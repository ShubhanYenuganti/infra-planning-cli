# Design Improvements

This document augments docs `01`–`04` with two design decisions that close gaps against the original spec and strengthen the agent-facing positioning:

1. **Agent integration surface** — how the CLI exposes machine-readable data to consuming agents.
2. **Discovery breadth** — how many IaC formats the CLI inspects and how detectors are organized.

These decisions feed directly into a revised implementation plan (`05-implementation-plan.md`).

## Goals

- Give consuming agents (Claude Code, Codex) structured access to plan metadata without weakening the human readability of the markdown artifact.
- Honor spec 04's "still supports compute and databases in the combined flow" by dispatching to per-domain plan builders instead of silently routing every request to a networking plan.
- Surface evidence from a broader set of AWS-first IaC formats so the discovery yield on real repos matches the project's "AWS-first" positioning.
- Keep scope tight: web evidence fetching targets `docs.aws.amazon.com` only (no community fallback); cost hints use hardcoded seed values with disclaimer notes; conflict detection covers multi-IaC repos only. All three are conservative choices matching the "no credentials, no side effects" posture.

## Non-goals

- CDK source-file parsing beyond detecting `cdk.json`.
- Pulumi support (`Pulumi.yaml`).
- Resource-level domain inference (e.g., classifying `aws_vpc` vs `aws_lambda_function` from repo evidence to feed back into the classifier).

---

## 1. CLI Contract

```
infra-plan plan REQUEST [--repo PATH] [--out PATH] [--json]
```

| Mode | Behavior |
|---|---|
| Default | Writes one `.md` file with YAML frontmatter to `docs/infra-planning-cli/YYYY-MM-DD-<slug>.md` under `--repo` (or to `--out` if given). Exit 0. |
| `--json` | Emits a single JSON object to stdout. **No file written.** Exit 0. Mutually exclusive with `--out`. |
| Clarification needed (default) | Prints the question to stdout, no file written. Exit 0. |
| Clarification needed (`--json`) | Emits payload with `clarification_needed: true` and empty `plan_md`. Exit 0. |

**Exit codes**
- `0` — success (plan written, JSON emitted, or clarification surfaced)
- `1` — hard error (bad flags, unreadable repo path, write failure)

**Rationale.** `--json` matches `kubectl -o json`, `terraform show -json`, and `gh ... --json` conventions. Spec 01's "agent-facing copilot" framing is best served by JSON as a first-class mode rather than a side channel.

## 2. Output Schema

Frontmatter and JSON share the same shape. JSON adds one extra field, `plan_md`, containing the rendered markdown body.

**Frontmatter form (default mode)**

```yaml
---
schema_version: 1
generated_at: 2026-05-19T01:23:45Z
request: "provision private networking for a new service"
domain: networking            # networking | compute | database | unknown
confidence: high              # low | medium | high
clarification_needed: false
clarification_question: ""
evidence:
  - source: repo              # repo | docs
    kind: aws-iac
    path: main.tf
    snippet: 'resource "aws_vpc" "main"'
  - source: docs
    kind: official-docs
    title: "Amazon VPC documentation"
    url: "https://docs.aws.amazon.com/vpc/"
conflicts: []                 # always present, may be empty
prerequisites:
  - "terraform >= 1.5"
---

# Provision Private Networking

## Problem summary
...
```

**JSON form (`--json` mode)**

```json
{
  "schema_version": 1,
  "generated_at": "2026-05-19T01:23:45Z",
  "request": "...",
  "domain": "networking",
  "confidence": "high",
  "clarification_needed": false,
  "clarification_question": "",
  "evidence": [/* same shape as frontmatter */],
  "conflicts": [],
  "prerequisites": [/* ... */],
  "plan_md": "# Provision Private Networking\n\n## Problem summary\n..."
}
```

**Confidence rule (MVP)**

| Inputs | Result |
|---|---|
| domain is `compute` or `database` (MVP stub builders) | `low` |
| domain is `networking`, repo has matching evidence | `high` |
| domain is `networking`, no repo evidence (relying on docs only) | `medium` |
| `conflicts` is non-empty | capped at `medium` |

**Conflicts field**

Always emitted (even when `[]`) so spec 02's conflict-surfacing requirement has a stable wire shape. Detection logic: if `discovery.RepoContext.DetectedTools` contains more than one entry, one `Conflict` record is emitted listing all detected tools with three resolution options. See §10.3 for the full design.

