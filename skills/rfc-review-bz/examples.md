# rfc-review-bz — verbatim style anchors

Drawn from real RFC and design-doc comment threads, with customer names,
people, and internal codenames made generic. Match this register.

## Reasoning left off the page

> what are the tradeoffs?

> what are those extra thoughts? shouldn't that be elucidated here? or at least why the difference makes a difference in this case

> how is this easier to reason about? and can't we have audit trail with the event stream also

> allow it based on what ?

> wdym by "tenant-scoped secrets"?

> wdym by env var?

> i assume some new code

> that's not clear from the wording

## Permissions and access

> is there consideration of permissions within a workspace? can everyone in a workspace see everything? nothing?

> do we support user level control? do we distinguish between ownership and access grants? i think we still need to flesh out our model tbh

> is there a split between config that we set and that a user or admin can set?
> eg we don't allow anyone to access app foo, or we just allow one company to access app bar. regardless of user setting if we don't allow the app they still don't have access

> feels like an essential requirement

> i think we really need permissions or roles before we can have config like this

## Data leaving its boundary

> this is tricky. is production config of a customer their data? can we exfil that locally ?

> how encrypt ? different than the usual encryption we have in the database? where is the key stored, who has permissions to read from what part and access the keys?

> it could be e.g. a key ID into a secret manager store

## New infrastructure or a special case

> even if we did use a different scheme than the event stream, that doesn't necessitate a new table. the database is designed to have multiple logical universes or "logical tables" in it. talk to @owner, this was already settled for the other event types.

> why? or at least we can give them special treatment but isn't this just another type of event? putting it in its own table feels wrong

> i don't think this is a good argument, you can segment the workflow such that it doesn't get in its own way. and making it a separate table doesn't buy anything extra that this couldn't

## Uncontrolled read paths and flexibility

> i would say that there should probably be a config reader that gets the config. that would be the place where overrides could be applied etc or version control can happen. letting anyone read the config from anywhere seems messy and prone to misbehavior

> i think you need a stricter def for workspace config.
> at the moment it's very app focused and doesn't allow misc config options not in an app. is that the idea going forward?

> i don't know if uber flexible is a feature

> i dont know if it makes sense to make copies of the workspace config and pass that around vs taking the actual workspace config and transforming that with edits into the specific config you need and passing that down.

## Legacy worth removing now

> are these enterprise roles here to stay? this seems like the time to reevaluate (and remove) them

> enterprise and workspace assignment should also be evaluated to be moved to the new paradigm

## Operability

> 1. are we thinking about idempotency (keys)
> 2. it's hard to follow the logic when anything can throw an "expected error" and return a flow early. similarly if business logic is split between handlers and actions (and other task-runner calls) and doesn't necessarily write standard end states back to the system, it's extremely hard to reason over.
> 3. logging and debug metrics are extremely important
> 4. good testing is a must. not ai generated 100% coverage and not just happy path. but really good edge case testing.

> are there defaults / how would (e2e) testing work or a fresh workspace work?

> +hard to know what the behavior of code is at a given time because it might depend on config

## Doc hygiene

> avoid document meta commentary

> i think this context section is a little verbose and out of place

> not relevant to the document, just talk about the external system and maybe leave a link to read/watch more

> if this document doesn't cover secrets then it should mention that up front and then mostly stop mentioning throughout

## Facts in the narrative

> wasn't this for config and NOT secrets? that was the whole point

> wdym? my comment was about the facts of your timeline. the PR you linked was about storing non-secret config

## Handing off

> elaborate or have @owner elaborate plz

> this is confusing since all of these configs are in production. what's the purpose of a non-production production config?
