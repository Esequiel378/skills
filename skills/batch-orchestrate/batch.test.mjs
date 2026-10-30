// Runs batch.js against a fake agent(): node --test skills/batch-orchestrate/batch.test.mjs
import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'

const src = readFileSync(new URL('./batch.js', import.meta.url), 'utf8').replace('export const meta', 'const meta')
const AsyncFunction = (async () => {}).constructor
const script = new AsyncFunction('args', 'agent', 'parallel', 'log', 'phase', src)

async function run(args, respond) {
  const calls = []
  const agent = async (prompt, opts) => { calls.push({ prompt, ...opts }); return respond(prompt, opts, calls) }
  const parallel = fns => Promise.all(fns.map(f => f()))
  const out = await script(args, agent, parallel, () => {}, () => {})
  return { ...out, calls }
}

const slice = (branch, extra = {}) => ({
  status: 'ready', reason: '', worktree: `/wt/${branch}`, branch, commits: [], decisions: [],
  unverified: [], checks: 'ok', pr_title: 't', pr_body_path: '/tmp/b.md', ...extra,
})
const ticket = (id, extra = {}) => ({ id, slug: `s${id}`, kind: 'feat', model: 'haiku', effort: 'low', spec: `spec ${id}`, ...extra })
const clean = { findings: [] }
const branchIn = p => p.match(/git switch -c (\S+)/)[1]

test('stacks each ticket on the previous branch in its lane', async () => {
  const { results, calls } = await run({ lanes: [{ name: 'a', tickets: [ticket(1), ticket(2)] }] }, (p, o) =>
    o.phase === 'Build' ? slice(branchIn(p)) : clean)
  assert.equal(results[1].status, 'ready')
  assert.equal(results[2].base, 'feat/1-s1')
  const builds = calls.filter(c => c.phase === 'Build')
  assert.match(builds[1].prompt, /git switch -c feat\/2-s2 feat\/1-s1/)
  assert.ok(calls.filter(c => c.phase === 'Review').every(c => c.model === 'sonnet'))
})

test('retries a failed build one tier up, then skips its dependents', async () => {
  const { results, calls } = await run({ lanes: [{ name: 'a', tickets: [ticket(1), ticket(2, { deps: [1] })] }] }, (p, o) =>
    o.phase === 'Build' ? { ...slice(branchIn(p)), status: 'failed', reason: 'boom' } : clean)
  assert.deepEqual(calls.filter(c => c.phase === 'Build').map(c => c.model), ['haiku', 'sonnet'])
  assert.equal(results[1].status, 'failed')
  assert.equal(results[2].status, 'skipped')
})

test('fixes material findings and keeps ticket-contradicting ones for the user', async () => {
  let round = 0
  const { results, calls } = await run({ lanes: [{ name: 'a', tickets: [ticket(1, { model: 'sonnet' })] }] }, (p, o) => {
    if (o.phase === 'Build' || o.phase === 'Fix') return slice('feat/1-s1')
    if (o.label.endsWith('correctness r1')) {
      round++
      return { findings: [
        { severity: 'high', file: 'a.go', line: 1, problem: 'bug', fix: 'x', contradicts_ticket: false },
        { severity: 'med', file: 'b.go', line: 2, problem: 'scope', fix: 'y', contradicts_ticket: true },
      ] }
    }
    return clean
  })
  assert.equal(round, 1)
  assert.equal(calls.filter(c => c.phase === 'Fix').length, 1)
  assert.equal(calls.find(c => c.phase === 'Fix').model, 'sonnet')
  assert.equal(results[1].status, 'ready')
  assert.deepEqual(results[1].open.map(f => f.file), ['b.go'])
})
