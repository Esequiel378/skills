# Style

Skill files are read by an agent, not a person. The failure mode isn't
misunderstanding — it's dropping a clause from an overloaded sentence, or
silently picking one meaning when a term drifts. These rules are the procedural
half of [ASD-STE100](https://www.asd-ste100.org/) (Simplified Technical English),
which exists for exactly that failure mode.

## Scope

Rules 1–4 apply to **procedural** text only: numbered steps, pipeline stages,
checklists, guardrails — anything the agent executes.

**Descriptive** text (rationale, "what good looks like", worked examples) is
exempt. Vivid phrasing there is steering signal, not noise. Keep it.

Good: `Run story-write to sharpen it. Block on the user's confirmation before
continuing.`

Bad: `Write the PR body via pr-describe, then ensure the branch is committed
and show the user the command that would open it.`

## Rules

1. **One instruction per sentence.** Imperative, active voice. Two verbs joined
   by "and" or ";" is two sentences.
2. **≤20 words per procedural sentence.** Descriptive prose gets 25, and nobody
   counts.
3. **One term, one meaning.** Use the glossary below. Never reach for a synonym
   for variety.
4. **Condition before instruction.** "If on `main`, create a branch" — not
   "create a branch, if on `main`". The agent acts on the first clause it reads.

Explicitly **not** adopted from STE: the approved-vocabulary dictionary, the
`-ing` ban, the ban on metaphor. They serve a human reader in a second language,
which is not the reader here.

## Skill names

`{entity}-{action}`, optionally `-{variant}`. The entity comes first so related
skills share a prefix and sort together: `pr-review`, `pr-review-po`,
`pr-describe`, `pr-comments-address`. The action is a bare verb — *review*,
*describe*, *address*, *improve*, *write* — never a noun (`pr-describe`, not
`pr-description`).

The directory name, the `name:` frontmatter field, and the slash command are the
same string. A skill named after an established methodology (`tdd`) keeps that
name; there is no entity to put in front of it.

Renaming a skill breaks four things beyond the directory: the `name:`
frontmatter, any `"/old-name"` trigger in a `description:`, cross-reference
links in other skills, and `.claude-plugin/plugin.json`. If the skill owns a
data directory under `~/.claude/`, add a `migrate_data_dir` line to
`scripts/link-skills.sh` so other machines fix themselves on `make install`.

## Verb glossary

| Verb | Means |
| --- | --- |
| **run** `<skill>` | Invoke another skill |
| **dispatch** | Spawn a subagent |
| **read** `<file>` | Load a file into context |
| **ask** | Put a question to the user and continue |
| **block** | Stop and wait for the user's answer |
| **show** | Print output to the user |
| **name** | Cite a concrete path, section, or line |

Not: *fire*, *drive through*, *chain*, *load* (for skills), *follow ... exactly*
(for files — that's **read**).

## Noun glossary

| Noun | Means |
| --- | --- |
| **ticket** | The unit of tracked work, in any tracker |
| **diff** | The change under review |
| **finding** | One defect or objection a review produces |

A tracked work item is a **ticket** — never a *story*, *user story*, or *issue*.
`story-write` keeps *User story* and *Technical story* only where they name the
two ticket **types**; everywhere else the noun is *ticket*.
