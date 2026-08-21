#!/usr/bin/env bash
set -euo pipefail

# Links every skill in this repo into ~/.claude/skills as a symlink, so the
# local Claude CLI picks them up and editing a skill here is instantly live.

REPO="$(cd "$(dirname "$0")/.." && pwd)"
DEST="$HOME/.claude/skills"

# If ~/.claude/skills is itself a symlink into this repo, linking would write
# the per-skill symlinks back into the repo. Bail out instead of polluting it.
if [ -L "$DEST" ]; then
  resolved="$(readlink "$DEST")"
  case "$resolved" in
    "$REPO"|"$REPO"/*)
      echo "error: $DEST is a symlink into this repo ($resolved)." >&2
      echo "Remove it (rm \"$DEST\") and re-run; it will be recreated as a real dir." >&2
      exit 1
      ;;
  esac
fi

mkdir -p "$DEST"

# --- migrations -------------------------------------------------------------
# Skills were renamed to a {entity}-{action} scheme. A machine that installed
# the old names keeps working symlinks to directories that no longer exist, and
# a data directory named after the old skill. Fix both here so `make install`
# is all any other machine needs to run.

# Data directories keyed by a skill name that changed.
migrate_data_dir() {
  old="$HOME/.claude/$1"
  new="$HOME/.claude/$2"
  [ -d "$old" ] || return 0
  if [ -d "$new" ]; then
    echo "warn: both $old and $new exist; leaving both, merge them by hand" >&2
    return 0
  fi
  mv "$old" "$new"
  echo "migrated $old -> $new"
}

migrate_data_dir improve-prompt prompt-improve

# Symlinks pointing into this repo at a path that no longer exists are skills
# that were renamed or deleted. Anything else in $DEST is left alone.
for link in "$DEST"/*; do
  [ -L "$link" ] || continue
  target="$(readlink "$link")"
  case "$target" in "$REPO"/*) ;; *) continue ;; esac
  if [ ! -e "$target" ]; then
    rm "$link"
    echo "pruned stale $(basename "$link") -> $target"
  fi
done
# --- end migrations ---------------------------------------------------------

find "$REPO/skills" -name SKILL.md -not -path '*/node_modules/*' -print0 |
while IFS= read -r -d '' skill_md; do
  src="$(dirname "$skill_md")"
  name="$(basename "$src")"
  target="$DEST/$name"

  # Replace a pre-existing real directory with a symlink to the repo copy.
  if [ -e "$target" ] && [ ! -L "$target" ]; then
    rm -rf "$target"
  fi

  ln -sfn "$src" "$target"
  echo "linked $name -> $src"
done
