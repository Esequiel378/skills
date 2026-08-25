# ADR format

ADRs live in `docs/adr/`, named sequentially: `0001-slug.md`. In a
multi-context repo, system-wide decisions go in the root `docs/adr/`;
context-specific ones go in that context's `docs/adr/`.

The value is in recording *that* a decision was made and *why*, not in filling
out sections. A single paragraph is a valid ADR.

## Minimal template

```md
# {Concise decision title}

{1–3 sentences: the context, the decision, and the rationale.}
```

## Optional sections

Add **Status**, **Considered options**, or **Consequences** only when they
genuinely clarify the decision. Most ADRs omit all three.

## What warrants an ADR

Strong candidates: architectural patterns, technology lock-ins, context
boundaries, deliberate deviations from convention, external constraints, and
non-obvious rejections of alternatives. The gate in
[SKILL.md](./SKILL.md#offer-adrs-sparingly) decides — all three conditions, or
no ADR.
