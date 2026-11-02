#!/usr/bin/env bash
set -euo pipefail

# Links every skill in this repo into the global skill dirs of Claude Code
# (~/.claude/skills) and Codex (~/.agents/skills) as symlinks, so both CLIs
# pick them up and editing a skill here is instantly live.

REPO="$(cd "$(dirname "$0")/.." && pwd)"
DESTS=("$HOME/.claude/skills" "$HOME/.agents/skills")

# If a skill dir is itself a symlink into this repo, linking would write the
# per-skill symlinks back into the repo. Bail out instead of polluting it.
for DEST in "${DESTS[@]}"; do
  [ -L "$DEST" ] || continue
  resolved="$(readlink "$DEST")"
  case "$resolved" in
    "$REPO"|"$REPO"/*)
      echo "error: $DEST is a symlink into this repo ($resolved)." >&2
      echo "Remove it (rm \"$DEST\") and re-run; it will be recreated as a real dir." >&2
      exit 1
      ;;
  esac
done

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
# that were renamed or deleted. Anything else in the skill dir is left alone.
prune_stale() {
  for link in "$1"/*; do
    [ -L "$link" ] || continue
    target="$(readlink "$link")"
    case "$target" in "$REPO"/*) ;; *) continue ;; esac
    if [ ! -e "$target" ]; then
      rm "$link"
      echo "pruned stale $(basename "$link") -> $target"
    fi
  done
}
# --- end migrations ---------------------------------------------------------

# Git Bash's plain `ln -s` silently copies the directory instead of linking it,
# and a native symlink needs Developer Mode or an admin shell. So on Windows,
# try a native symlink and fall back to a directory junction, which needs no
# privilege and reads as a symlink to [ -L ] and readlink.
link_dir() {
  src="$1"
  target="$2"
  case "$(uname -s)" in
    MINGW*|MSYS*|CYGWIN*)
      [ -L "$target" ] && rm "$target"
      if ! MSYS=winsymlinks:nativestrict ln -sn "$src" "$target" 2>/dev/null; then
        cmd //c mklink //J "$(cygpath -w "$target")" "$(cygpath -w "$src")" >/dev/null
      fi
      ;;
    *)
      ln -sfn "$src" "$target"
      ;;
  esac
}

for DEST in "${DESTS[@]}"; do
  mkdir -p "$DEST"
  prune_stale "$DEST"
  find "$REPO/skills" -name SKILL.md -not -path '*/node_modules/*' -print0 |
  while IFS= read -r -d '' skill_md; do
    src="$(dirname "$skill_md")"
    name="$(basename "$src")"
    target="$DEST/$name"

    # Replace a pre-existing real directory with a symlink to the repo copy.
    if [ -e "$target" ] && [ ! -L "$target" ]; then
      rm -rf "$target"
    fi

    link_dir "$src" "$target"
    echo "linked $name -> $target"
  done
done

# --- graphify ---------------------------------------------------------------
# graphify ships its own skill, so it is installed from PyPI rather than linked
# from here. `graphify install` writes that skill and the global CLAUDE.md
# section for each CLI; re-running upgrades both to the package's version.
if command -v uv >/dev/null; then
  uv tool install --upgrade graphifyy -q
  graphify install --platform claude >/dev/null
  graphify install --platform agents >/dev/null
  # `graphify claude install` and `graphify codex install` are per project:
  # they write ./CLAUDE.md or ./AGENTS.md, plus hooks. Point the same writers
  # at the user-level files so every project gets the rules. Claude's
  # PreToolUse hooks do nothing in a repo without graphify-out/. Codex's hook
  # is a no-op by design, so AGENTS.md is all Codex needs.
  uv tool run --from graphifyy python - <<'PY'
from pathlib import Path
from graphify import install as g
home = Path.home()
for md, marker, kind in ((home / ".claude" / "CLAUDE.md", g._CLAUDE_MD_MARKER, "claude-md"),
                         (home / ".codex" / "AGENTS.md", g._AGENTS_MD_MARKER, "agents-md")):
    text, nl = g._read_keep_newline(md) if md.exists() else ("", "\n")
    new = g._replace_or_append_section(text, marker, g._always_on(kind)) if text else g._always_on(kind)
    if new != text:
        md.parent.mkdir(parents=True, exist_ok=True)
        g._write_keep_newline(md, new, nl)
g._install_claude_hook(home)
PY
  echo "installed $(graphify --version)"
else
  echo "skipped graphify: uv not on PATH"
fi

# --- plugins ----------------------------------------------------------------
# Third-party plugins this setup depends on. Both commands are idempotent, so
# re-running is free. One line per plugin: "<repo> <plugin>@<marketplace>".
command -v claude >/dev/null || { echo "skipped Claude plugins: claude not on PATH"; exit 0; }
while read -r repo plugin; do
  [ -n "$repo" ] || continue
  claude plugin marketplace add "$repo" >/dev/null
  claude plugin install "$plugin" >/dev/null
  echo "installed $plugin"
done <<'PLUGINS'
anthropics/healthcare healthcare@healthcare
PLUGINS
