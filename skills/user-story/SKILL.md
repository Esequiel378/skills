---
name: user-story
description: Use when the user wants to turn a rough ticket idea into a great user story — triggers on "/user-story", "write a story for", "help me write this ticket", "make this ticket better", "split this into stories", or any half-formed task headed for Jira/Linear/GitHub Issues. Drafts in the team's house format, slices big asks into independently-shippable vertical slices, then offers to create the ticket in whatever tracker is connected.
---

# User Story

Turn a rough idea into a great story — one that a teammate could pick up cold and
ship as a single mergeable PR. Tracker-agnostic: draft as plain text, then offer
to create it wherever the user tracks work (Jira, Linear, GitHub Issues, …).

## What "great" means

A great story is **INVEST**: Independent, Negotiable, Valuable, Estimable, Small,
Testable. In practice that means it has four things:

1. **The right type.** Not every ticket is "As a user…".
   - **User story** — user-facing capability or goal. Use the narrative form.
   - **Technical story** — refactor, tech debt, infra, test coverage, migration.
     A persona narrative is a lie here; state *what changes + why + scope* instead.
2. **Value that's obvious.** Why does this matter, to whom? A story with no "so
   that" is a task, not a story.
3. **Acceptance criteria that are testable.** Each bullet is something you can
   demonstrably check off. "Works well" is not criteria; "returns 404 when the
   record is soft-deleted" is.
4. **Concrete anchors.** Real file paths, endpoints, permission keys, related
   tickets/PRs — so the reader doesn't reverse-engineer the ask. Never invent
   these (see step 2).

## Vertical slicing (the whole point of high PR throughput)

If the ask is bigger than **one mergeable PR**, split it — but slice **vertically**,
not horizontally.

- **Horizontal (bad):** "all the backend", then "all the frontend", then "wire it
  up". Nothing ships until the last slice; PRs pile up unmerged.
- **Vertical (good):** each slice cuts through every layer it needs and delivers
  something demonstrable on its own. Ship it, merge it, the trunk moves forward.

Test each proposed slice: *"Could this merge to main on its own and leave the app
working, even if the rest never happens?"* If no, it's a horizontal layer — reslice.
Ordering: earliest slice should be the thinnest end-to-end path (a walking skeleton);
later slices widen it. Feature-flag the incomplete UI rather than holding the PR.

## Process

1. **Read the ask.** Take everything the user gave. Don't ask for what you can infer.

2. **Find anchors (if a repo is present).** If run inside the relevant codebase,
   grep/read for the real files, endpoints, and keys the story touches and cite them.
   Not in the repo? Ask the user for the anchors or leave a clearly-marked
   `TODO(anchor)`. **Never guess a path or API** — a wrong anchor is worse than none.
   Fire a `context_research` call in parallel if institutional context (prior PRs,
   decisions, related tickets) would sharpen the story.

3. **Ask only the essentials you can't infer** — one message, batched. Usually:
   who's the user / what's the value, and any acceptance edge cases that aren't
   obvious. If it's already clear, skip this.

4. **Decide: one story or a slice set?** Apply the vertical-slicing test above.
   If splitting, present the ordered slices with a one-line rationale each before
   drafting them in full.

5. **Draft** in the template below. Match the team's house wording if you've seen
   their existing tickets (headings, label conventions); otherwise use this.

6. **Offer to create it.** Detect the connected tracker (look for an Atlassian /
   Jira, Linear, or GitHub MCP tool). Show the draft, ask which tracker + project,
   and **confirm before writing**. If no tracker is connected, just hand over the
   text. Never auto-create without a yes.

## Template

**User story:**
```
Title: <verb-first, specific — "Add v2 avatar endpoint", not "Avatar work">

As a <specific role>, I want <capability> so that <value>.

### Acceptance Criteria
- <testable outcome>
- <edge case / error path>
- <non-functional bar if it matters: perf, permission, audit>

### Scope
In:  <what this story delivers>
Out: <what it deliberately does NOT — link the follow-up slice>

### Anchors
- <path/to/file.ts>, <endpoint>, <permission key>, related: <TICKET-123>, <PR link>

### Open Questions
- <assumption that needs confirming, or "none">
```

**Technical story** — drop the persona line; lead with the change:
```
Title: <e.g. "Migrate avatar reads off the v1 endpoint">

<What changes and why. The user-facing behavior must not change / changes as follows.>

### Acceptance Criteria
- <observable, testable — "no call sites reference v1", "tests green", "flag removed">

### Scope / Anchors / Open Questions  (same as above)
```

## Guardrails

- A story you can't write acceptance criteria for is under-specified — say so and
  ask, don't paper over it with vague bullets.
- Don't gold-plate. The smallest story that's INVEST-complete beats an exhaustive one.
- If the "story" is really an epic (many slices, weeks of work), say so and propose
  the slice breakdown rather than one giant ticket.
