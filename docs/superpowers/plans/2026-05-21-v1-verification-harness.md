# v1 Verification Harness Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Ship an 8-layer composed verification harness driven by `catalog.yaml` that runs on every PR via a new `verify.yml` workflow, plus a `/verify-pr` slash command that wraps the local invocation for Claude-driven PR triage.

**Architecture:** Each layer is an independent Go test (or existing tool) invoked by a root `Makefile`. `catalog.yaml` is the single source of truth — every layer iterates it and emits one subtest per CLI. L1–L6 are repo-local hard checks; L7 is network advisory (soft fail); L8 is live cloud, build-tag gated, per-cloud auto-skip when GitHub OIDC creds for that cloud are absent.

**Tech Stack:** Go 1.20+ (existing `tests/golden/` and `tools/lint-conventions/` are go1.20 modules), `gopkg.in/yaml.v3`, `github.com/getkin/kin-openapi/openapi3` (L7 only), bash, GitHub Actions, `gh` CLI.

**Spec reference:** [docs/superpowers/specs/2026-05-21-v1-verification-design.md](../specs/2026-05-21-v1-verification-design.md). When in doubt about a design decision, the spec is authoritative.

---

## File Structure Overview

**New directories / files:**

```
Makefile                                          (NEW — root, drives every layer)
.github/workflows/verify.yml                      (NEW — Phase 3 / Phase 5 updates)
.claude/commands/verify-pr.md                     (NEW — Phase 6)
tests/fixtures/specs/<cli>.spec.json    × 6       (NEW — Phase 1)
tests/smoke/                                       (NEW Go module — Phase 2)
  go.mod
  helpers.go
  build_test.go
  cli_smoke_test.go
  recipe_test.go
tests/spec_check/                                  (NEW Go module — Phase 1 + Phase 4)
  go.mod
  helpers.go
  sync_test.go
  upstream_test.go
tests/live/                                        (NEW Go module — Phase 5)
  go.mod
  helpers.go
  creds.go
  live_test.go
```

**Modified files:**

```
library/<cloud>/<svc>/pressfile.yaml   × 6       (Phase 1 — replace spec_etag with spec_sha256+vendored_spec)
catalog.yaml                                      (Phase 5 — add live_smoke per CLI)
.gitignore                                        (Phase 6 — un-ignore .claude/commands/)
```

**Untouched (referenced by harness):**

```
catalog.yaml                                      (read by every layer)
tools/lint-conventions/                            (L1, already exists)
tests/golden/                                      (L3, already exists)
scripts/recipes/doctor-all.sh                      (L5 invokes this)
```

**Why these boundaries:**

- Each test module is independently `go test`-able from CI or Make.
- `tests/smoke/`, `tests/spec_check/`, `tests/live/` are separate Go modules to match the existing `tests/golden/` / `tools/lint-conventions/` pattern (each test/tool is its own module with its own go.mod).
- Each module duplicates the small (~30-line) catalog loader in its own `helpers.go` rather than introducing a shared module — DRY violation is intentional, keeps inter-module dependencies at zero.

---

# Phase 1 — L6 Foundation (Vendor Specs + Pressfile SHA + Sync Test)

After this phase: L6 spec-sync layer works end-to-end. PRs that hand-edit a vendored spec without updating its pressfile SHA fail. `make verify-spec-sync` is a working Make target.

## Task 1.1: Vendor the 6 OpenAPI 3.0 spec files

**Files:**
- Create: `tests/fixtures/specs/cloud-run-admin-pp-cli.spec.json`
- Create: `tests/fixtures/specs/cloud-functions-pp-cli.spec.json`
- Create: `tests/fixtures/specs/lambda-pp-cli.spec.json`
- Create: `tests/fixtures/specs/apprunner-pp-cli.spec.json`
- Create: `tests/fixtures/specs/functions-pp-cli.spec.json`
- Create: `tests/fixtures/specs/container-apps-pp-cli.spec.json`

- [ ] **Step 1: Create the fixtures directory**

Run from repo root:

```bash
mkdir -p tests/fixtures/specs
```

Expected: directory exists, no output.

- [ ] **Step 2: Vendor the 4 direct apis.guru specs (GCP × 2, AWS × 2)**

Run from repo root:

```bash
curl -fsSL https://api.apis.guru/v2/specs/googleapis.com/run/v2/openapi.json \
  -o tests/fixtures/specs/cloud-run-admin-pp-cli.spec.json
curl -fsSL https://api.apis.guru/v2/specs/googleapis.com/cloudfunctions/v2/openapi.json \
  -o tests/fixtures/specs/cloud-functions-pp-cli.spec.json
curl -fsSL https://api.apis.guru/v2/specs/amazonaws.com/lambda/2015-03-31/openapi.json \
  -o tests/fixtures/specs/lambda-pp-cli.spec.json
curl -fsSL https://api.apis.guru/v2/specs/amazonaws.com/apprunner/2020-05-15/openapi.json \
  -o tests/fixtures/specs/apprunner-pp-cli.spec.json
```

Expected: 4 files created. Each file's first line should contain `"openapi": "3.0`. Verify:

```bash
for f in tests/fixtures/specs/*.spec.json; do
  echo "=== $f ==="
  head -c 200 "$f"
  echo
done
```

Expected output includes `"openapi": "3.0"` somewhere in each file's first 200 bytes (key order may vary).

- [ ] **Step 3: Vendor the 2 Azure specs (post-conversion)**

Azure specs ship as OpenAPI 2.0 (Swagger). We must convert them with `scripts/spec-to-openapi3.sh` (already present) before vendoring. The vendored file is the post-conversion OpenAPI 3.0 output.

```bash
./scripts/spec-to-openapi3.sh \
  https://api.apis.guru/v2/specs/azure.com/web-WebApps/2018-11-01/swagger.json \
  tests/fixtures/specs/functions-pp-cli.spec.json

./scripts/spec-to-openapi3.sh \
  https://raw.githubusercontent.com/Azure/azure-rest-api-specs/main/specification/app/resource-manager/Microsoft.App/ContainerApps/stable/2024-03-01/ContainerApps.json \
  tests/fixtures/specs/container-apps-pp-cli.spec.json
```

Expected: 2 files created. Verify both start with `{"openapi":"3.0`:

```bash
head -c 50 tests/fixtures/specs/functions-pp-cli.spec.json
head -c 50 tests/fixtures/specs/container-apps-pp-cli.spec.json
```

Expected output: `{"openapi":"3.0...` in both.

- [ ] **Step 4: Confirm all 6 are present and sized reasonably**

```bash
ls -la tests/fixtures/specs/
```

Expected: 6 `.spec.json` files. Each between ~50KB and ~5MB. None empty.

## Task 1.2: Update each pressfile.yaml with spec_sha256 + vendored_spec

