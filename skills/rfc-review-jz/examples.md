# rfc-review-jz — verbatim style anchors

Drawn from real RFC and design-doc comment threads, with customer names,
people, and internal codenames made generic. Match this register.

## The mechanism is described wrong

> Technically it can still happen; Two emit tasks might run simultaneously before either result is written to the database. It is just much, much rarer, and wouldn't happen on rerunning of tasks

> If it were its own parent, wouldn't that be a loop? I would expect the first event to have no parent

> Why wouldn't it be valid to have a newer meta point at an older meta?

> Also, technically the writer doesn't inspect meta objects but asks the serializer for information about them

> Those would be in a different layer (the storage layer). In the writer's types they would be storage backend implementations

> This filter is to narrow down the events a projector sees, separate from the linked list

> The gateway IS where the auth happens/matters; everywhere else is just passing tokens along so that they can reach the gateway

## A primitive already covers it

> Can we use API keys?

> The auth flow through the gateway works this way; use a longer lived token to mint a scoped token for a specific bit of work. This token then gets passed down to the task runner

> Config projections could be where this logic would be centralized

> the events payloads already have a schema and hashing them + the parent key should be deterministically unique

> A feature we have discussed, but not yet added, is to have "local defers". The task only needs to restart if the other task actually goes through the queue. We could have certain tasks, like emitting events, which satisfy the task runner's calling semantics with reading from and writing to memory, that are run on the same worker almost like a regular function call.

> Is there a world in which we have validation on our structured output and if the llm output is structurally sound but fails validation, we ask the llm for a correction, chat style, given the error message?

## The requirement has no consumer

> Who would be waiting? I don't imagine a cron tick cares about the results?

> Even if this were true, projection speed at this scale is not a material concern. And if anything we can no longer stream our projection from event 0 to the end, but rather we must first traverse all events to track parents before yielding them in reverse order.

> This annoyance should only be felt if we don't use caching.

## Breaks what works today

> This would mean that you can't fork older records.

## Scale

> For a selector, I don't think radio works as UI (always imagine we have a million of them). Also, since the id by itself is not human friendly, users would generally always copy paste one, or need to filter by something else (e.g. tags)

> One reason to use path params rather than query params is to make caching more reliable/deterministic. query param order isn't guaranteed (unless we implement them that way) and doesn't convey as clearly (to user as well as author) how responses might be cached.

## Reversibility

> Worth noting that getting good at migrations makes it less important to make perfect decisions, which allows us to move faster.

> With that in mind, exposing anything to our customers hinders our ability to pivot/course correct, so let's be stringent on public facing data.

## Proposing an alternative

> What about xxxxxxxxxxx#yyyyyy#zzzzzz
> Where x is a timestamp (either unix numerical or human friendly), y is a machine id or some salt and z is a step counter (basically an unpacked snowflake id). Cons: larger, non-standard, potentially less efficient
> Pros: easier to understand, easier to dissect, easier to extract timestamp/machine info).
>
> The reason for suggesting this here is that these id specs introduce some "arbitrary complexity". If all we need is uniqueness, but we're not tied to a max size or compatibility with other systems, we don't need to tie ourselves to a foreign, less human friendly format.
>
> To make this call it would be good to understand the actual implications of string length in database usage.

## Correcting a quoted claim

> "No ~~Javascript~~ function can create a nanosecond precise timestamp relative to the rest of the world"
>
> I wouldn't worry about nanoseconds as a precision mechanism at all; in this use case they only serve as some sort of "salt" to make clashes less likely

## Security

> Mitigation is helpful, but a proper remediation should involve rotating secrets (and communicating it).
>
> One point of projections is to be able to access previous points-in-time so you couldn't hide value by overriding it.

> Are document dates considered sensitive? Could we use them as tags? (link to the legal-questions ticket)

## Naming and missing artifacts

> 1. Wouldn't call artifact refs ids
> 2. In the current model, on a single record, we wouldn't need to provide a from_ value, that is already encoded in the record.

> Could you include the projection schema?

> One nuance in my personal mental model of canonical vs special records is that canonical ones have a fixed id, and special ones have a random id but special treatment

## Priorities and concessions

> Honestly I personally don't care as much about the backing storage as I do about the interface to it
>
> For 99.99+ cases this should be a noop

> POC, sure

> More or less relevant depending on the scope of this RFC

## Humor, once

> One of us will need to volunteer for cryogenic suspension to deal with the catastrophe of our interplanetary consciousness
