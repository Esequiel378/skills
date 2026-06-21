# Weekly review — tune the improve-prompt skill

Goal: read the accumulated runs and make `SKILL.md` produce better rewrites
next week. This is the "get better over time" loop.

## What to do

1. Read `runs.jsonl` in this directory. Focus on entries since the last
   review — previous reviews append a `## Reviewed up to <date>` marker at
   the bottom of this file; start after the latest one.

2. For a sample of recent `raw` → `improved` pairs, ask:
   - Did the rewrite pin down the *right* things, or just add ceremony?
   - Recurring weakness? (e.g. never cites file paths, over-expands tiny
     asks, drops the user's constraints, redirects intent.)
   - Re-improve 1–2 of the raw prompts with the *current* `SKILL.md` and
     compare — has it drifted, improved, or regressed?

3. Propose **one** concrete, targeted edit to `SKILL.md` that fixes the most
   common weakness you found — a small change to Steps or Guardrails, not a
   rewrite.

4. Apply it, append `## Reviewed up to <date> — <one-line summary>` to the
   bottom of this file, and commit (or open a PR) titled
   `improve-prompt: weekly tune`.

Keep it lazy: one good improvement per week beats a redesign. If the log
shows nothing wrong, append the marker and change nothing.

## Schedule setup (run once)

Register a weekly routine with the `/schedule` skill, pointed at this repo:

> /schedule weekly Monday 9am — Follow skills/improve-prompt/weekly-review.md to tune the improve-prompt skill.

The routine reviews what's been **pushed** to the repo, so commit and push
`runs.jsonl` for it to have data to learn from.

<!-- ponytail: cloud routine learns from committed history. If you'd rather
     keep the run log out of git, point the routine at an external store
     (gist, object bucket) instead and update step 1's path. -->

## Review log

<!-- weekly reviews append their markers below -->
