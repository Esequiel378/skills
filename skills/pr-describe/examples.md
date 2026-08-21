# pr-describe — verbatim style anchors

Real PR bodies from the monorepo, with customer, vendor, and colleague names
genericised. Note what is absent throughout: headings, checklists, "Summary"
and "Testing" sections, and `Closes #`.

## No body at all

A third of PRs ship with an empty body. The kind is consistent across authors:
mechanical, low-risk, or fully described by the title — version bumps, key
rotations, CI timeouts, `build(nix)` work, small `chore:` housekeeping.

## Short — one line or a few bullets

> Remove the SPA-focused body height rules.

> Overload files-api to allow you to not pass infra, only stem context.

> was missing the test driver for factoid test. probably a casualty of the refactor a few months ago

> identical to #9656 which github mysteriously closed

> combines schema updates from #8652, #8735, and #8764

> see #8758 for real usage

> hoping this will address: https://github.com/<org>/<repo>/actions/runs/<id>/job/<id>?pr=8207

```
- generate flattened FHIR schemas in python
- move `__init__.py` and `py.typed` up one folder to enable absolute imports within the schema package
```

```
- replace FIXME in previous PR with simple nexus API call
- add request ID and idempotency key to meta passed through brrr calls
```

## Medium — bullets with a reason attached

```
various related cleanups to nexus sdk:
- do not export `getWebConfig` and `getWebInfra`, to guide users toward using only `createWebApp`
  - expose logger from return value of `createWebApp`
- drop "fromEnv" from name of `createWebAppFromEnv` -- no visible "env" in arguments
- small knock-on edits to apps that use this sdk
```

> some PDFs have images that require WASM for pdf.js to render them. our CSP does not allow for this.
>
> let's see if we can avoid complicated post-processing of raw PDFs with WASM-dependent images by simply disabling WASM. main tradeoff: slower PDF render times

> we had agreed on the more restful "faxes", but they've written "fax" for their fax results route... go with that for now
>
> https://<org>.slack.com/archives/<channel>/<id>

```
move <vendor> UI functionality to a separate app from the intake endpoint, for a few reasons (copied from Linear ticket):
1. the UI app can be moved to the nexus SDK without interfering with the intake app's auth story
2. the UI app can remain accessible by /y/<vendor> without mTLS
3. the intake app can remain behind mTLS at <vendor>.anterior.app
```

## Long — when the context earns it

> supports #9645 which creates a large schema based on a large enum
>
> allows for some schemas to be derived from other schemas, such that they all can both be used by downstream codegen and be persisted back into the repo

> see [nodejs github issue](https://github.com/nodejs/node/issues/<id>) and comment in code.
>
> this PR only "fixes" `serveStaticFile`, and other parts of the codebase may have the same pattern (returning a `Blob` via `Response` from a web endpoint).  I think this PR covers the vast majority of memory leaks seen in production, so I am inclined not to make this hacky "fix" in multiple places.
>
> repro'd locally by repeatedly fetching a static file and seeing process memory usage climb:
>
> `for i in $(seq 1 10000); do curl -X GET -s -H "Cookie: access_token=${ACCESS_TOKEN}" http://localhost:22171/static/vendor/scripts/htmx.min.js > /dev/null; done`
>
> validated in ephemeral

```
- adds the fax result PDF upload path so outbound PDF bytes can be written to the <vendor>'s `responserepository` blob location. emit the FaxResultDocumentUploadedV1 event -- i've intentionally kept this event minimal while we figure out how to represent HTTP interactions with arbitrary services (see: https://linear.app/<org>/issue/ENG-5636/generic-http-requestresponse-events)
- register upload_fax_result_document_v1 and subtask _upload_fax_result_to_azure_v1
- make blob I/O methods keyword-only and update the fetch download callsite and test doubles accordingly

NB: While <vendor>'s Azure PUT behavior is still being validated, HttpResponseError from upload is logged (PHI-safe) and not re-raised, but the stem still receives the upload event with the canonical URL; follow-up (ENG-5625) will replace this with explicit Azure error handling once expected upload behavior is confirmed
```

## Stack and merge order

> this must be merged AFTER #9566 or else <vendor> evals will break due to broken stem token logic

> depends on both https://github.com/<org>/<repo>/pull/9367 and https://github.com/<org>/<repo>/pull/9371

> this PR exists separately from #9599 to isolate 3000+ lines of package-lock.json changes

> create the skeleton for the new <vendor> UI app, in such a way that the next PR in this stack is as easy to review as possible (e.g., as many pure file renames as possible).  notably this PR contains 2,996 lines of `package-lock.json` changes.  the next PR in the stack actually moves functionality over.

> finishes revert of https://github.com/<org>/<repo>/pull/8450 and https://github.com/<org>/<repo>/pull/8370
> started here https://github.com/<org>/<repo>/pull/8455 but missed the actual package json

> EDIT: these two blocker PRs were merged. so far, no issues seen

## Evidence

Screenshots are pasted as raw `<img>`, straight from GitHub's uploader, usually
with no caption:

```
<img width="1269" height="783" alt="Screenshot 2026-02-17 at 23 16 44" src="https://github.com/user-attachments/assets/<hash>" />
```

One-line sign-offs used in place of a test plan:

> tested by checking that ephemeral web workers are up

> tested locally by running an IWE task

> validated in ephemeral

Pasted observability evidence for a bugfix:

> test run logs:
> https://app.datadoghq.com/logs?query=<deep-link>

## Open questions left for the reviewer

```
TODO:
- use either `/y/` or vanity domain
- should `createLegacyServer` wrap `createWebApp`?
- figure out plan for how to plug into `mountExpectedEndpoints` to serve redirection
  - this change would be in a follow-up PR, but this PR should be informed by it
```

> TODO: wait for #9514 which rearranges the nix tests

## Voice

Lowercase openers are the norm. `NB:` marks a caveat, after the main bullets.
Hedges and dry asides are normal and welcome:

> nb: does not include the output of the note template in the send response yet

> turns out this is ~~pretty~~ somewhat simple

> one of those PRs where thinking-to-lines-of-code ratio was very high

> i typically would have included dependencies/e2e tests but we need to figure out how to deal with the remote dependencies