**Files:**
- Modify: `library/gcp/cloud-run-admin/pressfile.yaml`
- Modify: `library/gcp/cloud-functions/pressfile.yaml`
- Modify: `library/aws/lambda/pressfile.yaml`
- Modify: `library/aws/apprunner/pressfile.yaml`
- Modify: `library/azure/functions/pressfile.yaml`
- Modify: `library/azure/container-apps/pressfile.yaml`

For each pressfile: replace the `spec_etag: ...` line with two new lines: `spec_sha256` (computed) and `vendored_spec` (repo-relative path to the vendored fixture).

- [ ] **Step 1: Compute SHA-256 of each vendored spec**

Run from repo root:

```bash
for f in tests/fixtures/specs/*.spec.json; do
  cli=$(basename "$f" .spec.json)
  sha=$(shasum -a 256 "$f" | awk '{print $1}')
  echo "$cli  $sha"
done
```

Expected: 6 lines, each with a 64-character hex SHA-256.

Record the output — you'll paste each SHA into the matching pressfile in Step 2.

- [ ] **Step 2: Edit each pressfile.yaml — example with cloud-run-admin**

Find the existing line (around line 8 of `library/gcp/cloud-run-admin/pressfile.yaml`):

```yaml
spec_etag: "e3c49030d6657c77f91ff3f4fbec5c97"
```

Replace it with two lines (substitute the actual SHA from Step 1):

```yaml
spec_sha256: <64-char hex SHA from step 1 for cloud-run-admin-pp-cli>
vendored_spec: tests/fixtures/specs/cloud-run-admin-pp-cli.spec.json
```

Repeat for all 6 pressfiles, substituting:

| CLI | Pressfile path | Vendored spec path |
|---|---|---|
| cloud-run-admin | `library/gcp/cloud-run-admin/pressfile.yaml` | `tests/fixtures/specs/cloud-run-admin-pp-cli.spec.json` |
| cloud-functions | `library/gcp/cloud-functions/pressfile.yaml` | `tests/fixtures/specs/cloud-functions-pp-cli.spec.json` |
| lambda | `library/aws/lambda/pressfile.yaml` | `tests/fixtures/specs/lambda-pp-cli.spec.json` |
| apprunner | `library/aws/apprunner/pressfile.yaml` | `tests/fixtures/specs/apprunner-pp-cli.spec.json` |
| functions | `library/azure/functions/pressfile.yaml` | `tests/fixtures/specs/functions-pp-cli.spec.json` |
| container-apps | `library/azure/container-apps/pressfile.yaml` | `tests/fixtures/specs/container-apps-pp-cli.spec.json` |

- [ ] **Step 3: Verify no `spec_etag` lines remain**

```bash
grep -r "spec_etag" library/ || echo "OK — no spec_etag remains"
grep -rl "spec_sha256" library/ | wc -l
grep -rl "vendored_spec" library/ | wc -l
```

Expected: first command prints `OK — no spec_etag remains`. The next two each print `6`.

## Task 1.3: Create `tests/spec_check/` Go module with sync test

**Files:**
- Create: `tests/spec_check/go.mod`
- Create: `tests/spec_check/helpers.go`
- Create: `tests/spec_check/sync_test.go`

- [ ] **Step 1: Initialize the Go module**

```bash
mkdir -p tests/spec_check
cd tests/spec_check
go mod init github.com/ShubhanYenuganti/infra-press/tests/spec_check
go get gopkg.in/yaml.v3
cd ../..
```

Expected: `tests/spec_check/go.mod` exists and declares the module path above with `gopkg.in/yaml.v3` as a dependency.

- [ ] **Step 2: Write `tests/spec_check/helpers.go`**

```go
package spec_check

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type CLI struct {
	Name    string `yaml:"name"`
	Cloud   string `yaml:"cloud"`
	Service string `yaml:"service"`
	Path    string `yaml:"path"`
	Binary  string `yaml:"binary"`
	Status  string `yaml:"status"`
	SpecURL string `yaml:"spec_url"`
}

type Catalog struct {
	Version int   `yaml:"version"`
	CLIs    []CLI `yaml:"clis"`
}

type Pressfile struct {
	SpecSHA256   string `yaml:"spec_sha256"`
	VendoredSpec string `yaml:"vendored_spec"`
	SpecURL      string `yaml:"spec_url"`
}

// RepoRoot walks up from the working directory looking for catalog.yaml.
func RepoRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "catalog.yaml")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

func LoadCatalog(repoRoot string) (*Catalog, error) {
	b, err := os.ReadFile(filepath.Join(repoRoot, "catalog.yaml"))
	if err != nil {
		return nil, err
	}
	var c Catalog
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func LoadPressfile(repoRoot, cliPath string) (*Pressfile, error) {
	b, err := os.ReadFile(filepath.Join(repoRoot, cliPath, "pressfile.yaml"))
	if err != nil {
		return nil, err
	}
	var p Pressfile
	if err := yaml.Unmarshal(b, &p); err != nil {
		return nil, err
	}
	return &p, nil
}
```

- [ ] **Step 3: Write the failing test — `tests/spec_check/sync_test.go`**

```go
package spec_check

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSpecSync(t *testing.T) {
	root, err := RepoRoot()
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	cat, err := LoadCatalog(root)
	if err != nil {
		t.Fatalf("load catalog: %v", err)
	}
	if len(cat.CLIs) == 0 {
		t.Fatal("catalog has no CLIs")
	}
	for _, cli := range cat.CLIs {
		cli := cli
		t.Run(cli.Name, func(t *testing.T) {
			pf, err := LoadPressfile(root, cli.Path)
			if err != nil {
				t.Fatalf("load pressfile: %v", err)
			}
			if pf.SpecSHA256 == "" {
				t.Fatal("pressfile missing spec_sha256")
			}
			if pf.VendoredSpec == "" {
				t.Fatal("pressfile missing vendored_spec")
			}

			specBytes, err := os.ReadFile(filepath.Join(root, pf.VendoredSpec))
			if err != nil {
				t.Fatalf("read vendored spec: %v", err)
			}

			// 1) SHA-256 matches recorded value
			sum := sha256.Sum256(specBytes)
			actualHex := hex.EncodeToString(sum[:])
			if actualHex != pf.SpecSHA256 {
				t.Fatalf("spec SHA mismatch:\n  vendored file: %s\n  pressfile:     %s",
					actualHex, pf.SpecSHA256)
			}

			// 2) Vendored spec parses as JSON
			var doc map[string]interface{}
			if err := json.Unmarshal(specBytes, &doc); err != nil {
				t.Fatalf("vendored spec not parseable JSON: %v", err)
			}

			// 3) Spec advertises OpenAPI 3.x
			v, ok := doc["openapi"].(string)
			if !ok {
				t.Fatal("vendored spec missing 'openapi' field")
			}
			if !strings.HasPrefix(v, "3.") {
				t.Fatalf("vendored spec not OpenAPI 3.x (got openapi=%q)", v)
			}
		})
	}
}
```

