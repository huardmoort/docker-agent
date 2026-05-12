---
name: review-pr
description: |
  Perform a complete pull-request review: fetch the PR and its diff, walk
  through the code-quality / testing / security / documentation / git-hygiene
  checklist, draft inline comments per file, post them via the GitHub MCP,
  and submit a verdict (Approve / Request changes / Comment). Use when the
  user asks to review a PR, asks for feedback on a PR, asks whether a PR is
  ready to merge, or types /review-pr.
context: fork
model: opus
---

# Review a Pull Request

This skill runs as a fork sub-agent on opus. The MCP traffic for fetching a
diff and posting many inline comments would otherwise crowd the parent
agent's context window; the fork keeps it isolated and folds back only the
final verdict.

The comment-style rules in `pr-comment-style` and the code-quality bar in
`code-quality-checklist` apply throughout. Read them if you haven't yet.

## 1. Identify the PR

Accept any of these inputs:

- a PR number (`#123`)
- a PR URL (`https://github.com/owner/repo/pull/123`)
- "the PR I'm working on" → use the current branch and resolve via
  `gh pr view --json number,url`

Confirm the resolved PR back to the user before fetching anything else.

## 2. Plan first

Output a one-screen plan and **stop**. Wait for explicit approval before
running any side-effecting MCP call (posting comments, changing labels,
submitting a verdict).

```
Review plan for PR #123 — "Add Bedrock streaming"
- Files changed: 7  (+312 / -41)
- Areas:        area/providers/bedrock, area/api
- Tests added:  yes (3 new test files)
- Approach:
  1. Fetch diff and CI status
  2. Walk each file (5 are .go, 1 is yaml example, 1 is docs)
  3. Apply checklist (style / errors / tests / security / docs / git)
  4. Draft inline comments + a global verdict
  5. STOP, show the comment plan, await approval
  6. After approval: post comments and submit verdict
```

## 3. Fetch context

Read-only calls — no approval needed for these:

```sh
gh pr view <n> --json number,title,author,baseRefName,headRefName,labels,files,statusCheckRollup,reviewDecision
gh pr diff <n>
gh pr checks <n>
```

Or via the GitHub MCP: `get_pull_request`, `get_pull_request_files`,
`get_pull_request_status`.

Capture: title, author, target branch, files changed (with patches), CI
status, existing reviews, linked issues.

## 4. Apply the review checklist

Walk these five categories. Don't write comments yet — collect findings.

### 4.1 Code quality

(See `code-quality-checklist`.)

- Style and patterns match surrounding code.
- Functions reasonable size, names clear.
- Error handling proper (wrap, don't swallow).
- No commented-out code, no dead code.
- No restate-the-code comments.

### 4.2 Testing

- New behaviour has new tests; bug fixes have a regression test.
- Edge cases and error paths covered.
- Test names describe the scenario.
- All CI tests pass (check `gh pr checks`).

### 4.3 Security

- No hardcoded credentials, tokens, or PII.
- Inputs validated at trust boundaries.
- No obvious injection vectors (SQL, shell, path traversal).
- Dependencies from trusted sources; no surprise additions.

### 4.4 Documentation

- Public APIs documented.
- Non-obvious logic has a `why` comment.
- README updated if user-visible behaviour changed.
- `docs/` updated if end-user surface changed.
- `examples/` reflect the change if relevant.
- `agent-schema.json` updated if the configuration model changed.
- CHANGELOG updated (if the project keeps one).

### 4.5 Git hygiene

- Conventional-commit subjects (see `git-commit-conventions`).
- Commits atomic; no "and" in subjects.
- `Assisted-By: docker-agent` trailer if AI-authored.
- Branch up to date with target; no merge conflicts.

## 5. Draft comments

Use the prefixes from `pr-comment-style`: **Blocking**, **Non-blocking**,
**Question**, **Nit**, **Praise**. One concern per comment.

- Prefer **inline** on the relevant line(s).
- Use **global** only for the verdict and cross-cutting concerns.

Group findings by file and by severity. The output of this step is a
"comment plan" — a bullet list of (file, line, prefix, body) — not yet
posted.

Also draft the **verdict comment** (one global post):

- **Approve** — all Blocking resolved, ready to merge.
- **Request changes** — list Blocking items by file/line.
- **Comment** — observations only, no merge decision.

Keep the verdict under 8 lines. Detail goes into the inline comments.

## 6. Show the plan, await approval

Print:

- Number of comments by prefix (`3 Blocking, 2 Non-blocking, 1 Nit`).
- The verdict.
- The comment plan as a markdown table.

Then **stop**. Wait for the user to say go before any posting call.

## 7. Post (after approval)

Post inline comments first, then the verdict, via the GitHub MCP:

- `add_pull_request_review_comment` (or `create_pending_pull_request_review`
  + `submit_pending_pull_request_review` if the MCP supports batch reviews
  — preferred so all comments land as one review).
- Final verdict: submit the review with state `APPROVE`,
  `REQUEST_CHANGES`, or `COMMENT`.

Optional, only if the user asked for it:

- Apply / adjust labels (use `triage-issue` rules — usually a `kind/*` and
  one or more `area/*`).
- Link related issues (`closes #...` belongs in the PR body, not in a
  review comment).

## 8. Report back

Return a concise summary to the parent agent:

```
Reviewed PR #123. Verdict: REQUEST_CHANGES.
- 3 Blocking: error wrapping in providers/bedrock/client.go (3 sites)
- 2 Non-blocking: extract helper, simplify retry loop
- Tests pass (CI green). Docs updated.
```

That's the only thing folded back into the parent conversation — the
inline comment text stays on GitHub.

## 9. Decide on follow-ups (suggest, don't do)

If your verdict was REQUEST_CHANGES, suggest handing off to `coder` with
the Blocking list. If APPROVE but docs are missing, suggest `doc_writer`.
Don't auto-handoff — surface the suggestion and let the parent decide.
