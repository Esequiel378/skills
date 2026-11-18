---
name: rfc-review
description: Use when the user is a reviewer on someone else's spec, design doc, RFC, or ADR and wants comments to post on it — triggers on "/rfc-review", "review this RFC", "review this spec", "review this design doc", "what should I push back on", "what would you descope". Reads the doc through every lens a senior reviewer applies, verifies each finding, and drafts comments for senior technical authors: pushback, descoping proposals, edge cases, gaps. Never posts. For a blunt verdict on an architecture doc for the user's own eyes, use architecture-review instead.
---

# rfc-review

Review a spec, design doc, or RFC that someone else wrote. Draft the comments
the user posts on it.

## Who reads the comments

The author is a tech lead, a head of engineering, or a CTO. They know the domain
better than you do. They have lived with this doc for weeks. They asked for
review because they want it attacked. Everything below follows from that.

- Skip the basics. Don't explain what idempotency is. Say where the design
  lacks it.
- Don't soften for seniority. A title changes nothing about whether the retry
  charges the customer twice.
- Show the reasoning in full. A senior author will find the missing step in an
  argument and dismiss the comment for it.
- Assume the obvious objection occurred to them. Look for it in the doc before
  you raise it. If they rejected it, answer their reason or drop the finding.
- A wrong comment costs the user credibility with the people whose trust they
  need most. One verified finding is worth more than five plausible ones.

## 0. Gather inputs

- If the user gave no doc, ask for the path, URL, or pasted text. Block.
- Read the whole doc before forming any finding.
- Read every linked doc you can reach.
- If you cannot open a link, list it under "Left out" in the report.
- If the doc names code, schemas, or configs in a repo you can read, read them.
- Find what the author asks of reviewers: approval, direction, or detail.
- If the doc does not say, ask the user. Continue as if the ask is approval.
- If the doc already has review comments, read them.

## 1. Restate the doc

Write three sentences: the problem, the proposed decision, and what it costs.

- If you cannot write them, make that the first finding. Name the sentence you
  could not write.
- Keep the restatement. It opens the summary in step 5.

The author reads the restatement first. A misread shows up there, before it
spreads into ten comments built on it.

## 2. Read through every lens

Take the lenses in order. Use every one. A lens that yields nothing goes on the
coverage list in the report, in one line.

1. **Problem and evidence.** Is the problem measured, or asserted? Whose problem
   is it? What does doing nothing cost, and over what period? Does each part of
   the solution trace back to a stated problem?
2. **Goals, non-goals, and success measures.** Can each goal be checked with a
   number or a yes/no test? What result would make the team stop or roll back?
   Does a non-goal quietly contradict a goal?
3. **Scope.** What is the smallest version that tests the riskiest assumption?
   Which parts are independent and bundled anyway? Which parts serve a
   requirement nobody stated? Does phase 1 deliver value if phase 2 never ships?
4. **Alternatives.** Is "do nothing" there? Is the boring option there: the
   existing system, a vendor, a cron job, a manual process? Is each alternative
   argued at its strongest, or set up to lose?
5. **Assumptions.** Which claims does the design depend on that the doc never
   states? Which stated ones are verified, and how? What is the cheapest test
   for each unverified one?
6. **Edge cases.** Walk the list at the end of this section against each flow.
7. **Failure modes.** For each step of each multi-step operation: what state is
   left if it fails here? Who notices? What does the user see? Can the
   operation be retried safely?
8. **Data and migration.** What happens to existing data that does not fit the
   new model? How long does the backfill take at real volume? Can the team roll
   back after new-shape data has been written? What about retention and
   deletion?
9. **Interfaces and one-way decisions.** Which consumers break? How are
   versions handled? Which decisions are expensive to reverse once shipped:
   public API shapes, storage formats, identifiers, vendor contracts? Those
   get the most scrutiny. Cheap-to-reverse decisions get the least.
10. **Rollout and rollback.** Is there a flag, a staged rollout, a kill switch?
    Where is the point of no return? Do old and new versions run at the same
    time, and does that work?
11. **Operations and cost.** What does on-call see when this breaks? Which
    alert fires? What does it cost at expected load and at ten times that?
    Redo every calculation in the doc. Check units.
12. **Security and privacy.** What new data is stored, and who can read it?
    Where are the trust boundaries, and what validates input at each? Does it
    touch PII, PHI, credentials, or audit requirements?
13. **People and sequencing.** Have the other teams this depends on agreed? Who
    owns it after launch? Does the timeline fit the scope and the named staff?
    What blocks what?
14. **Verification.** How will the team prove it works before launch? What can
    only be tested in production, and what limits the damage there?
15. **Consistency and clarity.** Do sections contradict each other? Does a term
    carry two meanings? Does the diagram match the text? Do the numbers agree
    across sections? Is an "open question" or "TBD" something the rest of the
    design depends on?

If the doc lacks a section a lens needs, that absence is a finding of kind
`gap`. Say what decision cannot be reviewed without it.

### Edge case list

Apply each line to every flow, entity, and interface in the doc. Keep the ones
the doc does not answer.

