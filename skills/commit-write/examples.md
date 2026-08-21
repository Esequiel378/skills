# commit-write — verbatim style anchors

Real commit subjects and PR titles from the monorepo, with customer and vendor
names genericised. The `(#NNNN)` suffix is GitHub's, appended at squash time —
never typed by the author.

## feat

```
feat: serve /.well-known/jwks from <vendor> app (#9540)
feat(py): filename access with test (#9672)
feat: use 2 GB for all web apps (#9541)
feat: tighten schemas used in nexus FHIR API
feat: build fax fetch blob URLs from configured azure host and request params
```

## fix

```
fix: useWasm false for PDFs with wasm-dependent images (#8273)
fix: move api class out of managed<Vendor>OutboundClient to allow reuse (#9585)
fix: brrr ctx - eng-6276 (#9624)
fix: memory leak in www due to nodejs (#9323)
fix: use 303 for login client cross-domain redirect
```

## prefactor

The house type: a refactor landed *ahead of* the feature it enables, so that
feature's diff stays reviewable.

```
prefactor: add docs folder (#6925)
prefactor(py-stores): cleanup meta wrapper and nano timestamp (#9550)
prefactor(lib-agent): make more fragile (#7810)
prefactor: narrow nexus sdk exports (#9601)
prefactor: make ANT_APP_NAME required (#9598)
prefactor: <vendor> brrr task stub
```

## refactor

Restructuring that prepares nothing in particular.

```
refactor: shared test builders module (#8036)
refactor(inspector): make event table collapsable (#7309)
refactor(factoid): rename tidbit (#9674)
refactor: simplify IWE task registry (#9445)
refactor(py): artifact store and event store from per request actions extra
```

## chore

```
chore: bump virtualisation.memorySize up to 2000 (#8228)
chore: no private ssh keys (#9160)
chore: rotate all anthropic llm keys (#9112)
chore(openapi): fix up openapi schemas (#8430)
```

## ci / build

```
ci(conventional-commit): no ch0re (#9510)
ci(merge): 45m timeout prem nix all mac (#9613)
build(sops): only one secrets file per env and create missing secrets.prod.json (#9110)
build(nix): jq system in scope (#8687)
build: avoid deprecated dynamo terraform attribute (#7661)
```

## test

```
test: run factoid test
test: skip flaky edges projector test (#9309)
```

## revert

Git's own boilerplate is left as-is, with no added explanation in the revert
itself:

```
revert: reduce thread count (#9435)
fix: Revert "prefactor: <vendor> brrr task stub" (#9705)
revert: revert: feat(py): filename access with test
```

## Scope

The scope names the code — the literal directory, package, or tool:

```
(py)  (py-stores)  (cdktf)  (nix)  (lib-agent)  (brrr-agent)
(inspector)  (factoid)  (openapi)  (merge)  (sops)
```

Omit it when the change spans several components, or when the summary already
names the component. `refactor(py): artifact store ...` scopes by package;
`feat: serve /.well-known/jwks from <vendor> app` needs no scope because the
summary already says where.

## Ticket on the PR title

Appended after the summary, lowercase — not in the scope slot, and not in
parentheses, which would collide with GitHub's `(#NNNN)`:

```
fix: brrr ctx - eng-6276
```

## NOMERGE / WIP markers

Replace the type on a PR that must not merge. Restore a real type before review:

```
NOMERGE(<vendor>) get-token script for dev
NOMERGE add users script
NOMERGE oncohealth schema update June 22
WIP feat: move <vendor> UI app to nexus
```

## What the house never writes

```
Rename Tidbit to Factoid across py/ (ENG-6276)   ← Title Case, no type, ticket in parens
Updated the exports.                              ← past tense, capitalized, trailing period
fix(py): fix filename access                      ← "fix" twice
chore: add retry to payment handler               ← behavior change hidden as a chore
feat: various fixes and improvements              ← names nothing
```
