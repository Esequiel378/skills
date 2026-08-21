---
name: prompt-improve-review
description: Use to run the periodic analysis + tuning pass over the prompt-improve run log — triggers on "/prompt-improve-review", "review my prompt log", "tune prompt-improve", "prompt-improve stats", or any request to analyze accumulated prompt rewrites and improve the skill. Prints stats, then proposes one targeted edit to the prompt-improve skill.
---

# Improve Prompt — Review

Reproducible analysis + improvement pass for the [prompt-improve](../prompt-improve/SKILL.md) skill.

## Steps

1. **Stats.** Run the quantitative view and read it:
   ```bash
   bash ~/.claude/skills/prompt-improve/stats.sh
   ```
   Note the `since last review` count, the expansion ratio (creeping >2.5x = over-expansion), and skim the recent notes for a repeated weakness.

2. **Review & tune.** Read [../prompt-improve/review.md](../prompt-improve/review.md) and follow it exactly:
   - Sample raw→improved pairs.
   - Find the most common weakness.
   - Propose **one** targeted edit to `prompt-improve/SKILL.md`.
   - Apply it.
   - Append the `## Reviewed up to <date>` marker.
   - Commit as `prompt-improve: tune`.

If stats and notes show nothing wrong, append the marker and change nothing. One good improvement per session beats a redesign.
