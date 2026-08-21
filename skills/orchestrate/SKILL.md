---
name: orchestrate
description: Use when the user wants to take a ticket all the way to a review-ready PR in one shot — triggers on "/orchestrate", "orchestrate this", "ticket to PR", "run the full workflow", or handing over a ticket/issue URL and asking for the finished PR. Gates on a user-story-shaped ticket, drives TDD, runs pr-review until clean, then prepares the PR.
---

# Orchestrate

One command from a ticket to a PR the user can review. This skill doesn't do the
work itself — it *sequences* the repo's other skills into a gated pipeline. The
end state is a described PR, staged and ready, waiting only on the user's
go-ahead to open.

Chains: [user-story](../user-story/SKILL.md) → [tdd](../tdd/SKILL.md) →
[pr-review](../pr-review/SKILL.md) → [pr-description](../pr-description/SKILL.md).

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
  scope), run **user-story** to sharpen it. Block on the user's confirmation
  before continuing.
- **Just an idea, no ticket** → run **user-story** to draft one. Block on the
  user's approval. If they want it tracked, create it in the tracker. No ticket,
  no start.

Then flag blockers as a single batched question set and **block until answered**:
- Unknowns in scope, acceptance criteria, or affected surfaces.
- Missing anchors — a path/endpoint/key the ticket assumes but you can't verify.
- Branch: confirm we are on a fresh branch off `main`. If on `main`, create one
  first. If the environment supports worktrees, run
  `superpowers:using-git-worktrees` to keep the ticket's work off the current
  workspace. Otherwise create a plain fresh branch.

Do not proceed to stage 1 until the ticket is user-story-complete and every
question is answered.

### 1. Implement — TDD

Run [tdd](../tdd/SKILL.md) on the ticket: red-green-refactor, one acceptance
criterion at a time. Don't skip the failing-test-first step.

### 2-3. Review and loop — hand off to pr-review

Run [pr-review](../pr-review/SKILL.md) over the working diff (`git diff
main...HEAD`) **in fix mode**. That skill owns the whole review stage: the
six-lens parallel fan-out, reconciling the findings into one ranked list, and
the fix-and-re-review loop with its thrash guard.

Orchestrate adds one thing on top: since the code was just written here, a
finding is not automatically a defect in someone else's work — it may be a gap in
the acceptance criteria. When a finding contradicts the ticket, that's a stage-0
question resurfacing, not a fix. Stop and ask.

Exit this stage when pr-review reports nothing material, or when it surfaces a
judgment call — pause and ask rather than looping on a subjective point.

### 4. PR — describe and stage, then confirm

Run [pr-description](../pr-description/SKILL.md) to write the PR body — the *why*,
the approach, the review guide, risk/rollout. `Closes <TICKET>`.

Then **prepare, don't auto-push**. Commit the branch. Show the user the final PR
description. Show the exact `gh pr create` command that would open it. Hand off
there. The user opens it. Honor the repo's git attribution rules — no
self-attribution in commits, PR title, or body.

## End state

A committed branch and a ready PR description, with the reviews already run and
their findings resolved — so the user's review is the *first human* review, not a
cleanup pass.

## Guardrails

- **Never start on assumptions.** Stage 0's gate is the whole point. If in doubt,
  ask.
- **Don't gold-plate the loop.** Clean means no *material* findings, not zero
  nitpicks. Ship when it's right, not when it's exhaustively polished.
- **The user opens the PR.** This skill stops at "ready to create", never pushes a
  PR without the final go-ahead.
