# Diátaxis — the four kinds of documentation

Diátaxis (https://diataxis.fr/) is the framework this project uses to
classify documentation. Every doc serves exactly **one** of these four
purposes; mixing them produces docs that fail at all four.

## The four quadrants

|              | Practical steps      | Theoretical knowledge   |
|--------------|----------------------|-------------------------|
| **Study**    | Tutorials            | Explanation             |
| **Work**     | How-to guides        | Reference               |

## 1. Tutorials (learning-oriented)

The reader is new and wants to learn the basics by doing.

- Lesson, not problem-solving. The reader follows your steps; you
  guarantee they end in a working state.
- Concrete and minimal. Skip alternatives, options, edge cases.
- One achievable outcome per tutorial.
- Voice: encouraging, hand-held.

> "First steps with docker-agent — your first multi-agent config in
>  10 minutes."

Anti-pattern: "Here are the 47 things you can configure" (that's
Reference).

## 2. How-to guides (task-oriented)

The reader has a specific problem and assumes some prior knowledge.

- Recipe-style: prerequisites, steps, result.
- Solves one problem; doesn't try to teach the underlying model.
- Skips background — links to Explanation when the reader needs more.
- Voice: direct, prerequisite-aware.

> "How to expose a tool over MCP."
> "How to bump a Go dependency safely."

Anti-pattern: open-ended walks through every option (Reference);
re-teaching the basics (Tutorial).

## 3. Reference (information-oriented)

The reader knows what they want to look up; they just need it
accurately.

- Exhaustive, structured, predictable layout.
- API surfaces, configuration fields, CLI flags, error codes.
- No tutorial-style narrative; no opinions.
- Voice: neutral, terse.

> "Agent configuration: every field of `pkg/config/latest`."
> "CLI reference."

Anti-pattern: storytelling, motivation, "why this is great"
(Explanation).

## 4. Explanation (understanding-oriented)

The reader wants to understand the *why* behind a design.

- Discursive, contextual, can include alternatives that were
  considered.
- Discusses tradeoffs, history, related concepts.
- Doesn't have to be runnable.
- Voice: thoughtful, opinionated where useful.

> "Why agents are configured in YAML, not code."
> "Tradeoffs between handoffs, fork-skills, and inline skills."

Anti-pattern: step-by-step procedure (How-to); exhaustive lists
(Reference).

## Decision rules

- If you find yourself writing "first, second, finally" → Tutorial or
  How-to. Pick one based on whether the reader knows the basics.
- If the doc is mostly a table of fields/flags/codes → Reference.
- If the doc is mostly prose discussing tradeoffs → Explanation.
- If you're tempted to do "all of the above" → split into multiple
  documents and cross-link.

## Project layout this maps to

| Diátaxis category | Lives under              |
|-------------------|--------------------------|
| Tutorial          | `docs/getting-started/`  |
| How-to guide      | `docs/guides/`           |
| Reference         | `docs/api/`, `docs/configuration/` |
| Explanation       | `docs/architecture/`, `docs/concepts/` |

`README.md` at the repo root is a special case: a *Tutorial* for the
absolute first encounter, plus index links into the four categories.
