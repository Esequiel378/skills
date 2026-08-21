---
name: commit-write
description: Use when writing a commit message or a PR title — before `git commit`, before `gh pr create`, when renaming an existing PR, or when checking whether a commit or PR title already written matches the house convention.
---

# Commit Write

The repo squash-merges. A PR title does not stay on the PR — it lands on `main`
as the commit subject, `<title> (#1234)`, and the branch's commit subjects become
that commit's body. Both are permanent `git log`. Write them for whoever runs
`git blame` in two years, not for the reviewer reading the PR today.

The PR **body** is the one part that never reaches git — [pr-describe](../pr-describe/SKILL.md)
owns that.

## The format

```
type(scope): summary
```

The scope is the ticket when the work has one, and the app or package when it
does not:

```
feat(eng-5526): generic render mustache template
feat(inspector): make event table collapsable
```

A commit and its PR title are the same string. Never type the `(#1234)` suffix —
GitHub appends it at squash time.

| Part | Rule |
|---|---|
| `type` | Required. From the table below. Lowercase. |
| `(scope)` | The ticket, else the app or package. See below. |
| `summary` | Lowercase, imperative, no trailing period. Aim for 50 characters. |

## Types

| Type | Use for |
|---|---|
| `feat` | New behavior a user or caller can observe |
| `fix` | Corrects behavior that was wrong |
| `prefactor` | A refactor landed *ahead of* the feature it enables, to keep that feature's diff small |
| `refactor` | Restructuring with no behavior change, not preparing anything specific |
| `chore` | Keys, bumps, config, housekeeping |
| `test` | Tests only |
| `ci` / `build` | Pipeline and build tooling |
| `docs` | Documentation only |
| `revert` | Backs out an earlier change |

`prefactor` is this repo's own type, not stock Conventional Commits. Reach for it
when the honest description is "this makes the next PR reviewable". If the change
prepares nothing in particular, it is a `refactor`.

## Scope

One slot, two cases. Check for a ticket first.

- **The work has a ticket** → the scope is the ticket, lowercase:
  `feat(eng-5526):`, `fix(eng-6276):`.
- **No ticket** → the scope is the app or package: `feat(inspector):`,
  `chore(openapi):`. Match the product's own casing — `feat(GuidingCare):`, not
  `feat(guidingcare):`.
- **Neither applies** → omit the scope. Do not invent a bucket.

Never put the ticket in trailing parentheses — `... (ENG-5526)` collides with
GitHub's `(#1234)`. The ticket goes in the scope or nowhere.

## Summary line

Name what changed, in the imperative, at the altitude a reader cares about.

- Good: `narrow nexus sdk exports`, `use 303 for login client cross-domain redirect`
- Bad: `Updated exports.` — past tense, capitalized, trailing period.
- Bad: `fix bug` — names nothing.
- Bad: `fix(py): fix filename access` — the type already said "fix".

One change, one commit. If the summary needs "and", consider two commits.

## Markers

Prefix a PR that must not merge with `NOMERGE` or `WIP`, in place of the type:
`NOMERGE(vendor) get-token script for dev`. Remove the marker and restore a real
type before the PR is ready.

## Reviewing an existing commit or title

Check it against this list, and report only what fails:

1. Type present, and the right one — a behavior change filed as `chore` hides it
   from anyone reading the log.
2. `prefactor` vs `refactor` used honestly.
3. Scope is the ticket when one exists, otherwise a real app or package.
4. Summary is lowercase, imperative, no trailing period, names something concrete.
5. No ticket in trailing parentheses.
6. No `NOMERGE` / `WIP` marker left on a PR that is ready.

A title that merely reads awkwardly is not a finding. Say so and move on.

## Guardrails

- **Never invent a ticket ID.** If no ticket is named, omit it and say you did.
- **Never invent the why.** Recover it from the diff, the branch name, or the
  ticket. If it is still not recoverable, ask the user one question.
- **Do not restate the diff.** A subject listing every file is a changelog.
- **Do not add attribution trailers.** No `Co-authored-by` for the agent, no
  generated-by footers. Follow the repo's git attribution rules.

Verbatim style anchors from the real corpus: [examples.md](examples.md).