- [ ] **Step 4: Run the test and verify it passes**

```bash
cd tests/spec_check
go test -v -run TestSpecSync .
cd ../..
```

Expected: `PASS` with 6 subtests, one per CLI. Output ends with `ok   github.com/ShubhanYenuganti/infra-press/tests/spec_check`.

If a SHA mismatch is reported, the SHA you wrote in the pressfile (Task 1.2) doesn't match the file you vendored. Re-compute and re-edit.

## Task 1.4: Create root Makefile with verify-spec-sync target

**Files:**
- Create: `Makefile`

- [ ] **Step 1: Write the initial Makefile**

```makefile
# infra-press verification harness
# See docs/superpowers/specs/2026-05-21-v1-verification-design.md

.PHONY: verify verify-fast verify-spec-sync verify-lint verify-build verify-golden verify-smoke verify-recipe verify-upstream verify-live

# Default target — full verify (will grow as later phases land).
verify: verify-spec-sync

verify-fast: verify-spec-sync

verify-spec-sync:
	@echo "==> L6 spec sync"
	cd tests/spec_check && go test -v -run TestSpecSync .
```

- [ ] **Step 2: Run `make verify` and confirm L6 passes through it**

```bash
make verify
```

Expected: prints `==> L6 spec sync` then runs `go test`, ends with `PASS` and a successful Make exit (`$?` is 0).

## Task 1.5: Commit Phase 1

- [ ] **Step 1: Stage and commit**

```bash
git add tests/fixtures/specs/ tests/spec_check/ library/*/pressfile.yaml Makefile
git status
```

Expected: 6 new spec fixtures, the new `tests/spec_check/` module, 6 modified pressfiles, new Makefile. Confirm no other files are staged.

```bash
git commit -m "$(cat <<'EOF'
feat(verify): L6 spec sync layer + vendor 6 specs

Vendors each CLI's OpenAPI 3.0 spec under tests/fixtures/specs/ (Azure
specs are vendored post-swagger2openapi conversion). Replaces the
upstream-controlled spec_etag in each pressfile.yaml with a
reproducible spec_sha256 plus a vendored_spec path field. New
tests/spec_check/ module asserts spec_sha256 matches the vendored
file's actual SHA-256, and that the vendored spec parses as OpenAPI
3.x. Root Makefile bootstrapped with verify-spec-sync target.

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

Expected: clean commit. `git log -1 --oneline` shows the new commit.

---

# Phase 2 — L2 Build + L4 Behavioral Smoke + L5 Recipe Smoke

After this phase: `make verify-fast` runs L1 → L2 → L3 → L4 → L5 → L6 with all hard-fail layers green.

## Task 2.1: Create `tests/smoke/` Go module + shared helpers

**Files:**
- Create: `tests/smoke/go.mod`
- Create: `tests/smoke/helpers.go`

- [ ] **Step 1: Initialize the module**

```bash
mkdir -p tests/smoke
cd tests/smoke
go mod init github.com/ShubhanYenuganti/infra-press/tests/smoke
go get gopkg.in/yaml.v3
cd ../..
```

Expected: `tests/smoke/go.mod` exists.

- [ ] **Step 2: Write `tests/smoke/helpers.go`**

```go
package smoke

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type CLI struct {
	Name          string   `yaml:"name"`
	Cloud         string   `yaml:"cloud"`
	Service       string   `yaml:"service"`
	Path          string   `yaml:"path"`
	Binary        string   `yaml:"binary"`
	Status        string   `yaml:"status"`
	DailyCommands []string `yaml:"daily_commands"`
}

type Catalog struct {
	Version int   `yaml:"version"`
	CLIs    []CLI `yaml:"clis"`
}

func RepoRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "catalog.yaml")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

