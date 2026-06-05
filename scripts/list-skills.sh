#!/usr/bin/env bash
set -euo pipefail

# Prints each skill's name and description, read from SKILL.md frontmatter.

REPO="$(cd "$(dirname "$0")/.." && pwd)"

find "$REPO/skills" -name SKILL.md -not -path '*/node_modules/*' -print0 |
sort -z |
while IFS= read -r -d '' skill_md; do
  name="$(basename "$(dirname "$skill_md")")"
  desc="$(awk '/^---/{c++; next} c==1 && /^description:/{sub(/^description:[[:space:]]*/,""); print; exit}' "$skill_md")"
  printf '%s\n  %s\n\n' "$name" "$desc"
done
