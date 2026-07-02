---
name: pr-description
description: Use when creating, opening, or writing a pull request, pushing a branch for review, or running `gh pr create` — any time a PR description or summary is about to be written. Produces a description that explains why the change exists and what a reviewer should focus on, instead of a per-file changelog that restates the diff.
---

# PR Description

The diff already says **what** changed, line by line. A PR description that
re-lists the files is wasted space — the reviewer can read the diff. Your job is
everything the diff *can't* show: **why this change exists, why it's built this
way, and where a reviewer should look hardest.**

Closely related to [user-story](../user-story/SKILL.md): the story says why the
work matters *before* it's built; the PR says why the code is *right* now that it
is. If a story/ticket exists, pull its "so that" straight into the Why.

## Ground it first — never invent the why

Motivation you make up is worse than none. Recover it from real signal:

- `git diff main...HEAD` and `git log main..HEAD` — what actually changed and the
  commit story.
- The linked issue/ticket and the branch name — the original problem.
- If the *why* still isn't recoverable, **ask the user one question** ("what
  problem does this solve?") rather than fabricate a plausible-sounding reason.

## The body IS these four parts, in this order

1. **Why** — the problem or goal. What was broken, missing, or slow; who felt it;
   link the issue/story. One or two sentences. This is the part reviewers and
   future archaeologists actually need.
2. **Approach & why this way** — how you solved it, in prose, and why *this* way
   over the obvious alternative (the trade-off you made). Not a file list.
3. **Review guide** — where to look hardest: the subtle/risky bits, decisions
   worth a second opinion, anything non-obvious. Point at specific spots
   (`file.go:42`) *only* where the diff isn't self-explanatory.
4. **Risk & rollout** — what could break, how it was verified, how to roll back,
   what to watch after deploy (a metric, a flag). Omit a part only if it's
   genuinely N/A — say so in one line rather than padding.

A per-file "## Changes" section is the default trap — it looks thorough but only
echoes the diff. Fold anything worth saying about a file into **Review guide** as
a reviewer would care about it, not as a changelog entry.

## Template

```markdown
## Why
<the problem/goal; link the issue or story>

## Approach
<how, and why this way over the alternative>

## Review guide
<where to look hardest; risky/subtle bits; file:line for non-obvious spots>

## Risk & rollout
<what could break · how it was verified · rollback · what to watch>

Closes <TICKET>
```

## Good vs. bad

**Bad** (echoes the diff — reviewer learns nothing new):
> ## Changes
> - `payment.go`: wrapped Stripe call in a 3s timeout, added retry
> - `checkout.yaml`: added `payment_timeout_ms`, `payment_max_retries`
> - `metrics.go`: added `payment_retry_total` counter

**Good** (says what the diff can't):
> ## Why
> Checkout failed intermittently at peak (PROJ-812): the Stripe call had no
> bounded timeout, so slow upstream responses blocked handlers and cascaded.
>
> ## Approach
> Bound the call to 3s and retry twice with backoff, so transient slowness
> recovers instead of failing the checkout. Limits live in config, not code, so
> we can tune against real p99 without a redeploy.
>
> ## Review guide
> Worst case is 3s × 3 attempts ≈ 9s — confirm the caller's deadline tolerates
> that. Retries reuse the same Stripe idempotency key (`payment.go:58`) to avoid
> double-charging when a slow first attempt actually succeeded — worth a close look.
>
> ## Risk & rollout
> Config defaults ship in the diff, no manual step. Watch `payment_retry_total`;
> a sustained climb means 3s is too tight. Revert = drop the retry wrapper.

## After drafting

Show the description. If the user asked to actually open the PR (or push for
review), create it with `gh pr create` using this body. Follow the repo's git
attribution rules.
