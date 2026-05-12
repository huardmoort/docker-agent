---
name: validate-go-change
description: |
  Run the project's standard validation sequence (`task build`, `task test`,
  `task lint`) before declaring a Go change done. Use after implementing or
  fixing code in this repository, before opening a PR, before committing
  anything beyond docs or trivial chores, or whenever a clean-build sanity
  check is required.
---

# Validate a Go Change

Run these three commands in order. Fix issues as they appear; don't move on
until the current step is green.

## 1. Build

```sh
task build
```

- Output binary lands at `./bin/docker-agent`.
- A failure here is a compile error. Read the message, fix, re-run.

## 2. Test

```sh
task test
```

- Clears API keys to keep tests deterministic.
- A failure here is either a real regression or a flaky test. Run the
  failing package alone to confirm:
  ```sh
  go test -run <TestName> -count=1 ./pkg/<path>/...
  ```
- New behaviour requires new tests (see `code-quality-checklist`).

## 3. Lint

```sh
task lint
```

- Uses `.golangci.yml`. Fix all new findings; don't introduce warnings.
- If a finding is genuinely incorrect, justify it in code with a focused
  `//nolint:` and a short reason — do not blanket-disable.
- For purely formatting-related findings:
  ```sh
  task format
  ```

## 4. Reporting

Report the outcome concisely:

- ✅ `build`, ❌ `test` (1 failure in pkg/foo) — paste the failing assertion
- ✅ all three — proceed

Never claim success without having run all three. Never run any of them in
the background and assume they passed.

## 5. When to skip

Skip this skill only when the change touches **none** of: Go source files,
`go.mod`/`go.sum`, taskfile, golangci config. Doc-only, comment-only, and
example-yaml-only changes do not require validation.
