---
name: code-quality-checklist
description: |
  Quality bar for source-code changes in this repository: clean self-documenting
  code, minimal but meaningful comments, proper error handling, edge cases, and
  matching tests. Use when writing new code, reviewing code, or judging whether
  a change is ready for review.
---

# Code Quality Checklist

Apply these standards to every code change you produce or review.

## 1. Style and structure

- Follow the existing code style and patterns in the area you touch. Match,
  don't reinvent.
- Functions are well-scoped and reasonably sized. Extract when a function does
  more than one thing.
- Names are clear and descriptive. Avoid abbreviations the rest of the file
  doesn't use.
- Keep code DRY, but don't over-abstract: one duplication is acceptable, three
  is a refactor.

## 2. Comments

- Self-documenting code first. Comments only when the code's purpose or logic
  is not immediately evident.
- Never write a comment that restates what the code does
  (e.g. `// increment counter` above `counter++`).
- Comments explain **why**, not what: non-obvious edge cases, why an algorithm
  was chosen, references to issues or specs.

## 3. Errors and edge cases

- Return errors, don't swallow them. Wrap with context when crossing layers.
- Validate inputs at the boundary; trust them inside.
- Think through the failure modes: empty inputs, nil pointers, timeouts,
  cancellation, concurrent access.

## 4. Tests

- New behaviour requires new tests. Bug fixes require a regression test.
- Cover the happy path **and** at least one error / edge case.
- Tests use the project's existing assertion library
  (`github.com/stretchr/testify/assert` and `require` for Go).
- Test names describe the scenario, not the function under test.

## 5. Maintainability

- The change is small enough that a reviewer can hold it all in their head.
- Public interfaces are stable and minimal. Don't expose what callers don't need.
- No commented-out code, no `TODO` without a tracking issue.

## 6. Self-check before declaring done

- [ ] Style matches surrounding files.
- [ ] Errors handled meaningfully.
- [ ] At least one new test per behavioural change.
- [ ] No restate-what-the-code-does comments.
- [ ] No dead code, no debug prints.

When validation tooling is needed (build / test / lint), use the
`validate-go-change` skill.
