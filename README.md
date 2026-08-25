# Skills

My personal [Claude Code](https://claude.com/claude-code) agent skills, centralized in one repo so they're version-controlled, shareable, and edited in one place.

Each skill is a self-contained folder under [`skills/`](./skills) with a `SKILL.md` (plus any supporting files the skill needs). Procedural sections follow [`STYLE.md`](./STYLE.md).

## Skills

| Skill | What it does |
| --- | --- |
| [`architecture-hexagonal`](./skills/architecture-hexagonal/SKILL.md) | Authoring/reviewing Go or Python service code in ports-and-adapters style — use cases, driven ports, driving adapters, stage-fit and promotion decisions. |
| [`architecture-review`](./skills/architecture-review/SKILL.md) | Blunt principal-architect review of an architecture doc / RFC — severity-tagged findings across 7 dimensions and a stage-fit verdict. |
| [`commit-write`](./skills/commit-write/SKILL.md) | House convention for commit messages and PR titles — `type(scope): summary` where the scope is the ticket or the app, the repo-invented `prefactor` type, and the squash-merge fact that a PR title becomes the permanent commit subject. Distilled from 240 real PRs. |
| [`domain-modeling`](./skills/domain-modeling/SKILL.md) | Builds and sharpens a project's domain model as you design — challenges terms against a `CONTEXT.md` glossary, stress-tests concepts with edge-case scenarios, captures resolved terms inline, and offers ADRs only for hard-to-reverse decisions. |
| [`plan-grill`](./skills/plan-grill/SKILL.md) | Stress-tests a plan or design by interrogating decisions branch-by-branch before implementation. |
| [`pr-comments-address`](./skills/pr-comments-address/SKILL.md) | Triages a PR's unresolved review comments — accept / descope / move (stack-aware) / reject with evidence — cross-validates contested ones with `pr-review`, fixes accepted ones in the working tree. Never commits, pushes, or replies. |
| [`pr-describe`](./skills/pr-describe/SKILL.md) | Writes and reviews PR *bodies* in the team's measured house style — no headings, no checklists, empty when the title suffices — with stack/merge-order coordination and pasted evidence instead of a test plan. |
| [`pr-orchestrate`](./skills/pr-orchestrate/SKILL.md) | Sequences the other skills into one gated pipeline — ticket to review-ready PR: gates on an INVEST-complete ticket, drives TDD, runs `pr-review` in fix mode until clean, then prepares the PR. |
| [`pr-review`](./skills/pr-review/SKILL.md) | Fans out six review lenses in parallel over one diff — `code-review` ∥ `ponytail-review` ∥ `security-review` ∥ `pr-review-po` ∥ `pr-review-sl` ∥ `pr-review-az` — then reconciles them into a single ranked list. Reports by default; fixes only when asked. |
| [`pr-review-az`](./skills/pr-review-az/SKILL.md) | Craft-and-correctness review lens — type-system lies, error values over exceptions, naming/comment hygiene, schema back-compat — each finding taught with the rule behind it. Approves with follow-up tickets instead of blocking. |
| [`pr-review-po`](./skills/pr-review-po/SKILL.md) | Review lens that hunts silent failure modes, tests that prove nothing, type lies, and drift risks. Cites evidence, hedges honestly, keeps blockers rare and explicit. |
| [`pr-review-sl`](./skills/pr-review-sl/SKILL.md) | Terse staff-architect review lens — interrogates the design rather than the syntax: wrong layer, wrong name, wrong schema, scope creep. Socratic, minimum words. |
| [`prompt-improve`](./skills/prompt-improve/SKILL.md) | Rewrites a raw, vague prompt into a sharp one grounded in the current project, shows the rewrite, then executes it. Logs every run; a periodic review tunes the skill over time. |
| [`prompt-improve-review`](./skills/prompt-improve-review/SKILL.md) | Analyzes the `prompt-improve` run log, judges which rewrites helped, and proposes one targeted edit to that skill. |
| [`tdd`](./skills/tdd/SKILL.md) | Test-driven development with the red-green-refactor loop. |
| [`ticket-write`](./skills/ticket-write/SKILL.md) | Turns a rough idea into a great, INVEST-complete ticket in the team's house format — slices big asks into independently-shippable vertical slices, then offers to create it in the connected tracker. |

## Install

### Local (symlink — recommended for the author)

Symlinks each skill into `~/.claude/skills/`, so the local Claude CLI picks them up and editing a skill here is instantly live:

```bash
make install          # or: ./scripts/link-skills.sh
```

Safe to re-run. It also prunes symlinks left behind by renamed or deleted skills, and migrates any `~/.claude/` data directory whose skill was renamed — so pulling and running `make install` is all another machine needs.

### As a Claude Code plugin

This repo ships a [`.claude-plugin/plugin.json`](./.claude-plugin/plugin.json), so it can be installed as a plugin (e.g. via a marketplace entry pointing at this repo, or by cloning and adding it as a local plugin).

## Utilities

```bash
./scripts/list-skills.sh   # print each skill's name + description
```

## License

[MIT](./LICENSE)
