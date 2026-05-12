---
name: write-docs
description: |
  Write or update project documentation: pick the right Diátaxis category
  (tutorial, how-to, reference, explanation), draft clear markdown with
  working examples, place files in the right directory, and validate code
  samples against the actual code. Use when the user asks for documentation,
  a README update, an API reference, a how-to guide, a tutorial, or types
  /write-docs.
context: fork
model: sonnet
---

# Write Documentation

Runs as a fork sub-agent on sonnet. Doc work is read-heavy (browsing the
code, the existing docs, related issues) and edit-heavy (markdown
drafting); isolating it keeps the parent agent's context clean.

The Diátaxis framework guides choices: see
`references/diataxis.md` for the full taxonomy. Read it the first time
you use this skill, or when in doubt about which category fits.

## 1. Identify what's needed

Ask only what the user hasn't said:

- **Audience**: end user, agent author, contributor, maintainer?
- **Outcome**: should the reader learn / accomplish a task / look
  something up / understand a concept?
- **Surface**: README, `docs/`, code comment, CHANGELOG, example YAML?

These determine the Diátaxis category.

## 2. Pick the category

| Reader's goal     | Category    | Voice                |
|-------------------|-------------|----------------------|
| Learn the basics  | Tutorial    | Lesson, hand-held    |
| Solve a problem   | How-to      | Recipe, prerequisite-aware |
| Look something up | Reference   | Structured, exhaustive |
| Understand        | Explanation | Discursive, contextual |

If the topic spans more than one category, **split** the document.
Mixing them produces docs that fail at all four jobs.

## 3. Plan first

Output a short plan and stop:

```
Docs plan for "exposing skills as slash commands"
- Category: Reference (with one How-to subsection)
- Location: docs/features/skills.md (existing — append)
- Examples: validated against pkg/skills/ + pkg/tui/commands/
- STOP — show outline, await approval before drafting prose
```

Wait for approval before drafting the body.

## 4. Locate the file

| Surface                               | Where                       |
|---------------------------------------|-----------------------------|
| End-user docs                         | `docs/`                     |
| API reference                         | `docs/api/`                 |
| How-to guides / tutorials             | `docs/guides/`              |
| Architecture & design                 | `docs/architecture/`        |
| Project entry point                   | `README.md` (root)          |
| Per-feature working examples          | `examples/<topic>.yaml`     |
| Agent / contributor guidelines        | `AGENTS.md`                 |

If none of the above fits, the doc may not belong in this repo —
flag it.

## 5. Draft

Markdown conventions for this project:

- Code blocks always declare a language (`go`, `yaml`, `sh`, …).
- Use mermaid for any non-trivial diagram.
- Active voice, present tense, second person ("you").
- No filler ("comprehensive", "robust", "leverage", "delve into",
  "we are excited to").
- Keep paragraphs short. Tables for any 3+ columns of structured
  comparison.
- Examples should be runnable as written.

## 6. Validate examples

Code samples that look right but don't run are worse than no samples.
For each fenced example:

- **Go**: paste into a scratch `_test.go` and compile, or build the
  containing package.
- **YAML config**: validate against `agent-schema.json`:
  ```sh
  jq empty examples/<file>.yaml
  ./bin/docker-agent run --dry-run examples/<file>.yaml  # if available
  ```
- **Shell**: run it, copy real output into the doc; don't fabricate.

## 7. Cross-link

- Link to related sections of the same doc set.
- Link from the project README to top-level new docs.
- Update any docs index / sidebar configuration the project maintains
  (e.g. Jekyll `_data/`, Mintlify `docs.json`).

## 8. Self-check

Before presenting:

- [ ] One Diátaxis category per document; mixing flagged.
- [ ] All examples validated against real code.
- [ ] No filler words.
- [ ] File placed in the correct directory.
- [ ] Cross-links to siblings and from the index.
- [ ] If a feature doc, an example YAML accompanies it.

## 9. Fold back

Return a concise summary to the parent:

```
Documented: configurable retry policy
- File: docs/features/retry.md (new, Reference)
- Example: examples/retry.yaml (new)
- Cross-links: docs/configuration/agents.md, README.md feature list
- Validation: jq empty passed; example loads with --dry-run
```

Suggest, but don't perform, next steps (e.g. handoff back to `root`).