**Rationale.**
- `schema_version` + `generated_at` let agents reason about staleness and evolution.
- `domain` + `confidence` directly surface the dispatch decision so an agent reading a compute plan sees `confidence: low` and knows to verify carefully.
- `evidence` preserves the structured records that `discovery.RepoContext` already produces.
- `clarification_needed` is a first-class boolean so agents do not text-parse stdout.

## 3. Plan Dispatch & Per-Domain Builders

A single dispatch function routes a classified `Request` to a domain-specific plan builder. Each builder inspects the relevant subset of `discovery.RepoContext`.

```
ClassifyRequest(raw)  →  Request{domain}
       │
       ▼
BuildPlan(request, repoCtx)
       ├── domain=networking → buildNetworkingPlan(...)   [full hero plan]
       ├── domain=compute    → buildComputePlan(...)      [minimal stub, low confidence]
       ├── domain=database   → buildDatabasePlan(...)     [minimal stub, low confidence]
       └── domain=unknown    → (never reached — clarification gate fires first)
```

**Networking builder.** Three concrete steps: inspect existing networking IaC, add/extend VPC resources, review blast radius. Verification with `terraform fmt -check -recursive`, `terraform validate`, `terraform plan`. Confidence `high` when matching repo evidence exists, otherwise `medium`.

**Compute builder (minimal stub).** Three coarse but correct steps:
1. Identify compute primitive (EC2 / ECS / EKS / Lambda) from request keywords.
2. Inspect existing compute IaC and IAM patterns in the repo.
3. Add the smallest viable resource definition with a least-privilege IAM role.

Verification: `terraform validate` exits 0; `terraform plan` shows only expected compute resources. Confidence always `low` in MVP.

**Database builder (minimal stub).** Three coarse steps:
1. Identify engine (RDS Postgres / RDS MySQL / DynamoDB / Aurora) from request keywords.
2. Inspect existing database IaC, subnet-group networking, and secrets handling.
3. Add the smallest viable instance, parameter group, and security group attachment.

Verification: same shape as compute. Confidence always `low` in MVP.

**Why `confidence: low` instead of a hard refusal.** Spec 04 says compute and database are "still supports... in the combined flow." A typed low-confidence signal honors that while being honest about MVP quality. Agents downstream can decide whether to (a) proceed with extra human review, (b) re-prompt the user, or (c) fall back to another tool.

## 4. Discovery Breadth — Detector Architecture

Replace the inline file-type switch in `discover.go` with a pluggable `Detector` interface. Each IaC format becomes one detector file in `internal/discovery/`.

```go
type Detector interface {
    Name() string                                            // "terraform", "cloudformation", ...
    Matches(repoRoot, path string, info fs.FileInfo) bool    // cheap: filename / extension / presence
    Extract(repoRoot, path string, content []byte) []Evidence
}

var defaultDetectors = []Detector{
    &terraformDetector{},
    &cloudformationDetector{},
    &samDetector{},
    &cdkDetector{},
    &serverlessFrameworkDetector{},
}
```

**Two-layer matching.** `Matches()` is cheap (no file reads); only files that match get `Extract()` called. Keeps the walk fast on large repos.

**`RepoContext` shape.**

```go
type RepoContext struct {
    Root          string
    Evidence      []Evidence
    DetectedTools []string  // distinct Detector.Name() values that fired
}
```

The old `HasTerraform bool` and `HasAWS bool` flags retire; `DetectedTools` plus a filter over `Evidence` cover the same use cases more flexibly.

**Walk hygiene (uniform across detectors).**
- Skipped directories: `.git`, `node_modules`, `vendor`, `dist`, `build`, `.terraform`, `cdk.out`.
- Symlinks not followed.
- Files larger than 1 MiB skipped.
- Snippet extraction: first non-empty line whose first non-whitespace character is not `#` or `/` (skips HCL `#` and `//` comments and YAML `#` comments; JSON has none), trimmed and truncated to 120 chars.

## 5. MVP Detector Set

| Detector | Filename match | Content trigger | Evidence kind |
|---|---|---|---|
| `terraform` | `*.tf`, `*.tfvars` | `aws_` prefix or `provider "aws"` | `aws-iac` |
| `cloudformation` | `*.yaml`, `*.yml`, `*.json` | line starts with `AWSTemplateFormatVersion` | `aws-cloudformation` |
| `sam` | `template.yaml`, `template.yml`, `template.json` | `Transform: AWS::Serverless` | `aws-sam` |
| `cdk` | `cdk.json` (anywhere in tree) | (no content check — file presence is sufficient) | `aws-cdk` |
| `serverless` | `serverless.yml`, `serverless.yaml` | `provider: aws` or `name: aws` under `provider:` | `aws-serverless` |

