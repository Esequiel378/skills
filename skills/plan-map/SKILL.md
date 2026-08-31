---
name: plan-map
description: Use when a chunk of work is too big to hold in one agent session and the way to the end isn't visible yet — triggers on "/plan-map", "chart a map", "this is too big for one session", "plan this out across sessions", or handing over a map issue and asking what's next. Charts the work as decision tickets on GitHub, then resolves them one per session.
---

# Plan Map

A loose idea has arrived, too big for one agent session, and wrapped in fog: the way from here to the **destination** isn't visible yet. This skill is about finding that way, not charging at the destination. It charts the way as a **shared map** of GitHub issues, then works its **decision tickets** — questions whose resolution is a decision, not slices of a build to execute — one at a time until the route is clear.

The destination varies per effort, and naming it is the first act of charting: it shapes every ticket. It might be a spec to hand off and iterate on, a decision to lock before planning starts, or a change made in place like a data-structure migration. The map is domain-agnostic: engineering work, course content, whatever fits the shape.

## Plan, don't do

This skill is **planning** by default: each ticket resolves a decision, and the map is done when the way is clear, with nothing left to decide before someone goes and does the thing. The pull to just do the work is usually the signal you've reached the edge of the map and it's time to hand off. An effort can override this in its **Notes**, carrying execution into the map itself. Absent that override, produce decisions, not deliverables.

## Refer by name

Every map and ticket is an issue, so it has a **name**: its title. Refer to it by that name in everything the user reads. Never refer to it by a bare id, number, or slug. A wall of `#42, #43, #44` is illegible; names read at a glance. The id and URL don't vanish; a name wraps its link, but they ride *inside* the name, never stand in for it.

## The map

The map is a single GitHub issue labelled `plan-map:map`. Its tickets are child issues of the map.

The map is an **index**, not a store. It lists the decisions made and points at the tickets that hold their detail. A decision lives in exactly one place, its ticket, so the map never restates it, only gists it and links.

### The map body

The whole map at low resolution, loaded once per session. Open tickets are **not** listed: they are open child issues, found by query.

```markdown
## Destination

<what reaching the end of this map looks like: the spec, decision, or change this effort is finding its way to. One or two lines; every session orients to it before choosing a ticket.>

## Notes

<domain; skills every session should run; standing preferences for this effort>

## Decisions so far

<!-- the index: one line per closed ticket, enough to judge relevance, then zoom the link for the detail the ticket holds -->

- [<closed ticket title>](link): <one-line gist of the answer>

## Not yet specified

<!-- see "Fog of war": in-scope fog you can't ticket yet; graduates as the frontier advances -->

## Out of scope

<!-- see "Out of scope": work ruled beyond the destination; closed, never graduates -->
```

### Tickets

Each ticket is a **child issue** of the map. The issue number is its identity. Its body is the question, sized to one 100K token agent session:

```markdown
## Question

<the decision or investigation this ticket resolves>
```

