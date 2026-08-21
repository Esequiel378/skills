# po-review — verbatim style anchors

Style anchors drawn from real production monorepo reviews, with domain specifics
genericised. Match this register exactly.

## Review summaries

> @author thanks! more small changes requested

> @teammate thanks, it's moving in the right direction. a couple more (small) blockers though. i think we should set the os env in Settings and use get_settings() in main for consistency with the other env vars

> left a few comments, non-blockers. thanks @teammate

> @teammate slack me and we can discuss! good stab but i see a few blockers worth discussion

> would be great to see some tests that use these improvements, but looks good 👍

> approving, but do we need slack-edge? i believe only the one bot depends on it, and i don't believe we have any plans for it

## Silent failure modes

> parameter.resource is FHIRResource: z.object({ resourceType, id }). zod keeps those two keys and drops the rest. App runs endpoint.body.parse before the handler, so this will happen silently (not 400)

> i'm not sure this `cs_` filter is doing anything. keys are typed CaseId, so if non-case keys can appear, the schema is wrong and this will silently drop data. if they can't appear, the filter doesn't serve a purpose

> nb: if meta.user_id is missing we silently skip, so the review disappears from filtered results silently

> this will happen silently. if this is exceptional, you may want to raise

## Tests that don't prove anything

> this cast skips the request schema, which is what is supposed to drop unknown keys. the handler only checks coverage-info === false, so this test suite would pass with no schema at all
>
> you should mount the app the way the discovery test already does and POST JSON to the real route. for the on/default cases, assert systemActions is present (not just cards.length > 0)

> this test asserts a throw, not an HTTP 400. the conformance suite and other clients never see ClientError, they see whatever the app mapped

> these three asserts don't prove we read the fields this PR exists to read. update + same id + "a coverage-information extension exists" still passes if orderCode only looks at one resource type

> these tests cover only pure helpers, but none of the async hydration + event read logic. hydrate_sample in particular could use coverage

## Type lies / parse at the edge

> this array is a type lie. it's typed as every FHIR resourceType, so `includes` compiles and doesn't narrow

> i mentioned something similar in another PR, but we should parse at the edge. this helper (and orderOf, as never, as unknown as in coverageInfoOf) asserts the JSON is a request instead of proving it

> this `as` cast could hide problems, eg silent filtering (if a projection has partial records, the response would omit them without any signal)

> @teammate we should avoid cast

## HTTP handlers fail with responses, not throws

> this function only returns a `Response` on the happy path (200). if the builder throws ClientError, this handler does not return (it throws) ... instead please return `new Response(..., { status: 400 })` here when the request can't be read (same as the reviews app). then this function always returns a Response, the test can strictEqual(res.status, 400), and we don't depend on "4xx as exception"

## Drift

> the problem is this same "succeeded or absent -> proceed" logic lives in different files in different forms, and these are likely to drift

> we should update the yaml example to mirror those (the spec documents the old keys). the codegen doesn't generate values from example (so the conformance suite won't fail) but the drift isn't ideal

> this comment is wrong (it's now in context.build_review_url)

## Scale scars

> this is the only getRangeProjection call site in the codebase with an unbounded {} range. so every request will: replay all review artifact_added_v1 events in the workspace from scratch, never read/write projection cache, then hydrate every case (n artifact fetches each)
>
> we hit this with adjudications -- i'd at least add a small range, but i think the real fix is to find a way to use cached projections

> you can try 30 days, but i won't be surprised if we hit gateway timeouts in deployed envs as the number of events grow. we were hitting them looking back just a week for one customer

## Concrete failure scenarios

> document_too_large is checked only after the fetch returns Valid, so a 501-page PDF is fully downloaded and written to the artifact store before we reject it. just confirming this is intended?

> gotcha: this makes byte-equality with the expected fixture depend on const key insertion order. a reorder of DECISION will fail the expected fixture contents with no type error

## Stack awareness

> this export doesn't appear used anywhere else yet -- assuming its in a higher PR in the stack? if not, drop the export. if you will be using in another PR in the stack, it can stay

> not blocking this PR so the demo route can stay canned. the unreadable request can still fall through until #1234. just flagging so #1235 doesn't only add config and forget to assert the builder works

## Comment noise

> this file has a comment on almost every constant ... we are working to avoid superfluous comment noise. comments are definitely justifiable when code is unsafe, or the code as written has a non-obvious function / reason for existing ... such comments will be forgotten and go stale very quickly

> i think we can remove this comment

## Naming

> should_process_documents is slightly misnamed i think, bc when it returns False, we still "process", we just skip extraction. maybe "has_processable_document"?

> then should it be called "case_completed_v1" or "case_reviews_completed_v1"?

## One-line nits

> nit, typo

> is this used?

> i think this is stale

> this indentation looks off / malformed

> looks like npm pollution?

> `Addres`?

> awkward return type? int | Invalid

> breaking your own rule? 😄

## PHI / trust boundaries

> i'm almost positive document_id cannot contain PHI, but we may need to formalize this with the downstream team before we start logging it

> it would be helpful to know for sure whether or not this is PHI

## Conceding

> ah, ok you're right. resolved 👍

> alright, i'm sold 👍

> > how will the downstream service learn that this has happened
>
> nevermind -- i see the DOCUMENT_DOWNLOAD_FAILED post to them