func LoadCatalog(repoRoot string) (*Catalog, error) {
	b, err := os.ReadFile(filepath.Join(repoRoot, "catalog.yaml"))
	if err != nil {
		return nil, err
	}
	var c Catalog
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func BinaryPath(repoRoot string, cli CLI) string {
	return filepath.Join(repoRoot, cli.Path, "cmd", cli.Binary, cli.Binary)
}
```

## Task 2.2: Implement L2 build matrix

**Files:**
- Create: `tests/smoke/build_test.go`

- [ ] **Step 1: Write the test**

```go
package smoke

import (
	"os/exec"
	"path/filepath"
	"testing"
)

func TestBuildMatrix(t *testing.T) {
	root, err := RepoRoot()
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	cat, err := LoadCatalog(root)
	if err != nil {
		t.Fatalf("load catalog: %v", err)
	}
	for _, cli := range cat.CLIs {
		cli := cli
		t.Run(cli.Name, func(t *testing.T) {
			cliDir := filepath.Join(root, cli.Path)

			t.Run("go-build-all", func(t *testing.T) {
				cmd := exec.Command("go", "build", "./...")
				cmd.Dir = cliDir
				if out, err := cmd.CombinedOutput(); err != nil {
					t.Fatalf("go build ./... failed in %s:\n%s", cli.Path, out)
				}
			})

			t.Run("binary-builds", func(t *testing.T) {
				bin := BinaryPath(root, cli)
				cmd := exec.Command("go", "build", "-o", bin, "./cmd/"+cli.Binary+"/")
				cmd.Dir = cliDir
				if out, err := cmd.CombinedOutput(); err != nil {
					t.Fatalf("binary build failed for %s:\n%s", cli.Name, out)
				}
			})
		})
	}
}
```

- [ ] **Step 2: Run the test**

```bash
cd tests/smoke
go test -v -run TestBuildMatrix .
cd ../..
```

Expected: PASS with 6 top-level subtests, each with 2 nested (`go-build-all`, `binary-builds`). All binaries are now present at `library/<cloud>/<svc>/cmd/<cli>/<cli>`.

## Task 2.3: Implement L4 behavioral smoke

**Files:**
- Create: `tests/smoke/cli_smoke_test.go`

- [ ] **Step 1: Write the test**

```go
package smoke

import (
	"encoding/json"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

var semverRE = regexp.MustCompile(`^v\d+\.\d+\.\d+`)

var requiredSubcommands = []string{
	"api", "auth", "doctor", "export", "import",
	"sync", "search", "sql", "version",
}

func TestCLISmoke(t *testing.T) {
	root, err := RepoRoot()
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	cat, err := LoadCatalog(root)
	if err != nil {
		t.Fatalf("load catalog: %v", err)
	}
	for _, cli := range cat.CLIs {
		cli := cli
		t.Run(cli.Name, func(t *testing.T) {
			bin := BinaryPath(root, cli)

			t.Run("version-is-semver", func(t *testing.T) {
				out, err := exec.Command(bin, "--version").Output()
				if err != nil {
					t.Fatalf("--version failed: %v", err)
				}
				if !semverRE.Match(out) {
					t.Fatalf("--version output not SemVer: %q", out)
				}
			})

			t.Run("help-mentions-required-subcommands", func(t *testing.T) {
				out, err := exec.Command(bin, "--help").Output()
				if err != nil {
					t.Fatalf("--help failed: %v", err)
				}
				body := string(out)
				for _, sub := range requiredSubcommands {
					if !strings.Contains(body, sub) {
						t.Errorf("--help missing required subcommand %q", sub)
					}
				}
			})

			t.Run("unknown-command-exits-2", func(t *testing.T) {
				cmd := exec.Command(bin, "asdf-not-a-real-command")
				err := cmd.Run()
				exitErr, ok := err.(*exec.ExitError)
				if !ok {
					t.Fatalf("expected non-zero exit, got %v", err)
				}
				if exitErr.ExitCode() != 2 {
					t.Fatalf("expected exit code 2, got %d", exitErr.ExitCode())
				}
			})

			t.Run("doctor-json-parseable", func(t *testing.T) {
				// Accept any exit code; assert stdout parses as JSON.
				out, _ := exec.Command(bin, "doctor", "--json").Output()
				if len(out) == 0 {
					t.Fatal("doctor --json produced no stdout")
				}
				var v interface{}
				if err := json.Unmarshal(out, &v); err != nil {
					t.Fatalf("doctor --json output not parseable JSON: %v\noutput: %s", err, out)
				}
			})

			t.Run("agent-flag-honored", func(t *testing.T) {
				cmd := exec.Command(bin, "--agent", "doctor")
				err := cmd.Run()
				if exitErr, ok := err.(*exec.ExitError); ok {
					// Exit 2 = unknown flag/command; means --agent was rejected.
					if exitErr.ExitCode() == 2 {
						t.Fatalf("--agent treated as unknown (exit 2)")
					}
				}
			})
		})
	}
}
```

- [ ] **Step 2: Run the test**

```bash
cd tests/smoke
go test -v -run TestCLISmoke .
cd ../..
```

Expected: PASS with 6 top-level subtests × 5 nested = 30 subtests. If `help-mentions-required-subcommands` fails for a CLI, that CLI's `--help` output is missing one of the 9 required subcommands — investigate before moving on (don't relax the list).

## Task 2.4: Implement L5 recipe smoke

**Files:**
- Create: `tests/smoke/recipe_test.go`

- [ ] **Step 1: Write the test**

```go
package smoke

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRecipeSmoke(t *testing.T) {
	root, err := RepoRoot()
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	cat, err := LoadCatalog(root)
	if err != nil {
		t.Fatalf("load catalog: %v", err)
	}

	// Build PATH with every CLI binary directory prepended.
	binDirs := make([]string, 0, len(cat.CLIs))
	for _, cli := range cat.CLIs {
		binDirs = append(binDirs, filepath.Join(root, cli.Path, "cmd", cli.Binary))
	}
	augmentedPath := strings.Join(binDirs, ":") + ":" + os.Getenv("PATH")

	script := filepath.Join(root, "scripts", "recipes", "doctor-all.sh")
	cmd := exec.Command("bash", script)
	env := os.Environ()
	for i, kv := range env {
		if strings.HasPrefix(kv, "PATH=") {
			env[i] = "PATH=" + augmentedPath
			break
		}
	}
	cmd.Env = env

	out, _ := cmd.Output() // Accept any exit code; assert structure of stdout.

	var report struct {
		CLIs []struct {
			CLI    string          `json:"cli"`
			Report json.RawMessage `json:"report"`
		} `json:"clis"`
	}
	if err := json.Unmarshal(out, &report); err != nil {
		t.Fatalf("recipe stdout not parseable JSON: %v\noutput: %s", err, out)
	}
	if got, want := len(report.CLIs), len(cat.CLIs); got != want {
		t.Fatalf("recipe report listed %d CLIs, expected %d", got, want)
	}
}
```

- [ ] **Step 2: Run the test**

```bash
cd tests/smoke
go test -v -run TestRecipeSmoke .
cd ../..
```

Expected: PASS. The test relies on Task 2.2 having built the binaries; if it can't find them, re-run `go test -run TestBuildMatrix` first.

## Task 2.5: Extend Makefile with verify-fast wiring all hard layers

**Files:**
- Modify: `Makefile`

- [ ] **Step 1: Replace the existing Makefile content with the expanded version**

```makefile
# infra-press verification harness
# See docs/superpowers/specs/2026-05-21-v1-verification-design.md

.PHONY: verify verify-fast verify-lint verify-build verify-golden verify-smoke verify-recipe verify-spec-sync verify-upstream verify-live

verify: verify-fast

verify-fast: verify-lint verify-build verify-golden verify-smoke verify-recipe verify-spec-sync
	@echo "==> verify-fast complete (L1-L6)"

verify-lint:
	@echo "==> L1 lint"
	@for path in $$(yq '.clis[].path' catalog.yaml); do \
		echo "  - $$path"; \
		(cd tools/lint-conventions && go run . "../../$$path/") || exit 1; \
	done

verify-build:
	@echo "==> L2 build matrix"
	cd tests/smoke && go test -v -run TestBuildMatrix .

verify-golden:
	@echo "==> L3 golden"
	cd tests/golden && go test -v ./...

verify-smoke:
	@echo "==> L4 behavioral smoke"
	cd tests/smoke && go test -v -run TestCLISmoke .

verify-recipe:
	@echo "==> L5 recipe smoke"
	cd tests/smoke && go test -v -run TestRecipeSmoke .

verify-spec-sync:
	@echo "==> L6 spec sync"
	cd tests/spec_check && go test -v -run TestSpecSync .
```

- [ ] **Step 2: Run `make verify-fast`**

```bash
make verify-fast
```

Expected: each layer prints its banner and completes. Final output ends with `==> verify-fast complete (L1-L6)`. If `yq` is not installed, the `verify-lint` step will fail — install via `brew install yq` (or document this in install.sh; deferred to a later task if it surfaces).

## Task 2.6: Commit Phase 2

- [ ] **Step 1: Stage and commit**

```bash
git add tests/smoke/ Makefile
git commit -m "$(cat <<'EOF'
feat(verify): L2/L4/L5 smoke package + verify-fast Makefile target

New tests/smoke/ Go module with three layers:
- L2 build matrix: go build per CLI dir + per-CLI binary target
- L4 behavioral smoke: --version/--help/unknown-cmd/doctor-json/--agent
- L5 recipe smoke: doctor-all.sh runs with all 6 binaries on PATH

verify-fast Make target now chains L1 lint, L2 build, L3 golden,
L4 smoke, L5 recipe, L6 spec sync as one entry point.

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

Expected: clean commit.

---

# Phase 3 — verify.yml CI Workflow

After this phase: every PR runs `make verify-fast` in GitHub Actions and posts the output to the workflow summary.

## Task 3.1: Create verify.yml workflow

**Files:**
- Create: `.github/workflows/verify.yml`

- [ ] **Step 1: Write the workflow**

```yaml
name: verify

on:
  pull_request:
  push:
    branches: [main]

permissions:
  contents: read
  pull-requests: write
  id-token: write   # required for OIDC in Phase 5; harmless here.

jobs:
  verify:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'

      - name: Install yq
        run: sudo snap install yq

      - name: Run make verify-fast
        id: verify
        run: |
          set -o pipefail
          make verify-fast 2>&1 | tee verify.log

      - name: Append verify.log to job summary
        if: always()
        run: |
          {
            echo "## Verification report"
            echo
            echo '```'
            tail -200 verify.log
            echo '```'
          } >> "$GITHUB_STEP_SUMMARY"
```

- [ ] **Step 2: Smoke-test the workflow locally (best-effort)**

GitHub Actions can't be run from local; instead, confirm syntactic validity:

```bash
# If `act` is installed, dry-run the workflow:
command -v act && act -j verify --dryrun || echo "act not installed; will validate on push"

# Otherwise just confirm the YAML parses:
python3 -c "import yaml; yaml.safe_load(open('.github/workflows/verify.yml'))" && echo OK
```

Expected: `OK` from the Python parse check.

## Task 3.2: Commit Phase 3

- [ ] **Step 1: Commit**

```bash
git add .github/workflows/verify.yml
git commit -m "$(cat <<'EOF'
feat(ci): verify.yml workflow runs make verify-fast on every PR

Runs L1-L6 on every pull_request and push to main. Tail of make
output is appended to the GitHub Actions job summary as markdown.

PR comment posting is added in Phase 6 via the verify-pr slash
command; this workflow keeps a CI-resident always-on signal.

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

Expected: clean commit. On the next PR push, the verify workflow runs.

---

# Phase 4 — L7 Upstream Advisory

After this phase: `make verify` (no longer just verify-fast) includes L7 — a network layer that fetches each spec from apis.guru, computes a ranked operation-set diff against the vendored fixture, and emits advisory rows. **L7 never fails the PR.**

## Task 4.1: Add kin-openapi dependency to spec_check

**Files:**
- Modify: `tests/spec_check/go.mod`
- Modify: `tests/spec_check/go.sum`

- [ ] **Step 1: Add the parser dependency**

```bash
cd tests/spec_check
go get github.com/getkin/kin-openapi/openapi3
cd ../..
```

Expected: `tests/spec_check/go.mod` gains a `require github.com/getkin/kin-openapi v0.x.x` line.

## Task 4.2: Write the upstream advisory test

**Files:**
- Create: `tests/spec_check/upstream_test.go`

- [ ] **Step 1: Write the test**

```go
package spec_check

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
)

type finding struct {
	cli      string
	priority string // "high" if op present in vendored, "info" if upstream-only.
	message  string
}

func TestUpstreamAdvisory(t *testing.T) {
	root, err := RepoRoot()
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	cat, err := LoadCatalog(root)
	if err != nil {
		t.Fatalf("load catalog: %v", err)
	}

	var (
		mu       sync.Mutex
		findings []finding
		wg       sync.WaitGroup
	)
	client := &http.Client{Timeout: 5 * time.Second}

	for _, cli := range cat.CLIs {
		cli := cli
		wg.Add(1)
		go func() {
			defer wg.Done()

			pf, err := LoadPressfile(root, cli.Path)
			if err != nil {
				mu.Lock()
				findings = append(findings, finding{cli.Name, "info",
					"could not load pressfile: " + err.Error()})
				mu.Unlock()
				return
			}
			if pf.SpecURL == "" {
				return
			}

			// Fetch upstream with timeout.
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			req, _ := http.NewRequestWithContext(ctx, "GET", pf.SpecURL, nil)
			resp, err := client.Do(req)
			if err != nil {
				mu.Lock()
				findings = append(findings, finding{cli.Name, "info",
					"upstream unreachable: " + err.Error() + " (advisory only)"})
				mu.Unlock()
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode == 404 {
				mu.Lock()
				findings = append(findings, finding{cli.Name, "high",
					"spec URL returned 404 — consider re-pinning"})
				mu.Unlock()
				return
			}
			if resp.StatusCode != 200 {
				mu.Lock()
				findings = append(findings, finding{cli.Name, "info",
					"upstream HTTP " + resp.Status + " (advisory only)"})
				mu.Unlock()
				return
			}

			// Load both specs through kin-openapi for operation-id awareness.
			loader := openapi3.NewLoader()
			loader.IsExternalRefsAllowed = true

			liveDoc, err := loader.LoadFromIoReader(resp.Body)
			if err != nil {
				mu.Lock()
				findings = append(findings, finding{cli.Name, "info",
					"upstream spec did not parse as OpenAPI 3.x (advisory only)"})
				mu.Unlock()
				return
			}

			vendoredBytes, err := os.ReadFile(filepath.Join(root, pf.VendoredSpec))
			if err != nil {
				mu.Lock()
				findings = append(findings, finding{cli.Name, "info",
					"vendored spec missing: " + err.Error()})
				mu.Unlock()
				return
			}
			vendoredDoc, err := loader.LoadFromData(vendoredBytes)
			if err != nil {
				mu.Lock()
				findings = append(findings, finding{cli.Name, "info",
					"vendored spec did not parse (this is a separate bug)"})
				mu.Unlock()
				return
			}

			// Compute operation-set diff.
			liveOps := operationSet(liveDoc)
			vendoredOps := operationSet(vendoredDoc)

			var removed, added []string
			for op := range vendoredOps {
				if _, ok := liveOps[op]; !ok {
					removed = append(removed, op)
				}
			}
			for op := range liveOps {
				if _, ok := vendoredOps[op]; !ok {
					added = append(added, op)
				}
			}
			sort.Strings(removed)
			sort.Strings(added)

			mu.Lock()
			defer mu.Unlock()
			for _, op := range removed {
				findings = append(findings, finding{cli.Name, "high",
					"operation removed upstream: " + op + " (CLI exposes it; consider regen)"})
			}
			if len(added) > 0 {
				findings = append(findings, finding{cli.Name, "info",
					"new endpoints available upstream (not exposed): " +
						strings.Join(added, ", ")})
			}
		}()
	}
	wg.Wait()

	// Sort: high-priority findings first, then info.
	sort.SliceStable(findings, func(i, j int) bool {
		if findings[i].priority != findings[j].priority {
			return findings[i].priority == "high"
		}
		return findings[i].cli < findings[j].cli
	})

	if len(findings) == 0 {
		t.Log("upstream advisory: no findings (clean)")
		return
	}

	t.Log("upstream advisory findings (ranked):")
	for _, f := range findings {
		t.Logf("  [%s] %s: %s", f.priority, f.cli, f.message)
	}
	// L7 NEVER fails the PR.
}

// operationSet builds a "METHOD path::operationId" set from an OpenAPI 3 doc.
func operationSet(doc *openapi3.T) map[string]struct{} {
	out := map[string]struct{}{}
	if doc == nil || doc.Paths == nil {
		return out
	}
	for path, item := range doc.Paths.Map() {
		for method, op := range item.Operations() {
			id := op.OperationID
			if id == "" {
				id = "<no-op-id>"
			}
			out[method+" "+path+"::"+id] = struct{}{}
		}
	}
	return out
}
```

- [ ] **Step 2: Run the test**

```bash
cd tests/spec_check
go test -v -run TestUpstreamAdvisory .
cd ../..
```

Expected: PASS regardless of network state. Test logs may include findings; that's OK and informational. Verify with no network:

```bash
# Verify graceful behavior on no network (best-effort — depends on OS):
GODEBUG=netdns=go HTTP_PROXY=http://0.0.0.0:1 \
  go test -v -run TestUpstreamAdvisory ./tests/spec_check/ 2>&1 | head -20
```

Expected: still PASS; findings say "upstream unreachable" or similar.

## Task 4.3: Wire L7 into Makefile

**Files:**
- Modify: `Makefile`

- [ ] **Step 1: Add verify-upstream target and chain it into verify**

Replace the `verify: verify-fast` line with:

```makefile
verify: verify-fast verify-upstream
```

Add a new target block (place it after `verify-spec-sync`):

```makefile
verify-upstream:
	@echo "==> L7 upstream advisory (advisory only, never fails PR)"
	cd tests/spec_check && go test -v -run TestUpstreamAdvisory .
```

- [ ] **Step 2: Run `make verify` and confirm L7 runs as the last layer**

```bash
make verify
```

Expected: L1-L6 run, then `==> L7 upstream advisory` runs and produces logs (PASS regardless of findings).

## Task 4.4: Update verify.yml to run full verify (with L7)

**Files:**
- Modify: `.github/workflows/verify.yml`

- [ ] **Step 1: Change the verify step to use `make verify` instead of `make verify-fast`**

In the existing workflow, find:

```yaml
      - name: Run make verify-fast
        id: verify
        run: |
          set -o pipefail
          make verify-fast 2>&1 | tee verify.log
```

Replace with:

```yaml
      - name: Run make verify (L1-L7)
        id: verify
        run: |
          set -o pipefail
          make verify 2>&1 | tee verify.log
```

## Task 4.5: Commit Phase 4

- [ ] **Step 1: Commit**

```bash
git add tests/spec_check/upstream_test.go tests/spec_check/go.mod tests/spec_check/go.sum Makefile .github/workflows/verify.yml
git commit -m "$(cat <<'EOF'
feat(verify): L7 upstream advisory (network, soft-fail)

Best-effort fetch each CLI's spec_url via 5s-timeout HTTP. Parses
both upstream and vendored specs via kin-openapi; computes an
operation-set diff and reports removed/added operations. Findings
are ranked: operations present in vendored (backing real CLI
subcommands) are 'high' priority; upstream-only operations are
'info'.

L7 never red-lights a PR. Network failure / 5xx / parse failure
all log as advisory and pass cleanly.

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

---

# Phase 5 — L8 Live Cloud + OIDC

After this phase: `make verify` includes L8 (live cloud), which auto-skips clouds without OIDC creds and hard-fails clouds where creds are present but the API call fails. The verify.yml workflow gains three optional `configure-creds` steps gated by repo secrets.

## Task 5.1: Add `live_smoke` field to catalog.yaml entries

**Files:**
- Modify: `catalog.yaml`

- [ ] **Step 1: Add `live_smoke.list_subcommand` to each CLI entry**

For each of the 6 CLI entries in `catalog.yaml`, insert a `live_smoke:` block immediately before the `daily_commands:` line. Use the following per-CLI mapping:

| CLI | `list_subcommand` |
|---|---|
| cloud-run-admin-pp-cli | `services-list` |
| cloud-functions-pp-cli | `functions-list` |
| lambda-pp-cli | `functions-list` |
| apprunner-pp-cli | `services-list` |
| functions-pp-cli | `sites-list` |
| container-apps-pp-cli | `managed-environments-list` |

Each block looks like:

```yaml
    live_smoke:
      list_subcommand: services-list
      args: []
```

(Substitute the per-CLI value.)

**Note:** the exact subcommand names depend on what press generated. If a particular subcommand doesn't exist for a CLI, leave `list_subcommand` empty (`""`) — L8 will fall back to running only the `doctor` reachability check for that CLI.

- [ ] **Step 2: Verify catalog still parses**

```bash
yq '.clis[].live_smoke.list_subcommand' catalog.yaml
```

Expected: 6 lines, one per CLI.

## Task 5.2: Create `tests/live/` Go module with creds helper

**Files:**
- Create: `tests/live/go.mod`
- Create: `tests/live/helpers.go`
- Create: `tests/live/creds.go`

- [ ] **Step 1: Initialize module**

```bash
mkdir -p tests/live
cd tests/live
go mod init github.com/ShubhanYenuganti/infra-press/tests/live
go get gopkg.in/yaml.v3
cd ../..
```

- [ ] **Step 2: Write `tests/live/helpers.go`** (duplicates catalog loader — intentional, see plan header)

```go
//go:build live

package live

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type LiveSmoke struct {
	ListSubcommand string   `yaml:"list_subcommand"`
	Args           []string `yaml:"args"`
}

type CLI struct {
	Name      string    `yaml:"name"`
	Cloud     string    `yaml:"cloud"`
	Path      string    `yaml:"path"`
	Binary    string    `yaml:"binary"`
	Status    string    `yaml:"status"`
	LiveSmoke LiveSmoke `yaml:"live_smoke"`
}

type Catalog struct {
	Version int   `yaml:"version"`
	CLIs    []CLI `yaml:"clis"`
}

func RepoRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "catalog.yaml")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

