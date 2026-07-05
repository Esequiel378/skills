#!/usr/bin/env bash
# Quantitative view over the improve-prompt run log. The qualitative judgment
# (which rewrites actually helped) is the reviewer's job — see review.md.
set -euo pipefail

LOG="${IMPROVE_PROMPT_LOG:-$HOME/.claude/improve-prompt/runs.jsonl}"
REVIEW="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/review.md"

[ -f "$LOG" ] || { echo "no log at $LOG"; exit 0; }

last_review=$(grep -oE '## Reviewed up to [0-9-]+' "$REVIEW" 2>/dev/null | tail -1 | grep -oE '[0-9-]+$' || true)

jq -rs --arg last "$last_review" '
  (map(select(.date > $last)) | length) as $new
  | "runs total:        \(length)",
    "date range:        \(min_by(.date).date) .. \(max_by(.date).date)",
    "since last review: \($new)   (marker: \(if $last=="" then "none" else $last end))",
    "avg expansion:     \(((map(.improved|length)|add) / (map(.raw|length)|add) * 10 | round / 10))x  (\((map(.raw|length)|add/length|round))->\((map(.improved|length)|add/length|round)) chars)",
    "shrank (improved<raw): \(map(select((.improved|length) < (.raw|length))) | length)",
    "",
    "recent notes (newest 8) — scan for recurring weakness:",
    (sort_by(.date) | reverse | .[:8][] | "  \(.date)  \(.note[:110])")
' "$LOG"
