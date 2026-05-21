# Scorecard — cloud-functions-pp-cli (press v1.3.2, first generation)

| Check | Status | Notes |
|---|---|---|
| pressfile-schema | PASS | All required fields present |
| subcommands-present | PASS | sync/search/sql stubs added via patches |
| required-flags | PASS | All required flags present in internal/cli/root.go |
| sqlite-path | PASS | Patched to use `~/.infra-press/cloud-functions-pp-cli.db` |
| agents-md-non-trivial | PASS | AGENTS.md written (≥8 lines) |
| binary builds | PASS | `go build ./...` clean after removing 2 unused imports |
| --help works | PASS | exit 0, all required flags present |
| golden: help-exits-zero | PASS | Binary exits 0 on --help |
| golden: unknown-command-exits-2 | PASS | Unknown command exits with code 2 via usageErr |
| golden: version-is-semver | PASS | `--version` outputs `v0.1.0` (SemVer) |

## Generated subcommands (from --help)

Press generated: `api`, `auth`, `completion`, `doctor`, `export`, `help`, `import`,
`resource-get-iam-policy`, `version`.

Also generated: `name-generate-download-url`, `resource-set-iam-policy`,
`resource-test-iam-permissions`.

Missing (convention-required, added as stubs): `sync`, `search`, `sql`.

## Skipped paths (press warnings)

| Path | Warning |
|---|---|
| `/v2/{name}` | could not derive resource name |
| `/v2/{name}/locations` | could not derive resource name |
| `/v2/{name}/operations` | could not derive resource name |
| `/v2/{parent}/functions` | could not derive resource name |
| `/v2/{parent}/functions:generateUploadUrl` | could not derive resource name |
| `/v2/{parent}/runtimes` | could not derive resource name |

All 6 skips are GCP hierarchical paths with `{name}` or `{parent}` parameters that press cannot resolve to a resource name. This is a press/GCP API mismatch — not a bug in this CLI. Out of scope for v1 patch.

Additional press warnings (non-blocking):
- `warning: skipping body field "policy": complex type not supported as CLI flag`
- 10 global query params filtered (access_token, alt, callback, fields, key, oauth_token, prettyPrint, quotaUser, uploadType, upload_protocol) — all present on 4/4 endpoints, correctly elided

## Patches applied

| Patch | File | Reason |
|---|---|---|
| Patch 6 | internal/cli/helpers.go | Removed unused `path/filepath`/`time` imports (press v1.3.2 bug) |
| Patch 1 | internal/cli/sync.go | Press does not generate sync; stub satisfies convention check |
| Patch 1 | internal/cli/search.go | Press does not generate search; stub satisfies convention check |
| Patch 1 | internal/cli/sql.go | Press does not generate sql; stub satisfies convention check |
| Patch 1 | internal/cli/root.go | Register sync/search/sql in root AddCommand |
| Patch 2 | internal/mcp/tools.go | Press emits `~/.local/share/…/data.db`; convention requires `~/.infra-press/<cli>.db` |
| Patch 3 | internal/cli/root.go | Version set to `v0.1.0` (SemVer); template stripped to bare `{{ .Version }}` |
| Patch 4 | internal/cli/root.go | Unknown-command error wrapped in `usageErr` for exit code 2 |
| Patch 5 | AGENTS.md | Press does not generate AGENTS.md; written per-CLI convention |

## Gaps and actions

| Gap | Type | Action |
|---|---|---|
| Unused imports (`path/filepath`, `time`) in helpers.go | Universal press v1.3.2 bug | Fixed inline (deleted lines) |
| `sync`/`search`/`sql` subcommands missing | Universal press gap | Stub patches; full SQLite backend in v1.1 |
| SQLite path wrong prefix | Universal press gap | Patch `internal/mcp/tools.go` to use `~/.infra-press/` |
| AGENTS.md missing | Universal press gap | Written as AGENTS.md patch artifact |
| 6 GCP hierarchical paths skipped by press | Press/GCP mismatch | Documented; out of scope for v1 patch |
| body field "policy" skipped | Press limitation | Complex nested type; API passthrough via `api` subcommand |

## Decisions

- Unused imports: fixed by deleting lines, not adding a separate patch file.
- sync/search/sql: stub implementations that return "not yet implemented" satisfy the
  convention check and preserve the subcommand contract for agents; full SQLite-backed
  implementation deferred to v1.1.
- SQLite path: used string concatenation (`home + "/.infra-press/..."`) per lint requirement
  that the literal substring `/.infra-press/` appears in source text.
