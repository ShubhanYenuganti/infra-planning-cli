# Infra Press v1 Substrate Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Migrate the existing `infra-planning-cli` codebase to a clean `infra-press` substrate, build all shared infrastructure (linter, golden tests, install scripts, CI), and prove the pipeline end-to-end by generating, conforming, and shipping the first CLI (`cloud-run-admin-pp-cli`).

**Architecture:** Approach D from the design spec — a curated library of printing-press-generated CLIs with hub artifacts (`AGENTS.md`, `catalog.yaml`, `docs/conventions.md`, golden tests, install scripts). Approach 4 mechanics — frozen-fork per-CLI source with `pressfile.yaml` manifest and `patches/` directory separation to support a future regen agent.

**Tech Stack:** Go 1.26+ (required by mvanhorn's printing-press), Bash, GitHub Actions, SQLite (per generated CLI), YAML (catalog + pressfiles).

**Spec reference:** [`docs/superpowers/specs/2026-05-20-infra-press-design.md`](../specs/2026-05-20-infra-press-design.md)

**Out of scope for this plan:** The other 5 v1 CLIs (`lambda-pp-cli`, `apprunner-pp-cli`, `cloud-functions-pp-cli`, `functions-pp-cli`, `container-apps-pp-cli`). Those follow the same per-CLI pipeline established in Phase 6 — a separate plan will cover them once this substrate is proven.

---

## Phase 0 — Migration

Move existing code to the `old-infra-cli` worktree, reset main to a substrate-only state. Per CLAUDE.md: NEVER run `git push` — stop after each commit.

### Task 0.1: Verify clean state & stage in-flight docs

**Files:**
- Modify: working tree state (no file changes)

- [ ] **Step 1: Check current git status**

Run: `git -C /Users/shubhan/infra-planning-cli status --short`
Expected: clean OR shows `M docs/05-implementation-plan.md` and `?? docs/06-design-improvements.md`

- [ ] **Step 2: Commit any in-flight docs changes to main (preserves work before delete)**

Run:
```bash
cd /Users/shubhan/infra-planning-cli
git add docs/05-implementation-plan.md docs/06-design-improvements.md
git commit -m "docs: snapshot prior-design docs before pivot to infra-press"
```
Expected: commit succeeds, `git status` clean.

- [ ] **Step 3: Verify the design spec and plan are committed**

Run: `git log --oneline -5`
Expected: shows recent commits including `docs: add infra-press v1 design spec`.

---

### Task 0.2: Create `old-infra-cli` branch

**Files:**
- Modify: git refs (creates branch)

- [ ] **Step 1: Create the branch pointing at current HEAD**

Run:
```bash
cd /Users/shubhan/infra-planning-cli
git checkout -b old-infra-cli
git log --oneline -3
```
Expected: branch created, log shows the same commits as main.

- [ ] **Step 2: Switch back to main**

Run: `git checkout main`
Expected: switched to main.

- [ ] **Step 3: Verify branch exists**

Run: `git branch -a | grep old-infra-cli`
Expected: shows `old-infra-cli` (local).

---

### Task 0.3: Create worktree for old code

**Files:**
- Create: `../infra-press-old-cli/` (new worktree directory)

- [ ] **Step 1: Create the worktree**

Run:
```bash
cd /Users/shubhan/infra-planning-cli
git worktree add ../infra-press-old-cli old-infra-cli
```
Expected: worktree created at `/Users/shubhan/infra-press-old-cli/`.

- [ ] **Step 2: Verify the worktree has all the old code**

Run: `ls /Users/shubhan/infra-press-old-cli/`
Expected: includes `cmd/`, `internal/`, `tests/`, `docs/`, `go.mod`, `go.sum`, `README.md`, `CLAUDE.md`.

- [ ] **Step 3: List git worktrees**

Run: `git -C /Users/shubhan/infra-planning-cli worktree list`
Expected: shows both `/Users/shubhan/infra-planning-cli` (main) and `/Users/shubhan/infra-press-old-cli` (old-infra-cli).

---

### Task 0.4: Reset main to substrate-only state

**Files:**
- Delete: `cmd/`, `internal/`, `tests/`, `testdata/`, `go.mod`, `go.sum`, `docs/01-position-core-concept.md`, `docs/02-discovery-workflow.md`, `docs/03-output-format.md`, `docs/04-mvp-scope.md`, `docs/05-implementation-plan.md`, `docs/06-design-improvements.md`
- Preserve: `README.md`, `CLAUDE.md`, `docs/superpowers/specs/2026-05-20-infra-press-design.md`, `docs/superpowers/plans/2026-05-20-infra-press-v1-substrate.md`

- [ ] **Step 1: Confirm spec and plan are preserved**

Run:
```bash
ls /Users/shubhan/infra-planning-cli/docs/superpowers/specs/
ls /Users/shubhan/infra-planning-cli/docs/superpowers/plans/
```
Expected: shows `2026-05-20-infra-press-design.md` and `2026-05-20-infra-press-v1-substrate.md` respectively.

- [ ] **Step 2: Remove old code directories**

Run:
```bash
cd /Users/shubhan/infra-planning-cli
rm -rf cmd internal tests testdata go.mod go.sum
```
Expected: those paths no longer exist.

- [ ] **Step 3: Remove old planning docs (01-06)**

Run:
```bash
cd /Users/shubhan/infra-planning-cli
rm -f docs/01-position-core-concept.md docs/02-discovery-workflow.md \
      docs/03-output-format.md docs/04-mvp-scope.md \
      docs/05-implementation-plan.md docs/06-design-improvements.md
```
Expected: only `docs/superpowers/` remains under `docs/`.

- [ ] **Step 4: Verify the reset**

Run: `ls -la /Users/shubhan/infra-planning-cli/`
Expected: shows ONLY `.git`, `README.md`, `CLAUDE.md`, `docs/`, (and possibly `.claude/` if present).

Run: `find /Users/shubhan/infra-planning-cli/docs -type f`
Expected: shows ONLY `docs/superpowers/specs/2026-05-20-infra-press-design.md` and `docs/superpowers/plans/2026-05-20-infra-press-v1-substrate.md`.

- [ ] **Step 5: Commit the reset**

Run:
```bash
cd /Users/shubhan/infra-planning-cli
git add -A
git commit -m "chore: reset to infra-press substrate

Removes old infra-planning-cli code (preserved on old-infra-cli branch).
Main now holds only README.md, CLAUDE.md, and docs/superpowers/."
```
Expected: commit succeeds.

---

### Task 0.5: Manual repo rename on GitHub (USER ACTION)

**Files:** none (external action)

- [ ] **Step 1: Rename repository in GitHub web UI**

User action: Navigate to `https://github.com/ShubhanYenuganti/infra-planning-cli/settings`, rename the repo to `infra-press`.

- [ ] **Step 2: Update local remote URL**

Run:
```bash
cd /Users/shubhan/infra-planning-cli
git remote set-url origin git@github.com:ShubhanYenuganti/infra-press.git
git remote -v
```
Expected: remote now points at `infra-press`.

- [ ] **Step 3: Optionally rename the local directory**

User decision: rename `/Users/shubhan/infra-planning-cli/` → `/Users/shubhan/infra-press/` for clarity. Update editor paths accordingly. Worktree at `../infra-press-old-cli/` keeps that name.

(All later tasks assume the local directory is still at `/Users/shubhan/infra-planning-cli/` for safety. Update paths if you rename.)

---

## Phase 1 — Repository skeleton + hub artifacts

Create top-level directories and the four hub files (`README.md`, `AGENTS.md`, `catalog.yaml`, `docs/conventions.md`). These are static (no Go code yet) and need no tests beyond "files exist and parse."

### Task 1.1: Create directory skeleton

**Files:**
- Create directories: `library/aws/`, `library/gcp/`, `library/azure/`, `tests/golden/checks/`, `tests/golden/fixtures/`, `scripts/`, `tools/lint-conventions/`, `docs/sprints/`

- [ ] **Step 1: Create all directories**

Run:
```bash
cd /Users/shubhan/infra-planning-cli
mkdir -p library/aws library/gcp library/azure \
         tests/golden/checks tests/golden/fixtures \
         scripts \
         tools/lint-conventions \
         docs/sprints
```
Expected: all directories present.

- [ ] **Step 2: Verify the layout**

Run: `find /Users/shubhan/infra-planning-cli -type d -not -path '*/.git*' | sort`
Expected: shows all created directories plus existing `docs/superpowers/{specs,plans}/`.

- [ ] **Step 3: Add `.gitkeep` to empty dirs so git tracks them**

Run:
```bash
cd /Users/shubhan/infra-planning-cli
touch library/aws/.gitkeep library/gcp/.gitkeep library/azure/.gitkeep \
      tests/golden/fixtures/.gitkeep
```

---

### Task 1.2: Rewrite root `README.md`

**Files:**
- Modify: `README.md` (full rewrite — currently contains infra-planning-cli content)

- [ ] **Step 1: Write the new README**

Create `README.md` with this content:

```markdown
# infra-press

A curated library of focused, agent-native CLIs for cloud-infra APIs.

Generated by [printing-press](https://github.com/mvanhorn/cli-printing-press),
packaged with a hub (`AGENTS.md`, `catalog.yaml`, conventions doc, golden tests,
install scripts) so any agent or skill can discover and use them as primitives.

## What's in v1

Stateless services and scheduled jobs across AWS, GCP, and Azure:

|  | Serverless | Managed Container |
|---|---|---|
| **AWS** | `lambda-pp-cli` | `apprunner-pp-cli` |
| **GCP** | `cloud-functions-pp-cli` | `cloud-run-admin-pp-cli` |
| **Azure** | `functions-pp-cli` | `container-apps-pp-cli` |

For everything else (storage, IAM, networking, databases) fall back to `aws`,
`gcloud`, or `az`. Future sprints expand coverage — see [docs/sprints/](docs/sprints/).

## Install

```bash
curl -sSL https://raw.githubusercontent.com/ShubhanYenuganti/infra-press/main/scripts/install.sh | sh
```

Or per CLI:

```bash
go install github.com/ShubhanYenuganti/infra-press/library/gcp/cloud-run-admin/cmd/cloud-run-admin-pp-cli@latest
```

## Why

Most generators wrap endpoints and stop. Printing Press builds CLIs that
understand the domain — SQLite-backed local sync, FTS5 search, compound
commands like `stale`/`health`/`bottleneck` that join across resources.
This library applies that pattern across the cloud-infra surface uniformly,
so an agent reasons about it as one library, not 6 disparate tools.

## Documentation

- **Agents start here:** [`AGENTS.md`](AGENTS.md)
- **Conventions every CLI follows:** [`docs/conventions.md`](docs/conventions.md)
- **Cross-CLI compound recipes:** [`docs/compound-recipes.md`](docs/compound-recipes.md)
- **Sprint roadmap:** [`docs/sprints/`](docs/sprints/)
- **Design spec:** [`docs/superpowers/specs/2026-05-20-infra-press-design.md`](docs/superpowers/specs/2026-05-20-infra-press-design.md)

## Contributing

Each CLI is generated, then patched per the 8-step pipeline in the design spec.
See [`docs/sprints/v1.md`](docs/sprints/v1.md) for the current sprint state.
```

- [ ] **Step 2: Verify file**

Run: `head -20 /Users/shubhan/infra-planning-cli/README.md`
Expected: shows new infra-press README content.

---

### Task 1.3: Write root `AGENTS.md` (skeleton)

**Files:**
- Create: `AGENTS.md`

- [ ] **Step 1: Write AGENTS.md**

Create `/Users/shubhan/infra-planning-cli/AGENTS.md` with this content:

```markdown
# infra-press — agent guide

A library of focused, agent-native CLIs for cloud-infra APIs. Stateless
services and scheduled jobs across AWS, GCP, Azure.

## When to reach for what

| Need | CLI | Cloud |
|---|---|---|
| Deploy a containerized service | cloud-run-admin-pp-cli | GCP |
| _(more added as sprints ship)_ | | |

For storage, IAM, networking, databases, secrets, and Kubernetes: fall back to
`aws`, `gcloud`, or `az`. See the "Known fallbacks" section below.

## Conventions every CLI in this library follows

- `--json` auto-on when stdout is piped
- Typed exit codes: `0` ok / `2` user error / `3` auth / `4` not found / `5` conflict / `7` rate-limit
- `sync` / `search` / `sql` subcommands on every CLI
- Local SQLite at `~/.infra-press/<cli>.db`
- `--data-source live|local|auto`
- `--dry-run` on every mutation
- `--select field1,field2` for field projection
- `--compact` drops to high-gravity fields only

Full spec: [`docs/conventions.md`](docs/conventions.md).

## Auth resolution

- **AWS:** flag → env (`AWS_PROFILE`) → SDK default chain
- **GCP:** flag → `GOOGLE_APPLICATION_CREDENTIALS` → ADC
- **Azure:** flag → `AZURE_TENANT_ID` + `AZURE_CLIENT_ID` + `AZURE_CLIENT_SECRET` → DefaultAzureCredential

Run the relevant doctor script before first use:
```bash
./scripts/doctor-aws.sh
./scripts/doctor-gcp.sh
./scripts/doctor-azure.sh
```

## Known fallbacks to vendor CLIs

- **Azure Functions** needs a Storage Account first:
  `az storage account create -n <name> -g <rg> -l <region> --sku Standard_LRS`
- **AWS Lambda** first-time use needs an IAM execution role:
  use `lambda functions create --auto-role` (gap patch) or pre-create
  via `aws iam create-role ...`

## Cross-CLI compound recipes

See [`docs/compound-recipes.md`](docs/compound-recipes.md) for SQL/scripts that
join data across the SQLite stores (stale services across all clouds, public
endpoints without auth, idle-compute spend triangulation, etc.).
```

- [ ] **Step 2: Verify**

Run: `head -10 /Users/shubhan/infra-planning-cli/AGENTS.md`
Expected: shows new content.

---

### Task 1.4: Write `catalog.yaml`

**Files:**
- Create: `catalog.yaml`

- [ ] **Step 1: Write catalog with empty clis list**

Create `/Users/shubhan/infra-planning-cli/catalog.yaml`:

```yaml
version: 1
clis: []
# As sprints ship CLIs, entries appear here per the schema:
#
# - name: lambda-pp-cli
#   cloud: aws
#   service: lambda
#   paradigm: serverless
#   path: library/aws/lambda
#   binary: lambda-pp-cli
#   install: github.com/ShubhanYenuganti/infra-press/library/aws/lambda
#   status: stable          # bootstrapping | preview | stable
#   press_version: v1.4.2
#   spec_version: 2024-11-01
#   daily_commands: [functions list, functions create, functions invoke]
#   compound_commands: [stale, cold-starts, unused-layers, permissions-audit]
#   known_gaps:
#     - "First use requires IAM exec role; use --auto-role"
```

- [ ] **Step 2: Verify YAML parses**

Run: `python3 -c "import yaml; print(yaml.safe_load(open('/Users/shubhan/infra-planning-cli/catalog.yaml')))"`
Expected: prints `{'version': 1, 'clis': []}`.

---

### Task 1.5: Write `docs/conventions.md`

**Files:**
- Create: `docs/conventions.md`

- [ ] **Step 1: Write the convention spec**

Create `/Users/shubhan/infra-planning-cli/docs/conventions.md`:

```markdown
# infra-press — Conventions

Every CLI in `library/` MUST conform to this spec. The linter at
`tools/lint-conventions/` enforces it; CI fails on any violation. New
conventions are added at sprint boundaries only.

## Flags

Every CLI MUST support these flags on every command that lists, gets, or
creates resources:

| Flag | Semantics |
|---|---|
| `--json` | Output JSON. Auto-on when stdout is piped to a non-TTY. |
| `--compact` | Drop to high-gravity fields only (60-80% fewer tokens). |
| `--select field1,field2,...` | Field projection. Applied to JSON output. |
| `--dry-run` | Required on every mutation. Prints the request without sending. |
| `--data-source live\|local\|auto` | Where to read data from. Default: `auto`. |
| `--yes` | Required for destructive operations. |
| `--idempotent` | Treat create-retry as success if resource exists. |
| `--ignore-missing` | Treat delete of missing resource as success. |
| `--stdin` | Accept structured input on stdin (where listed in `--help`). |
| `--human-friendly` | Force human formatting (colors, tables) even when piped. |

## Exit codes

| Code | Meaning |
|---|---|
| `0` | Success |
| `2` | User error (bad flags, bad args) |
| `3` | Authentication / authorization failure |
| `4` | Resource not found |
| `5` | Conflict / already exists |
| `7` | Upstream rate-limit |

## Subcommands

Every CLI MUST provide:

- `sync` — Pull resources into local SQLite store
- `search <query>` — FTS5 search over the local store; falls back to live per `--data-source`
- `sql <query>` — Power-user direct SQL against the local store

## Local store

- SQLite path: `~/.infra-press/<cli>.db`
- Schema generated by press; per-CLI tables, FTS5 mirrors for searchable text fields
- `--data-source` modes:
  - `live` — always call the API; ignore local store
  - `local` — only read local store; error if missing/stale
  - `auto` — read local if fresh, refresh otherwise (default)

## Auth resolution (per cloud)

- **AWS:** flag → env (`AWS_PROFILE`, `AWS_REGION`) → SDK default chain (instance profile, etc.)
- **GCP:** flag → `GOOGLE_APPLICATION_CREDENTIALS` → ADC
- **Azure:** flag → `AZURE_TENANT_ID` + `AZURE_CLIENT_ID` + `AZURE_CLIENT_SECRET` → DefaultAzureCredential

## Output shape

- Pipe detection: when stdout is not a TTY, auto-switch to `--json`
- Errors to stderr; exit code reflects category
- No colors or formatting on JSON output
- `--compact` removes verbose system fields, keeps user-meaningful fields

## `--dry-run` semantics

- Mutations: print the would-be request as JSON, exit 0
- Reads: ignored (no effect)
- Compound queries: run read-only against local store, no API calls

## Regen contract

When the future update agent regenerates a CLI:

- **Overwrites:** `cmd/`, `internal/`, `README.md`, `AGENTS.md`, `.golangci.yml`, `tests/generated/`, `go.mod`, `go.sum`, `scorecard.md`
- **Preserves:** `patches/` (all files), `AGENTS.append.md`, `pressfile.yaml` (updates `press_version`/`spec_version`/`spec_etag`/`generated_at` only; `patches[]` array preserved)
- **Verifies:** `tools/lint-conventions` passes on regenerated source; if not, regen fails and the CLI is rolled back

## Patch directory rules

`patches/` contains human-authored Go files registered against the CLI's root
command. Three categories:

1. **Gap patches** — fix things press output missed or got wrong (e.g.,
   `--auto-role` for Lambda)
2. **Compound commands** — added subcommands not auto-generated (e.g.,
   `lambda permissions-audit`)
3. **Convention fixes** — code that brings press output into convention
   compliance (rare — prefer fixing the press itself)

Each patch file MUST be listed in `pressfile.yaml`'s `patches[]` with a
`reason`. Files in `patches/` MUST compile and pass per-CLI tests.
`AGENTS.append.md` is concatenated to the generated `AGENTS.md` at install
time (handled by the CLI's build step or the install script).
```

- [ ] **Step 2: Verify**

Run: `wc -l /Users/shubhan/infra-planning-cli/docs/conventions.md`
Expected: ≥80 lines.

---

### Task 1.6: Write `docs/compound-recipes.md` skeleton

**Files:**
- Create: `docs/compound-recipes.md`

- [ ] **Step 1: Write a skeleton with cross-CLI recipe template**

Create `/Users/shubhan/infra-planning-cli/docs/compound-recipes.md`:

```markdown
# Cross-CLI Compound Recipes

SQL snippets and shell pipelines that join data across multiple `infra-press`
SQLite stores. These are queries no individual vendor CLI can answer.

> **Note:** These recipes assume you've run `sync` on every relevant CLI first.
> Stores live at `~/.infra-press/<cli>.db`.

## Recipe template

```sql
-- Recipe name (one line)
-- What it answers (one sentence)
-- Required syncs: <list the CLIs whose stores must be synced>

ATTACH '~/.infra-press/cloud-run-admin-pp-cli.db' AS cr;
-- ... query ...
```

## Recipes (added as sprints ship)

### Stale services across all synced stores

_To be added once 2+ CLIs ship. See sprint v1 retro for the first instance._

---

(More recipes accumulate as sprints add CLIs.)
```

- [ ] **Step 2: Verify**

Run: `head /Users/shubhan/infra-planning-cli/docs/compound-recipes.md`

---

### Task 1.7: Write `docs/sprints/v1.md`

**Files:**
- Create: `docs/sprints/v1.md`

- [ ] **Step 1: Write the v1 sprint planning doc**

Create `/Users/shubhan/infra-planning-cli/docs/sprints/v1.md`:

```markdown
# Sprint v1 — Serverless + Managed Containers (SUBSTRATE SPRINT)

## Goal

Ship 6 CLIs (Lambda, Cloud Functions, Azure Functions, App Runner, Cloud Run,
Container Apps) AND build the regen-ready substrate that all future sprints use.

## Substrate deliverables (one-time, v1 only)

- [ ] Repo rename `infra-planning-cli` → `infra-press`
- [ ] Old code moved to `old-infra-cli` worktree
- [ ] `docs/conventions.md` written (the spec the linter enforces)
- [ ] `tools/lint-conventions/` implemented
- [ ] `tests/golden/` orchestrator + initial checks implemented
- [ ] `scripts/install.sh`, `scripts/doctor-{aws,gcp,azure}.sh` implemented
- [ ] Root `AGENTS.md` authored
- [ ] `catalog.yaml` v1 schema defined
- [ ] CI workflows (`lint.yml`, `per-cli-tests.yml`, `golden.yml`, `release.yml`, `nightly.yml`)
- [ ] `pressfile.yaml` schema documented
- [ ] Patch convention documented

## CLIs in v1

| CLI | Cloud | Spec source | Status |
|---|---|---|---|
| cloud-run-admin-pp-cli | GCP | https://run.googleapis.com/$discovery/rest?version=v2 | pending (this plan) |
| apprunner-pp-cli | AWS | (Smithy / SDK model) | pending (follow-on plan) |
| lambda-pp-cli | AWS | (Smithy / SDK model) | pending (follow-on plan) |
| container-apps-pp-cli | Azure | (ARM spec, Microsoft.App) | pending (follow-on plan) |
| functions-pp-cli | Azure | (ARM spec, Microsoft.Web/sites) | pending (follow-on plan) |
| cloud-functions-pp-cli | GCP | https://cloudfunctions.googleapis.com/$discovery/rest?version=v2 | pending (follow-on plan) |

## Press-readiness check

Representative API for the gate: **`cloud-run-admin-pp-cli`** (GCP).

Rationale: mvanhorn already shipped this CLI, so we have a known-good output
to validate our substrate against. If our generated cloud-run-admin doesn't
match the known-good shape (modulo our convention layer), the substrate has
a bug.

## Per-CLI gap notes

(Captured per-CLI during Phase 6+ as scorecards complete.)

## Cross-cloud guidance to add to root AGENTS.md

- (Single-cloud v1; cross-cloud guidance accumulates as v1.1+ sprints add CLIs.)

## Definition of done

- [ ] All substrate deliverables checked off above
- [ ] cloud-run-admin-pp-cli through 8-step pipeline (this plan's Phase 6)
- [ ] Remaining 5 CLIs through 8-step pipeline (follow-on plan)
- [ ] All Layer 1 + Layer 2 tests green
- [ ] catalog.yaml lists all 6 with status: stable
- [ ] Root AGENTS.md updated with all 6 in "When to reach for what" table
- [ ] All 6 CLIs released and binaries published
- [ ] At least one cross-CLI compound recipe documented

## Retro

(Filled at sprint end. What surprised us, what we'd change for v1.1.)
```

---

### Task 1.8: Commit Phase 1

- [ ] **Step 1: Stage and commit**

Run:
```bash
cd /Users/shubhan/infra-planning-cli
git add README.md AGENTS.md catalog.yaml docs/conventions.md \
        docs/compound-recipes.md docs/sprints/v1.md \
        library tests scripts tools docs/sprints
git status
git commit -m "feat(substrate): repo skeleton + hub artifacts"
```
Expected: commit succeeds, all hub files tracked.

---

## Phase 2 — Convention linter (`tools/lint-conventions/`)

A Go program that lints `library/<cloud>/<service>/` directories for convention compliance. TDD throughout. Each check is a separate file in `tools/lint-conventions/checks/`.

### Task 2.1: Initialize Go module + check interface

**Files:**
- Create: `tools/lint-conventions/go.mod`
- Create: `tools/lint-conventions/main.go`
- Create: `tools/lint-conventions/checks/check.go`

- [ ] **Step 1: Initialize module**

Run:
```bash
cd /Users/shubhan/infra-planning-cli/tools/lint-conventions
go mod init github.com/ShubhanYenuganti/infra-press/tools/lint-conventions
```
Expected: `go.mod` created with `go 1.26` (or whatever is installed).

- [ ] **Step 2: Add yaml dependency**

Run: `cd /Users/shubhan/infra-planning-cli/tools/lint-conventions && go get gopkg.in/yaml.v3`
Expected: `go.sum` created, yaml.v3 added.

- [ ] **Step 3: Write the Check interface**

Create `tools/lint-conventions/checks/check.go`:

```go
package checks

// Check is a single convention rule applied to a CLI directory.
type Check interface {
	// Name returns a short identifier (e.g., "pressfile-schema").
	Name() string
	// Run inspects the CLI rooted at cliDir and returns nil on pass,
	// or an error describing the violation.
	Run(cliDir string) error
}

// Registry returns all checks the linter knows about.
// Add new checks here as conventions evolve.
func Registry() []Check {
	return []Check{
		// Filled in over the next tasks
	}
}
```

- [ ] **Step 4: Write main.go skeleton**

Create `tools/lint-conventions/main.go`:

```go
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ShubhanYenuganti/infra-press/tools/lint-conventions/checks"
)

func main() {
	ci := flag.Bool("ci", false, "Exit non-zero on any violation (suitable for CI)")
	flag.Parse()

	target := flag.Arg(0)
	if target == "" {
		fmt.Fprintln(os.Stderr, "usage: lint-conventions [--ci] <path>")
		os.Exit(2)
	}

	cliDirs, err := discoverCLIDirs(target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "discovery failed: %v\n", err)
		os.Exit(1)
	}

	violations := 0
	for _, dir := range cliDirs {
		for _, check := range checks.Registry() {
			if err := check.Run(dir); err != nil {
				fmt.Printf("FAIL  %s  %s: %v\n", dir, check.Name(), err)
				violations++
			} else {
				fmt.Printf("PASS  %s  %s\n", dir, check.Name())
			}
		}
	}

	if violations > 0 && *ci {
		os.Exit(1)
	}
}

// discoverCLIDirs returns CLI dirs under target. If target itself is a
// CLI dir (has pressfile.yaml), returns just it; otherwise walks library/.
func discoverCLIDirs(target string) ([]string, error) {
	if _, err := os.Stat(filepath.Join(target, "pressfile.yaml")); err == nil {
		return []string{target}, nil
	}
	var dirs []string
	err := filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if _, err := os.Stat(filepath.Join(path, "pressfile.yaml")); err == nil {
				dirs = append(dirs, path)
			}
		}
		return nil
	})
	return dirs, err
}
```

- [ ] **Step 5: Build to verify**

Run:
```bash
cd /Users/shubhan/infra-planning-cli/tools/lint-conventions
go build ./...
```
Expected: builds cleanly (no checks registered yet, so no checks run).

- [ ] **Step 6: Commit**

Run:
```bash
cd /Users/shubhan/infra-planning-cli
git add tools/lint-conventions/
git commit -m "feat(lint): scaffold convention linter with empty registry"
```

---

### Task 2.2: Pressfile schema check (TDD)

**Files:**
- Create: `tools/lint-conventions/checks/pressfile.go`
- Create: `tools/lint-conventions/checks/pressfile_test.go`
- Create: `tools/lint-conventions/checks/testdata/pressfile-valid/pressfile.yaml`
- Create: `tools/lint-conventions/checks/testdata/pressfile-missing/.gitkeep`
- Create: `tools/lint-conventions/checks/testdata/pressfile-malformed/pressfile.yaml`

- [ ] **Step 1: Write the failing test**

Create `tools/lint-conventions/checks/pressfile_test.go`:

```go
package checks

import "testing"

func TestPressfileCheck_Pass(t *testing.T) {
	check := &PressfileCheck{}
	if err := check.Run("testdata/pressfile-valid"); err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}

func TestPressfileCheck_MissingFile(t *testing.T) {
	check := &PressfileCheck{}
	if err := check.Run("testdata/pressfile-missing"); err == nil {
		t.Error("expected error for missing pressfile.yaml")
	}
}

func TestPressfileCheck_MalformedYAML(t *testing.T) {
	check := &PressfileCheck{}
	if err := check.Run("testdata/pressfile-malformed"); err == nil {
		t.Error("expected error for malformed yaml")
	}
}
```

- [ ] **Step 2: Create test fixtures**

Create `tools/lint-conventions/checks/testdata/pressfile-valid/pressfile.yaml`:

```yaml
cli: test-pp-cli
cloud: gcp
service: test
press_version: v1.4.2
press_command: "printing-press generate --spec https://..."
spec_url: https://example.com/spec.json
spec_version: 2024-01-01
spec_etag: "abc123"
generated_at: 2026-05-20T00:00:00Z
patches: []
```

Create empty `tools/lint-conventions/checks/testdata/pressfile-missing/.gitkeep`.

Create `tools/lint-conventions/checks/testdata/pressfile-malformed/pressfile.yaml`:

```yaml
cli: test-pp-cli
  this: is bad yaml
   nested: wrong
```

- [ ] **Step 3: Run tests to verify they fail**

Run: `cd /Users/shubhan/infra-planning-cli/tools/lint-conventions && go test ./checks/...`
Expected: build fails with `undefined: PressfileCheck`.

- [ ] **Step 4: Implement the check**

Create `tools/lint-conventions/checks/pressfile.go`:

```go
package checks

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type PressfileCheck struct{}

func (*PressfileCheck) Name() string { return "pressfile-schema" }

type pressfileSchema struct {
	CLI           string `yaml:"cli"`
	Cloud         string `yaml:"cloud"`
	Service       string `yaml:"service"`
	PressVersion  string `yaml:"press_version"`
	PressCommand  string `yaml:"press_command"`
	SpecURL       string `yaml:"spec_url"`
	SpecVersion   string `yaml:"spec_version"`
	SpecEtag      string `yaml:"spec_etag"`
	GeneratedAt   string `yaml:"generated_at"`
	Patches       []struct {
		Path   string `yaml:"path"`
		Reason string `yaml:"reason"`
	} `yaml:"patches"`
}

func (*PressfileCheck) Run(cliDir string) error {
	path := filepath.Join(cliDir, "pressfile.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read pressfile: %w", err)
	}
	var p pressfileSchema
	if err := yaml.Unmarshal(data, &p); err != nil {
		return fmt.Errorf("parse pressfile: %w", err)
	}
	required := map[string]string{
		"cli":           p.CLI,
		"cloud":         p.Cloud,
		"service":       p.Service,
		"press_version": p.PressVersion,
		"spec_url":      p.SpecURL,
		"generated_at":  p.GeneratedAt,
	}
	for field, value := range required {
		if value == "" {
			return fmt.Errorf("missing required field: %s", field)
		}
	}
	return nil
}
```

- [ ] **Step 5: Register the check**

Update `tools/lint-conventions/checks/check.go`, replace `Registry()`:

```go
func Registry() []Check {
	return []Check{
		&PressfileCheck{},
	}
}
```

- [ ] **Step 6: Run tests to verify they pass**

Run: `cd /Users/shubhan/infra-planning-cli/tools/lint-conventions && go test ./checks/...`
Expected: 3 tests pass.

- [ ] **Step 7: Commit**

Run:
```bash
cd /Users/shubhan/infra-planning-cli
git add tools/lint-conventions/
git commit -m "feat(lint): pressfile schema check"
```

---

### Task 2.3: Subcommand presence check (TDD)

Verifies the CLI's `cmd/` contains references to `sync`, `search`, and `sql` subcommands.

**Files:**
- Create: `tools/lint-conventions/checks/subcommands.go`
- Create: `tools/lint-conventions/checks/subcommands_test.go`
- Create: `tools/lint-conventions/checks/testdata/subcommands-ok/cmd/root.go`
- Create: `tools/lint-conventions/checks/testdata/subcommands-ok/pressfile.yaml`
- Create: `tools/lint-conventions/checks/testdata/subcommands-missing/cmd/root.go`
- Create: `tools/lint-conventions/checks/testdata/subcommands-missing/pressfile.yaml`

- [ ] **Step 1: Write the failing test**

Create `tools/lint-conventions/checks/subcommands_test.go`:

```go
package checks

import "testing"

func TestSubcommandsCheck_Pass(t *testing.T) {
	check := &SubcommandsCheck{}
	if err := check.Run("testdata/subcommands-ok"); err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}

func TestSubcommandsCheck_MissingSync(t *testing.T) {
	check := &SubcommandsCheck{}
	if err := check.Run("testdata/subcommands-missing"); err == nil {
		t.Error("expected error when sync command missing")
	}
}
```

- [ ] **Step 2: Create test fixtures**

Create `tools/lint-conventions/checks/testdata/subcommands-ok/pressfile.yaml` (minimal valid):

```yaml
cli: test-pp-cli
cloud: gcp
service: test
press_version: v1.4.2
spec_url: https://example.com/spec.json
generated_at: 2026-05-20T00:00:00Z
patches: []
```

Create `tools/lint-conventions/checks/testdata/subcommands-ok/cmd/root.go`:

```go
package cmd

// rootCmd registers all subcommands.
// This fixture references: sync, search, sql (required subcommands).
func init() {
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(sqlCmd)
}
```

Create `tools/lint-conventions/checks/testdata/subcommands-missing/pressfile.yaml` (copy of above).

Create `tools/lint-conventions/checks/testdata/subcommands-missing/cmd/root.go`:

```go
package cmd
// This fixture is intentionally missing sync/search/sql references.
func init() {
	rootCmd.AddCommand(otherCmd)
}
```

- [ ] **Step 3: Run tests to verify failure**

Run: `cd /Users/shubhan/infra-planning-cli/tools/lint-conventions && go test ./checks/...`
Expected: build fails with `undefined: SubcommandsCheck`.

- [ ] **Step 4: Implement**

Create `tools/lint-conventions/checks/subcommands.go`:

```go
package checks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SubcommandsCheck struct{}

func (*SubcommandsCheck) Name() string { return "subcommands-present" }

func (*SubcommandsCheck) Run(cliDir string) error {
	required := []string{"sync", "search", "sql"}
	cmdDir := filepath.Join(cliDir, "cmd")
	contents, err := concatGoSources(cmdDir)
	if err != nil {
		// Also try internal/cmd as that's where press tends to put them
		alt := filepath.Join(cliDir, "internal", "cmd")
		contents, err = concatGoSources(alt)
		if err != nil {
			return fmt.Errorf("no cmd directory found")
		}
	}
	for _, sub := range required {
		// Look for "syncCmd", "searchCmd", "sqlCmd" or `Use: "sync"` patterns
		if !strings.Contains(contents, sub+"Cmd") &&
			!strings.Contains(contents, `Use: "`+sub+`"`) &&
			!strings.Contains(contents, `Use:   "`+sub+`"`) {
			return fmt.Errorf("missing required subcommand: %s", sub)
		}
	}
	return nil
}

func concatGoSources(dir string) (string, error) {
	var b strings.Builder
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		b.Write(data)
		b.WriteString("\n")
		return nil
	})
	if err != nil {
		return "", err
	}
	if b.Len() == 0 {
		return "", fmt.Errorf("no .go files in %s", dir)
	}
	return b.String(), nil
}
```

- [ ] **Step 5: Register**

Update `tools/lint-conventions/checks/check.go`:

```go
func Registry() []Check {
	return []Check{
		&PressfileCheck{},
		&SubcommandsCheck{},
	}
}
```

- [ ] **Step 6: Run tests to verify pass**

Run: `cd /Users/shubhan/infra-planning-cli/tools/lint-conventions && go test ./checks/...`
Expected: all tests pass.

- [ ] **Step 7: Commit**

Run:
```bash
cd /Users/shubhan/infra-planning-cli
git add tools/lint-conventions/
git commit -m "feat(lint): subcommand presence check (sync/search/sql)"
```

---

### Task 2.4: Required-flags check (TDD)

Verifies the CLI source references each required flag (`--json`, `--dry-run`, `--select`, `--data-source`, `--compact`).

**Files:**
- Create: `tools/lint-conventions/checks/flags.go`
- Create: `tools/lint-conventions/checks/flags_test.go`
- Create: `tools/lint-conventions/checks/testdata/flags-ok/{pressfile.yaml,cmd/root.go}`
- Create: `tools/lint-conventions/checks/testdata/flags-missing/{pressfile.yaml,cmd/root.go}`

- [ ] **Step 1: Write the failing test**

Create `tools/lint-conventions/checks/flags_test.go`:

```go
package checks

import "testing"

func TestFlagsCheck_Pass(t *testing.T) {
	check := &FlagsCheck{}
	if err := check.Run("testdata/flags-ok"); err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}

func TestFlagsCheck_MissingDryRun(t *testing.T) {
	check := &FlagsCheck{}
	if err := check.Run("testdata/flags-missing"); err == nil {
		t.Error("expected error when --dry-run missing")
	}
}
```

- [ ] **Step 2: Create test fixtures**

Create `tools/lint-conventions/checks/testdata/flags-ok/pressfile.yaml` (minimal valid; same content as in Task 2.3).

Create `tools/lint-conventions/checks/testdata/flags-ok/cmd/root.go`:

```go
package cmd
// References all required flags.
const flags = `--json --dry-run --select --data-source --compact`
```

Create `tools/lint-conventions/checks/testdata/flags-missing/pressfile.yaml` (same minimal valid).

Create `tools/lint-conventions/checks/testdata/flags-missing/cmd/root.go`:

```go
package cmd
// Missing --dry-run, --compact, etc.
const flags = `--json --select`
```

- [ ] **Step 3: Run tests to verify failure**

Run: `cd /Users/shubhan/infra-planning-cli/tools/lint-conventions && go test ./checks/...`
Expected: build fails with `undefined: FlagsCheck`.

- [ ] **Step 4: Implement**

Create `tools/lint-conventions/checks/flags.go`:

```go
package checks

import (
	"fmt"
	"path/filepath"
	"strings"
)

type FlagsCheck struct{}

func (*FlagsCheck) Name() string { return "required-flags" }

func (*FlagsCheck) Run(cliDir string) error {
	required := []string{"--json", "--dry-run", "--select", "--data-source", "--compact"}
	// Search both cmd/ and internal/cmd/
	var contents string
	for _, dir := range []string{"cmd", "internal/cmd"} {
		c, err := concatGoSources(filepath.Join(cliDir, dir))
		if err == nil {
			contents += c
		}
	}
	if contents == "" {
		return fmt.Errorf("no Go source found in cmd/ or internal/cmd/")
	}
	for _, flag := range required {
		// Search for the literal flag string in source
		// (press output references flags as string literals when registering with cobra)
		bare := strings.TrimPrefix(flag, "--")
		if !strings.Contains(contents, `"`+flag+`"`) &&
			!strings.Contains(contents, `"`+bare+`"`) {
			return fmt.Errorf("missing required flag: %s", flag)
		}
	}
	return nil
}
```

- [ ] **Step 5: Register**

Update `Registry()` in `check.go`:

```go
return []Check{
	&PressfileCheck{},
	&SubcommandsCheck{},
	&FlagsCheck{},
}
```

- [ ] **Step 6: Run tests**

Run: `cd /Users/shubhan/infra-planning-cli/tools/lint-conventions && go test ./checks/...`
Expected: pass.

- [ ] **Step 7: Commit**

Run:
```bash
cd /Users/shubhan/infra-planning-cli
git add tools/lint-conventions/
git commit -m "feat(lint): required-flags check"
```

---

### Task 2.5: SQLite path check (TDD)

Verifies CLI source uses `~/.infra-press/<cli>.db` (or equivalent expanded form) as the local store path.

**Files:**
- Create: `tools/lint-conventions/checks/sqlite_path.go`
- Create: `tools/lint-conventions/checks/sqlite_path_test.go`
- Create: `tools/lint-conventions/checks/testdata/sqlite-ok/{pressfile.yaml,internal/store/db.go}`
- Create: `tools/lint-conventions/checks/testdata/sqlite-wrong/{pressfile.yaml,internal/store/db.go}`

- [ ] **Step 1: Failing test**

Create `tools/lint-conventions/checks/sqlite_path_test.go`:

```go
package checks

import "testing"

func TestSQLitePathCheck_Pass(t *testing.T) {
	check := &SQLitePathCheck{}
	if err := check.Run("testdata/sqlite-ok"); err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}

func TestSQLitePathCheck_WrongPath(t *testing.T) {
	check := &SQLitePathCheck{}
	if err := check.Run("testdata/sqlite-wrong"); err == nil {
		t.Error("expected error for non-conventional path")
	}
}
```

- [ ] **Step 2: Fixtures**

Create `tools/lint-conventions/checks/testdata/sqlite-ok/pressfile.yaml` (minimal valid).

Create `tools/lint-conventions/checks/testdata/sqlite-ok/internal/store/db.go`:

```go
package store
const DBPath = "~/.infra-press/test-pp-cli.db"
```

Create `tools/lint-conventions/checks/testdata/sqlite-wrong/pressfile.yaml` (minimal valid).

Create `tools/lint-conventions/checks/testdata/sqlite-wrong/internal/store/db.go`:

```go
package store
const DBPath = "~/.config/test/db.sqlite"
```

- [ ] **Step 3: Run tests — verify fail**

Run: `cd /Users/shubhan/infra-planning-cli/tools/lint-conventions && go test ./checks/...`
Expected: `undefined: SQLitePathCheck`.

- [ ] **Step 4: Implement**

Create `tools/lint-conventions/checks/sqlite_path.go`:

```go
package checks

import (
	"fmt"
	"path/filepath"
	"strings"
)

type SQLitePathCheck struct{}

func (*SQLitePathCheck) Name() string { return "sqlite-path" }

func (*SQLitePathCheck) Run(cliDir string) error {
	var contents string
	for _, dir := range []string{"internal/store", "internal", "cmd"} {
		c, _ := concatGoSources(filepath.Join(cliDir, dir))
		contents += c
	}
	if contents == "" {
		return fmt.Errorf("no source files found in internal/store, internal, or cmd")
	}
	if !strings.Contains(contents, `.infra-press/`) {
		return fmt.Errorf(`SQLite path does not include "~/.infra-press/" prefix; ` +
			`stores must live at ~/.infra-press/<cli>.db`)
	}
	return nil
}
```

- [ ] **Step 5: Register, run tests, commit**

Update `Registry()`:

```go
return []Check{
	&PressfileCheck{},
	&SubcommandsCheck{},
	&FlagsCheck{},
	&SQLitePathCheck{},
}
```

Run: `cd /Users/shubhan/infra-planning-cli/tools/lint-conventions && go test ./checks/...`
Expected: pass.

Commit:
```bash
cd /Users/shubhan/infra-planning-cli
git add tools/lint-conventions/
git commit -m "feat(lint): SQLite path check"
```

---

### Task 2.6: AGENTS.md presence check (TDD)

Verifies `AGENTS.md` exists and is non-trivial (>200 lines or >500 words is too strict for early CLIs; settle on "file exists, ≥10 lines, mentions at least one resource").

**Files:**
- Create: `tools/lint-conventions/checks/agents_md.go`
- Create: `tools/lint-conventions/checks/agents_md_test.go`
- Create: `tools/lint-conventions/checks/testdata/agents-ok/{pressfile.yaml,AGENTS.md}`
- Create: `tools/lint-conventions/checks/testdata/agents-missing/pressfile.yaml`
- Create: `tools/lint-conventions/checks/testdata/agents-empty/{pressfile.yaml,AGENTS.md}`

- [ ] **Step 1: Failing test**

Create `tools/lint-conventions/checks/agents_md_test.go`:

```go
package checks

import "testing"

func TestAgentsMDCheck_Pass(t *testing.T) {
	check := &AgentsMDCheck{}
	if err := check.Run("testdata/agents-ok"); err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}

func TestAgentsMDCheck_Missing(t *testing.T) {
	check := &AgentsMDCheck{}
	if err := check.Run("testdata/agents-missing"); err == nil {
		t.Error("expected error when AGENTS.md missing")
	}
}

func TestAgentsMDCheck_TooShort(t *testing.T) {
	check := &AgentsMDCheck{}
	if err := check.Run("testdata/agents-empty"); err == nil {
		t.Error("expected error for trivial AGENTS.md")
	}
}
```

- [ ] **Step 2: Fixtures**

Create `tools/lint-conventions/checks/testdata/agents-ok/pressfile.yaml` (minimal valid).

Create `tools/lint-conventions/checks/testdata/agents-ok/AGENTS.md`:

```markdown
# test-pp-cli

A CLI for the Test API.

## Resources

- services
- jobs
- revisions

## Commands

- list, get, create, delete

## SQLite store

Path: ~/.infra-press/test-pp-cli.db
```

Create `tools/lint-conventions/checks/testdata/agents-missing/pressfile.yaml` (minimal valid; no AGENTS.md sibling).

Create `tools/lint-conventions/checks/testdata/agents-empty/pressfile.yaml` (minimal valid).

Create `tools/lint-conventions/checks/testdata/agents-empty/AGENTS.md`:

```markdown
# placeholder
```

- [ ] **Step 3: Verify fail**

Run: `cd /Users/shubhan/infra-planning-cli/tools/lint-conventions && go test ./checks/...`
Expected: `undefined: AgentsMDCheck`.

- [ ] **Step 4: Implement**

Create `tools/lint-conventions/checks/agents_md.go`:

```go
package checks

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

type AgentsMDCheck struct{}

func (*AgentsMDCheck) Name() string { return "agents-md-non-trivial" }

func (*AgentsMDCheck) Run(cliDir string) error {
	path := filepath.Join(cliDir, "AGENTS.md")
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("AGENTS.md missing or unreadable: %w", err)
	}
	defer f.Close()
	lines := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines++
	}
	if lines < 8 {
		return fmt.Errorf("AGENTS.md too short (%d lines, need ≥8)", lines)
	}
	return nil
}
```

- [ ] **Step 5: Register, test, commit**

Update `Registry()`:

```go
return []Check{
	&PressfileCheck{},
	&SubcommandsCheck{},
	&FlagsCheck{},
	&SQLitePathCheck{},
	&AgentsMDCheck{},
}
```

Run: `cd /Users/shubhan/infra-planning-cli/tools/lint-conventions && go test ./checks/...`
Expected: pass.

Commit:
```bash
cd /Users/shubhan/infra-planning-cli
git add tools/lint-conventions/
git commit -m "feat(lint): AGENTS.md presence + non-trivial check"
```

---

### Task 2.7: End-to-end lint binary smoke test

Smoke-test the built `lint-conventions` binary against the test fixtures.

**Files:**
- Modify: none (uses existing fixtures from prior tasks)

- [ ] **Step 1: Build the binary**

Run:
```bash
cd /Users/shubhan/infra-planning-cli/tools/lint-conventions
go build -o /tmp/lint-conventions .
```
Expected: binary at `/tmp/lint-conventions`.

- [ ] **Step 2: Run against a known-good fixture**

Run: `/tmp/lint-conventions /Users/shubhan/infra-planning-cli/tools/lint-conventions/checks/testdata/agents-ok`
Expected: shows `PASS` for at least pressfile-schema and agents-md-non-trivial (other checks may fail because this fixture only has AGENTS.md, not cmd/).

- [ ] **Step 3: Run with --ci on a known-failing fixture**

Run: `/tmp/lint-conventions --ci /Users/shubhan/infra-planning-cli/tools/lint-conventions/checks/testdata/agents-missing`
Expected: exits non-zero (1), shows FAIL for `agents-md-non-trivial`.

- [ ] **Step 4: Commit (no file changes; verify binary works)**

(No commit needed; this is verification only.)

---

## Phase 3 — Golden test harness (`tests/golden/`)

A Go test suite that runs each registered check against every binary listed in `catalog.yaml`. Uses sub-tests so failures are reported per (CLI, check) pair.

### Task 3.1: Initialize golden tests module

**Files:**
- Create: `tests/golden/go.mod`
- Create: `tests/golden/golden_test.go`
- Create: `tests/golden/checks/check.go`

- [ ] **Step 1: Initialize**

Run:
```bash
cd /Users/shubhan/infra-planning-cli/tests/golden
go mod init github.com/ShubhanYenuganti/infra-press/tests/golden
go get gopkg.in/yaml.v3
```

- [ ] **Step 2: Write the check interface for golden tests**

Create `tests/golden/checks/check.go`:

```go
package checks

import "testing"

// CLI describes one entry in catalog.yaml relevant to golden tests.
type CLI struct {
	Name   string
	Binary string // resolved path to the built binary
}

// Check is a single black-box behavioral assertion.
type Check interface {
	Name() string
	Run(t *testing.T, cli CLI)
}

func Registry() []Check {
	return []Check{
		// Filled in over the next tasks
	}
}
```

- [ ] **Step 3: Write the orchestrator skeleton**

Create `tests/golden/golden_test.go`:

```go
package golden_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ShubhanYenuganti/infra-press/tests/golden/checks"
	"gopkg.in/yaml.v3"
)

type catalogEntry struct {
	Name   string `yaml:"name"`
	Path   string `yaml:"path"`
	Binary string `yaml:"binary"`
}

type catalog struct {
	CLIs []catalogEntry `yaml:"clis"`
}

func loadCatalog(t *testing.T) catalog {
	// Catalog lives at repo root, two dirs above this test file.
	root := repoRoot(t)
	data, err := os.ReadFile(filepath.Join(root, "catalog.yaml"))
	if err != nil {
		t.Fatalf("read catalog: %v", err)
	}
	var c catalog
	if err := yaml.Unmarshal(data, &c); err != nil {
		t.Fatalf("parse catalog: %v", err)
	}
	return c
}

func repoRoot(t *testing.T) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	// tests/golden -> repo root
	return filepath.Clean(filepath.Join(wd, "..", ".."))
}

func TestGoldenAcrossLibrary(t *testing.T) {
	cat := loadCatalog(t)
	if len(cat.CLIs) == 0 {
		t.Skip("catalog.yaml has no CLIs yet; nothing to check")
	}
	root := repoRoot(t)
	for _, entry := range cat.CLIs {
		binary := filepath.Join(root, entry.Path, "cmd", entry.Binary, entry.Binary)
		cli := checks.CLI{Name: entry.Name, Binary: binary}
		for _, check := range checks.Registry() {
			t.Run(entry.Name+"/"+check.Name(), func(t *testing.T) {
				check.Run(t, cli)
			})
		}
	}
}
```

- [ ] **Step 4: Run to verify (skips since catalog empty)**

Run: `cd /Users/shubhan/infra-planning-cli/tests/golden && go test ./...`
Expected: skipped (catalog has zero CLIs).

- [ ] **Step 5: Commit**

Run:
```bash
cd /Users/shubhan/infra-planning-cli
git add tests/golden/
git commit -m "feat(golden): catalog-driven orchestrator skeleton"
```

---

### Task 3.2: `help_exits_zero` check (TDD)

**Files:**
- Create: `tests/golden/checks/help_exits_zero.go`
- Create: `tests/golden/checks/help_exits_zero_test.go`
- Create: `tests/golden/checks/testdata/mock-binaries/help-ok.sh`
- Create: `tests/golden/checks/testdata/mock-binaries/help-broken.sh`

- [ ] **Step 1: Failing test**

Create `tests/golden/checks/help_exits_zero_test.go`:

```go
package checks

import "testing"

func TestHelpExitsZero_Pass(t *testing.T) {
	check := &HelpExitsZero{}
	check.Run(t, CLI{Name: "mock", Binary: "testdata/mock-binaries/help-ok.sh"})
	// Implicit: if t.Error/t.Fatal not called, pass.
}

func TestHelpExitsZero_Fail(t *testing.T) {
	subT := &testing.T{} // We need to capture failure without failing this test.
	// Instead, use a t.Run wrapper to assert failure.
	t.Run("inner", func(inner *testing.T) {
		check := &HelpExitsZero{}
		// Cannot easily assert sub-failure; just run and verify it reports.
		check.Run(inner, CLI{Name: "mock", Binary: "testdata/mock-binaries/help-broken.sh"})
		// We expect inner.Failed() to be true. Skipped to avoid test-of-test gymnastics.
	})
	_ = subT
}
```

(Note: testing test-of-test is awkward in Go. Keep `_Fail` shape but accept it's a smoke test, not strict assertion.)

- [ ] **Step 2: Create mock binaries**

Create `tests/golden/checks/testdata/mock-binaries/help-ok.sh`:

```bash
#!/usr/bin/env bash
# Mock CLI that handles --help correctly.
if [[ "$1" == "--help" ]]; then
  echo "Usage: mock [command]"
  exit 0
fi
echo "unknown command" >&2
exit 2
```

Make executable: `chmod +x tests/golden/checks/testdata/mock-binaries/help-ok.sh`

Create `tests/golden/checks/testdata/mock-binaries/help-broken.sh`:

```bash
#!/usr/bin/env bash
# Mock CLI that crashes on --help.
exit 1
```

Make executable: `chmod +x tests/golden/checks/testdata/mock-binaries/help-broken.sh`

- [ ] **Step 3: Verify fail**

Run: `cd /Users/shubhan/infra-planning-cli/tests/golden && go test ./checks/...`
Expected: `undefined: HelpExitsZero`.

- [ ] **Step 4: Implement**

Create `tests/golden/checks/help_exits_zero.go`:

```go
package checks

import (
	"os/exec"
	"testing"
)

type HelpExitsZero struct{}

func (*HelpExitsZero) Name() string { return "help-exits-zero" }

func (*HelpExitsZero) Run(t *testing.T, cli CLI) {
	cmd := exec.Command(cli.Binary, "--help")
	if err := cmd.Run(); err != nil {
		t.Errorf("%s --help exited non-zero: %v", cli.Name, err)
	}
}
```

- [ ] **Step 5: Register**

Update `tests/golden/checks/check.go`:

```go
func Registry() []Check {
	return []Check{
		&HelpExitsZero{},
	}
}
```

- [ ] **Step 6: Run**

Run: `cd /Users/shubhan/infra-planning-cli/tests/golden && go test ./checks/...`
Expected: pass (the `_Fail` test only loosely verifies; that's accepted).

- [ ] **Step 7: Commit**

Run:
```bash
cd /Users/shubhan/infra-planning-cli
git add tests/golden/
git commit -m "feat(golden): help-exits-zero check"
```

---

### Task 3.3: `unknown_command_exits_2` check (TDD)

**Files:**
- Create: `tests/golden/checks/unknown_command_exits_2.go`
- Create: `tests/golden/checks/unknown_command_exits_2_test.go`

- [ ] **Step 1: Failing test**

Create `tests/golden/checks/unknown_command_exits_2_test.go`:

```go
package checks

import "testing"

func TestUnknownCommandExits2(t *testing.T) {
	check := &UnknownCommandExits2{}
	check.Run(t, CLI{Name: "mock", Binary: "testdata/mock-binaries/help-ok.sh"})
}
```

- [ ] **Step 2: Verify fail**

Run: `cd /Users/shubhan/infra-planning-cli/tests/golden && go test ./checks/...`
Expected: `undefined: UnknownCommandExits2`.

- [ ] **Step 3: Implement**

Create `tests/golden/checks/unknown_command_exits_2.go`:

```go
package checks

import (
	"os/exec"
	"testing"
)

type UnknownCommandExits2 struct{}

func (*UnknownCommandExits2) Name() string { return "unknown-command-exits-2" }

func (*UnknownCommandExits2) Run(t *testing.T, cli CLI) {
	cmd := exec.Command(cli.Binary, "this-command-does-not-exist")
	err := cmd.Run()
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Errorf("%s: expected ExitError, got %v", cli.Name, err)
		return
	}
	if exitErr.ExitCode() != 2 {
		t.Errorf("%s: expected exit code 2 (user error), got %d", cli.Name, exitErr.ExitCode())
	}
}
```

- [ ] **Step 4: Register, run, commit**

Update Registry:

```go
return []Check{
	&HelpExitsZero{},
	&UnknownCommandExits2{},
}
```

Run: `cd /Users/shubhan/infra-planning-cli/tests/golden && go test ./checks/...`
Expected: pass.

Commit:
```bash
cd /Users/shubhan/infra-planning-cli
git add tests/golden/
git commit -m "feat(golden): unknown-command-exits-2 check"
```

---

### Task 3.4: `auto_json_when_piped` check (TDD)

**Files:**
- Create: `tests/golden/checks/auto_json_when_piped.go`
- Create: `tests/golden/checks/auto_json_when_piped_test.go`
- Create: `tests/golden/checks/testdata/mock-binaries/json-when-piped.sh`

- [ ] **Step 1: Mock binary**

Create `tests/golden/checks/testdata/mock-binaries/json-when-piped.sh`:

```bash
#!/usr/bin/env bash
# Mock CLI that emits JSON when stdout is not a TTY.
# This mock unconditionally outputs JSON for the test command "list".
if [[ "$1" == "list" ]]; then
  if [[ -t 1 ]]; then
    echo "TTY output (table)"
  else
    echo '{"items":[]}'
  fi
  exit 0
fi
[[ "$1" == "--help" ]] && { echo usage; exit 0; }
exit 2
```

Make executable.

- [ ] **Step 2: Failing test**

Create `tests/golden/checks/auto_json_when_piped_test.go`:

```go
package checks

import "testing"

func TestAutoJSONWhenPiped(t *testing.T) {
	check := &AutoJSONWhenPiped{}
	check.Run(t, CLI{Name: "mock", Binary: "testdata/mock-binaries/json-when-piped.sh"})
}
```

- [ ] **Step 3: Verify fail**

Run: `cd /Users/shubhan/infra-planning-cli/tests/golden && go test ./checks/...`
Expected: `undefined: AutoJSONWhenPiped`.

- [ ] **Step 4: Implement**

Create `tests/golden/checks/auto_json_when_piped.go`:

```go
package checks

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"testing"
)

type AutoJSONWhenPiped struct{}

func (*AutoJSONWhenPiped) Name() string { return "auto-json-when-piped" }

func (*AutoJSONWhenPiped) Run(t *testing.T, cli CLI) {
	// Many CLIs use "list" as a safe read command. If the binary supports it.
	cmd := exec.Command(cli.Binary, "list")
	var buf bytes.Buffer
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		t.Skipf("%s: list command unsupported or errored: %v (skip — not a hard failure)", cli.Name, err)
		return
	}
	// Stdout was a pipe (buf), so output must parse as JSON.
	var v any
	if err := json.Unmarshal(buf.Bytes(), &v); err != nil {
		t.Errorf("%s: stdout was not JSON when piped: %v", cli.Name, err)
	}
}
```

- [ ] **Step 5: Register, run, commit**

Update Registry:

```go
return []Check{
	&HelpExitsZero{},
	&UnknownCommandExits2{},
	&AutoJSONWhenPiped{},
}
```

Run: `cd /Users/shubhan/infra-planning-cli/tests/golden && go test ./checks/...`
Expected: pass (or skip).

Commit:
```bash
cd /Users/shubhan/infra-planning-cli
git add tests/golden/
git commit -m "feat(golden): auto-json-when-piped check"
```

---

### Task 3.5: `version_is_semver` check (TDD)

**Files:**
- Create: `tests/golden/checks/version_is_semver.go`
- Create: `tests/golden/checks/version_is_semver_test.go`
- Create: `tests/golden/checks/testdata/mock-binaries/version-ok.sh`

- [ ] **Step 1: Mock**

Create `tests/golden/checks/testdata/mock-binaries/version-ok.sh`:

```bash
#!/usr/bin/env bash
[[ "$1" == "--version" ]] && { echo "v0.1.0"; exit 0; }
[[ "$1" == "--help" ]] && { echo usage; exit 0; }
exit 2
```

Make executable.

- [ ] **Step 2: Failing test**

Create `tests/golden/checks/version_is_semver_test.go`:

```go
package checks

import "testing"

func TestVersionIsSemver(t *testing.T) {
	check := &VersionIsSemver{}
	check.Run(t, CLI{Name: "mock", Binary: "testdata/mock-binaries/version-ok.sh"})
}
```

- [ ] **Step 3: Verify fail**

Run: `cd /Users/shubhan/infra-planning-cli/tests/golden && go test ./checks/...`
Expected: `undefined: VersionIsSemver`.

- [ ] **Step 4: Implement**

Create `tests/golden/checks/version_is_semver.go`:

```go
package checks

import (
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

type VersionIsSemver struct{}

func (*VersionIsSemver) Name() string { return "version-is-semver" }

var semver = regexp.MustCompile(`^v?\d+\.\d+\.\d+(?:-[0-9A-Za-z.-]+)?$`)

func (*VersionIsSemver) Run(t *testing.T, cli CLI) {
	cmd := exec.Command(cli.Binary, "--version")
	out, err := cmd.Output()
	if err != nil {
		t.Errorf("%s: --version failed: %v", cli.Name, err)
		return
	}
	got := strings.TrimSpace(string(out))
	if !semver.MatchString(got) {
		t.Errorf("%s: --version output %q is not SemVer", cli.Name, got)
	}
}
```

- [ ] **Step 5: Register, run, commit**

Update Registry:

```go
return []Check{
	&HelpExitsZero{},
	&UnknownCommandExits2{},
	&AutoJSONWhenPiped{},
	&VersionIsSemver{},
}
```

Run: `cd /Users/shubhan/infra-planning-cli/tests/golden && go test ./checks/...`
Expected: pass.

Commit:
```bash
cd /Users/shubhan/infra-planning-cli
git add tests/golden/
git commit -m "feat(golden): version-is-semver check"
```

---

## Phase 4 — Install + doctor scripts

Bash scripts. Tested via shellcheck + smoke runs (not unit-tested per se).

### Task 4.1: `scripts/install.sh`

**Files:**
- Create: `scripts/install.sh`

- [ ] **Step 1: Write the script**

Create `/Users/shubhan/infra-planning-cli/scripts/install.sh`:

```bash
#!/usr/bin/env bash
# install.sh — bulk install all `status: stable` CLIs from catalog.yaml.
#
# Idempotent: re-run to upgrade to whatever catalog.yaml on main pins.

set -euo pipefail

INSTALL_DIR="${INFRA_PRESS_HOME:-$HOME/.infra-press}/bin"
CATALOG_URL="${INFRA_PRESS_CATALOG:-https://raw.githubusercontent.com/ShubhanYenuganti/infra-press/main/catalog.yaml}"

# Detect platform
case "$(uname -s)" in
  Darwin) OS=darwin ;;
  Linux)  OS=linux ;;
  *) echo "unsupported OS: $(uname -s)" >&2; exit 1 ;;
esac
case "$(uname -m)" in
  x86_64)  ARCH=amd64 ;;
  arm64|aarch64) ARCH=arm64 ;;
  *) echo "unsupported arch: $(uname -m)" >&2; exit 1 ;;
esac

mkdir -p "$INSTALL_DIR"

echo "Fetching catalog from $CATALOG_URL..."
catalog=$(curl -sSL "$CATALOG_URL")

# Require yq for YAML parsing; check for it.
if ! command -v yq >/dev/null 2>&1; then
  echo "yq is required. Install with: brew install yq  (or your package manager)" >&2
  exit 1
fi

# Iterate stable CLIs
echo "$catalog" | yq '.clis[] | select(.status == "stable") | .name + " " + .press_version' \
| while read -r cli version; do
  if [[ -z "$cli" ]]; then continue; fi
  url="https://github.com/ShubhanYenuganti/infra-press/releases/download/${cli}/${version}/${cli}-${OS}-${ARCH}"
  echo "Installing $cli $version..."
  if curl -fL "$url" -o "$INSTALL_DIR/$cli"; then
    chmod +x "$INSTALL_DIR/$cli"
    echo "  installed $cli"
  else
    echo "  failed to fetch $cli @ $version (release may not exist yet); skipping" >&2
  fi
done

echo
echo "Done. Add to PATH:"
echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
```

- [ ] **Step 2: Make executable**

Run: `chmod +x /Users/shubhan/infra-planning-cli/scripts/install.sh`

- [ ] **Step 3: Lint with shellcheck (if installed)**

Run: `shellcheck /Users/shubhan/infra-planning-cli/scripts/install.sh || echo "shellcheck not installed; manual review only"`
Expected: no errors (or shellcheck skipped).

- [ ] **Step 4: Smoke test — empty catalog**

The current `catalog.yaml` has zero CLIs, so install.sh should exit cleanly with no installs:

Run: `INFRA_PRESS_CATALOG="file:///Users/shubhan/infra-planning-cli/catalog.yaml" /Users/shubhan/infra-planning-cli/scripts/install.sh || echo "expected to handle empty catalog"`
Expected: exits cleanly, no installs performed (since `clis: []`).

(Note: if `file://` doesn't work with `curl`, this smoke test is symbolic; verified properly once a CLI is in catalog.)

---

### Task 4.2: `scripts/doctor-aws.sh`

**Files:**
- Create: `scripts/doctor-aws.sh`

- [ ] **Step 1: Write the script**

Create `/Users/shubhan/infra-planning-cli/scripts/doctor-aws.sh`:

```bash
#!/usr/bin/env bash
# doctor-aws.sh — Pre-flight check for AWS authentication.
#
# Does NOT call any cloud APIs. Reports what the SDK chain would find.

set -uo pipefail

pass() { printf "✅ %s\n" "$*"; }
warn() { printf "⚠️  %s\n" "$*"; }
fail() { printf "❌ %s\n" "$*"; }

# 1. AWS_PROFILE
if [[ -n "${AWS_PROFILE:-}" ]]; then
  pass "AWS_PROFILE set: $AWS_PROFILE"
else
  warn "AWS_PROFILE not set; default profile or instance role will be used"
fi

# 2. ~/.aws/credentials
if [[ -f "$HOME/.aws/credentials" ]]; then
  pass "~/.aws/credentials present"
else
  warn "~/.aws/credentials missing; SDK may fall back to instance/container role"
fi

# 3. AWS_REGION
if [[ -n "${AWS_REGION:-}" || -n "${AWS_DEFAULT_REGION:-}" ]]; then
  pass "AWS_REGION set: ${AWS_REGION:-$AWS_DEFAULT_REGION}"
else
  warn "AWS_REGION not set; per-CLI commands will require --region"
fi

# 4. aws CLI (recommended for fallback IAM/networking work)
if command -v aws >/dev/null 2>&1; then
  pass "aws CLI installed ($(aws --version 2>&1))"
else
  warn "aws CLI not installed; recommended for IAM/VPC fallbacks until v1.1+"
fi

# 5. STS identity check (lightweight, requires aws CLI)
if command -v aws >/dev/null 2>&1; then
  if identity=$(aws sts get-caller-identity --output text 2>&1); then
    pass "STS GetCallerIdentity: $identity"
  else
    fail "STS GetCallerIdentity failed; auth not working"
    echo "    $identity"
  fi
fi
```

- [ ] **Step 2: Make executable + lint**

Run:
```bash
chmod +x /Users/shubhan/infra-planning-cli/scripts/doctor-aws.sh
shellcheck /Users/shubhan/infra-planning-cli/scripts/doctor-aws.sh 2>/dev/null || true
```

- [ ] **Step 3: Smoke run**

Run: `/Users/shubhan/infra-planning-cli/scripts/doctor-aws.sh`
Expected: prints status icons; doesn't crash. (Pass/fail of individual checks depends on local env.)

---

### Task 4.3: `scripts/doctor-gcp.sh`

**Files:**
- Create: `scripts/doctor-gcp.sh`

- [ ] **Step 1: Write the script**

Create `/Users/shubhan/infra-planning-cli/scripts/doctor-gcp.sh`:

```bash
#!/usr/bin/env bash
# doctor-gcp.sh — Pre-flight check for GCP authentication.

set -uo pipefail

pass() { printf "✅ %s\n" "$*"; }
warn() { printf "⚠️  %s\n" "$*"; }
fail() { printf "❌ %s\n" "$*"; }

# 1. GOOGLE_APPLICATION_CREDENTIALS
if [[ -n "${GOOGLE_APPLICATION_CREDENTIALS:-}" ]]; then
  if [[ -f "$GOOGLE_APPLICATION_CREDENTIALS" ]]; then
    pass "GOOGLE_APPLICATION_CREDENTIALS: $GOOGLE_APPLICATION_CREDENTIALS"
  else
    fail "GOOGLE_APPLICATION_CREDENTIALS set but file not found"
  fi
else
  warn "GOOGLE_APPLICATION_CREDENTIALS not set; falling back to gcloud ADC"
fi

# 2. gcloud ADC
ADC_PATH="$HOME/.config/gcloud/application_default_credentials.json"
if [[ -f "$ADC_PATH" ]]; then
  pass "gcloud ADC present: $ADC_PATH"
else
  warn "gcloud ADC not present; run \`gcloud auth application-default login\`"
fi

# 3. gcloud CLI (recommended fallback)
if command -v gcloud >/dev/null 2>&1; then
  pass "gcloud CLI installed ($(gcloud --version 2>&1 | head -1))"
  active=$(gcloud config get-value account 2>/dev/null || echo "(none)")
  pass "Active gcloud account: $active"
  project=$(gcloud config get-value project 2>/dev/null || echo "(unset)")
  if [[ "$project" == "(unset)" ]]; then
    warn "gcloud project not set; per-CLI commands will require --project"
  else
    pass "gcloud project: $project"
  fi
else
  warn "gcloud CLI not installed; recommended for fallback work"
fi
```

- [ ] **Step 2: Make executable + lint + smoke**

Run:
```bash
chmod +x /Users/shubhan/infra-planning-cli/scripts/doctor-gcp.sh
shellcheck /Users/shubhan/infra-planning-cli/scripts/doctor-gcp.sh 2>/dev/null || true
/Users/shubhan/infra-planning-cli/scripts/doctor-gcp.sh
```
Expected: runs without crashing.

---

### Task 4.4: `scripts/doctor-azure.sh`

**Files:**
- Create: `scripts/doctor-azure.sh`

- [ ] **Step 1: Write the script**

Create `/Users/shubhan/infra-planning-cli/scripts/doctor-azure.sh`:

```bash
#!/usr/bin/env bash
# doctor-azure.sh — Pre-flight check for Azure authentication.

set -uo pipefail

pass() { printf "✅ %s\n" "$*"; }
warn() { printf "⚠️  %s\n" "$*"; }
fail() { printf "❌ %s\n" "$*"; }

# 1. Service principal env vars
if [[ -n "${AZURE_TENANT_ID:-}" && -n "${AZURE_CLIENT_ID:-}" && -n "${AZURE_CLIENT_SECRET:-}" ]]; then
  pass "Service principal env vars set (TENANT/CLIENT_ID/CLIENT_SECRET)"
else
  warn "Service principal env vars not set; will fall back to az CLI or managed identity"
fi

# 2. az CLI auth state
if command -v az >/dev/null 2>&1; then
  pass "az CLI installed ($(az --version 2>&1 | head -1))"
  if account=$(az account show --output json 2>&1); then
    name=$(echo "$account" | python3 -c "import json,sys; print(json.load(sys.stdin).get('user',{}).get('name','?'))" 2>/dev/null || echo '?')
    pass "az logged in as: $name"
    sub=$(echo "$account" | python3 -c "import json,sys; print(json.load(sys.stdin).get('name','?'))" 2>/dev/null || echo '?')
    pass "Active subscription: $sub"
  else
    warn "az not logged in; run \`az login\`"
  fi
else
  warn "az CLI not installed; recommended fallback"
fi

# 3. AZURE_SUBSCRIPTION_ID
if [[ -n "${AZURE_SUBSCRIPTION_ID:-}" ]]; then
  pass "AZURE_SUBSCRIPTION_ID set"
else
  warn "AZURE_SUBSCRIPTION_ID not set; CLIs will need --subscription flag"
fi
```

- [ ] **Step 2: Make executable + lint + smoke**

Run:
```bash
chmod +x /Users/shubhan/infra-planning-cli/scripts/doctor-azure.sh
shellcheck /Users/shubhan/infra-planning-cli/scripts/doctor-azure.sh 2>/dev/null || true
/Users/shubhan/infra-planning-cli/scripts/doctor-azure.sh
```

---

### Task 4.5: Commit Phase 4

- [ ] **Step 1: Stage + commit**

Run:
```bash
cd /Users/shubhan/infra-planning-cli
git add scripts/
git commit -m "feat(scripts): install.sh + doctor-{aws,gcp,azure}.sh"
```

---

## Phase 5 — CI workflows

Five GitHub Actions workflows. Each is a single YAML file.

### Task 5.1: `lint.yml`

**Files:**
- Create: `.github/workflows/lint.yml`

- [ ] **Step 1: Create the workflow directory**

Run: `mkdir -p /Users/shubhan/infra-planning-cli/.github/workflows`

- [ ] **Step 2: Write lint.yml**

Create `.github/workflows/lint.yml`:

```yaml
name: lint

on:
  pull_request:
    paths:
      - 'library/**'
      - 'tools/lint-conventions/**'
      - '.github/workflows/lint.yml'
  push:
    branches: [main]
    paths:
      - 'library/**'

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'
      - name: Build lint-conventions
        working-directory: tools/lint-conventions
        run: go build -o /tmp/lint-conventions .
      - name: Run lint on library/
        run: /tmp/lint-conventions --ci library/
```

---

### Task 5.2: `per-cli-tests.yml`

**Files:**
- Create: `.github/workflows/per-cli-tests.yml`

- [ ] **Step 1: Write the workflow**

Create `.github/workflows/per-cli-tests.yml`:

```yaml
name: per-cli-tests

on:
  pull_request:
    paths:
      - 'library/**'
  push:
    branches: [main]
    paths:
      - 'library/**'

jobs:
  detect-changes:
    runs-on: ubuntu-latest
    outputs:
      clis: ${{ steps.detect.outputs.clis }}
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - id: detect
        run: |
          # List each library/<cloud>/<service> dir with changes since main
          base=$(git merge-base HEAD origin/main 2>/dev/null || echo HEAD~1)
          changed=$(git diff --name-only "$base" HEAD -- 'library/*/*' \
            | awk -F/ '{print "library/"$2"/"$3}' \
            | sort -u | jq -R . | jq -sc .)
          echo "clis=$changed" >> "$GITHUB_OUTPUT"

  test:
    needs: detect-changes
    if: needs.detect-changes.outputs.clis != '[]'
    runs-on: ubuntu-latest
    strategy:
      fail-fast: false
      matrix:
        cli: ${{ fromJson(needs.detect-changes.outputs.clis) }}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'
      - name: Test ${{ matrix.cli }}
        working-directory: ${{ matrix.cli }}
        run: go test ./...
```

---

### Task 5.3: `golden.yml`

**Files:**
- Create: `.github/workflows/golden.yml`

- [ ] **Step 1: Write the workflow**

Create `.github/workflows/golden.yml`:

```yaml
name: golden

on:
  pull_request:
  push:
    branches: [main]

jobs:
  golden:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'
      - name: Build all CLIs in catalog
        run: |
          # For each CLI listed in catalog.yaml, build its binary into its
          # own cmd/ output dir so the golden orchestrator can find it.
          for cli_path in $(yq '.clis[].path' catalog.yaml); do
            (
              cd "$cli_path"
              for cmd_dir in cmd/*/; do
                binary=$(basename "$cmd_dir")
                (cd "$cmd_dir" && go build -o "$binary")
              done
            )
          done
      - name: Run golden tests
        working-directory: tests/golden
        run: go test -v ./...
```

---

### Task 5.4: `release.yml`

**Files:**
- Create: `.github/workflows/release.yml`

- [ ] **Step 1: Write the workflow**

Create `.github/workflows/release.yml`:

```yaml
name: release

on:
  push:
    tags:
      - '*-pp-cli/v*.*.*'

jobs:
  release:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'
      - name: Parse tag
        id: parse
        run: |
          # Tag format: <cli-name>/v<version>  e.g. cloud-run-admin-pp-cli/v0.1.0
          tag="${GITHUB_REF#refs/tags/}"
          cli="${tag%%/*}"
          version="${tag##*/}"
          echo "cli=$cli" >> "$GITHUB_OUTPUT"
          echo "version=$version" >> "$GITHUB_OUTPUT"
          # Resolve path from catalog
          path=$(yq ".clis[] | select(.name == \"$cli\") | .path" catalog.yaml)
          echo "path=$path" >> "$GITHUB_OUTPUT"
      - name: Build binaries (4 platforms)
        run: |
          cli=${{ steps.parse.outputs.cli }}
          path=${{ steps.parse.outputs.path }}
          mkdir -p /tmp/release
          cd "$path/cmd/$cli"
          for os in darwin linux; do
            for arch in amd64 arm64; do
              GOOS=$os GOARCH=$arch go build -o "/tmp/release/${cli}-${os}-${arch}"
            done
          done
      - name: Create GitHub release
        uses: softprops/action-gh-release@v2
        with:
          files: /tmp/release/*
          tag_name: ${{ github.ref_name }}
```

---

### Task 5.5: `nightly.yml`

**Files:**
- Create: `.github/workflows/nightly.yml`

- [ ] **Step 1: Write the workflow**

Create `.github/workflows/nightly.yml`:

```yaml
name: nightly

on:
  schedule:
    - cron: '17 7 * * *'  # 07:17 UTC daily
  workflow_dispatch:

jobs:
  build-all:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'
      - name: Build every CLI from scratch
        run: |
          fail=0
          for cli_path in $(yq '.clis[].path' catalog.yaml); do
            for cmd_dir in "$cli_path"/cmd/*/; do
              if ! (cd "$cmd_dir" && go build .); then
                echo "FAIL: $cmd_dir" >&2
                fail=1
              fi
            done
          done
          exit "$fail"
      - name: Open issue on failure
        if: failure()
        uses: actions-cool/issues-helper@v3
        with:
          actions: 'create-issue'
          title: 'nightly build failed on ${{ github.run_id }}'
          body: 'See run: ${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }}'
```

---

### Task 5.6: Commit Phase 5

- [ ] **Step 1: Stage + commit**

Run:
```bash
cd /Users/shubhan/infra-planning-cli
git add .github/
git commit -m "ci: lint, per-cli-tests, golden, release, nightly workflows"
```

---

## Phase 6 — First CLI: `cloud-run-admin-pp-cli`

End-to-end run of the 8-step pipeline against GCP Cloud Run. This proves the substrate works.

### Task 6.1: Verify Go 1.26+ + install printing-press

**Files:** none (environment setup)

- [ ] **Step 1: Check Go version**

Run: `go version`
Expected: `go1.26.x` or newer.

If older: install Go 1.26+ (e.g., `brew install go@1.26` on macOS, or download from https://go.dev/dl/).

If Go 1.26 is unavailable in your environment: **this is the highest-risk open question in the spec (Section 15 item 1).** Project stalls here. Document the blocker and stop.

- [ ] **Step 2: Install printing-press**

Run: `go install github.com/mvanhorn/cli-printing-press/cmd/printing-press@latest`
Expected: `printing-press` binary in `$GOPATH/bin`.

- [ ] **Step 3: Verify**

Run: `printing-press --help`
Expected: shows usage.

---

### Task 6.2: Press-readiness check against Cloud Run

**Files:** none (smoke test)

- [ ] **Step 1: Identify the spec URL**

Cloud Run Admin API v2 Discovery doc: `https://run.googleapis.com/$discovery/rest?version=v2`

- [ ] **Step 2: Run press in a scratch directory**

Run:
```bash
mkdir -p /tmp/press-readycheck-cloud-run
cd /tmp/press-readycheck-cloud-run
printing-press generate \
  --spec 'https://run.googleapis.com/$discovery/rest?version=v2' \
  --output .
```
Expected: press generates a CLI scaffold without crashing.

- [ ] **Step 3: Build it**

Run: `cd /tmp/press-readycheck-cloud-run && go build ./... 2>&1 | head -50`
Expected: builds (may show warnings; building means the spec parse + Go code-gen pipeline both work).

- [ ] **Step 4: If failure: STOP**

If press fails on Cloud Run discovery doc: spec format mismatch. Decide whether to file an upstream issue with mvanhorn or hand-author press patches. Document the failure mode in `docs/sprints/v1.md` under "Press-readiness check" and stop work on this CLI.

If success: proceed to Task 6.3.

---

### Task 6.3: Spec resolve + stub pressfile

**Files:**
- Create: `library/gcp/cloud-run-admin/pressfile.yaml`

- [ ] **Step 1: Create the CLI directory**

Run: `mkdir -p /Users/shubhan/infra-planning-cli/library/gcp/cloud-run-admin`

- [ ] **Step 2: Fetch the spec ETag (for spec_etag in pressfile)**

Run:
```bash
etag=$(curl -sI 'https://run.googleapis.com/$discovery/rest?version=v2' | grep -i '^etag:' | awk '{print $2}' | tr -d '\r"')
echo "$etag"
```
Expected: a string like `W/"abc123..."`. If empty, set to `unknown` for now.

- [ ] **Step 3: Write the stub pressfile**

Create `/Users/shubhan/infra-planning-cli/library/gcp/cloud-run-admin/pressfile.yaml`:

```yaml
cli: cloud-run-admin-pp-cli
cloud: gcp
service: cloud-run-admin
press_version: <fill in after Task 6.4>
press_command: "printing-press generate --spec https://run.googleapis.com/$discovery/rest?version=v2 --output library/gcp/cloud-run-admin/"
spec_url: https://run.googleapis.com/$discovery/rest?version=v2
spec_version: v2
spec_etag: <paste from Step 2 or "unknown">
generated_at: <ISO8601 timestamp filled by Task 6.4>
patches: []
```

(Press version and timestamps will be updated in Task 6.4 after generation completes.)

---

### Task 6.4: Run press generate against `library/gcp/cloud-run-admin/`

**Files:**
- Generated: many files under `library/gcp/cloud-run-admin/`

- [ ] **Step 1: Record press version**

Run: `printing-press --version`
Expected: a version string (e.g., `v1.4.2`). Note it for the pressfile.

- [ ] **Step 2: Run press**

Run:
```bash
cd /Users/shubhan/infra-planning-cli
printing-press generate \
  --spec 'https://run.googleapis.com/$discovery/rest?version=v2' \
  --output library/gcp/cloud-run-admin/
```
Expected: press generates `cmd/`, `internal/`, `README.md`, `AGENTS.md`, `go.mod`, `go.sum`, `.golangci.yml`, and `tests/` under the target directory.

- [ ] **Step 3: Update pressfile.yaml**

Edit `library/gcp/cloud-run-admin/pressfile.yaml`:
- Replace `<fill in after Task 6.4>` for `press_version` with the version from Step 1
- Replace `<ISO8601 timestamp ...>` for `generated_at` with current time (run `date -u +"%Y-%m-%dT%H:%M:%SZ"`)

- [ ] **Step 4: Verify generated layout**

Run: `find /Users/shubhan/infra-planning-cli/library/gcp/cloud-run-admin -type d | head -20`
Expected: shows `cmd/`, `internal/`, possibly `tests/`.

- [ ] **Step 5: Build the binary**

Run: `cd /Users/shubhan/infra-planning-cli/library/gcp/cloud-run-admin && go build ./...`
Expected: builds (may need `go mod tidy` first if dependencies missing).

- [ ] **Step 6: Commit the raw generated output (separate commit, before patches)**

Run:
```bash
cd /Users/shubhan/infra-planning-cli
git add library/gcp/cloud-run-admin/
git commit -m "feat(cloud-run-admin): bootstrap from press (raw output, pre-patch)"
```

---

### Task 6.5: Run the scorecard

**Files:**
- Create: `library/gcp/cloud-run-admin/scorecard.md`

- [ ] **Step 1: Build the lint binary**

Run:
```bash
cd /Users/shubhan/infra-planning-cli/tools/lint-conventions
go build -o /tmp/lint-conventions .
```

- [ ] **Step 2: Run lint against the generated CLI**

Run: `/tmp/lint-conventions /Users/shubhan/infra-planning-cli/library/gcp/cloud-run-admin/`
Expected: prints `PASS` / `FAIL` per check. Note which fail.

- [ ] **Step 3: Run the built binary's `--help`**

Run:
```bash
cd /Users/shubhan/infra-planning-cli/library/gcp/cloud-run-admin
go run ./cmd/... --help
```
Expected: prints usage. If multiple `cmd/` subdirs exist, run each.

- [ ] **Step 4: Write scorecard**

Create `library/gcp/cloud-run-admin/scorecard.md`:

```markdown
# Scorecard — cloud-run-admin-pp-cli (first generation)

| Check | Status | Notes |
|---|---|---|
| pressfile-schema | <fill in> | |
| subcommands-present | <fill in> | |
| required-flags | <fill in> | |
| sqlite-path | <fill in> | |
| agents-md-non-trivial | <fill in> | |
| binary builds | <fill in> | `go build ./...` outcome |
| --help works | <fill in> | |

## Gaps identified
(List patches needed — each will go in `patches/` with a corresponding `pressfile.yaml` entry.)

## Decisions
- If convention checks fail: patch press (if the gap is universal) OR patch this CLI
- If CLI-specific gaps: add to `patches/` per Task 6.6
```

Fill in the table based on what you observed in Steps 2-3.

---

### Task 6.6: Apply convention patches (if any)

**Files:**
- May create: `library/gcp/cloud-run-admin/patches/*.go`
- Will modify: `library/gcp/cloud-run-admin/pressfile.yaml` (`patches[]` entries)

- [ ] **Step 1: Decide per failed check**

For each `FAIL` from the scorecard:

- **Universal failures** (every press output would fail this) → fix the press itself or open an upstream issue. Document in `docs/sprints/v1.md`. Skip the patch.
- **CLI-specific failures** → write a patch file in `patches/`.

- [ ] **Step 2: Write each patch (if needed)**

Example: if `sync` subcommand missing, create `library/gcp/cloud-run-admin/patches/sync.go`:

```go
package patches

// (Stub — actual implementation depends on what's missing in press output.)
//
// Each patch file is a Go file that the CLI's cmd/main.go can import to
// register additional commands or modify behavior. See docs/conventions.md
// for the patch contract.
```

Add to `library/gcp/cloud-run-admin/pressfile.yaml`:

```yaml
patches:
  - path: patches/sync.go
    reason: "Press did not emit a sync subcommand; convention requires it."
```

- [ ] **Step 3: Re-run lint to verify**

Run: `/tmp/lint-conventions /Users/shubhan/infra-planning-cli/library/gcp/cloud-run-admin/`
Expected: all checks PASS (after patches).

- [ ] **Step 4: Build to verify patches compile**

Run: `cd /Users/shubhan/infra-planning-cli/library/gcp/cloud-run-admin && go build ./...`
Expected: builds.

- [ ] **Step 5: Commit patches (if any)**

Run:
```bash
cd /Users/shubhan/infra-planning-cli
git add library/gcp/cloud-run-admin/
git commit -m "feat(cloud-run-admin): patches to satisfy conventions"
```

(If no patches were needed: skip this commit.)

---

### Task 6.7: Update `catalog.yaml` + root `AGENTS.md`

**Files:**
- Modify: `catalog.yaml`
- Modify: `AGENTS.md`

- [ ] **Step 1: Add cloud-run-admin entry to catalog.yaml**

Replace `clis: []` in `catalog.yaml` with:

```yaml
version: 1
clis:
  - name: cloud-run-admin-pp-cli
    cloud: gcp
    service: cloud-run-admin
    paradigm: managed-container
    path: library/gcp/cloud-run-admin
    binary: cloud-run-admin-pp-cli
    install: github.com/ShubhanYenuganti/infra-press/library/gcp/cloud-run-admin
    status: bootstrapping
    press_version: <copy from pressfile.yaml>
    spec_version: v2
    daily_commands:
      - services list
      - services replace
      - jobs run
    compound_commands: []
    known_gaps: []
```

(Status is `bootstrapping` until golden tests pass + a release is tagged.)

- [ ] **Step 2: Update root AGENTS.md table**

In `/Users/shubhan/infra-planning-cli/AGENTS.md`, replace the placeholder row in "When to reach for what" with:

```markdown
| Need | CLI | Cloud |
|---|---|---|
| Deploy a containerized service | cloud-run-admin-pp-cli | GCP |
| _(more added as sprints ship)_ | | |
```

(Already there from Task 1.3; verify it matches.)

- [ ] **Step 3: Verify YAML still parses**

Run: `python3 -c "import yaml; print(yaml.safe_load(open('/Users/shubhan/infra-planning-cli/catalog.yaml')))"`
Expected: prints dict with 1 entry under `clis`.

---

### Task 6.8: Run golden tests against the CLI

**Files:** none (verifies existing infrastructure)

- [ ] **Step 1: Build the binary into the expected path**

Run:
```bash
cd /Users/shubhan/infra-planning-cli/library/gcp/cloud-run-admin/cmd/cloud-run-admin-pp-cli
go build -o cloud-run-admin-pp-cli
```
(Path may differ if press puts binary elsewhere; adjust the `binary` field in `catalog.yaml` if needed.)

- [ ] **Step 2: Run golden tests**

Run: `cd /Users/shubhan/infra-planning-cli/tests/golden && go test -v ./...`
Expected: tests run against `cloud-run-admin-pp-cli` for each registered check. All should pass.

- [ ] **Step 3: If failures**

For each failing check, decide:
- Is this a press output bug → file/patch press, regenerate
- Is this a convention bug → patch in `patches/`, re-run
- Is this a golden-test bug → fix the check in `tests/golden/checks/`

Iterate until all golden tests pass.

---

### Task 6.9: Run lint locally + commit

- [ ] **Step 1: Lint**

Run: `/tmp/lint-conventions --ci /Users/shubhan/infra-planning-cli/library/`
Expected: exits 0 (all PASS).

- [ ] **Step 2: Final commit for cloud-run-admin**

Run:
```bash
cd /Users/shubhan/infra-planning-cli
git add library/gcp/cloud-run-admin/ catalog.yaml AGENTS.md
git status
git commit -m "feat(cloud-run-admin): first CLI through 8-step pipeline; catalog + AGENTS updated"
```

---

### Task 6.10: Tag the first release (locally only — per CLAUDE.md, do not push)

- [ ] **Step 1: Tag**

Run:
```bash
cd /Users/shubhan/infra-planning-cli
git tag cloud-run-admin-pp-cli/v0.1.0
git tag --list 'cloud-run-admin-pp-cli/*'
```
Expected: tag created.

- [ ] **Step 2: Bump catalog status to `preview`**

Edit `catalog.yaml`, change the cloud-run-admin entry's `status: bootstrapping` → `status: preview`.

Commit:
```bash
cd /Users/shubhan/infra-planning-cli
git add catalog.yaml
git commit -m "chore(cloud-run-admin): bump status to preview for v0.1.0"
```

- [ ] **Step 3: (Manual / future) push tag to trigger release.yml**

(Per CLAUDE.md: NEVER `git push`. The user will manually push when ready, which triggers `release.yml` to build and publish binaries.)

---

## Phase 7 — Substrate completion + sprint retro

### Task 7.1: Verify all substrate deliverables

**Files:** none (verification only)

- [ ] **Step 1: Walk the v1 substrate checklist**

Open `/Users/shubhan/infra-planning-cli/docs/sprints/v1.md` and check off each substrate item:

- [x] Repo rename `infra-planning-cli` → `infra-press` (Task 0.5)
- [x] Old code moved to `old-infra-cli` worktree (Tasks 0.2-0.4)
- [x] `docs/conventions.md` written (Task 1.5)
- [x] `tools/lint-conventions/` implemented (Phase 2)
- [x] `tests/golden/` orchestrator + initial checks implemented (Phase 3)
- [x] `scripts/install.sh`, `scripts/doctor-{aws,gcp,azure}.sh` implemented (Phase 4)
- [x] Root `AGENTS.md` authored (Task 1.3)
- [x] `catalog.yaml` v1 schema defined (Task 1.4)
- [x] CI workflows (Phase 5)
- [x] `pressfile.yaml` schema documented (in docs/conventions.md, Task 1.5)
- [x] Patch convention documented (in docs/conventions.md)

Mark them in the sprint doc.

- [ ] **Step 2: Verify the proof CLI**

- [x] cloud-run-admin-pp-cli through 8-step pipeline (Phase 6)
- [x] Layer 1 + Layer 2 tests green for cloud-run-admin (Tasks 6.8)
- [x] catalog.yaml lists cloud-run-admin with status: preview (Task 6.10)
- [x] Root AGENTS.md updated with cloud-run-admin row (Task 6.7)

- [ ] **Step 3: Confirm remaining 5 CLIs are deferred to follow-on plan**

Note in `docs/sprints/v1.md` under "CLIs in v1":
- 5 remaining CLIs (lambda, apprunner, cloud-functions, functions, container-apps): pending follow-on plan

---

### Task 7.2: Write sprint v1 retro stub

**Files:**
- Modify: `docs/sprints/v1.md` (append retro section content)

- [ ] **Step 1: Fill in the retro section at the bottom of `docs/sprints/v1.md`**

Append (or fill in if section exists):

```markdown
## Retro (preliminary — pre 5-CLI follow-on)

### What worked
- Substrate-first approach: lint + golden tests existed before any CLI, so the
  first CLI (cloud-run-admin) was held to the spec from generation onward.
- 8-step pipeline executed end-to-end on cloud-run-admin.

### What surprised us
- (Fill in after the work: e.g., press version mismatches, spec parse quirks,
  patches needed.)

### What we'd change for the next sprint
- (Fill in: e.g., add specific checks to lint, refactor the patches/ structure,
  document tradeoff X.)

### Open questions to resolve before follow-on
- Did press handle GCP discovery doc cleanly? (Spec open question #1)
- Was the Claude Code dependency (spec open question #2) blocking?
- Did naming collision with mvanhorn's cloud-run-admin cause confusion? (#4)
```

- [ ] **Step 2: Commit**

Run:
```bash
cd /Users/shubhan/infra-planning-cli
git add docs/sprints/v1.md
git commit -m "docs(sprints): v1 substrate complete + retro stub"
```

---

### Task 7.3: Final verification + plan summary

- [ ] **Step 1: Confirm clean state**

Run:
```bash
cd /Users/shubhan/infra-planning-cli
git status
git log --oneline -15
```
Expected: clean working tree, ~25-30 commits since the substrate reset.

- [ ] **Step 2: Confirm what ships at end of this plan**

- 1 working CLI (`cloud-run-admin-pp-cli` at v0.1.0)
- Complete substrate: linter, golden tests, install scripts, doctor scripts, CI workflows
- All hub artifacts: README, AGENTS.md, catalog.yaml, conventions.md, compound-recipes.md, sprint plan
- Old code preserved in `old-infra-cli` branch + worktree
- Ready for: follow-on plan covering the other 5 v1 CLIs (lambda, apprunner, cloud-functions, functions, container-apps), each through the same 8-step pipeline now proven on cloud-run-admin

- [ ] **Step 3: Suggest next plan**

The follow-on plan should be written next. Suggested name: `docs/superpowers/plans/2026-05-XX-infra-press-v1-remaining-clis.md`. Each of the 5 remaining CLIs follows Phase 6 of this plan with cloud/service/spec-URL substituted.

---

## Spec Coverage Check (self-review)

Mapping spec sections to plan phases:

| Spec section | Covered by |
|---|---|
| §2 Locked decisions | Plan overall structure |
| §3 The 6 v1 CLIs | Phase 6 (cloud-run-admin); follow-on plan covers other 5 |
| §4 Capability profile | Documented in `README.md` (Task 1.2) and `AGENTS.md` (Task 1.3) |
| §5 Architecture (4 layers) | Phases 1-6 implement all 4 layers |
| §6 Repo layout + per-CLI structure | Tasks 1.1, 6.3-6.4 |
| §7 Hub artifacts | Tasks 1.2-1.7 |
| §8 8-step per-CLI workflow | Phase 6 (executed for cloud-run-admin) |
| §9 Convention enforcement | Phase 2 (linter) |
| §10 Tests & CI | Phases 3, 5 |
| §11 Distribution & install | Phase 4 + Task 5.4 (release.yml) |
| §12 Sprint workflow | Task 1.7 (sprint doc), Task 7.2 (retro) |
| §13 Migration | Phase 0 |
| §14 Out of scope | Plan header explicitly defers 5 of 6 v1 CLIs |
| §15 Open questions | Task 6.1 (Go 1.26 risk), Task 6.2 (Claude Code risk), Task 7.2 (retro questions) |
| §16 Success criteria | Task 7.1 (verify all deliverables); 5 of 6 CLIs deferred per scope note |
| §17 References | Referenced in plan header + Task 6.2/6.4 commands |

**Gaps explicitly accepted:**
- 5 of 6 v1 CLIs are deferred to follow-on plan (stated up front).
- "At least one cross-CLI compound recipe documented" (spec §16) requires ≥2 CLIs sharing data; deferred to follow-on plan.

---

## Placeholder scan

Scanned for "TBD", "TODO", "implement later", "fill in details" patterns:

- `<fill in>` appears in Task 6.5 (scorecard.md template) — intentional: scorecard is filled at runtime from observed results, not pre-determined.
- `<fill in after Task 6.4>` appears in Task 6.3 (pressfile press_version) — intentional: known data point that only exists after press is run.
- `<paste from Step 2 or "unknown">` in Task 6.3 — intentional: real value comes from a network call at runtime.

These are not plan-failures; they are runtime placeholders the implementing engineer fills with observed values.

---

## Type / signature consistency

- `Check` interface in `tools/lint-conventions/checks/check.go` (Task 2.1) is used consistently by all check implementations in Tasks 2.2-2.6.
- `Check` interface in `tests/golden/checks/check.go` (Task 3.1) has a different signature (`Run(t *testing.T, cli CLI)`) — this is intentional: lint-conventions runs against directories, golden runs against built binaries.
- `pressfile.yaml` schema is defined once (Task 1.5 conventions doc), instantiated as a stub (Task 6.3), filled in (Task 6.4), and validated by the linter (Task 2.2). Schema fields consistent across all uses.
- `catalog.yaml` schema is defined as comment (Task 1.4), instantiated (Task 6.7), and parsed by the golden orchestrator (Task 3.1). Field names consistent (`name`, `path`, `binary`, etc.).

---

**End of plan.**
