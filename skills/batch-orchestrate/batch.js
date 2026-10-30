export const meta = {
  name: 'batch-orchestrate',
  description: 'Build a batch of tickets in stacked lanes: one worktree agent per ticket, three review lenses, fix rounds; never pushes',
  whenToUse: 'Run by the batch-orchestrate skill, which plans the args. Args {lanes: [{name, tickets: [{id, slug, kind, model, effort, spec, deps?}]}], base?: "main", rules?: string, maxLanes?: 4}. kind is "fix" or "feat". model is "opus", "sonnet" or "haiku". spec is the filled spec-write template for that ticket. Each ticket builds on the branch of the previous ticket in its lane, so tickets in one lane build one after another and lanes run in parallel. deps name tickets in the same lane; a ticket whose dep failed is skipped. A build that fails is retried once one model tier up. Nothing is pushed: the lead verifies and publishes after the run.',
  phases: [
    { title: 'Build', detail: 'one worktree agent per ticket runs pr-orchestrate in worker mode' },
    { title: 'Review', detail: 'three lenses in parallel on Sonnet' },
    { title: 'Fix', detail: 'fix material findings, rerun Done-when, review again; at most 2 rounds' },
  ],
}

const A = args || {}
if (!A.lanes || !A.lanes.length) throw new Error('batch-orchestrate: args need lanes')
const BASE = A.base || 'main'
const MAX_LANES = A.maxLanes || 4
const MAX_ROUNDS = 2
const RULES = A.rules || ''
const UP = { haiku: 'sonnet', sonnet: 'opus', opus: 'opus' }

const DECISION = {
  type: 'object',
  properties: { question: { type: 'string' }, chosen: { type: 'string' }, why: { type: 'string' } },
  required: ['question', 'chosen', 'why'],
}
const SLICE = {
  type: 'object',
  properties: {
    status: { type: 'string', enum: ['ready', 'blocked', 'failed'] },
    reason: { type: 'string', description: 'why, when not ready; empty otherwise' },
    worktree: { type: 'string', description: 'absolute path of your worktree' },
    branch: { type: 'string' },
    commits: { type: 'array', items: { type: 'string' }, description: 'sha + subject' },
    decisions: { type: 'array', items: DECISION, description: 'every choice the ticket and spec did not settle' },
    unverified: { type: 'array', items: { type: 'string' }, description: 'what no automated check covers' },
    checks: { type: 'string', description: 'each Done-when command and its actual outcome' },
    pr_title: { type: 'string' },
    pr_body_path: { type: 'string', description: 'absolute path of the PR body file, outside the repo tree' },
  },
  required: ['status', 'reason', 'worktree', 'branch', 'commits', 'decisions', 'unverified', 'checks', 'pr_title', 'pr_body_path'],
}
const FINDINGS = {
  type: 'object',
  properties: {
    findings: {
      type: 'array',
      items: {
        type: 'object',
        properties: {
          severity: { type: 'string', enum: ['high', 'med', 'low'] },
          file: { type: 'string' },
          line: { type: 'integer' },
          problem: { type: 'string', description: 'concrete failure scenario' },
          fix: { type: 'string' },
          contradicts_ticket: { type: 'boolean', description: 'true if fixing it would go against the ticket or the spec' },
        },
        required: ['severity', 'file', 'line', 'problem', 'fix', 'contradicts_ticket'],
      },
    },
  },
  required: ['findings'],
}
const LENSES = [
  { key: 'correctness', prompt: 'CORRECTNESS and security: logic bugs, concurrency, edge cases, hostile input, whether the code does what its names say.' },
  { key: 'po-az', prompt: 'FAILURE MODES and craft (pr-review-po + pr-review-az): silent breakage, tests that pass on broken code, ignored errors, type integrity, naming.' },
  { key: 'sl', prompt: 'DESIGN and scope (pr-review-sl + ponytail-review): layering, placement, naming, scope creep, over-engineering, anything the ticket requires that is missing.' },
]

// Nobody can answer a worker mid-run, so these rules replace pr-orchestrate's human gates.
const unattended = t => `Unattended rules: nobody can answer you. Never ask; never wait.
- If the ticket or spec leaves a choice open, take the ticket's default, or else the most conservative option. Record it in "decisions".
- Commit with commit-write: \`${t.kind}(${t.id}): ...\`. Never add a Co-Authored-By trailer, a generated-by footer or any mention of Claude.
- Never push, open a PR, merge, or touch ${BASE}.${RULES ? `\nRepo batch rules:\n${RULES}` : ''}`

// A desktop-app session bound to its own worktree refuses Edit and Write outside it (Tickforge #287, #310).
const foreign = dir => `Work only in ${dir} (absolute paths or \`git -C\`). If Edit or Write refuses a path there, write it with Bash instead. Never return "blocked" for that refusal.`

const branchOf = t => `${t.kind}/${t.id}-${t.slug}`

