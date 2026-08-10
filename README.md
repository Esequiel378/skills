# Skills

My personal [Claude Code](https://claude.com/claude-code) agent skills, centralized in one repo so they're version-controlled, shareable, and edited in one place.

Each skill is a self-contained folder under [`skills/`](./skills) with a `SKILL.md` (plus any supporting files the skill needs). Procedural sections follow [`STYLE.md`](./STYLE.md).

## Skills

| Skill | What it does |
| --- | --- |
| [`architecture-review`](./skills/architecture-review/SKILL.md) | Blunt principal-architect review of an architecture doc / RFC — severity-tagged findings across 7 dimensions and a stage-fit verdict. |
| [`grill-me`](./skills/grill-me/SKILL.md) | Stress-tests a plan or design by interrogating decisions branch-by-branch before implementation. |
| [`hexagonal-architecture`](./skills/hexagonal-architecture/SKILL.md) | Authoring/reviewing Go or Python service code in ports-and-adapters style — use cases, driven ports, driving adapters, stage-fit and promotion decisions. |
| [`improve-prompt`](./skills/improve-prompt/SKILL.md) | Rewrites a raw, vague prompt into a sharp one grounded in the current project, shows the rewrite, then executes it. Logs every run; a weekly review tunes the skill over time. |
| [`orchestrate`](./skills/orchestrate/SKILL.md) | Sequences the other skills into one gated pipeline — ticket to review-ready PR: gates on a user-story ticket, drives TDD, runs `code-review` ∥ `ponytail-review` ∥ `security-review` in parallel, loops until clean, then prepares the PR. |
| [`pr-description`](./skills/pr-description/SKILL.md) | Writes PR descriptions that explain *why* the change exists and where a reviewer should look — grounded in the diff, log, and linked ticket — instead of a per-file changelog that restates the diff. Fires on "create/open a PR" or "push for review". |
| [`tdd`](./skills/tdd/SKILL.md) | Test-driven development with the red-green-refactor loop. |
| [`user-story`](./skills/user-story/SKILL.md) | Turns a rough ticket idea into a great, INVEST-complete user story in the team's house format — slices big asks into independently-shippable vertical slices, then offers to create it in the connected tracker. |

## Install

### Local (symlink — recommended for the author)

Symlinks each skill into `~/.claude/skills/`, so the local Claude CLI picks them up and editing a skill here is instantly live:

```bash
./scripts/link-skills.sh
```

### As a Claude Code plugin

This repo ships a [`.claude-plugin/plugin.json`](./.claude-plugin/plugin.json), so it can be installed as a plugin (e.g. via a marketplace entry pointing at this repo, or by cloning and adding it as a local plugin).

## Utilities

```bash
./scripts/list-skills.sh   # print each skill's name + description
```

## License

[MIT](./LICENSE)
