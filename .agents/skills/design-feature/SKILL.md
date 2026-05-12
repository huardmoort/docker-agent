---
name: design-feature
description: |
  Produce a technical design for a new feature: gather requirements, analyse
  the existing codebase, propose architecture and interfaces, identify risks
  and tests, and write a specification document. Does not write
  implementation code. Use when the user asks to design a feature, plan an
  enhancement, write a technical spec, or types /design-feature.
context: fork
model: sonnet
---

# Design a Feature

Runs as a fork sub-agent on sonnet. Designs involve a lot of reading
(codebase exploration, prior-art research) and drafting; isolating that
work keeps the parent agent's context clean and folds back only the
final spec.

This skill is **planning only — do not write production code**. The fix
or implementation is a separate step, handled by the parent agent
(usually `coder`) after the spec is approved.

Composes with: `research-topic` (when prior art is needed).

## 1. Gather requirements

Ask only the questions the user hasn't answered yet:

- **Problem**: what user-visible problem does this solve?
- **Acceptance criteria**: what does "done" look like?
- **Scope**: what's explicitly in / out of scope?
- **Constraints**: backwards compatibility, perf, schema versioning,
  config compatibility, API stability.
- **Audience**: end users, library consumers, agent authors, all of
  the above?

If something fundamental is unclear, **stop and ask once** — don't
over-design around uncertainty.

## 2. Plan first

Output a one-screen plan and stop. Wait for explicit approval before
the longer exploration in step 3.

```
Design plan for "configurable retry policy"
- Phase 1: read pkg/runtime + pkg/config retry sites (read-only)
- Phase 2: survey similar features (sessions config, providers config)
- Phase 3: draft spec (architecture, types, migration, tests)
- Phase 4: STOP — present spec, await review before any code change
```

## 3. Analyse the codebase

Read first, write later. The aim is to discover the patterns the
project already uses, not to invent new ones.

- Identify the relevant packages and their public interfaces.
- Find similar existing features and follow their pattern.
- Note constraints: schema version (`pkg/config/latest`), agent config
  surface (`agent-schema.json`), examples conventions
  (`examples/*.yaml`).
- Check for existing prior art — issues, PRs, design docs.

If outside research is needed (libraries, articles, RFCs), invoke the
`research-topic` skill rather than browsing inside this skill.

## 4. Draft the specification

Write the spec to `.docker-agent/designs/<feature-slug>.md` (the
working directory; never committed). Required sections:

1. **Overview** — one paragraph: what, why, who for.
2. **Goals & non-goals** — bulleted, explicit.
3. **User-facing surface** — config YAML examples, CLI flags, MCP
   tool names — whatever the user sees first.
4. **Architecture** — components, data flow, interfaces. A mermaid
   diagram is recommended for any non-trivial design.
5. **Interfaces & types** — Go interfaces, JSON schema fragments, or
   protobuf — at the level of detail a reviewer can spot mistakes in.
6. **Schema and compatibility** — if `pkg/config/latest` changes,
   call out whether `agent-schema.json` and `examples/` need updates,
   and whether a config-version freeze is implied.
7. **Test strategy** — unit, integration, e2e. Name the new tests.
8. **Risks & mitigations** — security, perf, backwards compat,
   complexity. One sentence each.
9. **Implementation plan** — milestones with rough effort
   (`effort:small` / `medium` / `large`).
10. **Open questions** — explicit list. The parent agent (or user)
    answers these before implementation starts.

## 5. Self-check

Before presenting:

- [ ] Spec follows existing project patterns (you cited at least one).
- [ ] Backwards compatibility addressed (schema, config, examples).
- [ ] Test strategy is concrete (named tests, not "we'll add tests").
- [ ] Risks are honest, not boilerplate.
- [ ] Open questions clearly distinguish "needs user input" from
      "decide during implementation".

## 6. Present, await approval, fold back

Show the spec location and a 1-page summary in the response. Stop.
Wait for the user to approve, request changes, or close out open
questions.

When folding back to the parent:

```
Designed: configurable retry policy
- Spec: .docker-agent/designs/retry-policy.md
- Approach: per-toolset retry config in pkg/config/latest, read by
  pkg/runtime/dispatcher; behavioural default unchanged.
- Schema: agent-schema.json updated (new `retry` object on toolsets).
- Tests: TestRetry_Defaults, TestRetry_PerToolset, TestRetry_Backoff.
- Open: should retry apply to MCP tool calls too? (deferred)
```

Suggest, but don't perform, the next step:
- After approval, hand off to `coder` to implement, with this skill's
  spec path as the entry point.

## 7. When this skill is the wrong tool

- For pure bug fixes → `diagnose-bug`.
- For one-line config tweaks → just do them; don't write a spec.
- For changes to a frozen `pkg/config/v<N>` package → not a design
  question, it's "you can't"; explain instead.
- For schema-version freezing → use the `bump-config-version` skill.