**Overlap rule.** SAM templates ARE CloudFormation templates; both detectors may fire on the same file. That is correct — the evidence list reflects both facts, and agents can reason from the `kind` values.

**Test strategy.** Each detector gets a `<name>_test.go` with table-driven cases: matching file / non-matching file / empty file / oversized file / file with no AWS markers. One fixture under `testdata/repos/<format>-sample/` per detector.

**Deliberately excluded from MVP.**
- CDK deep parsing (TypeScript/Python/Go source files under `lib/`, `bin/`).
- Pulumi (`Pulumi.yaml`).
- GitHub Actions / `buildspec.yml` / `appspec.yml`.
- Resource-level domain tagging (e.g., classifying `aws_vpc` as networking vs `aws_lambda_function` as compute).

## 6. Impact on the Implementation Plan

The current 9-task plan in `05-implementation-plan.md` stays mostly intact. Three tasks expand and two new tasks slot in. Numbering shifts.

| # | Task | Status | Change |
|---|---|---|---|
| 1 | Scaffold Go CLI Module | unchanged | — |
| 2 | AWS Request Classification | unchanged | — |
| 3 | Repo Discovery | **expanded** | Refactor inline switch into `Detector` interface; add `terraform`, `cloudformation`, `sam`, `cdk`, `serverless` detectors. Retire `HasTerraform`/`HasAWS` flags in favor of `DetectedTools`. |
| 4 | Curated AWS Docs Evidence | unchanged | — |
| 5 | Clarification Gate | unchanged | — |
| 6 | Plan Model & Markdown Renderer | **expanded** | `Plan` gains `Metadata` (with `CostHints`, `AlreadySatisfied`); `Evidence` gains `FetchedAt`/`Status`; renderer prepends YAML frontmatter via `gopkg.in/yaml.v3`. Adds `CostHint` type. |
| 7 | Plan Dispatch & Per-Domain Builders | **expanded** | `BuildPlan` calls `DetectConflicts`, `DetectSatisfaction`, `CostHintsForDomain`; `computeConfidence` caps at `medium` when conflicts exist; stubs for Tasks 13–15 created here. |
| 8 | JSON Output Renderer | new | `internal/render/json.go` assembles the payload from `Plan` + `Metadata` + rendered markdown body. |
| 9 | CLI Orchestration | **expanded** | Adds `--json` and `--no-web` flags; validates mutual exclusion; creates `web.Fetcher` (unless `--no-web`); passes fetcher to `BuildPlan`. |
| 10 | Update README | **expanded** | Documents `--json`, `--no-web`, frontmatter schema, confidence levels, cost hints, already-satisfied flag, conflict detection, detector list. |
| 11 | Verify End-to-End Demo | **expanded** | Adds `--json`, `--no-web`, compute/database stub, multi-detector, multi-IaC conflict, cost-hints, and already-satisfied smoke tests. |
| 12 | **NEW: Web Evidence Fetcher** | new | `internal/web/fetcher.go` — 5s timeout, 1 retry on 5xx, `docs.aws.amazon.com` allowlist, per-run in-memory cache. Updates `combineEvidence` to attach `FetchedAt`/`Status`. |
| 13 | **NEW: Cost Hints** | new | `internal/planner/costs.go` — replaces stub with `CostHintsForDomain`; seed values for networking, compute, database with disclaimer notes. |
| 14 | **NEW: "Already Satisfied" Detection** | new | `internal/planner/satisfaction.go` — replaces stub; heuristic snippet match prepends a confirmation step and sets `AlreadySatisfied`. |
| 15 | **NEW: Multi-IaC Conflict Detection** | new | `internal/planner/conflicts.go` — replaces stub; fires when `DetectedTools` has more than one entry; `testdata/repos/aws-mixed/` fixture. |

**New types added in Task 6**