func LoadCatalog(repoRoot string) (*Catalog, error) {
	b, err := os.ReadFile(filepath.Join(repoRoot, "catalog.yaml"))
	if err != nil {
		return nil, err
	}
	var c Catalog
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func BinaryPath(repoRoot string, cli CLI) string {
	return filepath.Join(repoRoot, cli.Path, "cmd", cli.Binary, cli.Binary)
}
```

- [ ] **Step 3: Write `tests/live/creds.go`**

```go
//go:build live

package live

import "os"

// HasCreds returns true if env vars indicating OIDC creds for the given
// cloud are present. The vars are set by the canonical GitHub Actions
// configure-creds actions for each cloud.
func HasCreds(cloud string) bool {
	switch cloud {
	case "aws":
		return os.Getenv("AWS_ROLE_ARN") != "" && os.Getenv("AWS_WEB_IDENTITY_TOKEN_FILE") != ""
	case "gcp":
		return os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") != "" ||
			os.Getenv("GCP_WORKLOAD_IDENTITY_PROVIDER") != ""
	case "azure":
		return os.Getenv("AZURE_CLIENT_ID") != "" && os.Getenv("AZURE_TENANT_ID") != ""
	default:
		return false
	}
}
```

## Task 5.3: Implement L8 generic test (build-tag gated)

**Files:**
- Create: `tests/live/live_test.go`

- [ ] **Step 1: Write the test**

```go
//go:build live

package live

import (
	"encoding/json"
	"os/exec"
	"testing"
)

func TestLiveCloud(t *testing.T) {
	root, err := RepoRoot()
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	cat, err := LoadCatalog(root)
	if err != nil {
		t.Fatalf("load catalog: %v", err)
	}
	for _, cli := range cat.CLIs {
		cli := cli
		t.Run(cli.Name, func(t *testing.T) {
			if !HasCreds(cli.Cloud) {
				t.Skipf("no OIDC creds for %s in env", cli.Cloud)
			}
			bin := BinaryPath(root, cli)

			t.Run("doctor-reaches-cloud", func(t *testing.T) {
				out, err := exec.Command(bin, "doctor", "--json").Output()
				if err != nil {
					t.Fatalf("doctor --json failed: %v\noutput: %s", err, out)
				}
				var report map[string]interface{}
				if err := json.Unmarshal(out, &report); err != nil {
					t.Fatalf("doctor --json not parseable: %v\noutput: %s", err, out)
				}
				// Convention: doctor reports "status": "ok" on full reachability.
				if status, _ := report["status"].(string); status != "ok" {
					t.Fatalf("doctor reports non-ok status: %v", report)
				}
			})

			if cli.LiveSmoke.ListSubcommand == "" {
				t.Log("no list_subcommand declared; skipping list-endpoint assertion")
				return
			}

			t.Run("list-subcommand-returns-json", func(t *testing.T) {
				args := append([]string{cli.LiveSmoke.ListSubcommand, "--json"}, cli.LiveSmoke.Args...)
				out, err := exec.Command(bin, args...).Output()
				if err != nil {
					t.Fatalf("%s failed: %v\noutput: %s", cli.LiveSmoke.ListSubcommand, err, out)
				}
				var v interface{}
				if err := json.Unmarshal(out, &v); err != nil {
					t.Fatalf("%s output not parseable JSON: %v", cli.LiveSmoke.ListSubcommand, err)
				}
			})
		})
	}
}
```

- [ ] **Step 2: Confirm the test does NOT compile under default build**

```bash
cd tests/live
go test -v . 2>&1 | head -10
cd ../..
```

Expected: `no Go files in tests/live` or `build constraints exclude all Go files` — i.e., default build skips the package entirely. This is the design guarantee that `go test ./...` from anywhere can never accidentally invoke live cloud APIs.

- [ ] **Step 3: Run with `-tags=live` (will skip every subtest locally because no OIDC creds)**

```bash
cd tests/live
go test -tags=live -v .
cd ../..
```

Expected: PASS with every subtest marked `SKIP` — reason `no OIDC creds for aws/gcp/azure in env`.

## Task 5.4: Wire L8 into Makefile

**Files:**
- Modify: `Makefile`

- [ ] **Step 1: Update verify-live target and chain into verify**

Replace the `verify` target line with:

```makefile
verify: verify-fast verify-upstream verify-live
```

Append (or update) the `verify-live` target:

```makefile
verify-live:
	@echo "==> L8 live cloud (skips clouds without OIDC creds)"
	cd tests/live && go test -tags=live -v .
```

- [ ] **Step 2: Run full `make verify`**

```bash
make verify
```

Expected: L1-L7 run as before, then L8 runs and SKIPs every cloud (no creds locally). Final exit is 0.

## Task 5.5: Update verify.yml with OIDC configure-creds steps

**Files:**
- Modify: `.github/workflows/verify.yml`

- [ ] **Step 1: Insert three optional auth steps before `make verify`**

Edit `.github/workflows/verify.yml`. After the `actions/setup-go` step and before `Install yq`, insert these three steps:

```yaml
      - name: Configure AWS OIDC (optional)
        if: ${{ vars.AWS_OIDC_ROLE_ARN != '' }}
        uses: aws-actions/configure-aws-credentials@v4
        with:
          role-to-assume: ${{ vars.AWS_OIDC_ROLE_ARN }}
          aws-region: ${{ vars.AWS_REGION || 'us-east-1' }}

      - name: Configure GCP OIDC (optional)
        if: ${{ vars.GCP_WORKLOAD_IDENTITY_PROVIDER != '' }}
        uses: google-github-actions/auth@v2
        with:
          workload_identity_provider: ${{ vars.GCP_WORKLOAD_IDENTITY_PROVIDER }}
          service_account: ${{ vars.GCP_SERVICE_ACCOUNT }}

      - name: Configure Azure OIDC (optional)
        if: ${{ vars.AZURE_CLIENT_ID != '' }}
        uses: azure/login@v2
        with:
          client-id: ${{ vars.AZURE_CLIENT_ID }}
          tenant-id: ${{ vars.AZURE_TENANT_ID }}
          subscription-id: ${{ vars.AZURE_SUBSCRIPTION_ID }}
```

- [ ] **Step 2: Validate YAML**

```bash
python3 -c "import yaml; yaml.safe_load(open('.github/workflows/verify.yml'))" && echo OK
```

Expected: `OK`.

## Task 5.6: Commit Phase 5

- [ ] **Step 1: Commit**

```bash
git add catalog.yaml tests/live/ Makefile .github/workflows/verify.yml
git commit -m "$(cat <<'EOF'
feat(verify): L8 live cloud + OIDC plumbing

New tests/live/ module gated by //go:build live so default test runs
cannot accidentally hit cloud APIs. Per-CLI catalog-driven generic
test asserts doctor --json reports status:ok and (if declared) the
list_subcommand returns parseable JSON.

verify.yml gains three optional configure-creds steps, each gated by
a repo Variable being set (AWS_OIDC_ROLE_ARN, GCP_WORKLOAD_
IDENTITY_PROVIDER, AZURE_CLIENT_ID). Missing variable = step skipped
= env vars unset = L8 auto-skips that cloud. Zero configuration
required to opt out of any cloud.

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

Expected: clean commit.

---

# Phase 6 — Slash Command + Gitignore Fix

After this phase: `/verify-pr <num>` checks out a PR, runs `make verify`, parses the output, and posts a markdown summary as a PR comment using a stable marker so re-runs update in place.

## Task 6.1: Allow `.claude/commands/` through .gitignore

**Files:**
- Modify: `.gitignore`

- [ ] **Step 1: Update gitignore to track project slash commands**

Replace the current `.claude/` line:

```
# local Claude Code config
.claude/
```

With a more selective pattern:

```
# local Claude Code config — ignore everything except project-shared slash commands
.claude/*
!.claude/commands/
.claude/commands/.DS_Store
```

- [ ] **Step 2: Verify**

```bash
mkdir -p .claude/commands
touch .claude/commands/.gitkeep
git status --short .claude/
```

Expected: `.claude/commands/.gitkeep` shows as untracked-but-allowed (not ignored). Other `.claude/*` paths still ignored.

## Task 6.2: Create the verify-pr slash command

**Files:**
- Create: `.claude/commands/verify-pr.md`

- [ ] **Step 1: Write the slash command**

```markdown
---
description: Check out a PR, run the v1 verification harness, and post a ranked markdown report as a PR comment.
argument-hint: <PR-number>
---

# /verify-pr — run the verification harness against a PR

Use this command when reviewing a PR locally or when CI's verify workflow has run and you want a fuller human-readable report.

## What this command does

1. **Check out the PR locally** via `gh pr checkout $1`.
2. **Run `make verify`** from the repo root. This invokes L1 through L8 (L8 auto-skips clouds without OIDC creds in the local environment).
3. **Capture stdout to `verify.log`** in the repo root.
4. **Parse the log** into per-layer pass/fail/info counts.
5. **Render a markdown report** with these sections:
   - Header: `## Verification report — PR #$1 — commit <sha>`
   - Table: one row per layer (L1-L8) with status + counts.
   - L7 advisory subsection: rank "high" findings first, "info" findings collapsed under a `<details>` block.
   - L8 subsection: per-cloud SKIP / PASS / FAIL breakdown.
6. **Post the report** as a PR comment using a stable marker so re-runs update in place:

```bash
gh pr comment "$1" --body-file /tmp/verify-pr-report.md \
  --edit-last  # if a prior comment with the marker exists, update it
```

Use the marker `<!-- verify-pr -->` at the top of every report so subsequent runs can detect and replace the previous comment.

## Important caveats

- Running this locally requires `make`, `go`, `yq`, and (for Azure CLIs to build) network access to apis.guru. Confirm via `which make go yq` before invoking.
- L8 requires real cloud credentials in the environment (`aws sts get-caller-identity`, `gcloud auth list`, `az account show`). Without them, L8 SKIPs every cloud — that's expected and not a failure.
- The command will checkout the PR branch — make sure your working tree is clean first (`git status` shows clean) or stash before invoking.

## Output format

The posted comment should look like:

```markdown
<!-- verify-pr -->
## Verification report — PR #123 — commit abc1234

| Layer | Status | Detail |
|---|---|---|
| L1 Lint | PASS 30/30 | |
| L2 Build | PASS 6/6 | |
| L3 Golden | PASS 24/24 | |
| L4 Smoke | PASS 30/30 | |
| L5 Recipe | PASS | doctor-all.sh: 6 CLIs reported |
| L6 Spec sync | PASS 6/6 | |
| L7 Upstream | INFO | 2 high-priority findings |
| L8 Live cloud | SKIP | OIDC not configured for AWS/GCP/Azure |

### L7 advisory findings (ranked)

**High priority:**
- `lambda-pp-cli`: operation removed upstream: GET /functions::ListFunctions
- `cloud-run-admin-pp-cli`: arg `pageSize` changed from optional to required

<details>
<summary>Informational (3 findings)</summary>

- `cloud-functions-pp-cli`: 5 new endpoints available upstream (not exposed)
- ...
</details>

**Result: PASS** (L1–L7 green; L8 not configured)
```
```

## Task 6.3: Commit Phase 6

- [ ] **Step 1: Commit**

```bash
git add .gitignore .claude/commands/verify-pr.md
git commit -m "$(cat <<'EOF'
feat(claude): /verify-pr slash command + un-ignore .claude/commands/

The slash command wraps gh pr checkout + make verify + parsing +
gh pr comment posting. Uses a <!-- verify-pr --> marker so re-runs
update the existing PR comment in place rather than appending.

.gitignore refined to allow .claude/commands/ through while keeping
local-only config (.claude/settings.local.json, .claude/cache, etc.)
ignored.

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

---

# Final Verification

After all 6 phases land:

- [ ] **Run the full harness locally**

```bash
make verify
```

Expected: L1-L7 PASS hard; L8 SKIPs every cloud locally; exit 0.

- [ ] **Push to a feature branch and open a PR**

The verify.yml workflow runs automatically. Confirm:
- The workflow appears as a required-ish status check on the PR (depending on branch protection rules — set those separately).
- The workflow's job summary contains the tail of `verify.log`.

- [ ] **Run `/verify-pr <num>` against the same PR**

Confirm a PR comment is posted with the markdown report and that running it again **edits** the existing comment rather than appending a new one.

- [ ] **Confirm the final layer table is documented**

The Phase tags in commits 1-6 plus the spec at `docs/superpowers/specs/2026-05-21-v1-verification-design.md` constitute the complete documentation. No additional README is needed.

---

# Notes for the executing engineer

- **Phase ordering is sticky.** Each phase depends on prior phases. Phase 4's `kin-openapi` dependency wouldn't make sense without Phase 1's vendored fixtures.
- **Don't run `git push`.** Per repo convention. Each phase commits locally; the user pushes manually.
- **Catalog is the source of truth.** Every layer reads `catalog.yaml`. If a CLI is added/removed in the future, only `catalog.yaml` + the matching pressfile/fixture changes — no harness code edits.
- **L7 findings are not bugs.** Drift findings reported by L7 indicate the upstream changed; they do NOT mean the PR is broken. Handle them in a separate "regen the CLI" workflow (deferred — not part of this plan).
- **If a `live_smoke.list_subcommand` value is wrong** (i.e., the CLI doesn't actually have that subcommand), L4's `help-mentions-required-subcommands` will pass (it only checks the 9 universal subcommands) but L8 will fail when OIDC is configured. Update the catalog value to a real subcommand name and re-run.
- **If kin-openapi fails to parse the converted Azure specs**, the L7 test will log "vendored spec did not parse" advisory — that's a separate finding worth investigating, but it does not red-light the PR. The L6 sync check uses JSON-level parsing (no schema validation), so it stays green regardless.
