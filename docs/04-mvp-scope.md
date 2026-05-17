# MVP Scope

## One-sentence MVP

Turn a vague infra request into a detailed, step-by-step markdown execution plan.

## First demo target

- Provisioning X
- Hero example focused on networking
- Still supports compute and databases in the combined flow

## Target environment

- Generic across repos/services
- Current repo only for discovery
- Optimized for AWS first
- Helpful for both Claude Code and Codex

## Behavior constraints

- Plan-only for now.
- Conservative when unsure.
- Ask the minimum clarifying question.
- Never ask for credentials.
- Prefer repo context, then docs, then web/community.
- Output should be the best possible plan, even if it takes longer.

## Summary framing

An agent-facing CLI that turns vague infra asks into a high-confidence, repo-specific execution plan.
