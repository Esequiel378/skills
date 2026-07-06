---
name: orchestrate
description: Use when the user wants to take a ticket all the way to a review-ready PR in one shot — triggers on "/orchestrate", "orchestrate this", "ticket to PR", "run the full workflow", or handing over a ticket/issue URL and asking for the finished PR. Gates on a user-story-shaped ticket, drives TDD, runs the reviews in parallel, loops until clean, then prepares the PR.
---

# Orchestrate

One command from a ticket to a PR the user can review. This skill doesn't do the
work itself — it *sequences* the repo's other skills into a gated pipeline and
runs the two reviews in parallel. The end state is a described PR, staged and
ready, waiting only on the user's go-ahead to open.

Chains: [user-story](../user-story/SKILL.md) → [tdd](../tdd/SKILL.md) →
`code-review` ∥ `ponytail:ponytail-review` ∥ `security-review` →
[pr-description](../pr-description/SKILL.md).

## The one rule: no blind starts

**The workflow cannot begin while anything material is unknown.** Before writing a
single line of implementation code, surface *every* open question in one batched
message and stop. Guessing here is what makes the same prompt yield different
results each time — this skill exists to kill that variance. A wrong assumption
made fluent is still wrong.

## Pipeline

### 0. Gate — resolve the ticket, or refuse

The input must resolve to a ticket that follows the
[user-story](../user-story/SKILL.md) structure (INVEST, testable acceptance
criteria, real anchors).

- **Ticket URL / existing ticket** → read it. If it's already user-story-shaped
  with clear acceptance criteria, proceed. If it's thin (no criteria, vague
  scope), run **user-story** to sharpen it and confirm with the user before
  continuing.
- **Just an idea, no ticket** → run **user-story** to draft one, get the user's
  approval, and (if they want) create it in the tracker. No ticket, no start.

Then flag blockers as a single batched question set and **block until answered**:
- Unknowns in scope, acceptance criteria, or affected surfaces.
- Missing anchors — a path/endpoint/key the ticket assumes but you can't verify.
- Branch: are we on a fresh branch off `main`? If on `main`, create one first.
  Prefer an isolated **git worktree** when the environment supports it (run the
  `superpowers:using-git-worktrees` skill) — it keeps the ticket's work off the
  current workspace. Fall back to a plain fresh branch if worktrees aren't
  available.

Do not proceed to stage 1 until the ticket is user-story-complete and every
question is answered.

### 1. Implement — TDD

Drive the ticket through [tdd](../tdd/SKILL.md): red-green-refactor, one acceptance
criterion at a time. Don't skip the failing-test-first step.

### 2. Review — three lenses, in parallel

Dispatch **parallel sub-agents** (superpowers `dispatching-parallel-agents`) — one
per lens, all reading the same working diff (`git diff main...HEAD`):

- `code-review` — correctness, bugs, reuse/simplification.
- `ponytail:ponytail-review` — over-engineering: what to delete or replace with
  stdlib/native.
- `security-review` — vulnerabilities in the pending changes.

Collect all findings before acting on any — a ponytail "delete this" and a
code-review "fix this" can target the same code; reconcile them together.

### 3. Loop — fix, re-review, until clean

Address the findings (auto-fix), then **re-run the three reviews on the updated
diff**. Repeat until all three come back with nothing material. Each round:
- Fix, keeping tests green (re-run the suite).
- Re-review in parallel.
- Stop when clean, or when a finding is a genuine judgment call the user should
  make — pause and ask rather than loop forever on a subjective point.

Guardrail against thrash: if a finding survives two fix attempts, stop fixing it
and surface it to the user instead of oscillating.

### 4. PR — describe and stage, then confirm

Write the PR body via [pr-description](../pr-description/SKILL.md) — the *why*, the
approach, the review guide, risk/rollout. `Closes <TICKET>`.

Then **prepare, don't auto-push**: ensure the branch is committed, show the user
the final PR description and the exact `gh pr create` command that would open it.
Hand off there. The user opens it. Honor the repo's git attribution rules — no
self-attribution in commits, PR title, or body.

## End state

A committed branch and a ready PR description, with the reviews already run and
their findings resolved — so the user's review is the *first human* review, not a
cleanup pass.

## Guardrails

- **Never start on assumptions.** Stage 0's gate is the whole point. If in doubt,
  ask.
- **Parallel only where it's safe.** The three reviews read the same diff and don't
  mutate — safe to fan out. Fixes are sequential (they touch the same tree).
- **Don't gold-plate the loop.** Clean means no *material* findings, not zero
  nitpicks. Ship when it's right, not when it's exhaustively polished.
- **The user opens the PR.** This skill stops at "ready to create", never pushes a
  PR without the final go-ahead.
