# Compound recipes

Cross-CLI workflows that combine multiple infra-press CLIs.

## v1.0 recipes

| Recipe | Purpose | Status |
|---|---|---|
| `doctor-all.sh` | Run `doctor --json` against all 6 v1 CLIs and emit a unified JSON report | smoke test |

`doctor-all.sh` is a smoke-test recipe. It demonstrates that every v1 binary is
installed and can run its doctor command. It does not require live cloud
credentials to be useful: missing auth appears inside each CLI's doctor report.

## Adding a recipe

- Keep recipes shell-only unless a real parser is needed.
- Emit JSON or JSONL so agents can consume the result.
- Prefer read-only checks; mutations must support `--dry-run`.
- Document required cloud credentials or mocked fixtures.
