# Project Agent Guidelines

This document is loaded by every agent in this team's configuration via
`add_prompt_files`. Rules below apply to all agents.

---

## Plan-First Operating Rule

**Always present a plan and wait for explicit user approval before any
side-effecting action.**

Side-effecting actions (require approval):

- creating, modifying, or deleting files
- running build, test, lint, or format commands
- committing, pushing, branching, opening or merging pull requests
- calling MCP tools that mutate state (GitHub issue/PR edits, etc.)

Read-only actions do **not** require approval:

- reading files, searching, listing directories
- inspecting git or GitButler state
- fetching documentation, web searches
- thinking and drafting plans

Plan format:

1. A numbered list of concrete steps with the exact files or commands.
2. The expected outcome of each step.
3. End with `STOP. Awaiting approval.`
4. After approval, execute one step at a time and report after each.

This rule overrides any conflicting guidance in agent instructions or
skills.

---

## Communication Style

- Direct and efficient. No filler phrases, no excessive affirmations.
- Avoid LLM clichés: "comprehensive", "robust", "leverage", "delve into",
  "I'd be happy to", "great question".
- One clear point per message; split if you have two.
- Report results concisely; don't restate what the user just said.
- Use markdown tables for any structured comparison of three or more
  items.
- When generating code, do not paste it back in the chat — write it to
  a file. Status updates only in chat.

---

## Skill Catalogue

The team maintains a catalogue of skills under `.agents/skills/`. Prefer
invoking a skill over re-deriving its procedure inline. Skills surface as
slash commands in the TUI (`/skill-name`).

**Shared recipes (inline, short)**

| Skill | Use when |
|---|---|
| `code-quality-checklist` | Writing, reviewing, or judging readiness of code |
| `git-commit-conventions` | Writing or reviewing commit messages |
| `validate-go-change` | After a Go change, before declaring done |
| `pr-comment-style` | Drafting PR review comments |

**Workflow skills (fork sub-agents)**

| Skill | Use when |
|---|---|
| `triage-issue` | Classifying a single GitHub issue or PR |
| `triage-backlog` | Batch-triaging many open issues at once |
| `review-pr` | Performing a full PR review end-to-end |
| `diagnose-bug` | Investigating and fixing a reported bug |
| `design-feature` | Producing a technical design / spec |
| `write-docs` | Writing or updating project documentation |
| `research-topic` | Searching the web for prior art, docs, examples |
| `bump-config-version` | Freezing the latest config schema and bumping the version |
| `bump-go-dependencies` | Updating Go module dependencies one at a time |

When a skill applies, the agent should invoke it instead of inlining the
procedure. Skills compose: many of them reference each other.

---

## Code Quality Standards

The full quality bar lives in the `code-quality-checklist` skill. The
short version:

- Self-documenting code; comments only when the *why* isn't obvious.
- Never write a comment that restates what the code does.
- Proper error handling, considered edge cases, matching tests.
- Match the surrounding style.

---

## Working Approach

- Use tools to gather information rather than guessing.
- Examine existing code before making changes.
- Validate every change before declaring complete.
- Ask clarifying questions only when truly necessary.
- Call independent tools concurrently when possible — it's faster.

---

## Validation Requirements

For Go changes, run the `validate-go-change` skill (or directly):

- [ ] `task build` succeeds
- [ ] `task test` passes
- [ ] `task lint` reports no new issues
- [ ] Acceptance criteria met
- [ ] Project patterns followed
- [ ] Edge cases considered

For non-code changes, validate appropriately (e.g. `jq empty` for YAML,
markdown link checks for docs).

---

# Development Commands

## Build and Development

- `task build` — Build the application binary (outputs to `./bin/docker-agent`)
- `task test` — Run Go tests (clears API keys to ensure deterministic tests)
- `task lint` — Run golangci-lint (uses `.golangci.yml` configuration)
- `task format` — Format code using golangci-lint fmt
- `task dev` — Run lint, test, and build in sequence

## Docker and Cross-Platform Builds

- `task build-local` — Build binary for local platform using Docker Buildx
- `task cross` — Build binaries for multiple platforms (linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64, windows/arm64)
- `task build-image` — Build Docker image tagged as `docker/docker-agent`
- `task push-image` — Build and push multi-platform Docker image to registry

## Running docker-agent

- `./bin/docker-agent run <config.yaml>` — Run agent with configuration (launches TUI by default)
- `./bin/docker-agent run <config.yaml> -a <agent_name>` — Run specific agent from multi-agent config
- `./bin/docker-agent run agentcatalog/pirate` — Run agent directly from OCI registry
- `./bin/docker-agent run --exec <config.yaml>` — Execute agent without TUI (non-interactive)
- `./bin/docker-agent new` — Generate new agent configuration interactively
- `./bin/docker-agent new --model openai/gpt-5` — Generate with specific model
- `./bin/docker-agent share push ./agent.yaml namespace/repo` — Push agent to OCI registry
- `./bin/docker-agent share pull namespace/repo` — Pull agent from OCI registry
- `./bin/docker agent serve mcp ./agent.yaml` — Expose agents as MCP tools
- `./bin/docker agent serve a2a <config.yaml>` — Start agent as A2A server
- `./bin/docker agent serve api` — Start docker-agent API server

## Debug and Development Flags

- `--debug` or `-d` — Enable debug logging (logs to `~/.cagent/cagent.debug.log`)
- `--log-file <path>` — Specify custom debug log location
- `--otel` or `-o` — Enable OpenTelemetry tracing
- Example: `./bin/docker-agent run config.yaml --debug --log-file ./debug.log`

# Testing

- Tests are located alongside source files (`*_test.go`)
- Run `task test` to execute the full test suite
- E2E tests live in the `e2e/` directory
- Test fixtures and data live in `testdata/` subdirectories
- Use `github.com/stretchr/testify/assert` and `require` for assertions
- Cover edge cases and error conditions
- Mock external dependencies for unit tests

# Agent Config YAML

- Agent config files follow a strict schema: `./agent-schema.json`
- The schema is **versioned**
- `./pkg/config/v0`, `./pkg/config/v1`, ... packages handle older versions of the config
- `./pkg/config/latest` package handles the current, work-in-progress config format
- When adding new features to the config, **only add them to the latest config**
- Older config types are **frozen** — do not modify them
- When adding new features to the config:
  - Update `./agent-schema.json` accordingly
  - Create an example YAML that demonstrates the new feature

# Git Practices

This project uses GitButler (`but`) for all version-control write
operations. Never run `git commit`, `git add`, `git push`, `git checkout`,
`git rebase`, `git stash`, or `git merge` directly. Read-only `git log` /
`git blame` / `git show --stat` are fine.

The full conventions live in the `git-commit-conventions` skill. The
short version:

- Conventional Commits subjects: `feat:`, `fix:`, `docs:`, `chore:`,
  `refactor:`, `test:`, `perf:`, `security:`.
- Atomic commits — one logical change per commit, no "and" in subjects.
- Body explains the *why*, not the *what*.
- Agent-authored commits end with `Assisted-By: docker-agent`.
- Branches focused on a single feature or fix; up to date before submit.
