---
name: pr-describe
description: Use when writing or reviewing the body of a pull request — after `gh pr create`, when filling in an empty PR description, or when asked whether a PR body says enough. Covers the body only; the PR title and commit messages are commit-write.
---

# PR Description

The PR body is the only part of a change that never reaches git. The title
becomes the commit subject on `main`; the body stays on GitHub for the reviewer
in front of you. Write it for that reviewer, and for nobody else.

House style, measured over 240 real PRs: **no headings, no checklists, no
template.** A third ship with no body at all. The median body is 234 characters.

## Does this PR need a body?

**No** — if the title already says everything. Mechanical, low-risk, or
self-evident work ships bodyless: version bumps, key rotations, CI timeouts,
formatting, a rename the title names. An empty body on a `chore:` is correct,
not lazy.

**Yes** — if any of these is true:

- The *why* is not obvious from the diff.
- The PR is part of a stack, or must merge in a particular order.
- A reviewer would otherwise trip on something — generated noise, a known gap,
  a deliberate stopgap.
- The change is visual.
- You want a specific thing looked at.

If none of those hold, leave it empty and say so.

## What a body is made of

Flat dash-bullets, or one short paragraph. Lowercase openers are normal. No
`## Why` / `## Approach` headings — the corpus has two headings in 157 bodies,
both one-offs. No `- [ ]` checklists; there are none in the corpus at all.

Write only the parts that apply, in whatever order fits:

| Part | When |
|---|---|
| What changed, as bullets | The diff spans several things worth naming |
| Why, in a sentence | The reason is not visible in the code |
| The tradeoff you took | You picked a stopgap, or a slower path, on purpose |
| Stack and merge order | Another PR must land first, or this one enables the next |
| Reviewer warnings | Generated noise, a big lockfile, a deliberately deferred gap |
| Evidence | A screenshot, a repro, a log link, one "tested by …" line |
| Open questions | A `TODO:` block, when the answer belongs to the reviewer |

## Stack and merge order

Coordinate stacks in the body — the corpus does this in prose, never with a
tool. Be explicit about direction and consequence:

> this must be merged AFTER #9566 or else <vendor> evals will break

> this PR exists separately from #9599 to isolate 3000+ lines of package-lock.json changes

> supports #9645 which creates a large schema based on a large enum

Warn about noise you knowingly introduced. Naming "2,996 lines of
`package-lock.json`" up front saves the reviewer the discovery.

## Evidence beats a test plan

There is no "Testing" section in this house. Evidence is pasted inline instead:

- Visual change → paste the screenshot. Raw `<img>` from GitHub's uploader.
- Bug fix → the repro, or the log/observability link that proves it.
- Otherwise → one line: `tested by checking that ephemeral web workers are up`.

## Tickets

Reference tickets as a bare `ENG-1234` or a plain Linear URL. **Do not use
`Closes #` / `Fixes #`** — the corpus has five across 240 PRs, and auto-close
fires on merge whether or not the ticket is actually done. Link sibling PRs as
`#9645`.

## Ground it — never invent the why

- Read `git diff main...HEAD` and `git log main..HEAD`.
- Read the ticket and the branch name.
- If the why is still not recoverable, ask the user one question. Do not
  fabricate a plausible reason.

## Reviewing a body

Report only real gaps:

1. Empty on a PR that needed a body — name which trigger above it hit.
2. Bullets that restate the diff file by file, adding nothing.
3. A stack dependency the body does not mention.
4. A visual change with no screenshot.
5. `Closes #` used, or headings and checklists that do not match the house.

A short body is not a finding. Most bodies here are short on purpose.

## Guardrails

- **Never invent a ticket, a link, or a test that was run.**
- **Do not pad.** Adding a heading to a two-bullet body makes it worse.
- **Do not restate the title.** The reviewer already read it.
- Follow the repo's git attribution rules. No generated-by footers.

Verbatim style anchors from the real corpus: [examples.md](examples.md).
