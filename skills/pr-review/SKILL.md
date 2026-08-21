---
name: pr-review
description: Use when the user wants a diff or PR reviewed from every angle at once — triggers on "/pr-review", "review this PR", "review my changes", "full review", or handing over a PR URL/number and asking what's wrong with it. Fans out six review lenses in parallel, reconciles their findings into one ranked list, and reports. Fixes only when asked.
---

# PR Review

One command, six lenses, one ranked list. This skill doesn't review by itself —
it *fans out* the repo's review skills in parallel over the same diff, then
reconciles what comes back into a single verdict the user can act on.

Fans out: `code-review` ∥ `ponytail:ponytail-review` ∥ `security-review` ∥
[pr-review-po](../pr-review-po/SKILL.md) ∥ [pr-review-sl](../pr-review-sl/SKILL.md) ∥
[pr-review-az](../pr-review-az/SKILL.md).

## 0. Resolve the target, and read it

Establish exactly what diff is under review before dispatching anything:

- **PR number / URL** → `gh pr view <n> --json …` for metadata, `gh pr diff <n>`
  for the patch. Also pull existing review comments (`gh api
  repos/<o>/<r>/pulls/<n>/comments`) — a lens should not re-raise what a human
  already raised.
- **Branch** → `git diff <base>...HEAD`, base from the PR or the repo default.
- **Nothing given** → the working diff (`git diff HEAD` plus staged).

Two checks that are cheap here and expensive later:

- **Is the local branch what the PR actually shows?** If a PR was named, compare
  its head SHA to local `HEAD`. A stale or diverged local branch means the lenses
  review code the reviewer never saw — say so before proceeding.
- **Is the changed code reachable?** Trace at least one entry point. Code behind a
  route that still points elsewhere, a flag that's off, or an unexported symbol
  changes every finding's severity from "bug" to "not yet live".

## 1. Fan out — six lenses, in parallel

Dispatch **parallel sub-agents** (superpowers `dispatching-parallel-agents`), one
per lens, each reading the same resolved diff. They don't mutate, so this is safe
to run wide.

| Lens | Catches |
|---|---|
| `code-review` | correctness, bugs, reuse, simplification |
| `ponytail:ponytail-review` | over-engineering — what to delete or replace with stdlib/native |
| `security-review` | vulnerabilities in the pending changes |
| `pr-review-po` | silent failure modes, tests that prove nothing, type lies, drift |
| `pr-review-sl` | wrong layer, wrong name, wrong schema, scope creep |
| `pr-review-az` | type-system lies, undefined smuggling, exceptions as control flow, comment/naming hygiene, schema back-compat |

The last three are not stylistic garnish — they cover classes the mechanical
lenses systematically miss. `code-review` will pass a test suite that asserts a
throw no client can observe; `pr-review-po` is the lens that catches it. `code-review`
will pass a correct function living in the wrong package; `pr-review-sl` catches
that. `code-review` will pass an event rename that breaks every already-persisted
event; `pr-review-az` catches that.

**Collect everything before acting on any of it.** A ponytail "delete this" and a
code-review "fix this" routinely target the same lines.

## 2. Reconcile — one list, not five

Merge the raw findings into a single ranked list:

- **Dedupe** — the same defect surfaces in different words from different lenses.
  Merge them and note it was found more than once; that is a confidence signal.
- **Resolve conflicts** — when two lenses disagree (add a guard vs. delete the
  branch), decide, and say which lens you sided with and why. Never emit both.
- **Verify before reporting.** Every finding names a file, a line, and a concrete
  failure scenario — inputs → wrong output. A finding you can't make concrete is a
  vibe; drop it. Check each claim against the actual code, not the diff alone: the
  surrounding lines a hunk doesn't show are where false positives come from.
- **Rank by severity**, not by lens. Blockers first, then material, then nits.
- **Drop what a human already said** in the PR's existing comments, unless the
  point is still unaddressed and worth reinforcing.

## 3. Report

```
## review — <target> [APPROVE | COMMENTS | BLOCKERS]
<one line: what this change does, and whether it's safe to merge>

### Blockers
- <file>:<line> — <defect> → <failure scenario> [lens]

### Material
- ...

### Nits
- ...
```

Say plainly when nothing material came back. A clean review reported as clean is
more useful than a manufactured finding.

## 4. Fix — only when asked

Reporting is the default. Fix only on an explicit `--fix` / "fix them" / "apply
these", and then:

- Fix, keep tests green (re-run the suite — not "should still pass").
- Re-run the lenses on the updated diff.
- Repeat until nothing material comes back.
- **Thrash guard:** if a finding survives two fix attempts, stop and surface it.
  Oscillating on a judgment call is the user's decision to make, not yours.

## Guardrails

- **Never post and never push.** No `gh pr comment`, no `gh pr review`, no
  inline comments, no commits, no push — unless the user asks for that
  specifically, in that message. Reviewing a PR is not permission to write to it.
- **Verify before claiming.** Run the tests and the typechecker on the code as it
  actually is; report what the output said. "Tests pass" without having run them
  is the one failure mode that discredits the whole review.
- **Report what you skipped.** A lens that errored, a file too large to read, a
  check you couldn't run — say so. Silent partial coverage reads as full coverage.
- **Severity is not politeness.** A blocker stays a blocker when the author is in
  a hurry; a nit stays a nit when you found it interesting.
