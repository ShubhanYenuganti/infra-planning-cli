# Scorecard - container-apps-pp-cli (press v1.3.2, first generation)

| Check | Status | Notes |
|---|---|---|
| Spec URL check | PASS | Original plan path 404; corrected to `Microsoft.App/ContainerApps/stable/2024-03-01/ContainerApps.json` |
| OpenAPI conversion | PASS | Azure REST API Specs Swagger 2.0 converted to OpenAPI `3.0.0` with external refs resolved |
| pressfile-schema | PASS | All required fields present, including `spec_pipeline: ["swagger2openapi"]` |
| subcommands-present | PASS | sync/search/sql stubs added |
| required-flags | PASS | All required flags present in internal/cli/root.go |
| sqlite-path | PASS | Patched to use `~/.infra-press/container-apps-pp-cli.db` |
| agents-md-non-trivial | PASS | AGENTS.md written and documents the 2.0 to 3.0 conversion step |
| binary builds | PASS | `go build ./...` clean after removing 2 unused imports |
| --help works | PASS | exit 0, generated command list available |
| golden: help-exits-zero | PASS | Binary exits 0 on --help |
| golden: unknown-command-exits-2 | PASS | Unknown command exits with code 2 via usageErr |
| golden: version-is-semver | PASS | `--version` outputs `v0.1.0` (SemVer) |

## Generated subcommands (from --help)

Press generated: `api`, `auth`, `completion`, `doctor`, `export`, `help`,
`import`, `providers`, and `version`.

Missing (convention-required, added as stubs): `sync`, `search`, `sql`.

## Spec conversion

Planned input spec:
`https://raw.githubusercontent.com/Azure/azure-rest-api-specs/main/specification/app/resource-manager/Microsoft.App/stable/2024-03-01/ContainerApps.json`

Result: 404. The active stable path includes the service folder:

`https://raw.githubusercontent.com/Azure/azure-rest-api-specs/main/specification/app/resource-manager/Microsoft.App/ContainerApps/stable/2024-03-01/ContainerApps.json`

Converted spec:
`/tmp/azure-container-apps-openapi3.json`

Result: `openapi=3.0.0`, `paths=8`, no top-level Swagger 2.0 marker, no
remaining `./CommonDefinitions.json` external refs.

## Press warnings

Press generated source after `swagger2openapi -r` resolved sibling JSON refs.
Non-blocking warnings recorded during generation:

- 10 complex body fields skipped (`identity`, `properties`, `systemData`, `tags`, `extendedLocation`)
- 1 global query parameter filtered: `api-version`
- Initial press validation failed on unused `path/filepath` and `time` imports in `helpers.go`; fixed by Patch 6

## Patches applied

| Patch | File | Reason |
|---|---|---|
| C2 blocker fix | scripts/spec-to-openapi3.sh | Added `swagger2openapi -r` so Azure sibling JSON refs are resolved before press reads the spec |
| Patch 6 | internal/cli/helpers.go | Removed unused `path/filepath` and `time` imports |
| Patch 1 | internal/cli/sync.go | Press does not generate sync; stub satisfies convention check |
| Patch 1 | internal/cli/search.go | Press does not generate search; stub satisfies convention check |
| Patch 1 | internal/cli/sql.go | Press does not generate sql; stub satisfies convention check |
| Patch 1 | internal/cli/root.go | Register sync/search/sql in root AddCommand |
| Patch 2 | internal/mcp/tools.go | Press emits `~/.local/share/.../data.db`; convention requires `~/.infra-press/<cli>.db` |
| Patch 3 | internal/cli/root.go | Version set to `v0.1.0` (SemVer); template stripped to bare `{{ .Version }}` |
| Patch 4 | internal/cli/root.go | Unknown-command error wrapped in `usageErr` for exit code 2 |
| Patch 5 | AGENTS.md | Press does not generate AGENTS.md; written per-CLI convention |

## Gaps and actions

| Gap | Type | Action |
|---|---|---|
| Plan URL omitted `ContainerApps/` path segment | Spec URL drift | Corrected and documented; catalog pins working URL |
| External sibling JSON refs in Azure spec | Azure spec layout | Resolved with `swagger2openapi -r` in wrapper |
| Azure ARM complex request bodies simplified | Press/Azure limitation | Documented; use `api` for full raw endpoint coverage |
| `sync`/`search`/`sql` backend not implemented | Universal v1 gap | Stub patches; full SQLite backend in v1.1 |
| SQLite path wrong prefix | Universal press gap | Patched MCP path to use `~/.infra-press/` |
| AGENTS.md missing | Universal press gap | Written as AGENTS.md patch artifact |

## Decisions

- The working Azure REST API Specs URL is pinned in `catalog.yaml` and
  `pressfile.yaml`; it differs from the original plan URL because the upstream
  repo moved specs under a `ContainerApps/` service folder.
- `scripts/spec-to-openapi3.sh` resolves external refs by default. This remains
  compatible with the Azure Functions conversion and is required for Container
  Apps.
- `sync`/`search`/`sql` intentionally return "not yet implemented" for v1.0,
  matching the other v1 CLIs.