Each ticket carries a `plan-map:<type>` label, one of `research`, `prototype`, `grilling`, `task`. See [Ticket types](#ticket-types).

A session **claims** a ticket by assigning it to the user driving the map. Claim first, before any work, so concurrent sessions skip it. That assignee *is* the claim: an open, unassigned ticket is unclaimed.

Blocking uses GitHub's native issue dependencies. Native blocking renders the frontier visually in GitHub's own UI, so the user sees what's takeable without opening the map. A ticket is **unblocked** when every ticket blocking it is closed. The **frontier** is the open, unblocked, unclaimed children: the edge of the known.

The answer isn't part of the body. Record it on resolution, per [Work through the map](#work-through-the-map). Link assets created while resolving a ticket from the issue. Never paste them into it.

## Ticket types

Every ticket is either **HITL** (human in the loop, worked *with* a human who speaks for themselves) or **AFK**, driven by the agent alone. A HITL ticket only resolves through that live exchange. Never stand in for the human's side of it — a grilling that answers its own questions has broken this.

- **Research** (AFK): Reading documentation, third-party APIs, or local resources like knowledge bases to surface a fact a decision waits on. Dispatch a subagent to resolve it. Use when knowledge outside the current working directory is required.
- **Prototype** (HITL): Raise the fidelity of the discussion with a cheap, rough, concrete artifact to react to: an outline, a rough take, a stub, or throwaway UI/logic code. Build the roughest version that answers the question. Link it from the issue as an asset. Use when "how should it look" or "how should it behave" is the key question.
- **Grilling** (HITL): Conversation. The default case. Run plan-grill.
- **Task** (HITL or AFK): Manual work that must happen before a *decision* can be made: nothing to decide, prototype, or research, but the discussion is blocked until it's done. Signing up for a service so its API can be judged, provisioning access, moving data so its shape can be seen. This is the one type that *does* rather than decides, and it earns its place by unblocking a decision, not by delivering the destination. Drive it alone where you can (AFK). Otherwise show the user a precise checklist (HITL). Resolve it when the work is done. Record what was done, plus any resulting facts later tickets depend on: credentials location, new URLs, row counts.

## Fog of war

The map is *deliberately* incomplete: don't chart what you can't yet see. Beyond the live tickets lies the **fog of war**: the dim view of decisions and investigations you can tell are coming but can't yet pin down, because they hang on questions still open. Resolving a ticket clears the fog ahead of it, graduating whatever's now specifiable into fresh tickets, one at a time, until the way to the destination is clear and no tickets remain.

The map's **Not yet specified** section is where that dim view is written down: the suspected question, the area to revisit later. It's the undiscovered frontier *toward* the destination: everything here is in scope, just not sharp enough to ticket. Write as loosely or as fully as the view allows; it doubles as a signpost for collaborators reading where the effort is headed.

**Fog or ticket?** The test is whether you can state the question precisely now, *not* whether you can answer it now.

- **Ticket when** the question is already sharp, even if it's blocked and you can't act on it yet.
- **Not yet specified when** you can't yet phrase it that sharply. Don't pre-slice the fog into ticket-sized pieces: it's coarser than a ticket, and one patch may graduate into several tickets, or none, once the frontier reaches it.

**Not yet specified** excludes what's already decided (Decisions so far), what's already a live ticket, and what's out of scope.

## Out of scope

Fog only ever gathers *toward* the destination. The destination fixes the scope, so work beyond it is **out of scope**: it isn't fog, and it doesn't belong in **Not yet specified**. It gets its own **Out of scope** section on the map: work you've consciously ruled out of *this* effort. Scope, not sharpness, lands it here.

Out-of-scope work never graduates, because the frontier stops at the destination. It returns only if the destination is redrawn, and then as a fresh effort, not a resumption.

Ruling something out of scope is a scoping act, not a step on the route. If an existing ticket turns out to sit past the destination, close it. A closed ticket is unambiguously off the frontier. Leave one line in the **Out of scope** section: the gist, why it's out of scope, and a link to the closed ticket. Keep it out of **Decisions so far**, which records the route actually walked; a scope boundary isn't a step on it.

## Invocation

Two modes. In both, resolve at most one ticket per session. Research tickets are the exception.

### Chart the map

The user invokes with a loose idea.

1. Name the destination. Run plan-grill and domain-modeling. Pin down what this map is finding its way to: the spec, decision, or change. The destination fixes the scope, so settle it first.
2. Map the frontier. Grill again, breadth-first this time. Fan out across the whole space rather than deep on any one thread. Surface the open decisions and the first steps takeable now.
3. If this surfaces no fog, you don't need a map. The way to the destination is already clear, and the journey fits one session. Stop and ask the user how they'd like to proceed.
4. Create the map issue with the label `plan-map:map`. Fill in Destination and Notes. Leave Decisions so far empty. Sketch the fog into **Not yet specified**.
5. Create the tickets you can specify now as child issues of the map.
6. Wire the blocking edges in a second pass. Issues need numbers before they can reference each other.
7. Leave everything you can't yet specify in **Not yet specified**.
8. For each `research` ticket you just created, dispatch a subagent to resolve it in parallel. Capture its findings on a throwaway `research/<name>` branch. Add a link to those findings on the ticket.
9. Stop. Charting is one session's work. It hand-resolves nothing.

### Work through the map

The user invokes with a map, by URL or number. A ticket is optional. Without one, you pick the next decision, not the user.

1. Read the map issue. Load the low-res view, not every ticket body.
2. If the user named a ticket, use it. Otherwise take the first frontier ticket in map order.
3. Claim the ticket. Assign it to the user before any work.
4. Resolve it per its type. Zoom as needed: fetch the full body of any related or closed ticket on demand.
5. Run whichever skills the map's `## Notes` block names. If in doubt, run plan-grill and domain-modeling.
6. Post the answer as a resolution comment on the ticket.
7. Close the ticket.
8. Append a one-line gist plus link to the map's Decisions so far.
9. Create any newly-surfaced tickets, then wire their blocking edges.
10. Graduate any fog the answer made specifiable. Delete each graduated patch from **Not yet specified** so it lives only as its new ticket.
11. If the answer reveals a ticket sits beyond the destination, rule it out of scope. Don't resolve it on the route.
12. If the decision invalidates other parts of the map, update or delete those tickets.

The user may run unblocked tickets in parallel. Expect other sessions to edit the tracker concurrently.

## GitHub operations

All operations use the `gh` CLI. `gh` infers the repo from the clone you're in.

- **Create the map**: `gh issue create --label plan-map:map --title "..." --body "..."`. Use a heredoc for the body.
- **Read an issue**: `gh issue view <n> --comments`.
- **Child ticket**: create the issue, then link it to the map as a GitHub sub-issue via `gh api` on the sub-issues endpoint. If sub-issues aren't enabled on the repo, add the child to a task list in the map body. Put `Part of #<map>` at the top of the child body.
- **Blocking**: `gh api --method POST repos/<owner>/<repo>/issues/<child>/dependencies/blocked_by -F issue_id=<blocker-db-id>`. Get the blocker's numeric database id with `gh api repos/<owner>/<repo>/issues/<n> --jq .id`. It is not the `#number` and not the `node_id`. If dependencies aren't available, fall back to a `Blocked by: #<n>, #<n>` line at the top of the child body.
- **Frontier query**: list the map's open children with `gh issue list --state open`, scoped to the map's sub-issues or task list. Drop any with an open blocker (`issue_dependencies_summary.blocked_by > 0`, or an open issue in the `Blocked by` line). Drop any with an assignee. First in map order wins.
- **Claim**: `gh issue edit <n> --add-assignee @me`. This is the session's first write.
- **Resolve**: `gh issue comment <n> --body "<answer>"`, then `gh issue close <n>`.

GitHub shares one number space across issues and PRs, so a bare `#42` may be either. Resolve with `gh issue view 42` first.

---

Adapted from [`wayfinder`](https://github.com/mattpocock/skills/blob/main/skills/engineering/wayfinder/SKILL.md) by Matt Pocock (MIT).