- Zero, one, many, and the maximum. Empty input. The largest tenant.
- The same request twice. A retry after a timeout that succeeded.
- Events out of order. Events arriving late. Events that never arrive.
- Two writers on the same record at once.
- Failure halfway through a multi-step operation.
- A dependency that is slow, down, or returns well-formed garbage.
- Clock skew, time zones, daylight saving, a job that runs over midnight.
- Old clients against the new server. New clients against the old server.
- Existing records created before the new rules applied.
- A user who cancels midway, submits twice, or goes back.
- Permission boundaries: another tenant, a revoked user, a deactivated account.
- Quotas, rate limits, and payload size limits, ours and the vendor's.
- A user or caller acting in bad faith.

## 3. Verify each finding

A finding the author can dismiss in one line hurts every other comment in the
review. Check each one before it goes in.

- Quote the sentence in the doc the finding is about.
- Search the whole doc for an answer to the finding. If the doc answers it,
  drop the finding.
- If the finding is a failure, write the scenario: the inputs, the sequence,
  and the wrong result.
- If the finding is about code or data you can read, check it there. Name the
  file and line.
- If the finding is about a number, show your calculation.
- If you cannot verify a finding, turn it into kind `question`. Say which
  answer would make it a problem.
- If you cannot make a finding concrete, drop it.
- If an existing review comment already makes the point, drop the finding.

Never invent a finding to fill a lens. A lens reported as clean is useful to
the user. A manufactured comment is not.

## 4. Classify, merge, and rank

Give each finding one kind:

| Kind | Means |
| --- | --- |
| `pushback` | You disagree with a decision the doc makes |
| `descope` | You propose cutting or deferring part of the work |
| `edge case` | A concrete scenario the design does not handle |
| `gap` | Information a reviewer needs is missing |
| `question` | You don't know, and the answer changes the verdict |
| `inconsistency` | The doc contradicts itself |

Give each finding one severity:

- **blocker.** Approval should wait. The design fails in a scenario that will
  happen, or a one-way decision lacks justification.
- **material.** It changes the design, the scope, or the plan. The author could
  still proceed and fix it on the way.
- **nit.** Wording, naming, a typo, a clearer diagram.

Then:

- Merge findings that share a root cause into one. Name the cause, then list
  the symptoms.
- Rank by severity, not by lens.
- Collect all nits into a single comment.
- Pick the verdict: `APPROVE`, `APPROVE WITH COMMENTS`, `CHANGES NEEDED`, or
  `NEEDS A CONVERSATION`.

Use `NEEDS A CONVERSATION` when you dispute the problem or the direction.
Comment threads do not settle that kind of disagreement. A call does.

There is no cap on the number of findings. Every finding that survived step 3
goes in, ranked, so the author can stop reading where they choose.

## 5. Draft the comments

The user posts these by hand. Never post them yourself.

### How to push back

- State the position. Write "I disagree with X because Y", not "have you
  considered Y?".
- Ask a question only when you don't know the answer. Say which answer changes
  your view.
- Give the scenario with real values: this request, then this failure, then
  this result.
- Propose the alternative. State what the alternative costs.
- Say what evidence would change your mind.
- When the doc gives the author's reason, answer that reason, not a weaker one.
- Hedge once, where you are unsure, and say what you are unsure about.
- Leave out praise that pads. Name a strength only if it tells the author
  where not to spend time.

A `descope` comment names four things: what to cut, what the cut saves, what it
loses, and when to bring the cut part back.

Worked examples of the voice.

Weak:

> Great doc! I was wondering if we've considered what happens with retries? It
> might be worth thinking about idempotency here.

Pushback:

> **4.2 Payment worker, "the worker retries on timeout"** [edge case, blocker]
> The charge call carries no idempotency key. If the provider accepts the
> charge and the response times out, the retry bills the customer a second
> time. I'd send the order ID as the key. If the provider already dedupes on
> another field, say so in this section and I'll withdraw the comment.

Descope:

> **6 Phasing, "Phase 1: rule engine, rule editor, audit log"** [descope, material]
> Only the rule engine tests the claim in section 2 that rules cut manual
> review by 40%. I'd ship the engine alone with rules in a config file, measure
> for two weeks, and build the editor if the number holds. That takes phase 1
> from ten weeks to about four. The cost is that ops needs a deploy to change a
> rule during those weeks.

### Report shape

```
## Review of <doc title>: <VERDICT>

<The restatement from step 1, in one or two sentences.>
<The two or three findings that decide the verdict.>

### Blockers

**<section heading>, "<short quote>"** [<kind>, blocker]
<comment>

### Material

**<section heading>, "<short quote>"** [<kind>, material]
<comment>

### Nits

**Nits**
- <section heading>: <nit>

---
Not for posting

### Coverage
- <lens>: <n findings | read, nothing to raise | not applicable, with reason>

### Left out
- <dropped finding or unreachable link>: <reason>
```

- Omit a severity heading that has no findings.
- If the doc is a file in a pull request, head each comment with the file path
  and line instead. Then pr-comments-draft can load the comments.
- If no finding survived step 3, say so in the summary and approve.

## 6. Unslop the draft

Humans read every word of this. Before you show the report:

- Read [unslop](../unslop/SKILL.md).
- Apply its rules to the summary and to every comment.
- If an unslop rule conflicts with the report shape, keep the shape.

The shape keeps its bold anchors, bracketed tags, and quotes. Unslop governs
the words inside the comments. The rules that fail most often in review
comments are 22 (sycophantic tone), 24 (hedging), 27 (name the mechanism or the
number), and 29 (active voice).

Then show the report.
