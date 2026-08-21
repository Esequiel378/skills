# sl-review — verbatim style anchors

Style anchors drawn from real production monorepo reviews, with domain specifics
genericised. Match this register exactly. Note: across 200 recent PRs, zero
review-summary bodies and zero request-changes — everything below is an inline
comment.

## Wrong layer / wrong context

> This code is vendor specific, so wouldn't live in `lib-core`. I think the design is more like part of a class that wraps the Vendor Outbound API Client to expose things like `updateAuthorizationWithApproval()` -- this isn't clear in the ticket, which is confusing -- it really needs a design doc, considered together with the other vendor tickets.

> These endpoints should be case-scoped, not workspace-scoped.

> you should never create a case. Just declare a case scoped endpoint

> Don't create a case from workspace context and then act on it. Make your endpoint in case context and use `emit` from there.

> I think this is an anti pattern, you shouldn't be able to project one case from another case context. cc @owner for comment.

> I don't understand this one either. Mentioning Artifacts sounds like a leaky abstraction?

## Generic → specific

> Don't use ArtifactAdded -- make an event that actually describes what is happening, e.g., PatientInfoAdded

> Can we use a specific event to say what is being added here? DocFileAdded or something? Use of ArtifactAdded in application logic is a smell.

> Better to have an enum expressing what was discovered, not an assumed goal of running the action, e.g., `DOCUMENTS_WITHIN_TIME_RANGE` / `DOCUMENTS_OUTSIDE_TIME_RANGE`

> surely this has more structure than `string`?

> avoid boolean types, use an Enum for forward compatibility

> this is such a generic name, where I assume it's meant to be in the context of "action selected for an eval run against a single action?" -- probably should change the name to be more specific to its intended use

## Schema-driven

> Why isn't this a JSON schema

> This should be YAML in `/schemas`

> We are schema-driven, so start with the types in `/schemas/org` in YAML, then these types will be auto-generated into `@org/schemas/evals`

> Note: please do not edit the openapi spec directly. Customer provided schemas should be left in the codebase in the state they were received.

> this is in `/schemas` import it from `@org/schemas`

> Can we have description fields explaining what this schema represents? Why does it have `case` in the title for example?

## Naming / domain language

> s/UsersEndpoints/LegacyUserInfoEndpoints/

> s/iff/if/

> `s/ah/acmehealth` -- Acme Health slug form should be `acmehealth` everywhere

> "Basic" as a descriptor always sucked. If this is a shim, it can just be called Case. If it will itself be replaced, it could be called DynamoEvents which is what it is or something.

> Patient, not member, surely?

> And then, why not use Patient FHIR Resource?

> `REQUESTED_SERVICE_IDENTIFIED` is more specific, and will be more identifiable.

> started rather than intended maybe?

## PR scope / completeness

> **Every PR should enter the codebase such that if future work stops at that point, there is no trace of half-done work.**

> This is not an argument. We may make tech debt tradeoffs for time constraints, (if we do, a lead or myself will make the call. You should just default to writing good code).

> Don't mix the web app and the package in 1 PR

> Also, let's separate out the jq edits to the openapi schema into their own PR.

> This function should be a separate PR, contributed further upstream into the SDK. cc @owner

> That doesn't make sense to me. Let's make this PR complete with the env vars resolving to the correct values.

> I don't think we need the `ponytail` style comments. Code in master shouldn't refer to upcoming potential work.

> I'd rather this was a separate "refresh tag projection list" button I think? It's obviously worse user experience, but it's more honest until we get the actual behavior working, and I don't think inspector should "fake" system behavior for nicer UX.

## Reinvented machinery

> don't create your own tag store, use `artifact.tags`

> this projector format is confusing -- you are bundling everything into `update` and doing your own `if event_type` -- there is a projector pattern that handles this for you

> This functionality already exists on the case-tags page

> `/health` is provided by the platform SDK to all apps, gratis. (Is that right @owner ?)

> better yet, don't do any of the case setup here, that functionality exists elsewhere. Make these endpoints case-scoped and assume a case.

## The Socratic barrage

> I need help to understand what's happening in this file. Why so many new schemas (why can't we use the schemas you created in `/schemas`?) Why so many endpoints? Why `PATCH`? Why is it minting an event timestamp?

> Why does the sampling "logic" (this random "score" -- badly named surely, as what are you "scoring"?) belong in the data layer? Why can't you pull the cases and sample at QA-time?

> Is an eval run not a case?

> Are evaluator functions not themselves also actions?

> Is this an error? It doesn't seem like something that should raise. Isn't it just, no doc_files found, so didn't have to do work.

> What does the review_date mean exactly? Couldn't that be inferred from review events?

> I don't understand this get endpoint. If I have the artifact, I already have the review? Surely I would want to get the review by some identifier, like a tag?

## PHI

> That sounds bad... it could contain PHI. Better to add the stacktrace as an artifact and link to it

> see comment above - could contain PHI. Events should be ids and enums.

## Cruft / comments

> console.log

> console.log etc. remove these before review pls

> the console.log is still there tho?

> What is "the golden"? What does this comment mean?

> This comment is confusing. What is the reference implementation? Where is it? What is the `FIXTURE_001` sandbox capture?

> What's the purpose of this log message? We have tracing.

> Add a comment explaining why these endpoints exist.

## One-liners

> why?

> ?

> y tho

> Return a value.

> _really_ required

> "RI"?

> a... what?

> This sounds bad?

> then maybe it shouldn't be?

> can be fixed now, no?

> Platform.

## Citing sources

> https://www.linfo.org/rule_of_silence.html

> https://build.fhir.org/patient.html

> How do these document types relate to LOINC like https://hl7.org/fhir/R4/valueset-c80-doc-typecodes.html ?

## Dry humor

> By all means add a comment to explain what the acmehealth mock is doing in the sandbox app. This would of course be clearer if the sandbox app were called 'mock-customer', but where would be the fun in that?

> Sgtm, did TypeScript grow up suddenly?

> Please complain loudly if your memory causes issues -- 24GB sounds .... low

## Conceding

> ohhh yeah you're right there could be tasks in there. keep the web dir

> That one is my fault! I had the unnecessary `usage()` in my scripts I pointed a teammate to. Let's just remove it and do a one-liner.

> 2c I prefer the other approach. It's cleaner. Us assuming that the performance hit of unlucky requests is significant is a premature optimization.
