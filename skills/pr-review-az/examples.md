# pr-review-az — verbatim style anchors

Style anchors drawn from real production monorepo reviews, with domain specifics
genericised. Match this register exactly.

## Verdicts / summaries

> seems legit, thanks

> minor nits, overall looks good

> approved to unblock, but a bit confused about the ephemeral secrets

> I don't 100% agree philosophically, but YAGNI applies here, and it's trivial to separate it back out later, so I approve

> a few comments, main one is the comment about the comment.  overall looks reasonable -- feel free to ping me whenever you update the PR and I'll promptly take another look

## Type-system lies

> avoid `as T` whenever possible -- unlike C, casting in typescript does not result in a new value of the correct type, but rather tells the linter "trust me this is the correct type", which is basically lying to the linter and defeats the point of the linter

> trailing exclamation point is somewhat dangerous -- `!` does not do anything at runtime (unlike `?`, confusingly), so now the type checker thinks that `coding` has a real type, but it could actually be null/undefined.  then that blows up at runtime elsewhere, and it's hard to debug what happened.

> `as StemId` is a smell, can `collectResults` be modified such that its keys already satisfy `StemId`?

> leverage the type system!   there's a lot of defensive code in `countLeafFields`, which is a smell -- the code clearly expects an object of a certain type

> I assume this is nullable because only successful fetches will know the page count.  however this allows invalid state to be represented -- nothing except our own diligence prevents a successful fetch from lacking a page count, or a failed fetch from having a page count.  can you make the invalid state irrepresentable?

> typescript is powerful enough to accommodate a type like non-empty string

## null vs undefined

> prefer returning null instead of undefined here -- convention here is that `null` is an empty value, whereas `undefined` is a missing array/object access

> as a rule of thumb, we should only be checking whether something is `undefined`, never assigning `undefined` to a value

> prefer one of these, don't allow users to set a value to `undefined`

## Exceptions

> `Exception` is very broad, and `except Exception` is a massive anti-pattern.  do we need a catch at all?  if we expect this to always succeed, I wouldn't bother with a catch -- too defensive.

> I strongly prefer to return error values rather than raise exceptions, for all the standard reasons, especially this deep in the codebase

> a good principle is exception should only be raised when a retry would fix it; that does not seem to be the case.  happy to pair on this

> this recursion is scary since there is no guarantee that we hit a base case.

> if we don't know why this code failed, I'd prefer to let the exception bubble up, be logged, and cause a 5xx response

## Comment hygiene

> comments shouldn't record history -- the code works as-is, this explanation is not relevant to help someone looking at the code a month from now

> when leaving TODO comments, also leave your name and date

> this kind of comment goes out of date very quickly

> comments are used by future maintainers, not by the people writing/reading the code today.  if a future maintainer reads this comment, they will not learn where to find the design document or who to ask about it.

> is this comment necessary?  i.e., if someone adds a new event to this test, the test will immediately fail, which seems sufficient to me

## Naming / conventions

> rename to `selectionEnabled` to reflect boolean; `selectionMode` sounds like an enum

> camel case for typescript-native types, snake case for generated types

> in our typescript code, it's pretty nice that `snake_case` is always a zod schema (or at least generated code), and `camelCase` is typescript-local variables

> "column 2" is odd since the other columns seem to have non-numeric names

> nit: `source` is short enough that I prefer it; generally I bias towards full words instead of abbreviations

## Keyword args / defaults

> use keyword arguments here, easy to accidentally swap string arguments

> having multiple positional parameters of the same type is dangerous because it's easy to accidentally swap the order.  you can use `*` to force callers to provide the parameters as keywords

> avoid default value, pass in explicitly, especially since there only seem to be two valid options

## Layering / reuse / scope

> the logic of this handler does not seem specific to <vendor>, so why does it live in `app/<vendor>`?

> if we start changing the abstraction when we don't need to, the abstractions are no longer stable, and we open ourselves up to problems.

> see below -- this code has already been written once, better to reuse it

> should this change be in #<other PR>?  (it's small enough where I'm okay leaving it here, but curious why)

> since these new functions are not used elsewhere in this PR (other than tests), I would rather move them to the stacked PR that first uses them; otherwise it's hard to glean their intent

> this is not backwards-compatible, so event readers will fail to read existing events with this schema; I think this can be a new version of the event

> I don't remember who said it first, but good to "avoid coding aspirationally"

## Evidence

> `time_ns` is a lie
> ```
> In [1]: import time
> In [2]: for _ in range(10): a, b = time.time_ns(), time.time_ns(); print(b - a)
> 0
> 1000
> 0
> ```

> `OrderedDict` docs say: … we don't need the extra methods, and the fact that `dict` is ordered is a guarantee of the standard library (not just an implementation detail), so I strongly prefer plain `dict`

## PHI / logging

> I'd be careful about including filename in the error message -- is it guaranteed to never contain PHI?

> since `environment` is a free-form string, I don't think it's PHI-safe to log it, including in exception messages

> noisy log -- `warn` means failure, but we're returning a value here, so use `debug` instead

> rule of thumb is one `log.info` per task function call

## Tests

> I think this is testing implementation detail instead of overall behavior.  could the test use `with pytest.raises` and check whether `str(exc)` contains the marker?

> a better test would test the intended behavior, which is not "the temperature parameter is set to zero" but rather "the LLM is deterministic by default".  such a test may look like "run the LLM five times on the same prompt and verify the output is the same".

> add a test for when: indication ID and protocol ID are both present; there is an override for the same indication ID but a different protocol ID.  basically make sure we actually use the combo

## Voice / concession / praise

> idk if this is better but maybe??  very loosely held

> I did recommend plain `<NAME>` earlier.  on further thought I am persuaded in favor of `<PREFIXED_NAME>`; clarity of the prefix is well worth four extra characters

> nvm I can't read

> having a nullable promise type seems strange; can't tell if genius or madness (to borrow a phrase)

> 😍 I love this -- extremely clear and information-dense variable name

> this function is beautiful

> some folks here say typos are good because they are proof an LLM didn't write the code.  I like their spirit but I still also like to not have typos haha
