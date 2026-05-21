# infra-press v1 — Verification Harness Design

**Date:** 2026-05-21
**Status:** Design — ready for implementation plan
**Driving question:** How do we verify that every PR against an infra-press CLI produces a stable, correct CLI — without depending on upstream third-party state for the gating signal?

---

## 1. Goal and non-goals

### Goal

Every PR is verified by a layered harness that gives a single GREEN/RED signal in CI and an inline markdown report as a PR comment. The signal is stable: a PR can only go RED because of something inside the repo or because a live cloud rejected a real call against real credentials. Upstream weather (apis.guru flakiness, transient cloud unavailability, spec drift) never red-lights a PR.

The same harness is runnable locally via a single `make verify` command and is wrapped by a `.claude/commands/verify-pr.md` slash command for Claude-driven PR triage.

### Non-goals (v1.0 of the harness)

- **Press regen-diff drift detection.** Deferred. Regen is dep-driven (upstream spec changes or `press_version` bumps), not PR-driven, and belongs in a separate workflow.
- **Live cloud write operations.** L8 only reads. Create / describe / delete roundtrips are deferred to a later tier.
- **Mocked-credential fixtures.** No httptest servers or VCR cassettes. The plan already defers this to v1.1.
- **Performance / load testing.** Out of scope.

---

## 2. Layered scope

Eight layers. Layers 1–7 always run on PR. Layer 8 runs per-cloud only when GitHub OIDC credentials for that cloud are present in the workflow environment; otherwise it auto-skips.

| # | Layer | Network | Hard fail? | What it verifies |
|---|---|---|---|---|
| L1 | Lint | no | yes | 5 convention checks per CLI: sync/search/sql subcommands present in source, SQLite path uses `~/.infra-press/`, version is SemVer, AGENTS.md ≥ 8 non-trivial lines, unknown-cmd routes through `usageErr`. |
| L2 | Build matrix | no | yes | `go build ./...` clean in every `library/<cloud>/<svc>/`; per-CLI binary builds via `go build -o cmd/<cli>/<cli> ./cmd/<cli>/`. |
| L3 | Golden | no | yes | Existing checks: help-exits-zero, unknown-command-exits-2, version-is-semver, auto-json-when-piped (SKIP allowed). |
| L4 | Behavioral smoke | no | yes | Per CLI: `--version` exits 0 with semver; `--help` exits 0 and mentions required subcommands (api, auth, doctor, export, import, sync, search, sql, version); `unknown-cmd` exits 2; `doctor --json` produces parseable JSON (any exit code); `--agent` flag honored. |
| L5 | Recipe smoke | no | yes | `scripts/recipes/doctor-all.sh` runs with all 6 binaries on PATH and produces parseable JSON with 6 entries. |
| L6 | Spec sync | no | yes | Each vendored spec at `tests/fixtures/specs/<cli>.spec.json` parses as OpenAPI 3.0 and its SHA-256 matches the `spec_sha256` recorded in the CLI's `pressfile.yaml`. |
| L7 | Upstream advisory | yes (5s timeout) | no — advisory only | Best-effort fetch each CLI's `spec_url`. Three checks: (a) liveness/parseability of the response, (b) byte-level / Last-Modified delta against vendored, (c) operation-set diff — parse both specs and report any `paths × operationId` removed-from-live, args-changed, or new-in-live. Findings are **ranked by relevance**: ops present in both upstream and vendored come first (these back real CLI subcommands), with upstream-only "new endpoints available" collected at the bottom under a subheader. Emits drift / dead-link / removed-endpoint / new-endpoint advisories in the PR comment but never red-lights the PR. |
| L8 | Live cloud | yes | yes per cloud where OIDC is configured; auto-SKIP otherwise | Per CLI: `doctor` reports cloud reachability; if `live_smoke.list_subcommand` is declared in catalog, that subcommand returns parseable JSON. |

**Target PR-time runtime budget:** ~90 seconds for L1–L7. L8 adds another ~30–60 seconds per configured cloud.

---

## 3. Execution model

A single source of truth — `catalog.yaml` — drives every layer. Every layer iterates the catalog entries and produces one subtest per CLI. Adding a 7th CLI requires a catalog entry, a vendored spec file, and (for L8 coverage) one `live_smoke.list_subcommand` field. No new Go test code is required.

### A PR run, end to end