```go
// internal/planner/plan.go
type CostHint struct {
    ResourceType       string  `yaml:"resource_type"        json:"resource_type"`
    MonthlyUSDEstimate float64 `yaml:"monthly_usd_estimate" json:"monthly_usd_estimate"`
    Notes              string  `yaml:"notes"                json:"notes"`
}

type Metadata struct {
    SchemaVersion         int        `yaml:"schema_version"         json:"schema_version"`
    GeneratedAt           time.Time  `yaml:"generated_at"           json:"generated_at"`
    Request               string     `yaml:"request"                json:"request"`
    Domain                Domain     `yaml:"domain"                 json:"domain"`
    Confidence            Confidence `yaml:"confidence"             json:"confidence"`
    ClarificationNeeded   bool       `yaml:"clarification_needed"   json:"clarification_needed"`
    ClarificationQuestion string     `yaml:"clarification_question" json:"clarification_question"`
    Evidence              []Evidence `yaml:"evidence"               json:"evidence"`
    Conflicts             []Conflict `yaml:"conflicts"              json:"conflicts"`
    Prerequisites         []string   `yaml:"prerequisites"          json:"prerequisites"`
    CostHints             []CostHint `yaml:"cost_hints"             json:"cost_hints"`
    AlreadySatisfied      bool       `yaml:"already_satisfied"      json:"already_satisfied"`
}

type Confidence string
const (
    ConfidenceLow    Confidence = "low"
    ConfidenceMedium Confidence = "medium"
    ConfidenceHigh   Confidence = "high"
)

type Evidence struct {
    Source    string    `yaml:"source"               json:"source"`
    Kind      string    `yaml:"kind"                 json:"kind"`
    Path      string    `yaml:"path,omitempty"       json:"path,omitempty"`
    Snippet   string    `yaml:"snippet,omitempty"    json:"snippet,omitempty"`
    Title     string    `yaml:"title,omitempty"      json:"title,omitempty"`
    URL       string    `yaml:"url,omitempty"        json:"url,omitempty"`
    FetchedAt time.Time `yaml:"fetched_at,omitempty" json:"fetched_at,omitempty"`
    Status    int       `yaml:"status,omitempty"     json:"status,omitempty"`
}

type Conflict struct {
    Description        string   `yaml:"description"         json:"description"`
    Options            []string `yaml:"options,omitempty"   json:"options,omitempty"`
    ResolutionQuestion string   `yaml:"resolution_question" json:"resolution_question"`
}
```

The existing `discovery.Evidence` (internal to the discovery layer) and `planner.Evidence` (public output shape) are distinct types. A small adapter function in Task 7 converts between them.

**New dependency.** `gopkg.in/yaml.v3` for frontmatter marshaling.

## 7. Spec Coverage Check

| Spec doc | Requirement | Covered? |
|---|---|---|
| 01 | Agent-facing | ✅ — `--json` mode + structured frontmatter |
| 01 | Plan-only | ✅ — unchanged |
| 01 | Repo discovery first | ✅ — broader detector set strengthens this |
| 01 | Minimum clarifying question | ✅ — clarification is now machine-readable too |
| 02 | Evidence priority: repo → docs → community | ✅ — `evidence` records carry `source: repo | docs` |
| 02 | Conflict surfacing | ✅ — implemented via multi-IaC detector (Task 15); `conflicts` field populated when `DetectedTools` has more than one entry |
| 03 | Single `.md` file per plan | ✅ — `--json` is a non-file alternative mode, not a second artifact |
| 03 | Mandatory sections | ✅ — unchanged |
| 03 | Citation-free prose | ✅ — citations live in frontmatter/JSON metadata, not inline in the prose |
| 03 | "Structure should remain flexible" | ✅ — `schema_version` reserves the evolution hook |
| 04 | AWS-first | ✅ — five AWS-native detectors |
| 04 | Networking hero | ✅ — full builder |
| 04 | Compute / database supported in combined flow | ✅ — minimal stub builders with `confidence: low` |

## 9. Web Evidence Fetching

Closes spec 02's "Search official/vendor docs on the web" gap. The curated docs URL map (Task 4) already points to `docs.aws.amazon.com`; this adds real HTTP verification so `Evidence` records carry a live status code and timestamp.

**Package:** `internal/web/` — `fetcher.go` + `fetcher_test.go`.

**HTTP client behaviour:**
- Total timeout: 5 seconds (set on `http.Client`).
- User-Agent: `infra-planning-cli/0.1`.
- Retry policy: 1 retry on 5xx, no retry on 4xx.
- Cache: per-run in-memory `map[string]result` keyed by URL. No persistent cache in MVP.

**Allowlist:** only URLs whose host is `docs.aws.amazon.com` are fetched. Any other host returns an error that callers silently ignore (graceful degradation). No HTML body is parsed or stored in MVP — only the HTTP status code and fetch timestamp are recorded.

