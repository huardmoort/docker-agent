---
name: pr-comment-style
description: |
  How to write pull-request review comments in this repository: prefer inline
  over global, prefix with Blocking/Non-blocking/Question/Nit/Praise, be
  factual and concise, critique the code not the author. Use when leaving
  feedback on a PR, drafting review comments, or coaching another reviewer.
---

# PR Comment Style

## 1. Inline first, global second

- **Inline comment on the relevant line(s)** for anything that targets
  specific code. The reader sees the comment in context.
- **Global comment** only for cross-cutting concerns or the overall verdict.

## 2. Prefixes

Start each comment with one of:

| Prefix | Meaning |
|---|---|
| **Blocking** | Must be fixed before merge. State the concrete required change. |
| **Non-blocking** | Worth doing but not gating. The PR can merge without it. |
| **Question** | You don't yet have an opinion; you need information. |
| **Nit** | Style or formatting; explicitly optional. |
| **Praise** | Highlight a good solution. Use sparingly so it stays meaningful. |

## 3. Voice and length

- One clear point per comment. Split if you have two.
- Factual and direct. No filler ("just wanted to mention", "I was
  wondering if maybe…").
- No LLM verbosity. No "comprehensive", "robust", "leverage", "delve into".
- Polite by being concise: critique the code, not the author.
- For Blocking comments, state the *required* change, not three vague
  suggestions.

## 4. Templates

Blocking:

> **Blocking** — `foo()` returns the unwrapped error from `os.Open`; wrap
> with `fmt.Errorf("read config: %w", err)` so callers can `errors.Is` it.

Non-blocking:

> **Non-blocking** — this loop is O(n²); a map keyed by `ID` would make it
> linear. Fine as-is for current dataset sizes.

Question:

> **Question** — is there a reason this bypasses the `Validator` interface?
> Curious whether the validator path was considered.

Nit:

> **Nit** — variable spelled `recieved`; should be `received`.

Praise:

> **Praise** — nice use of `t.Cleanup` here, much clearer than the previous
> defer pattern.

## 5. Verdict comment (global)

When you finish a review, leave one global comment summarising the verdict:

- **Approve** — all Blocking resolved, ready to merge.
- **Request changes** — list the Blocking items by file/line.
- **Comment** — observations only; no merge decision yet.

Keep this summary under 8 lines. Detail belongs in the inline comments.