```
1. PR opened / pushed.
2. .github/workflows/verify.yml fires.
3. Checkout PR branch.
4. Setup Go.
5. Configure cloud creds via OIDC — three independent optional steps,
   each gated by presence of a repo secret (AWS_OIDC_ROLE_ARN,
   GCP_WORKLOAD_IDENTITY_PROVIDER, AZURE_CLIENT_ID). Missing secret =
   step skipped = env vars not set = L8 auto-skips that cloud.
6. Run `make verify`. Each layer is a separate `go test` invocation:
     go test ./tools/lint-conventions/... -catalog=catalog.yaml   # L1
     go test ./tests/smoke/...    -run BuildMatrix                # L2
     go test ./tests/golden/...                                   # L3
     go test ./tests/smoke/...    -run CLISmoke                   # L4
     go test ./tests/smoke/...    -run RecipeSmoke                # L5
     go test ./tests/spec_check/... -run Sync                     # L6
     go test ./tests/spec_check/... -run Upstream                 # L7
     go test -tags=live ./tests/live/...                          # L8
7. Each test reads catalog.yaml and emits one subtest per CLI. JSON
   report is written to tests/verify-report.json.
8. The slash command (or a small CI step) renders the report into a
   markdown comment and posts it to the PR with a `<!-- verify-pr -->`
   marker so subsequent runs update in place rather than appending.
```

### How L7 and L8 stay uniform across heterogeneous CLIs

Both layers have CLI-specific *inputs* but CLI-agnostic *logic*. The catalog supplies the inputs:

- **L7 input:** each CLI's `spec_url` (already in catalog). L7's Go code iterates catalog, fetches each URL in parallel, parses both upstream and vendored specs, and computes the operation-set diff. Findings are sorted so that operations present in both specs (the ones backing real CLI subcommands) appear first; upstream-only operations are collected under a "New endpoints available upstream" subheader at the end. This keeps the most actionable rows at the top of the PR comment.
- **L8 input:** each CLI's `cloud` (existing field) and `live_smoke.list_subcommand` (new field). L8's Go code is one generic function that for each catalog entry:
  - Calls `creds.HasCreds(cli.Cloud)` — if false, `t.Skip()`.
  - Always: runs the CLI's `doctor` against the real cloud, asserts OK.
  - If `list_subcommand` declared: runs `<cli> <list_subcommand> --json` and asserts parseable JSON.

If a CLI later wants a deeper per-CLI test (e.g., a create/describe/delete roundtrip for Lambda only), it can drop a `tests/live/lambda_test.go` alongside the generic one. The baseline stays catalog-driven.

---

## 4. Harness shape and file layout

```
tests/
├── golden/                   (existing — L3)
├── fixtures/
│   └── specs/                (NEW — vendored OpenAPI 3.0 specs)
│       ├── cloud-run-admin-pp-cli.spec.json
│       ├── cloud-functions-pp-cli.spec.json
│       ├── lambda-pp-cli.spec.json
│       ├── apprunner-pp-cli.spec.json
│       ├── functions-pp-cli.spec.json
│       └── container-apps-pp-cli.spec.json
├── smoke/                    (NEW — L2, L4, L5 share one Go package)
│   ├── build_test.go
│   ├── cli_smoke_test.go
│   ├── recipe_test.go
│   └── helpers.go            (shared catalog loader, binary invocation)
├── spec_check/               (NEW — L6, L7)
│   ├── sync_test.go          (hard, repo-local)
│   └── upstream_test.go      (soft, network with timeout)
└── live/                     (NEW — L8, build-tag gated `//go:build live`)
    ├── creds.go              (HasCreds(cloud) — env var detection)
    ├── live_test.go          (generic catalog-driven baseline)
    └── (optional per-CLI files for deeper tests)

tools/lint-conventions/       (existing — L1)
scripts/recipes/              (existing — L5 source)

Makefile                      (NEW)
.github/workflows/
  └── verify.yml              (NEW)
.claude/commands/
  └── verify-pr.md            (NEW)