function buildPrompt(t, base) {
  return `${t.spec}

Run the pr-orchestrate skill in worker mode on this spec. The caller reviews, so skip stages 2-3.
${unattended(t)}
- You are in a fresh git worktree. Run \`git switch -c ${branchOf(t)} ${base}\`. ${base} is a local branch; don't fetch it from origin.
- Write the PR title and body with commit-write and pr-describe. Put the body in a file outside the repo.
- If the change can't be built on ${base}, return status "blocked" with the reason.`
}

function reviewPrompt(t, r, base, lens) {
  return `Review branch ${r.branch} for ticket ${t.id}. Read \`git -C "${r.worktree}" diff ${base}...HEAD\`. Do not edit anything.
Lens: ${lens.prompt}
The spec it was built from:
${t.spec}
Decisions already taken: ${r.decisions.map(d => `${d.question} -> ${d.chosen}`).join('; ') || 'none'}. A finding that would go against the ticket or spec gets contradicts_ticket=true.
Report only real defects, each with file, line, a concrete failure scenario and a fix.`
}

function fixPrompt(t, r, material, forUser) {
  return `Continue ticket ${t.id} in the existing worktree ${r.worktree} on branch ${r.branch}. ${foreign(r.worktree)}
${unattended(t)}
Fix these findings, with tests where they apply. Then run every Done-when command from the spec below and fix what fails. If one still fails, return status "failed" with its trimmed output as the reason.
Update the PR body at ${r.pr_body_path}: add "Review found and fixed", and "Open for the user" with:
${forUser.map(f => `- [${f.severity}] ${f.file}:${f.line} ${f.problem}`).join('\n') || '- none'}
Findings to fix:
${material.map(f => `- [${f.severity}] ${f.file}:${f.line} ${f.problem} -> ${f.fix}`).join('\n')}
Spec:
${t.spec}
Return the full slice report, updated.`
}

const ready = r => r && r.status === 'ready'

async function build(t, base) {
  let model = t.model
  for (let attempt = 0; attempt < 2; attempt++) {
    const r = await agent(buildPrompt(t, base), {
      label: `build ${t.id} (${model})`, phase: 'Build', schema: SLICE, agentType: 'general-purpose',
      isolation: 'worktree', model, effort: t.effort,
    })
    if (ready(r) || (r && r.status === 'blocked')) return { r, model }
    log(`${t.id}: build ${r ? `${r.status}: ${r.reason}` : 'died'}; ${attempt ? 'giving up' : `retrying on ${UP[model]}`}`)
    model = UP[model]
  }
  return { r: null, model }
}

async function finish(t, r, base, model) {
  const forUser = []
  let open = []
  for (let round = 1; round <= MAX_ROUNDS; round++) {
    const found = (await parallel(LENSES.map(l => () => agent(reviewPrompt(t, r, base, l), {
      label: `review ${t.id} ${l.key} r${round}`, phase: 'Review', schema: FINDINGS, model: 'sonnet', effort: 'medium',
    })))).filter(Boolean).flatMap(x => x.findings)
    for (const f of found.filter(f => f.contradicts_ticket)) {
      if (!forUser.some(o => o.file === f.file && o.line === f.line && o.problem === f.problem)) forUser.push(f)
    }
    const material = found.filter(f => !f.contradicts_ticket && f.severity !== 'low')
    open = [...forUser, ...found.filter(f => !f.contradicts_ticket && f.severity === 'low')]
    log(`${t.id} round ${round}: ${material.length} material, ${forUser.length} for the user`)
    if (!material.length) return { ...r, status: 'ready', open }
    const fixed = await agent(fixPrompt(t, r, material, forUser), {
      label: `fix ${t.id} r${round}`, phase: 'Fix', schema: SLICE, agentType: 'general-purpose', model, effort: t.effort,
    })
    if (!ready(fixed)) return { ...(fixed || r), status: 'failed', reason: (fixed && fixed.reason) || 'a fix round died', open: material }
    r = fixed
  }
  return { ...r, status: 'review_exhausted', open }
}

const results = {}

async function runLane(lane) {
  let base = BASE
  for (const t of lane.tickets) {
    const bad = (t.deps || []).find(d => results[d] && !['ready', 'review_exhausted'].includes(results[d].status))
    if (bad) { results[t.id] = { status: 'skipped', reason: `depends on ${bad}, which did not build`, lane: lane.name }; continue }
    const { r, model } = await build(t, base)
    if (!ready(r)) { results[t.id] = { status: r ? r.status : 'failed', reason: r ? r.reason : 'build died twice', lane: lane.name, worktree: r && r.worktree }; continue }
    const done = await finish(t, r, base, model)
    results[t.id] = { ...done, lane: lane.name, base, model }
    if (['ready', 'review_exhausted'].includes(done.status)) base = done.branch
  }
}

phase('Build')
const pending = [...A.lanes]
const running = new Set()
while (pending.length || running.size) {
  while (running.size < MAX_LANES && pending.length) {
    const p = runLane(pending.shift()).finally(() => running.delete(p))
    running.add(p)
  }
  await Promise.race(running)
}
return { results }
