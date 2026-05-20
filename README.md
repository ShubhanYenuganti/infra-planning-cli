# infra-plan

An AWS-first CLI that turns a natural-language infrastructure request into a detailed, repo-aware markdown execution plan (with YAML frontmatter) or a structured JSON payload — ready for Claude Code, Codex, or another coding agent to consume.

The tool never provisions resources and never asks for credentials. It only analyzes your repository and generates a plan.

## Installation

```bash
go install github.com/ShubhanYenuganti/infra-planning-cli/cmd/infra-plan@latest
```

Or build from source:

```bash
git clone https://github.com/ShubhanYenuganti/infra-planning-cli
cd infra-planning-cli
go build -o infra-plan ./cmd/infra-plan
```

## Usage

```
infra-plan plan REQUEST [--repo PATH] [--out PATH] [--json]
```

### Flags

| Flag | Default | Description |
|---|---|---|
| `--repo` | `.` | Path to the repository to inspect |
| `--out` | _(stdout)_ | Write markdown plan to this file instead of stdout |
| `--json` | `false` | Emit a JSON payload to stdout instead of markdown |

### Examples

Generate a markdown plan for adding a VPC:

```bash
infra-plan plan "Add a VPC with public and private subnets" --repo .
```

Write the plan to a file:

```bash
infra-plan plan "Set up an RDS Postgres instance" --repo . --out plan.md
```

Emit JSON for agent consumption:

```bash
infra-plan plan "Deploy a Lambda behind API Gateway" --repo . --json
```

### Behavior by mode

| Mode | Behavior |
|---|---|
| Default | Prints a markdown plan with YAML frontmatter to stdout (or `--out` file). Exit 0. |
| `--json` | Emits a single JSON object to stdout. No file written. Exit 0. |
| Clarification needed | Prints the clarification question to stderr. No plan output. Exit 0. |
| Hard error | Exit 1 (bad flags, unreadable repo, write failure). |

## Output format

### Markdown (default)

Plans are emitted as markdown with a YAML frontmatter block:

```markdown
---
schema_version: 1
generated_at: 2026-05-19T17:00:00Z
request: "Add a VPC with public and private subnets"
domain: networking
confidence: high
clarification_needed: false
evidence:
  - source: repo
    kind: terraform
    path: main.tf
    snippet: "resource \"aws_vpc\" ..."
  - source: docs
    title: Amazon VPC documentation
    url: https://docs.aws.amazon.com/vpc/latest/userguide/
prerequisites:
  - AWS account with appropriate IAM permissions
  - Terraform >= 1.0 installed
---

# Provision Private Networking

## Problem summary
...

## Steps
1. **Define the VPC resource** ...

## Verification
- Run `terraform plan` and confirm no unexpected changes
```

### JSON (--json)

The JSON output mirrors the same structure with snake_case keys:

```json
{
  "metadata": {
    "schema_version": 1,
    "generated_at": "2026-05-19T17:00:00Z",
    "domain": "networking",
    "confidence": "high",
    ...
  },
  "title": "Provision Private Networking",
  "problem_summary": "...",
  "steps": [...],
  "verification": [...]
}
```

## Architecture

```
infra-plan plan REQUEST
    │
    ├── discovery.DiscoverRepo(--repo)
    │       Walks the repo, runs pluggable detectors:
    │       Terraform · CloudFormation · SAM · CDK · Serverless Framework
    │
    ├── planner.ClassifyRequest(REQUEST)
    │       Token-based domain classification:
    │       networking · compute · database · unknown
    │
    ├── planner.BuildPlan(req, repo)
    │       ├── ClarificationQuestion  — gate for ambiguous requests
    │       ├── DetectConflicts        — multi-IaC framework detection
    │       ├── DetectSatisfaction     — already-provisioned detection
    │       ├── CostHintsForDomain     — per-domain cost estimates
    │       ├── combineEvidence        — repo + curated AWS docs links
    │       └── per-domain builder     — networking / compute / database
    │
    └── render.MarkdownPlan or render.JSONPlan
```

### IaC detectors

| Detector | Matches | Evidence trigger |
|---|---|---|
| Terraform | `*.tf`, `*.tfvars` | `aws_` resource or `provider "aws"` |
| CloudFormation | `*.yaml`, `*.yml`, `*.json` | `AWSTemplateFormatVersion` line |
| SAM | `template.yaml/yml/json` | `Transform: AWS::Serverless` |
| CDK | `cdk.json` | always |
| Serverless Framework | `serverless.yml/yaml` | `provider: aws` |

### Confidence levels

| Domain | Condition | Confidence |
|---|---|---|
| networking | evidence found, no conflicts | high |
| networking | no evidence | medium |
| networking | evidence + conflicts | medium |
| compute | any | low |
| database | any | low |

## Development

```bash
# Run all tests
go test ./...

# Build the binary
go build -o infra-plan ./cmd/infra-plan

# Run against current directory
./infra-plan plan "Add an S3 bucket with versioning" --repo .
```

### Project layout

```
cmd/infra-plan/         binary entrypoint
internal/
  cli/                  Cobra command wiring
  discovery/            repo walker + pluggable detectors
  docs/                 curated AWS docs URL map
  planner/              request classification, plan model, dispatch
  render/               markdown and JSON renderers
  web/                  HTTP docs fetcher (allowlisted to docs.aws.amazon.com)
testdata/repos/         fixture repos for detector tests
```
