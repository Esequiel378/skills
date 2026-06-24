# Review — tune the improve-prompt skill

Goal: read the accumulated runs and make `SKILL.md` produce better rewrites.
This is the "get better over time" loop — run it whenever you see fit.

## What to do

1. Read `~/.claude/improve-prompt/runs.jsonl`. Focus on entries since the last
   review — previous reviews append a `## Reviewed up to <date>` marker to the
   Review log below; start after the latest one.

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
   Review log below, and commit titled `improve-prompt: tune`.

Keep it lazy: one good improvement per session beats a redesign. If the log
shows nothing wrong, append the marker and change nothing.

The log lives at `~/.claude/improve-prompt/runs.jsonl` (outside this repo, so
it's never committed). Run this review manually — no schedule.

## Review log

<!-- reviews append their markers below -->
