---
name: session-search
description: Use before re-investigating something that may have been done in a past session — triggers on "did we already…", "search past sessions", "what did we decide about", "what was I doing on <date>", or before re-deciding a question that earlier work may have settled. Searches every archived Claude Code and Codex transcript on zx with `cass`.
---

# Session Search

Every Claude Code and Codex transcript from the Mac (source `mac`) and Windows (source `win`) is indexed on `zx`. Call `cass` over ssh. Search the transcripts before you repeat work.

## When to search

- If the task needs a non-trivial investigation, search first.
- If you are about to re-decide something, search first.
- If the user asks "did we already…" or "what did we decide about…", search first.
- If a search finds nothing, continue as normal.

## Search, then expand

1. Search for a distinctive term.

```bash
ssh zx cass search "<term>" --robot --limit 10 --fields summary
```

2. Read the hits. Pick the one that matches.

3. Expand that hit with the values from the hit.

```bash
ssh zx cass expand "<source_path>" --source "<source_id>" --message-index <line_number> -C 5 --json
```

Pass `line_number` to `--message-index` unchanged. Never pass it to `--line`: `--line` is a physical file line.

Narrow a search with these flags.

- `--source mac` or `--source win`
- `--agent <slug>`
- `--days <N>` or `--since <date>`
- `--max-tokens <N>` to cap the output size
- `--robot-format sessions` to list only the matching `source_path` values

## Result fields

| Field | Use |
| --- | --- |
| `source_path` | First argument to `expand` |
| `line_number` | Value for `--message-index` |
| `source_id` | Value for `--source` on `expand` |
| `agent` | Which tool wrote the session (Claude Code or Codex) |

## What was I doing around a date

Run `timeline` to list sessions in a window.

```bash
ssh zx cass timeline --since 2026-11-01 --until 2026-11-03 --group-by day --json
```

Add `--source mac` or `--source win` to filter by machine. Search or expand inside a session that looks relevant.

## Limits

- The index refreshes nightly at 02:30. Sessions from the last day may be missing.
- Results contain raw transcripts. They may contain secrets.
- Never paste a secret value into a reply, ticket, commit, or file. Name where it appeared instead.
