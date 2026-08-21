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
type(component): summary
```

The PR title carries the ticket; branch commits do not need it:

```
type(component): summary - eng-1234
```

Never type the `(#1234)` suffix. GitHub appends it at squash time.

| Part | Rule |
|---|---|
| `type` | Required. From the table below. Lowercase. |
| `(component)` | Optional. The literal directory, package, or tool touched. |
| `summary` | Lowercase, imperative, no trailing period. Aim for 50 characters. |
| `- eng-1234` | PR title only. Lowercase. Omit when no ticket exists. |

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

The scope names **the code**, not the ticket and not a layer. Use the literal
directory, package, or tool: `(py)`, `(cdktf)`, `(nix)`, `(py-stores)`,
`(lib-agent)`, `(inspector)`.

- If the change spans several components, omit the scope. Do not invent a bucket.
- If the summary already names the component, omit the scope. Never say it twice.

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
3. Scope names real code, or is absent.
4. Summary is lowercase, imperative, no trailing period, names something concrete.
5. PR title carries the ticket when one exists.
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
