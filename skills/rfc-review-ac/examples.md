# rfc-review-ac — verbatim style anchors

Drawn from real RFC and design-doc comment threads, with customer names,
people, and internal codenames made generic. Match this register.

## Solution without requirements

> what are the requirements of a "configuration store"?  e.g., what is the version management story, from a user perspective, and why is it faster/better/etc. to build it outside of the rest of the platform?

The author replied "The requirements are those of a document store." Follow-up:

> > The requirements are those of a document store.
>
> what are the functional requirements that the ops team need for the business deliverables we have promised to the customer for the next month?

## Options without criteria

> related to above question about what the requirements are — how should we evaluate these options against each other?

> but this proposal includes a new table?  not necessarily asking for Postgres, but have we considered all the options?

## The center of the design is thin

> isn't the backing storage central to this design doc?  I assume this requires us to create new infra, so we need alignment on what that new infra is and its cost to maintain.  otherwise I'd like to understand how we're avoiding new infra, because that would imply awesome and dangerous things about our platform.

> this seems like the real meat of the proposal, and I think it could use more color about how the table is configured and guarded

> admin tooling is at least theoretically scoped to the user DB, and this seems separate

## New infra without an owner

> is the infra team on board with adding and maintaining a new table?

> how long does this take?  why not a DB with logging?

## Claims that don't match the description

> do we expect our agent to reduce manual outreach?  i.e., if a human needs to reach out manually for more info, I would expect that our agent would have to escalate too, rather than try to resolve on its own

> asking partially because the bold text doesn't seem to match the description here

> we currently don't send source references in JSON for the existing import; would that be a requirement for this one?

> current workflow?

## Custom where generic exists

> The fields (name, address, and account number) seem the same as in the generic task, can we reuse?

> @owner is there a generic version of this task?  I don't think we've used it in prior POCs

> I think we can create an artifact for `category_descriptions` and ensure it is picked up by the generic task, so we would not need a custom task

> and this might already exist?

> spoke with @teammate late yesterday — potentially we could reuse the v1 validation task.  unlike the more recent tasks, it extracts account identifiers from the document *first*, before comparing extracted values to expected values.

> @owner can you shed some light on why we moved from extracting account identifiers *before* comparison to extracting *at the same time*?  and does reusing the v1 task in this context make sense to you?

> Agree that we'd only be reusing the first half of the v1 task, I'm thinking that can be pulled out and refactored into an account identifier extract task.  I am hesitant to jump all the way to a generic extract task because of how much context is required for personal names.
>
> I see the motivation for moving from v1 to v3, but I don't think it applies for this import flow, since we don't know the expected account details.  Given that, do you think the proposed task would work here?

## Two problems merged into one

> > deferred building user management in case we move everything to SSO anyway
>
> SSO can handle authentication, but wouldn't we still handle the bulk of authorization?
>
> and either way, it's probably good to scope out the auth/auth APIs, and assess if we can take a shim approach.  SSO and roles are related but distinct problems, and I'm not even sure it's that much more overall work to tackle them separately.

## Interface ergonomics

> if we're creating these stores via UI, then the code should not directly reference the human-readable name.  either the end user passes in a name (or ID) that feeds into some variable in the code, or the code has some way to find available stores, like a `store.list` function

> why not `ctx.workspace.docs.getStandaloneStore`?  `standalone` feels odd at this high of a level

> maybe `/records/ids/html`?  currently with `/records/projections/html` and `/records/html`, it feels like `projections` and `html` are on the same level

> step 2 implies that some projections are unavailable to some records, is this expected?  if not, should the projection and ID dropdowns be independent?

> noob question: can we control the headers such that /latest is never cached but everything else is?

## Upstream causes, flagged out of scope

> not in scope for this project: why does the upstream form do this?  how feasible is it for the customer to improve the form?

## Joining someone else's pushback

> > this page is this workspace and only this workspace
>
> @teammate if links auto-redirect to the correct workspace (with a brief landing page), how user-friendly or not is that experience?
>
> with the caveat that other open tabs in the prior workspace will stop working…

> @owner I agree with @teammate's point that I don't see what problem is being solved by removing the workspace switcher

## Scope, priority, and closing

> imo radio button is fine for now, soon enough it can be improved into a legit picker with tag filtering, etc.

> priority atm is "how do I list them" not "how do I filter down my list"

> I think this is correct for now; would be curious to see what the end result looks like, and figure out what can be abstracted to a common task

> I think the contract has already been agreed on

> looks great, thanks!
