---
name: pr-comments-draft
description: Use when the user wants the review comments from this chat loaded onto the PR as a pending review — triggers on "/pr-comments-draft", "draft these on the PR", "add them as pending comments", "put the comments on GitHub but don't submit". Reads the latest drafts from the conversation, usually a pr-review run plus edits, and adds them to the user's pending review. Never submits the review.
---

# Draft PR Comments

The chat already holds the review. Usually a [pr-review](../pr-review/SKILL.md)
run, then some back and forth on which findings stay and how they read. This
skill moves that final set onto the PR as a **pending review**. A pending
review is visible only to the user until they press "Submit review" in GitHub.
That press is theirs. This skill never makes it.

## 0. Collect the drafts from the chat

- Find the most recent comment drafts in this conversation.
- If the user edited, dropped, or added a comment afterwards, use the latest
  version. The last word in the chat wins over the original draft.
- If a finding was listed as left out, leave it out.
- If the chat holds only a pr-review agent report, rewrite each finding in the
  pr-review comment-drafts voice first.
- If the chat holds no review of a PR, stop. Tell the user to run pr-review
  first.

Each draft needs a path, a line, and a body. The review-body summary is
optional.

## 1. Resolve the PR and its head

- If the chat names the PR, use it. Otherwise use the current branch's PR.
- Run `gh pr view <n> --json number,headRefOid,url`.
- If the head SHA moved since the review, block. Say what changed with
  `git log <old>..<new> --oneline`. Ask whether to re-verify the findings
  first.

## 2. Map every draft to a diff line

GitHub only accepts inline comments on lines inside the PR diff.

- Convert each path to repo-root-relative. pr-review headings may be
  package-relative.
- Fetch the patch with `gh pr diff <n>`.
- Check that each line falls inside a `+` or context line of a hunk, on the
  new side.
- If a line falls outside every hunk, don't drop the comment. Move it into the
  review body as `` `path` line N: <body> ``.

## 3. Add the comments to the pending review

GitHub allows one pending review per user per PR. Check for one first:

```bash
gh api --paginate repos/<o>/<r>/pulls/<n>/reviews \
  --jq '.[] | select(.state=="PENDING") | .node_id'
```

If none exists, create it with every comment in one call. Omit `event`. That
omission is what keeps the review pending.

```bash
gh api repos/<o>/<r>/pulls/<n>/reviews --input - <<'EOF'
{
  "commit_id": "<headRefOid>",
  "body": "<review-body summary>",
  "comments": [
    {"path": "<repo-root path>", "line": 142, "side": "RIGHT", "body": "<comment>"}
  ]
}
EOF
```

If one exists, add each comment to it. Leave its existing comments alone.

```bash
gh api graphql -f query='mutation($id:ID!,$path:String!,$line:Int!,$body:String!){
  addPullRequestReviewThread(input:{pullRequestReviewId:$id,path:$path,
    line:$line,side:RIGHT,body:$body}){thread{id}}}' \
  -f id=<node_id> -f path=<path> -F line=<line> -f body=<comment>
```

If the pending review already exists, don't touch its body. Show the
review-body summary for the user to paste at submit time.

Build the payload with `jq -n --arg`. Never write the JSON by hand. Comment
bodies hold backticks, quotes, and newlines that break hand-built strings.

If a call returns 422, show the error and the draft it rejected. Continue
with the rest.

## 4. Report

- Link the PR's "Files changed" tab.
- List each comment added, as `path` line N.
- List each comment moved into the review body, and why.
- List each comment GitHub rejected, with its draft text.
- Remind the user that nothing is visible to others until they submit.

## Guardrails

- **Never submit.** Never pass `event`. Never call the
  `/reviews/<id>/events` endpoint or the `submitPullRequestReview` mutation.
- **Never delete** a pending review or a comment the user already had in it.
- **Never invent a comment.** Every body comes from the chat. If a draft is
  ambiguous, ask.
