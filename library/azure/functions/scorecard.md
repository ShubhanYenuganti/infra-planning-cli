# Scorecard - functions-pp-cli (press v1.3.2, first generation)

| Check | Status | Notes |
|---|---|---|
| OpenAPI conversion | PASS | apis.guru Swagger 2.0 converted to OpenAPI `3.0.0` with `scripts/spec-to-openapi3.sh` |
| pressfile-schema | PASS | All required fields present, including `spec_pipeline: ["swagger2openapi"]` |
| subcommands-present | PASS | sync/search/sql stubs present |
| required-flags | PASS | All required flags present in internal/cli/root.go |
| sqlite-path | PASS | Patched to use `~/.infra-press/functions-pp-cli.db` |
| agents-md-non-trivial | PASS | AGENTS.md written and documents the 2.0 to 3.0 conversion step |
| binary builds | PASS | `go build ./...` clean |
| --help works | PASS | exit 0, generated command list available |
| golden: help-exits-zero | PASS | Binary exits 0 on --help |
| golden: unknown-command-exits-2 | PASS | Unknown command exits with code 2 via usageErr |
| golden: version-is-semver | PASS | `--version` outputs `v0.1.0` (SemVer) |

## Generated subcommands (from --help)

Press generated 48 top-level commands including `api`, `auth`, `doctor`,
`export`, `import`, `functions`, `slots`, `config`, `deployments`,
`sourcecontrols`, `webjobs`, and `workflow`.

Convention-required commands present as stubs: `sync`, `search`, `sql`.

## Spec conversion

Input spec:
`https://api.apis.guru/v2/specs/azure.com/web-WebApps/2018-11-01/swagger.json`

Converted spec:
`/tmp/azure-functions-openapi3.json`

Result: `openapi=3.0.0`, `paths=252`, no top-level Swagger 2.0 marker.

## Press warnings

Press completed generation and built the temporary readycheck binary. Non-blocking
warnings recorded during generation:

- 95 complex body fields skipped, mostly `properties`
- 1 global query parameter filtered: `api-version`

These are Azure ARM/WebApps schema-shape limitations in press v1.3.2. Raw API
coverage remains available through the `api` command.

## Patches applied

| Patch | File | Reason |
|---|---|---|
| Patch 1 | internal/cli/sync.go | Replaced generated experimental sync with v1.1 placeholder stub |
| Patch 1 | internal/cli/search.go | Replaced generated experimental search with v1.1 placeholder stub |
| Patch 1 | internal/cli/sql.go | Added missing sql placeholder stub |
| Patch 1 | internal/cli/root.go | Register sync/search/sql after version command |
| Patch 2 | internal/mcp/tools.go | Press emits `~/.local/share/.../data.db`; convention requires `~/.infra-press/<cli>.db` |
| Patch 2 | internal/cli/helpers.go | CLI defaultDBPath adjusted to `~/.infra-press/<cli>.db` |
| Patch 3 | internal/cli/root.go | Version set to `v0.1.0` (SemVer); version template stripped to bare `{{ .Version }}` |
| Patch 4 | internal/cli/root.go | Unknown-command error wrapped in `usageErr` for exit code 2 |
| Patch 5 | AGENTS.md | Press does not generate AGENTS.md; written per-CLI convention |

## Gaps and actions

| Gap | Type | Action |
|---|---|---|
| Azure ARM complex request bodies simplified | Press/Azure limitation | Documented; use `api` for full raw endpoint coverage |
| `sync`/`search`/`sql` backend not implemented | Universal v1 gap | Stub patches; full SQLite backend in v1.1 |
| SQLite path wrong prefix | Universal press gap | Patched MCP and CLI helper paths to use `~/.infra-press/` |
| AGENTS.md missing | Universal press gap | Written as AGENTS.md patch artifact |

## Decisions

- The Azure WebApps spec is the source for Azure Functions operations in this
  sprint because apis.guru publishes it as Swagger 2.0.
- `sync`/`search`/`sql` intentionally return "not yet implemented" for v1.0,
  matching the other v1 CLIs.
- SQLite paths use `~/.infra-press/functions-pp-cli.db` for both CLI and MCP
  helper paths.
