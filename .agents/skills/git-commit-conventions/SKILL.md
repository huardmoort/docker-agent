---
name: git-commit-conventions
description: |
  Conventional Commits format and atomic-commit rules for this repository:
  feat/fix/docs/chore/refactor/test prefixes, small focused commits with a
  one-line subject and an optional body, and the `Assisted-By: docker-agent`
  trailer for agent-authored work. Use when about to commit, when writing or
  reviewing commit messages, when squashing or reordering history, or when
  preparing a PR description.
---

# Commit Conventions

This repository uses [Conventional Commits](https://www.conventionalcommits.org/).

## 1. Subject line

```
<type>(<scope>): <imperative summary, lowercase, no trailing period>
```

Allowed types:

| Type | Use for |
|---|---|
| `feat` | New user-visible functionality |
| `fix` | Bug fix |
| `docs` | Documentation only |
| `refactor` | Code change with no behavioural impact |
| `chore` | Build, tooling, deps, CI, repo plumbing |
| `test` | Test-only changes |
| `perf` | Performance improvement |
| `security` | Security fix or hardening |

`<scope>` is optional but encouraged when the change is localised
(e.g. `feat(skills): ...`, `fix(runtime): ...`, `chore(docker-agent): ...`).

Keep the subject under ~72 characters.

## 2. Body

- Blank line after subject.
- Wrap at ~80 columns.
- Explain **why** the change is needed and any non-obvious decision.
- Don't restate the diff.

## 3. Trailer

Agent-authored commits end with:

```
Assisted-By: docker-agent
```

## 4. Atomicity

- One logical change per commit. If the message needs an "and", split.
- A commit must build and pass tests on its own (bisect-friendly).
- Group truly inseparable changes; never bundle unrelated fixes.

## 5. Using `but` (GitButler)

This project uses `but` for all write operations — never `git commit`/`add`/
`push`/`checkout`/`rebase`/`stash`/`merge` directly. Read-only `git log` /
`git blame` / `git show --stat` are fine.

Minimal flow:

```sh
but status -fv                                # gather change IDs
but commit <branch> -c -m "<message>" \
  --changes <id1>,<id2> --status-after        # -c creates the branch if new
```

When the message has a body, pass the whole multi-line message to a single
`-m` flag (a heredoc works well); `but commit` does **not** accept repeated
`-m` flags.

See the `but` skill for the full GitButler reference.

## 6. Pre-commit self-check

- [ ] Subject line in conventional format, imperative, no trailing period.
- [ ] Body explains the *why*, not the *what*.
- [ ] `Assisted-By: docker-agent` trailer present (agent commits).
- [ ] Commit is atomic (no "and" in the subject).
- [ ] `validate-go-change` passed if Go code changed.