```

### Makefile targets

```
verify         → runs L1–L8 in sequence. L8 auto-skips clouds without OIDC.
verify-fast    → runs L1–L7 only. Always skips L8 entirely.
verify-live    → runs L8 only. Requires `-tags=live` and at least one cloud's OIDC env vars.
verify-<layer> → individual layer targets for local triage (verify-lint, verify-build, etc.).
```

### Build-tag gating for L8

Files under `tests/live/` start with:

```go
//go:build live
```

So `go test ./...` from the repo root does **not** compile `tests/live/`. Only `make verify-live` (which adds `-tags=live`) does. This gives a strong guarantee that running the default test suite never hits cloud APIs.

### Slash command flow

`.claude/commands/verify-pr.md` documents a flow Claude executes when given a PR number:

1. `gh pr checkout <num>`
2. `make verify`
3. Parse `tests/verify-report.json`
4. Render markdown comment
5. `gh pr comment <num> --body-file <rendered>` — using a stable `<!-- verify-pr -->` marker so re-runs update the existing comment rather than spam.

---

## 5. Data flow and contracts

### catalog.yaml — adds one block per CLI

```yaml
- name: lambda-pp-cli
  cloud: aws
  # ... existing fields unchanged ...
  live_smoke:
    list_subcommand: list-functions   # required for L8 list-endpoint check
    args: []                          # optional extra args
```

If `live_smoke` is omitted, L8 still runs the generic doctor-reachability check for that CLI; only the list-endpoint assertion is skipped. CLIs can opt into deeper L8 coverage incrementally.

### pressfile.yaml — adds two fields per CLI

```yaml
spec_sha256: 7a3b...                                         # hex SHA-256 of vendored spec
vendored_spec: tests/fixtures/specs/lambda-pp-cli.spec.json  # repo-relative path
```

These replace the existing `spec_etag` placeholder. SHA-256 is chosen over ETag because it is reproducible from a local file, whereas ETag is upstream-controlled and not verifiable offline.

### tests/fixtures/specs/ — one vendored OpenAPI 3.0 file per CLI

For Azure CLIs, the vendored file is the **converted** (post-`swagger2openapi`) OpenAPI 3.0 output, not the upstream 2.0 swagger. This is the actual artifact the CLI was built from. Total size across all 6 is ~5–15 MB.

### JSON report shape (tests/verify-report.json)

```json
{
  "commit": "abc1234",
  "result": "PASS|FAIL|PARTIAL",
  "layers": [
    {
      "id": "L1",
      "name": "Lint",
      "status": "PASS|FAIL|INFO|SKIP",
      "subtests": [
        {"cli": "lambda-pp-cli", "status": "PASS", "detail": "5/5"},
        ...
      ]
    },
    {
      "id": "L7",
      "name": "Upstream advisory",
      "status": "INFO",
      "advisory": [
        {"cli": "lambda-pp-cli", "priority": "high", "message": "operation ListFunctions: arg `MaxItems` changed from optional to required"},
        {"cli": "lambda-pp-cli", "priority": "high", "message": "upstream modified 2026-05-19 (vendor: 2026-05-15)"},
        {"cli": "cloud-functions-pp-cli", "priority": "info", "message": "5 new endpoints available upstream (not yet exposed)"}
      ]
    },
    {
      "id": "L8",
      "name": "Live cloud",
      "status": "PARTIAL",
      "subtests": [
        {"cli": "lambda-pp-cli", "status": "PASS"},
        {"cli": "cloud-functions-pp-cli", "status": "SKIP", "reason": "no GCP OIDC creds in env"}
      ]
    }
  ]
}
```

### Rendered PR comment

```
## Verification report — commit abc1234

| Layer            | Status      | Detail                                                |
|------------------|-------------|-------------------------------------------------------|
| L1 Lint          | PASS 30/30  |                                                       |
| L2 Build         | PASS 6/6    |                                                       |
| L3 Golden        | PASS 24/24  |                                                       |
| L4 Smoke         | PASS 30/30  |                                                       |
| L5 Recipe        | PASS        | doctor-all.sh: 6 CLIs reported                        |
| L6 Spec sync     | PASS 6/6    |                                                       |
| L7 Upstream      | INFO        | 2 CLIs show upstream drift (cloud-functions, lambda)  |
| L8 Live cloud    | SKIP        | OIDC not configured for AWS/GCP/Azure                 |

Result: PASS (L1–L7 green; L8 not configured)
```

---

## 6. Failure modes and invariants

| Scenario | Result |
|---|---|
| L1–L6 assertion fails | Test fails, `make verify` exits non-zero, PR red, report shows failing subtest. |
| L7: apis.guru returns 200 with different bytes | Test passes; advisory logged (`upstream modified <date>`). PR green. |
| L7: apis.guru times out or 5xx | Test passes; INFO logged (`couldn't reach upstream`). PR green. |
| L7: apis.guru returns 404 | Test passes; WARN surfaced prominently in comment (`spec URL dead — consider re-pinning`). PR green. |
| L8: cloud OIDC env vars present, real call fails | Test fails, PR red. This is the breakage L8 exists to catch. |
| L8: cloud OIDC env vars absent | `t.Skip()`, report shows SKIPPED with reason. PR green. |
| L8: cloud OIDC present but creds rejected by cloud | Test fails with the cloud's error envelope. PR red. |
| Go toolchain failure | `make verify` halts before any layer runs. CI shows toolchain error, no report posted. Layers that did not run show as "not run" in any partial report. |

