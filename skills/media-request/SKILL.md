---
name: media-request
description: Use when the user wants a movie or show on Jellyfin, or asks about one — triggers on "I want <title>", "get me <title>", "download <title>", "is <title> ready", "how's <title> going", "delete <title>", "free up space". Requests through Seerr, reports progress, and deletes only through prompting tools. Needs the `media` MCP server.
---

# Media Request

Turn "I want Dune" into a Seerr request, then report status when asked. Seerr is
the front door. Radarr, Sonarr and qBittorrent are for progress and repairs.

Every tool below is `mcp__media__<tool>`. Jellyfin is at https://jellyfin.litoral.lat.

## Reply style

Keep replies short. Name the title with its year, then say what happened.

Good: `Requested *Dune: Part Two* (2024) in 1080p. Radarr is searching.`

## Flow

1. **Parse.** Read the title, the optional year, the type hint ("the show"), the seasons, and the quality. Never ask for anything the user did not leave open.
   - If the user asks for status ("is <title> ready", "how's <title> going") or a delete, check the library first.
   - Run `radarr_list_movies` and `sonarr_list_series`. If exactly one item matches the title, use it. Skip steps 2 and 3.
2. **Search.** Run `seerr_search` with the title.
   - Each result shows the title, year, `tmdbId`, `type` (`movie` or `tv`), and `status` — Seerr's `mediaInfo.status` (1 unknown, 2 pending, 3 processing, 4 partially available, 5 available), or `not requested`.
3. **Disambiguate.**
   - If exactly one result has the same title and, when the user gave a year, the same year, go ahead without asking.
   - Otherwise show up to 3 lines: `Title (year), movie|show, one-line overview`. Ask which one. Block on the answer.
   - Never guess between a movie and a show that share a title.
4. **Check the library.**
   - Use the matched result's `status` first:
     - `5` (available): say "Already in Jellyfin". Stop.
     - `4` (partially available, show only): fall through to the episode-file check below to name the seasons present.
     - `2` or `3` (pending/processing): go to step 7.
     - `not requested` or `1` (unknown): fall through to the Radarr/Sonarr check below — the status field can be stale or the item may already sit in the library outside of Seerr.
   - Fallback — for a movie, run `radarr_list_movies`. Match on title and year.
   - For a show, run `sonarr_list_series`. Match on title and year.
   - If the movie has a file (✅), or the show has every episode, say "Already in Jellyfin". Stop.
   - If the show has some episodes, run `sonarr_list_episode_files` for the series. Name the seasons present. Offer the missing seasons.
   - If the item is in the library but has no file, go to step 7.
   - If the item is not in the library, run `seerr_list_requests` with filter `pending`, then `processing`. If one row has the same `tmdbId`, go to step 7.
5. **Seasons** (show only).
   - If the user named seasons, use them.
   - Otherwise ask: "All N seasons, or just S1?" Block on the answer.
   - Add a size hint to the question. Use the season count you know and say "about".
   - Size arithmetic: ~3 GB per episode, ~12 GB per movie, 500 GB quota.
   - If the show is in Sonarr, take the counts from `sonarr_get_series`.
   - Add current usage when it is cheap: sum the GB column of `radarr_list_movies` and `sonarr_list_series`. Otherwise omit it.
6. **Request.**
   - For a movie, run `seerr_request_movie` with `tmdb_id` and `quality`.
   - For a show, run `seerr_request_tv` with `tmdb_id`, `seasons` (`"all"` or `"1,2"`) and `quality`.
   - Set `quality` to `1080p` unless the user named another profile, such as `4K`.
   - If the tool answers "Unknown quality", show the valid names it lists. Ask which one. Block on the answer.
   - Do not ask for confirmation. The request tools run without a prompt.
   - Reply with the title, the quality, and the request id. Say that Radarr or Sonarr is searching.
7. **Status.** Run it on demand. Never poll and never loop.
   - Find the request id with `seerr_list_requests` (filter `all`), matched on `tmdbId`. For a movie, `radarr_get_movie` shows the TMDB id.
   - Run `seerr_request_status`. Read the `Media status` code.
   - Codes: 1 unknown, 2 pending, 3 processing, 4 partially available, 5 available.
   - If the code is 5, say "Ready in Jellyfin". Stop.
   - If the library already has the file, trust the library over the code. Seerr lags behind Radarr and Sonarr. Say "Ready in Jellyfin". Stop.
   - If the code is 4, name the seasons present and offer the missing ones.
   - If the code is 2 or 3, run `radarr_queue` (movie) or `sonarr_queue` (show).
   - Find the torrent in `qbt_list_torrents` by release name. Report its percent, its speed, and the time left (remaining GB divided by speed).
   - If the queue is empty, say "no release yet, Radarr keeps searching". Offer `radarr_releases` for a manual pick.
   - If the torrent state is `stalledDL`, run `qbt_torrent_details`. If it has 0 seeds, offer `radarr_delete_queue_item` or `sonarr_delete_queue_item` with `blocklist` true, then a new search.
   - Run those repair tools only after the user agrees. They prompt.
8. **Delete.** Triggers: "delete <title>", "free up space".
   - For "delete <title>", find the item with `radarr_list_movies` or `sonarr_list_series`. Show its size in GB.
   - Then run `radarr_delete_movie` or `sonarr_delete_series` with `delete_files` true. These tools prompt. Let them prompt.
   - After a delete, run `seerr_delete_request` for the same title. A later request then works.
   - For "free up space", list the 10 largest items from both libraries, largest first. Ask which to delete. Block on the answer.
   - Delete one item per tool call. Never delete in bulk.
   - Never delete from a vague reply. If the user did not pick an item, ask again.
