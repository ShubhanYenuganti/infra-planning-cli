# Output Format

The CLI should generate a single markdown file per plan, stored in the repo under:

- `docs/<cli-name>/`

## File naming

Use a timestamp plus slug.

Example:

- `docs/<cli-name>/2026-05-15-provision-networking.md`

## Mandatory sections

- Problem summary
- Recommended path
- Steps
- Verification

## Style

- Superpowers-style.
- Structured and directive.
- Human-readable first, but suitable for Claude Code / Codex to execute.
- Very detailed.
- Include concrete command snippets / IaC examples for each step.
- Keep the final plan clean and citation-free.
- The structure should remain flexible for future augmentations, but for now it stays plan-only.

## Verification section

Include:

- exact commands
- expected pass criteria

This creates a clean artifact that a human can review and an agent can execute from.
