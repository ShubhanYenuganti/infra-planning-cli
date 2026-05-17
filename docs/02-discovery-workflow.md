# Discovery Workflow

This CLI should run a conservative, correctness-first discovery flow before writing a plan.

## Flow

1. Start in the current working directory.
2. Systematically inspect the repo for:
   - infra code paths
   - config and IaC definitions
   - existing patterns and examples
   - tests/examples only after stronger signals are checked
3. Search official/vendor docs on the web.
4. If repo + docs still leave gaps, search for high-signal community examples.
5. Ask the minimum clarifying question needed if a blocker remains.
6. Never ask for credentials.
7. Once the key unknowns are resolved, write the plan.

## Evidence priority

1. Repo code paths and patterns
2. Repo config / IaC
3. Official docs / vendor docs
4. Community examples

## Conflict rule

If repo evidence, metadata, and web results disagree, surface the conflict and ask a clarifying question.

## Stopping rule

Stop discovery once the key unknowns for the task are answered.

## Behavioral posture

- Conservative when unsure.
- Best possible plan over fastest possible plan.
- Minimum clarifying question only when blocked.
- Plan quality and trustworthiness matter most.
