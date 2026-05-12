---
name: diagnose-bug
description: |
  Diagnose and fix a bug end-to-end: reproduce the issue, narrow down the
  root cause (logs, stack traces, git bisect when needed), propose the
  smallest fix that addresses the root cause, write a regression test, and
  validate the change. Use when the user reports a bug, asks to investigate
  a failure, points at an error message or stack trace, or types
  /diagnose-bug.
context: fork
model: opus
---

# Diagnose a Bug

Runs as a fork sub-agent on opus. Bug investigation can involve many
exploratory tool calls (greps, file reads, test runs, git history
walks); isolating that traffic keeps the parent agent's context clean
and only the final diagnosis + fix folds back.

Composes with: `code-quality-checklist`, `validate-go-change`,
`git-commit-conventions`.

## 1. Capture the report

Pin down the facts before exploring. Ask only the questions the user
hasn't answered yet — don't re-ask things you can read off the report.

- **Symptom**: what was observed (error message, wrong output, crash).
- **Expected behaviour**: what should have happened.
- **Reproduction**: exact steps, command line, input, environment.
- **Scope**: always / sometimes / first time after upgrade / specific OS.
- **Artefacts**: stack trace, log excerpt, screenshot, failing test.

If the report is missing critical info, **stop and ask once** — don't
guess. Use the `triage-issue` rules to mark it `needs-info` if that's
where the request originated.

## 2. Plan first

Output a one-screen plan and stop. No file edits, no commits, no fix
attempts before approval. Read-only investigation in step 3 is fine.

```
Diagnosis plan for "TUI freezes on /clear"
- Phase 1 (read-only): reproduce locally, capture stack trace
- Phase 2 (read-only): narrow to a package, read the implicated code
- Phase 3 (read-only): form hypothesis, confirm with a focused test
- Phase 4 (write):     propose minimal fix + regression test
- STOP after each phase to report findings
```

## 3. Reproduce

Reproduction is the single most important step — without it you're
guessing.

- Run the exact steps from the report. Capture stdout / stderr.
- If the report mentions a test, run **that test alone**:
  ```sh
  go test -run <TestName> -count=1 -v ./pkg/<path>/...
  ```
- For TUI / runtime bugs, enable debug logging:
  ```sh
  ./bin/docker-agent run <config> --debug --log-file /tmp/repro.log
  ```
- If you cannot reproduce, **do not** continue to a fix. Report what
  you tried and stop.

## 4. Narrow the root cause

Use a layered approach — cheapest hypothesis first.

### 4.1 Read the error message

The stack trace usually points at the right package. Read it carefully.
Note the deepest non-test frame in the project's own code.

### 4.2 Search the codebase

```sh
# References to the symbol or message
grep -rn "<symbol-or-message>" --include="*.go" .
# Recent changes near the suspect file
git log -p --since="14 days ago" -- pkg/<path>/...
# Who last touched the failing line
git blame -L <start>,<end> pkg/<path>/<file>.go
```

### 4.3 Check recent changes

A regression that worked yesterday is almost always a recent commit.

```sh
git log --since="14 days ago" --oneline pkg/<path>/...
```

If a regression window is suspected and the test is deterministic,
**git bisect** is fastest:

```sh
git bisect start <bad-ref> <known-good-ref>
git bisect run go test -run <TestName> -count=1 ./pkg/<path>/...
```

### 4.4 Add focused logging

When the cause is non-obvious, add temporary `slog.Debug` lines or
`fmt.Printf` (clearly marked `// TEMP` so they're easy to remove)
around the suspect code path. Run the repro. Read the output. Form
the next hypothesis.

### 4.5 Form and confirm a hypothesis

State the cause in one sentence:

> "`session.Resume` panics because the `Sessions` map is read without
>  a read lock when `WithCleanup` runs concurrently with `Save`."

Confirm by reading the code (don't fix yet). If confirmation requires
running code, write a **failing test** that captures the bug — this
test becomes the regression test in step 6.

## 5. Plan the fix

Stop and present:

- Root cause (one sentence).
- Smallest change that addresses it (file paths, ~lines).
- Regression test (where it goes, what it asserts).
- Side effects considered (what else uses this code path).

Wait for approval before editing.

## 6. Implement

- Apply the **minimal** change — don't refactor adjacent code in the
  same commit. (If a refactor is warranted, do it in a separate PR.)
- Match surrounding style (`code-quality-checklist`).
- Write the regression test next to the existing tests for that
  package, named after the scenario, not the function.
- Remove any `// TEMP` debug statements added during step 4.

## 7. Validate

Run the full validation chain (`validate-go-change`):

```sh
task build
task test
task lint
```

All three must pass. Don't skip lint because "the change was small";
silent linter regressions are how style debt accumulates.

Also re-run the original repro from step 3 — confirm the symptom is
gone.

## 8. Document and report back

Fold this back to the parent agent:

```
Fixed: TUI freezes on /clear
- Root cause: race in pkg/session/manager.go:resumeWithCleanup
- Fix:        guard Sessions map with the existing RWMutex (5 lines)
- Test:       TestResumeWithConcurrentCleanup in manager_test.go
- Validation: task build / test / lint all green
```

Suggest, but do not perform, the next step:

- Commit message draft (using `git-commit-conventions`):
  `fix(session): guard Sessions map during Resume+Cleanup race`
- If the bug came from a GitHub issue, suggest `triage-issue` to update
  status and link the fix.

## 9. When to escalate, not fix

Stop and escalate to the user if:

- The fix requires an architectural change (escalate to `architect`).
- The bug spans more than ~3 files of non-trivial change (it's not a
  bug fix anymore — it's a refactor, and deserves a design).
- Reproduction requires credentials or infra you don't have.
- The repro cannot be created — never ship a "fix" without a repro.
