---
name: improve-prompt
description: Use when the user hands you a rough, vague, or messy prompt and wants it sharpened before acting — triggers on "improve this prompt", "fix my prompt", "/improve-prompt", "clean up this request", or any raw instruction handed over for rewriting. Rewrites the prompt in the context of the current project, shows the rewrite, then executes it.
---

# Improve Prompt

Take a raw, ugly prompt, rewrite it into a sharp one **in the context of this project**, show the rewrite, then carry it out.

## Steps

1. **Read the raw prompt** — everything the user passed after the skill name.

2. **Ground it in the project.** Spend a *cheap* pass on the context that actually changes the rewrite — not a full audit:
   - What is this repo? (README, manifest, top-level dirs)
   - Which files/areas does the prompt touch? Name them.
   - Existing conventions or utilities the task should reuse.

   Stop as soon as you know enough to rewrite well. Don't boil the ocean.

3. **Rewrite the prompt.** A good rewrite makes the implicit explicit:
   - **Goal** — one sentence, unambiguous.
   - **Scope** — what's in, what's out.
   - **Concrete anchors** — real file paths, function names, commands found in step 2.
   - **Constraints & success criteria** — how we'll know it's done.

   Cut filler. Preserve the user's actual intent — improve it, don't redirect it.

4. **Show it.** Print exactly:
   ```
   Improved prompt:
   <the rewrite>

   Why: <one line — what was vague and what you pinned down>
   ```

5. **Log it.** Append one JSON line to `runs.jsonl` in this skill's own directory
   (create the file if missing) using the Bash tool:
   ```bash
   printf '%s\n' '{"date":"YYYY-MM-DD","raw":"...","improved":"...","note":"..."}' \
     >> "<this-skill-dir>/runs.jsonl"
   ```
   This feeds the weekly review — see [weekly-review.md](weekly-review.md).

6. **Execute** the improved prompt — proceed to actually do the work, following
   any other applicable skills.

## Guardrails

- If the raw prompt is already clear, say so and run it as-is (still log it).
- If grounding reveals the prompt rests on a wrong assumption (file doesn't
  exist, feature already present), surface that *before* executing — a better
  prompt isn't a wrong one made fluent.
- Never invent file paths or APIs in the rewrite. Anchors must be ones you
  verified in step 2.

## Getting better over time

Every run is logged to `runs.jsonl`. A weekly review reads that log, judges
which rewrites helped, and proposes edits to this very file. Setup +
instructions: [weekly-review.md](weekly-review.md).