### Core invariant

> A PR can only go RED because of something inside the repo or because a live cloud rejected a real call against real credentials. Upstream weather never causes a red.

---

## 7. Implementation order

Six commits, each independently mergeable. After step 3 the PR pipeline is already functional; later steps add coverage.

1. **L6 foundation: vendor specs + pressfile sha256.**
   - Create `tests/fixtures/specs/` with the 6 spec files. For Azure CLIs, vendor the post-`swagger2openapi` converted output, not the upstream 2.0 swagger.
   - Add `spec_sha256` and `vendored_spec` to each `pressfile.yaml`. Remove the existing literal placeholder strings (current pressfiles contain unexpanded shell commands like `spec_etag: "<curl -sI ... | grep etag | ...>"` — these are dead text that must be replaced).
   - Implement `tests/spec_check/sync_test.go`.
   - Add `verify-spec-sync` Makefile target.
2. **L2, L4, L5: smoke package.**
   - Create `tests/smoke/` with build, CLI smoke, and recipe tests.
   - Add `verify-build`, `verify-smoke`, `verify-recipe` Makefile targets.
   - Wire `verify-fast` to run L1 → L2 → L3 → L4 → L5 → L6.
3. **verify.yml CI workflow.**
   - Workflow runs `make verify-fast` (since L7/L8 not built yet).
   - Output goes to GitHub Actions job summary as markdown. No PR comment yet.
4. **L7: upstream advisory.**
   - Implement `tests/spec_check/upstream_test.go` with 5s timeout, soft-fail.
   - Update `verify` Makefile target to include L7.
   - Update workflow to run `make verify` (now includes L7).
5. **L8: live cloud + OIDC.**
   - Implement `tests/live/` with build-tag gating, `HasCreds` detection, generic catalog-driven baseline test.
   - Add `live_smoke.list_subcommand` to catalog entries that opt in.
   - Workflow gains three optional OIDC configure-creds steps gated by repo secrets.
6. **Slash command.**
   - `.claude/commands/verify-pr.md` wraps `gh pr checkout` + `make verify` + JSON-report parsing + `gh pr comment`.
   - Uses `<!-- verify-pr -->` marker for in-place comment updates.

Each commit lands behind a green CI signal. Stopping after step 3 still leaves a working PR verification pipeline; step 6 is the developer/Claude UX layer.

---

## 8. Deferred / open questions

- **Regen-diff drift detection.** Will be designed separately when needed. Out of scope here.
- **Mocked-credential fixtures for offline live-API testing.** Deferred (existing plan defers to v1.1).
- **Cloud write operations.** L8 only reads. Create/describe/delete roundtrips are a later tier.
- **Per-CLI deeper L8 tests.** The design leaves room (drop `tests/live/<cli>_test.go` alongside the generic file), but none are required for v1.0.
- **Harness self-test.** The harness has no tests of its own. A small `tests/smoke/helpers_test.go` covering the catalog loader and binary-invocation helpers is worth adding but deferred to v1.1.
- **GitHub repo rename (`infra-planning-cli` → `infra-press`).** Tracked separately; not blocking this harness.

---

## Appendix — Why these choices

- **`spec_sha256` over `spec_etag`.** SHA is reproducible from a local file; ETag is upstream-controlled. L6 needs offline verifiability.
- **One Go package for L2 + L4 + L5.** All three layers shell out to binaries; they share a catalog loader and an `os/exec` helper.
- **Build-tag gating for L8.** Stronger than env-var gating: default `go test ./...` cannot accidentally compile or run live tests.
- **Catalog as single source of truth.** No per-CLI test code means linear scaling — 6 CLIs today, 20 in v2.0, same code path.
- **Soft-fail / advisory for L7.** Without this, any apis.guru hiccup would red-light unrelated PRs, exactly the failure mode the design exists to prevent.
- **Per-cloud OIDC independence.** Lets you wire up clouds one at a time without an all-or-nothing rollout.
