# Position / Core Concept

This CLI is an agent-facing infra planning copilot for Claude Code, Codex, and similar agentic developer workflows.

## What it does

- Takes a vague infra request from an engineer or agent.
- Automatically discovers the current repo context.
- Inspects repo code and config first.
- Checks official/vendor docs next.
- Falls back to web/community examples only when needed.
- Asks the minimum clarifying question when blocked.
- Writes a single, polished markdown execution plan into the repo.
- Stops there for now.

## What it is not

- Not a human-facing tutorial tool.
- Not an executor yet.
- Not a generic search tool.
- Not a multi-plan system.

## Main value

- Turns “I need to provision X” or “I need to connect A to B” into an agent-ready, step-by-step plan.
- Packages the workflow as a reusable CLI instead of ad hoc prompting.
- Produces a clean `.md` artifact that Claude Code / Codex can review and execute.

## Recommended default shape

- The current working directory is the primary repo context.
- If the repo does not look right, the CLI asks the user to confirm or point elsewhere.
- If confidence is low, the CLI asks the smallest possible clarifying question.
- The output is plan-only for now.
- The plan is detailed, concrete, correctness-first, and optimized for eventual agent execution.
