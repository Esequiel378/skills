---
name: pr-comments-address
description: Use when the user wants a PR's review comments handled — triggers on "/pr-comments-address", "address the review comments", "handle the reviewer feedback", "apply the review suggestions", or pointing at a PR with unresolved comments and asking to deal with them. Output is local working-tree changes plus a per-comment verdict report — never commits, pushes, or replies.
---

# Address PR Comments

Reviewer comments in, verdicts and local fixes out. Every comment gets judged —
not obeyed — and the accepted ones get fixed in the working tree. **This skill
never writes to the PR.** Committing, pushing, and replying are the user's
moves; the deliverable here is a dirty tree and a report they can act on.

## 0. Resolve the PR, fetch every unresolved comment

- **PR number / URL given** → use it. **Nothing given** → the current branch's
  PR: `gh pr view --json number,headRefOid,baseRefName,url`.
- **Is the local branch what the PR shows?** Compare `git rev-parse HEAD` to the
  PR's `headRefOid`. Comments anchor to lines in the PR's head; a diverged local
  branch means fixing code the reviewer never commented on — stop and say so.
- Fetch all three comment surfaces (each misses the others):
  - **Inline threads, with resolution state** — REST can't tell you what's
    resolved, so use GraphQL:
    ```bash
    gh api graphql -f query='query($o:String!,$r:String!,$n:Int!,$after:String){
      repository(owner:$o,name:$r){pullRequest(number:$n){author{login}
        reviewThreads(first:100,after:$after){pageInfo{hasNextPage endCursor}
          nodes{id isResolved isOutdated path line
            comments(first:20){pageInfo{hasNextPage endCursor}
              nodes{author{login} body}}}}}}}' \
      -f o=<owner> -f r=<repo> -F n=<pr>
    ```
    - If `reviewThreads.pageInfo.hasNextPage` is true, repeat the query with
      `-f after=<endCursor>`. Stop when it is false.
    - If a thread's `comments.pageInfo.hasNextPage` is true, query that thread
      by its node id. Page with `after:` until `hasNextPage` is false.
      ```bash
      gh api graphql -f query='query($t:ID!,$after:String){node(id:$t){
        ... on PullRequestReviewThread{comments(first:20,after:$after){
          pageInfo{hasNextPage endCursor} nodes{author{login} body}}}}}' \
        -f t=<thread id>
      ```
  - **Review bodies** — `gh api repos/<o>/<r>/pulls/<n>/reviews`
  - **Conversation comments** — `gh api repos/<o>/<r>/issues/<n>/comments`
- Work the **unresolved** threads; skip resolved ones and bot chatter unless the
  user says otherwise. Outdated-but-unresolved threads still count — check
  whether the concern survived the code that outdated them.
- **Thread state.** Read each thread's last comment before triage.
  - If the reviewer conceded in it, skip the thread.
  - If the PR author already answered in it, skip the thread.
  - If the user asks for all threads, skip nothing.
  - Record each skipped thread, with the reason, for the report.

## 1. Triage — every comment gets exactly one verdict

Read the code each comment anchors to — the file as it is, not just the diff
hunk; the surrounding lines are where wrong verdicts come from. Then assign:

| Verdict | Meaning |
|---|---|
| **accept** | Valid, and this PR's job → gets fixed in step 3 |
| **descope** | Valid, but not this change's job → name the follow-up it belongs in |
| **move** | Valid, but belongs to a different PR → name which one |
| **reject** | Wrong or moot → cite the evidence (code, test, doc) that says so |

- **Stacks:** list the PRs stacked on this one with
  `gh pr list --base <head branch>`, then follow the chain, as
  [pr-review](../pr-review/SKILL.md) does. A comment left on this PR about code
  introduced in a *parent* layer is a **move** — fixing it here buries the fix
  in the wrong diff and dooms the restack. Name the layer it belongs to.
- **Scope:** "while you're here, also…" comments that widen the PR beyond its
  stated intent default to **descope**. A PR that grows in review is a PR that
  never merges.
- **reject needs receipts.** A reviewer is a human who read the code; overruling
  them on vibes is how valid feedback gets lost. No evidence → the verdict is
  accept or descope, not reject.

## 2. Cross-validate contested verdicts with pr-review

Before finalizing any **reject**, and for any **accept** whose fix is invasive,
run the [pr-review](../pr-review/SKILL.md) skill scoped to the code in question.

- Collect every contested comment first. Run pr-review once for all of them.
- Pass `--lenses=<keys>`. Pick the keys that match each contested comment's
  topic, using pr-review's lens keys. Take the union across comments.
- Example: a comment about an unchecked input maps to `security`. A comment
  about a test that proves nothing maps to `po`.

The lenses are the second opinion: one independently re-raising the reviewer's
point flips a reject to accept; all of them passing the code silently is the
evidence a rejection can cite. Never reject a human's comment on your own
authority alone.

## 3. Fix the accepted ones — working tree only

- Fix the **root cause**, not just the line the comment names — grep the
  sibling callers; a comment about one call site is usually about all of them.
- Run the tests and report the actual output — not "should still pass".
- Leave everything **uncommitted**. A dirty tree is the deliverable.

## 4. Report

````
## comments — PR #<n> (<addressed>/<total> unresolved threads)

| comment (file:line — author) | verdict | why / what changed |
|---|---|---|
| ... | accept | fixed in <file:line>: <one line> |
| ... | descope | belongs in <follow-up>: <one line> |
| ... | move | belongs in PR #<m> / stack layer <branch> |
| ... | reject | <evidence> |

Skipped threads (reviewer conceded or author already answered):
- <file:line — author>: <reason>

Tests: <command> → <actual result>

## reply drafts

<file:line>

```markdown
<reply>
```

Next moves are yours: review the diff, commit, push, reply.
````

Every unresolved comment appears in the table or in the skipped list — a
comment silently skipped reads as a comment addressed.

### Reply drafts

Write one reply draft per judged thread. Skipped threads get none.

- Put the `file:line` heading outside the box. Put the reply in its own
  ` ```markdown ` box.
- If a reply holds a code fence, open its box with four backticks.
- Write each reply in the voice of pr-review's
  [comment drafts](../pr-review/SKILL.md#comment-drafts).
- If the verdict is **accept**, say what changed.
- If the verdict is **descope** or **move**, say where it goes.
- If the verdict is **reject**, give the evidence.
- Never post a draft.

[pr-comments-draft](../pr-comments-draft/SKILL.md) can't post thread replies
yet. It adds new comments to a pending review. The user pastes each reply
into its thread.

## Guardrails

**Forbidden — no exceptions, regardless of urgency or convenience:**

- `git commit`, `git push` — in any form
- `gh pr comment`, `gh pr review`, reply endpoints, resolving threads,
  re-requesting reviewers, reactions — anything that speaks on the PR

| Excuse | Reality |
|---|---|
| "Reviewers are waiting — I'll just push so they see it" | The user reviews the fixes before anyone else does. Speed pressure is a reason to triage faster, not to publish. |
| "One quick reply so they know it's handled" | A reply speaks *as the user* without their consent. Forbidden even as a courtesy. |
| "A local commit is harmless" | The user may want to reshape the branch history. The dirty tree *is* the output. |

- **Report what you skipped** — a thread you couldn't parse, a file too large, a
  test you couldn't run. Silent partial coverage reads as full coverage.