**New CLI flag:** `--no-web` (default `false`). When set, the fetcher is never created and `combineEvidence` skips all HTTP calls.

**Failure mode:** if a fetch fails or times out, the `Evidence` record is emitted without `fetched_at` / `status` fields. Plan generation is never blocked on a web failure.

**Evidence field additions** (to `planner.Evidence`):
```go
FetchedAt time.Time `yaml:"fetched_at,omitempty" json:"fetched_at,omitempty"`
Status    int       `yaml:"status,omitempty"     json:"status,omitempty"`
```

`planner.combineEvidence` is updated to accept a `*web.Fetcher` (nil = no fetch) and attach `Status` + `FetchedAt` to each `source: docs` evidence record when the fetch succeeds.

---

## 10. Plan Quality Enhancements

Three lightweight signals added to `planner.Metadata` to improve the plan's value for both humans and agents.

### 10.1 Cost Hints

**File:** `internal/planner/costs.go`.

```go
type CostHint struct {
    ResourceType       string  `yaml:"resource_type"        json:"resource_type"`
    MonthlyUSDEstimate float64 `yaml:"monthly_usd_estimate" json:"monthly_usd_estimate"`
    Notes              string  `yaml:"notes"                json:"notes"`
}

func CostHintsForDomain(domain Domain) []CostHint
```

Seed values (hardcoded; each carries `"estimate; verify against current AWS pricing"` as `notes`):

| Domain | Resource | $/mo |
|---|---|---|
| networking | NAT Gateway | 32.40 |
| networking | VPC endpoint (Interface) | 7.20 |
| compute | Lambda (1M req/mo) | 0.20 |
| compute | Fargate task (0.25 vCPU, 24/7) | 8.90 |
| database | RDS db.t3.micro | 12.40 |
| database | DynamoDB on-demand baseline | 0.00 |

`Metadata.CostHints []CostHint` is always emitted as a non-nil slice.

### 10.2 "Already Satisfied" Detection

**File:** `internal/planner/satisfaction.go`.

```go
func DetectSatisfaction(req Request, repo discovery.RepoContext) (satisfied bool, evidencePath string)
```

Heuristic: check whether any `Evidence.Snippet` contains a domain-specific keyword:

| Domain | Keywords |
|---|---|
| networking | `aws_vpc`, `AWS::EC2::VPC` |
| compute | `aws_lambda_function`, `aws_ecs_service`, `aws_instance`, `AWS::Serverless::Function`, `AWS::Lambda::Function` |
| database | `aws_db_instance`, `aws_rds_cluster`, `aws_dynamodb_table`, `AWS::RDS::DBInstance`, `AWS::DynamoDB::Table` |

When satisfied, `BuildPlan` prepends a step:

```go
Step{
    Title: "Confirm whether existing resource already satisfies this request",
    Body:  fmt.Sprintf("Detected matching resource at `%s`. Verify with the human before generating new IaC.", evidencePath),
}
```

`Metadata.AlreadySatisfied bool` is set accordingly. Confidence is unchanged — the prepended step and the flag are the signal. Plan generation always continues.

### 10.3 Conflict Surfacing — Multi-IaC Detection

**File:** `internal/planner/conflicts.go`.

```go
func DetectConflicts(repo discovery.RepoContext) []Conflict
```

MVP rule: if `repo.DetectedTools` contains more than one entry, emit one `Conflict`:

```go
Conflict{
    Description: "Multiple IaC tools detected: " + strings.Join(sortedTools, ", "),
    Options: []string{
        "Use the tool with the most existing evidence",
        "Pick one tool and migrate the others",
        "Document the boundary between tools",
    },
    ResolutionQuestion: "Which IaC tool should this plan target?",
}
```

`BuildPlan` populates `Metadata.Conflicts` via `DetectConflicts` (replacing the always-empty default). `computeConfidence` is updated to cap at `ConfidenceMedium` when `len(conflicts) > 0`.

---

## 8. Deferred Follow-ups

These were deliberately excluded to keep scope tight. Each is a candidate for its own follow-up brainstorm and design doc:

- **CDK deep parsing** (TypeScript/Python/Go source files under `lib/`, `bin/`).
- **Pulumi support** (`Pulumi.yaml`).
- **Resource-level domain inference** (classify `aws_vpc` vs `aws_lambda_function` from repo evidence to feed back into the classifier).
