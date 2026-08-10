---
name: improve-prompt-review
description: Use to run the periodic analysis + tuning pass over the improve-prompt run log — triggers on "/improve-prompt-review", "review my prompt log", "tune improve-prompt", "improve-prompt stats", or any request to analyze accumulated prompt rewrites and improve the skill. Prints stats, then proposes one targeted edit to the improve-prompt skill.
---

# Improve Prompt — Review

Reproducible analysis + improvement pass for the [improve-prompt](../improve-prompt/SKILL.md) skill.

## Steps

1. **Stats.** Run the quantitative view and read it:
   ```bash
   bash ~/.claude/skills/improve-prompt/stats.sh
   ```
   Note the `since last review` count, the expansion ratio (creeping >2.5x = over-expansion), and skim the recent notes for a repeated weakness.

2. **Review & tune.** Read [../improve-prompt/review.md](../improve-prompt/review.md) and follow it exactly:
   - Sample raw→improved pairs.
   - Find the most common weakness.
   - Propose **one** targeted edit to `improve-prompt/SKILL.md`.
   - Apply it.
   - Append the `## Reviewed up to <date>` marker.
   - Commit as `improve-prompt: tune`.

If stats and notes show nothing wrong, append the marker and change nothing. One good improvement per session beats a redesign.
