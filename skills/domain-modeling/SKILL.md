---
name: domain-modeling
description: Use when building or sharpening a project's domain model — triggers on "/domain-modeling", discussing codebase terminology, naming a new concept, resolving what a term means, writing or editing a CONTEXT.md or CONTEXT-MAP.md, or recording an architecture decision (ADR).
---

# Domain modeling

Actively build and sharpen the project's domain model as you design: challenge
terms, invent edge-case scenarios, and write the glossary and decisions down
the moment they crystallise. Merely *reading* `CONTEXT.md` for vocabulary is
not this skill — that's a one-line habit any skill can do. This skill is for
changing the model, not consuming it.

## File structure

Most repos have a single context:

```
/
├── CONTEXT.md
├── docs/
│   └── adr/
│       ├── 0001-event-sourced-orders.md
│       └── 0002-postgres-for-write-model.md
└── src/
```

If a `CONTEXT-MAP.md` exists at the root, the repo has multiple contexts. The
map points to where each one lives:

```
/
├── CONTEXT-MAP.md
├── docs/
│   └── adr/                          ← system-wide decisions
├── src/
│   ├── ordering/
│   │   ├── CONTEXT.md
│   │   └── docs/adr/                 ← context-specific decisions
│   └── billing/
│       ├── CONTEXT.md
│       └── docs/adr/
```

Create files lazily. If no `CONTEXT.md` exists, create one when the first term
resolves. If no `docs/adr/` exists, create it when the first ADR is needed.

## During the session

### Challenge against the glossary

If the user's term conflicts with `CONTEXT.md`, name the conflict immediately.
"Your glossary defines 'cancellation' as X, but you seem to mean Y. Which is
it?"

### Sharpen fuzzy language

If a term is vague or overloaded, propose one precise canonical term. "You're
saying 'account': do you mean the Customer or the User? Those are different
things."

### Discuss concrete scenarios

When domain relationships are discussed, stress-test them with specific
scenarios. Invent scenarios that probe edge cases. Force precision about the
boundaries between concepts.

### Cross-reference with code

When the user states how something works, check whether the code agrees. If
the code disagrees, show the contradiction. "Your code cancels entire Orders,
but you just said partial cancellation is possible. Which is right?"

### Update CONTEXT.md inline

When a term resolves, update `CONTEXT.md` immediately. Do not batch updates.
Read [CONTEXT-FORMAT.md](./CONTEXT-FORMAT.md) for the format.

Keep `CONTEXT.md` free of implementation details. It is a glossary — not a
spec, a scratch pad, or a home for implementation decisions.

### Offer ADRs sparingly

If all three conditions hold, offer to create an ADR:

1. **Hard to reverse** — the cost of changing your mind later is meaningful.
2. **Surprising without context** — a future reader will wonder "why did they
   do it this way?"
3. **A real trade-off** — genuine alternatives existed and one was picked for
   specific reasons.

If any condition is missing, skip the ADR. Read
[ADR-FORMAT.md](./ADR-FORMAT.md) for the format.
