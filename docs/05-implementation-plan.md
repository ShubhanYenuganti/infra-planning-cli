# Infra Planning CLI Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build an AWS-first, agent-facing CLI that turns a vague infra request into one detailed, repo-specific markdown execution plan (with YAML frontmatter) or a structured JSON payload, for Claude Code, Codex, or another coding agent to consume.

**Architecture:** Implement a Go CLI with a deterministic pipeline: repo discovery (pluggable per-format detectors) → AWS/domain classification → docs evidence capture + optional HTTP fetch of docs URLs → clarification gate → per-domain plan dispatch (with cost hints, already-satisfied detection, and multi-IaC conflict detection) → frontmatter markdown OR JSON rendering. The MVP is plan-only: never provisions resources, never asks for credentials. See `docs/06-design-improvements.md` for the design rationale.

**Tech Stack:** Go 1.24+, Cobra for CLI commands, `gopkg.in/yaml.v3` for frontmatter marshaling, standard-library `net/http` for docs fetching, standard-library filesystem/process/encoding utilities, table-driven Go tests.

---

## File Structure

All paths are relative to the project root (`infra-planning-cli/`).

**CLI layer**
- `go.mod` — module declaration.
- `cmd/infra-plan/main.go` — binary entrypoint.
- `internal/cli/root.go` — Cobra root command.
- `internal/cli/plan.go` — `plan` subcommand orchestration with `--json` and `--no-web` modes.
- `internal/cli/plan_test.go` — CLI integration tests.

**Discovery layer (pluggable detectors)**
- `internal/discovery/context.go` — `RepoContext` and `Evidence` types.
- `internal/discovery/detector.go` — `Detector` interface.
- `internal/discovery/discover.go` — walker + detector registry.
- `internal/discovery/discover_test.go` — walker integration tests.
- `internal/discovery/terraform.go` + `_test.go`
- `internal/discovery/cloudformation.go` + `_test.go`
- `internal/discovery/sam.go` + `_test.go`
- `internal/discovery/cdk.go` + `_test.go`
- `internal/discovery/serverless.go` + `_test.go`

**Planner layer**
- `internal/planner/request.go` — request model and `ClassifyRequest`.
- `internal/planner/request_test.go`
- `internal/planner/clarify.go` — `ClarificationQuestion`.
- `internal/planner/clarify_test.go`
- `internal/planner/plan.go` — `Plan`, `Step`, `Metadata`, `Confidence`, `Evidence`, `Conflict`, `CostHint` types.
- `internal/planner/dispatch.go` — `BuildPlan` dispatcher + per-domain builders + confidence rule.
- `internal/planner/dispatch_test.go`
- `internal/planner/costs.go` — `CostHintsForDomain` with seed values per domain (stub in Task 7; real in Task 13).
- `internal/planner/costs_test.go`
- `internal/planner/satisfaction.go` — `DetectSatisfaction` heuristic (stub in Task 7; real in Task 14).
- `internal/planner/satisfaction_test.go`
- `internal/planner/conflicts.go` — `DetectConflicts` for multi-IaC repos (stub in Task 7; real in Task 15).
- `internal/planner/conflicts_test.go`

**Web layer**
- `internal/web/fetcher.go` — HTTP fetcher (5 s timeout, 1 retry on 5xx, `docs.aws.amazon.com` allowlist, per-run in-memory cache).
- `internal/web/fetcher_test.go`

**Render layer**
- `internal/render/markdown.go` — `MarkdownPlan` (frontmatter + body) and `MarkdownBody` (body only).
- `internal/render/markdown_test.go`
- `internal/render/json.go` — `JSONPlan`.
- `internal/render/json_test.go`

**Docs evidence**
- `internal/docs/aws.go` — curated AWS docs URL mapping per domain.
- `internal/docs/aws_test.go`

**Fixtures**
- `testdata/repos/aws-terraform/main.tf`
- `testdata/repos/aws-cloudformation/template.yaml`
- `testdata/repos/aws-sam/template.yaml`
- `testdata/repos/aws-cdk/cdk.json`
- `testdata/repos/aws-serverless/serverless.yml`
- `testdata/repos/aws-mixed/main.tf` — Terraform + CDK side-by-side for conflict-detection tests.
- `testdata/repos/aws-mixed/cdk.json`

**Top-level**
- `README.md` — install and usage (already exists; will be updated).

## MVP CLI Contract

```bash
infra-plan plan REQUEST [--repo PATH] [--out PATH] [--json] [--no-web]
```

| Mode | Behavior |
|---|---|
| Default | Writes one `.md` file with YAML frontmatter to `docs/infra-planning-cli/YYYY-MM-DD-<slug>.md` under `--repo` (or to `--out` if given). Exit 0. |
| `--json` | Emits a single JSON object to stdout. **No file written.** Exit 0. Mutually exclusive with `--out`. |
| `--no-web` | Skips HTTP fetches for docs evidence; `fetched_at` and `status` remain zero-valued. |
| Clarification needed (default) | Prints question to stdout, no file written. Exit 0. |
| Clarification needed (`--json`) | Emits payload with `clarification_needed: true` and empty `plan_md`. Exit 0. |
| Hard error | Exit 1 (bad flags, unreadable repo, write failure). |

---

### Task 1: Scaffold the Go CLI Module

**Files:**
- Create: `go.mod`
- Create: `cmd/infra-plan/main.go`
- Create: `internal/cli/root.go`
- Create: `internal/cli/plan.go`

- [ ] **Step 1: Create the module file**

Create `go.mod`:

```
module github.com/ShubhanYenuganti/infra-planning-cli

go 1.24

require (
	github.com/spf13/cobra v1.8.0
	gopkg.in/yaml.v3 v3.0.1
)
```

- [ ] **Step 2: Create the binary entrypoint**

Create `cmd/infra-plan/main.go`:

```go
package main

import (
	"os"

	"github.com/ShubhanYenuganti/infra-planning-cli/internal/cli"
)

func main() {
	if err := cli.NewRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
```

- [ ] **Step 3: Create root command wiring**

Create `internal/cli/root.go`:

```go
package cli

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "infra-plan",
		Short: "Generate agent-ready AWS infrastructure execution plans",
	}
	cmd.AddCommand(newPlanCommand())
	return cmd
}
```

- [ ] **Step 4: Create a stub plan command**

Create `internal/cli/plan.go`:

```go
package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

type planOptions struct {
	repo string
	out  string
	json bool
}

func newPlanCommand() *cobra.Command {
	opts := &planOptions{}

	cmd := &cobra.Command{
		Use:   "plan REQUEST",
		Short: "Turn an infra request into a markdown execution plan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "request: %s\nrepo: %s\nout: %s\njson: %v\n", args[0], opts.repo, opts.out, opts.json)
			return nil
		},
	}

	cmd.Flags().StringVar(&opts.repo, "repo", ".", "repository path to inspect")
	cmd.Flags().StringVar(&opts.out, "out", "", "optional output markdown path")
	cmd.Flags().BoolVar(&opts.json, "json", false, "emit JSON payload to stdout instead of writing a file")
	return cmd
}
```

- [ ] **Step 5: Download dependencies and verify the stub**

Run:

```bash
go mod tidy
go run ./cmd/infra-plan plan "provision networking" --repo .
```

Expected output:

```text
request: provision networking
repo: .
out:
json: false
```

- [ ] **Step 6: Commit**

```bash
git add go.mod go.sum cmd/infra-plan/main.go internal/cli/root.go internal/cli/plan.go
git commit -m "feat: scaffold infra planning CLI"
```

---

### Task 2: Add AWS Request Classification

**Files:**
- Create: `internal/planner/request.go`
- Create: `internal/planner/request_test.go`

- [ ] **Step 1: Write classification tests**

Create `internal/planner/request_test.go`:

```go
package planner

import "testing"

func TestClassifyRequest(t *testing.T) {
	tests := []struct {
		name string
		text string
		want Domain
	}{
		{"vpc", "provision a vpc with private subnets", DomainNetworking},
		{"lambda", "deploy a lambda worker", DomainCompute},
		{"rds", "create a postgres rds database", DomainDatabase},
		{"unknown", "improve our platform", DomainUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassifyRequest(tt.text)
			if got.PrimaryDomain != tt.want {
				t.Fatalf("PrimaryDomain = %q, want %q", got.PrimaryDomain, tt.want)
			}
		})
	}
}
```

- [ ] **Step 2: Run the test and verify failure**

```bash
go test ./internal/planner -run TestClassifyRequest -v
```

Expected: FAIL because `Domain` and `ClassifyRequest` are undefined.

- [ ] **Step 3: Implement minimal classification**

Create `internal/planner/request.go`:

```go
package planner

import "strings"

type Domain string

const (
	DomainNetworking Domain = "networking"
	DomainCompute    Domain = "compute"
	DomainDatabase   Domain = "database"
	DomainUnknown    Domain = "unknown"
)

type Request struct {
	Raw           string
	PrimaryDomain Domain
}

func ClassifyRequest(raw string) Request {
	text := strings.ToLower(raw)

	for _, token := range []string{"vpc", "subnet", "network", "networking", "route table", "security group", "nat gateway"} {
		if strings.Contains(text, token) {
			return Request{Raw: raw, PrimaryDomain: DomainNetworking}
		}
	}
	for _, token := range []string{"ec2", "ecs", "eks", "lambda", "compute", "container", "worker", "service"} {
		if strings.Contains(text, token) {
			return Request{Raw: raw, PrimaryDomain: DomainCompute}
		}
	}
	for _, token := range []string{"rds", "postgres", "postgresql", "mysql", "dynamodb", "database", "db"} {
		if strings.Contains(text, token) {
			return Request{Raw: raw, PrimaryDomain: DomainDatabase}
		}
	}
	return Request{Raw: raw, PrimaryDomain: DomainUnknown}
}
```

- [ ] **Step 4: Run the test and verify pass**

```bash
go test ./internal/planner -run TestClassifyRequest -v
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/planner/request.go internal/planner/request_test.go
git commit -m "feat: classify AWS infra requests"
```

---

### Task 3: Repo Discovery with Pluggable Detectors

**Files:**
- Create: `internal/discovery/context.go`
- Create: `internal/discovery/detector.go`
- Create: `internal/discovery/discover.go`
- Create: `internal/discovery/discover_test.go`
- Create: `internal/discovery/terraform.go` + `terraform_test.go`
- Create: `internal/discovery/cloudformation.go` + `cloudformation_test.go`
- Create: `internal/discovery/sam.go` + `sam_test.go`
- Create: `internal/discovery/cdk.go` + `cdk_test.go`
- Create: `internal/discovery/serverless.go` + `serverless_test.go`
- Create: `testdata/repos/aws-terraform/main.tf`
- Create: `testdata/repos/aws-cloudformation/template.yaml`
- Create: `testdata/repos/aws-sam/template.yaml`
- Create: `testdata/repos/aws-cdk/cdk.json`
- Create: `testdata/repos/aws-serverless/serverless.yml`

- [ ] **Step 1: Add the types**

Create `internal/discovery/context.go`:

```go
package discovery

type Evidence struct {
	Path    string
	Kind    string
	Snippet string
}

type RepoContext struct {
	Root          string
	Evidence      []Evidence
	DetectedTools []string
}
```

- [ ] **Step 2: Add the Detector interface**

Create `internal/discovery/detector.go`:

```go
package discovery

import "io/fs"

type Detector interface {
	Name() string
	Matches(repoRoot, path string, info fs.FileInfo) bool
	Extract(repoRoot, path string, content []byte) []Evidence
}
```

- [ ] **Step 3: Implement the walker (no detectors yet)**

Create `internal/discovery/discover.go`:

```go
package discovery

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const maxFileBytes = 1 << 20 // 1 MiB

var skipDirs = map[string]bool{
	".git":         true,
	"node_modules": true,
	"vendor":       true,
	"dist":         true,
	"build":        true,
	".terraform":   true,
	"cdk.out":      true,
}

func DiscoverRepo(root string) (RepoContext, error) {
	return DiscoverRepoWith(root, defaultDetectors())
}

func DiscoverRepoWith(root string, detectors []Detector) (RepoContext, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return RepoContext{}, err
	}

	ctx := RepoContext{Root: abs}
	fired := map[string]bool{}

	walkErr := filepath.WalkDir(abs, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if skipDirs[entry.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return nil
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}
		if info.Size() > maxFileBytes {
			return nil
		}

		var content []byte
		for _, d := range detectors {
			if !d.Matches(abs, path, info) {
				continue
			}
			if content == nil {
				content, err = os.ReadFile(path)
				if err != nil {
					return nil
				}
			}
			records := d.Extract(abs, path, content)
			if len(records) == 0 {
				continue
			}
			if !fired[d.Name()] {
				fired[d.Name()] = true
				ctx.DetectedTools = append(ctx.DetectedTools, d.Name())
			}
			for _, ev := range records {
				if rel, err := filepath.Rel(abs, ev.Path); err == nil && !strings.HasPrefix(rel, "..") {
					ev.Path = rel
				}
				ctx.Evidence = append(ctx.Evidence, ev)
			}
		}
		return nil
	})
	if walkErr != nil {
		return RepoContext{}, walkErr
	}

	return ctx, nil
}

func defaultDetectors() []Detector {
	return []Detector{
		&terraformDetector{},
		&cloudformationDetector{},
		&samDetector{},
		&cdkDetector{},
		&serverlessDetector{},
	}
}

func firstCodeLine(content []byte) string {
	for _, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
			continue
		}
		if len(trimmed) > 120 {
			trimmed = trimmed[:120]
		}
		return trimmed
	}
	return ""
}
```

- [ ] **Step 4: Create all five fixtures**

Create `testdata/repos/aws-terraform/main.tf`:

```hcl
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "private" {
  vpc_id     = aws_vpc.main.id
  cidr_block = "10.0.1.0/24"
}
```

Create `testdata/repos/aws-cloudformation/template.yaml`:

```yaml
AWSTemplateFormatVersion: '2010-09-09'
Description: Plain CloudFormation VPC
Resources:
  VPC:
    Type: AWS::EC2::VPC
    Properties:
      CidrBlock: 10.0.0.0/16
```

Create `testdata/repos/aws-sam/template.yaml`:

```yaml
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  HelloFunction:
    Type: AWS::Serverless::Function
    Properties:
      Runtime: nodejs20.x
      Handler: index.handler
```

Create `testdata/repos/aws-cdk/cdk.json`:

```json
{
  "app": "npx ts-node --prefer-ts-exts bin/app.ts",
  "context": {}
}
```

Create `testdata/repos/aws-serverless/serverless.yml`:

```yaml
service: hello-service
provider:
  name: aws
  runtime: nodejs20.x
functions:
  hello:
    handler: handler.hello
```

- [ ] **Step 5: Write the integration test**

Create `internal/discovery/discover_test.go`:

```go
package discovery

import (
	"sort"
	"testing"
)

func TestDiscoverRepoTerraformFixture(t *testing.T) {
	ctx, err := DiscoverRepo("../../testdata/repos/aws-terraform")
	if err != nil {
		t.Fatalf("DiscoverRepo error: %v", err)
	}
	assertDetected(t, ctx, "terraform")
	assertEvidencePath(t, ctx, "main.tf")
}

func TestDiscoverRepoCloudFormationFixture(t *testing.T) {
	ctx, err := DiscoverRepo("../../testdata/repos/aws-cloudformation")
	if err != nil {
		t.Fatalf("DiscoverRepo error: %v", err)
	}
	assertDetected(t, ctx, "cloudformation")
	assertEvidencePath(t, ctx, "template.yaml")
}

func TestDiscoverRepoSAMFixture(t *testing.T) {
	ctx, err := DiscoverRepo("../../testdata/repos/aws-sam")
	if err != nil {
		t.Fatalf("DiscoverRepo error: %v", err)
	}
	// SAM templates are CloudFormation templates with a Transform; both detectors fire.
	tools := append([]string{}, ctx.DetectedTools...)
	sort.Strings(tools)
	want := []string{"cloudformation", "sam"}
	if len(tools) != 2 || tools[0] != want[0] || tools[1] != want[1] {
		t.Fatalf("DetectedTools = %v, want %v", tools, want)
	}
}

func TestDiscoverRepoCDKFixture(t *testing.T) {
	ctx, err := DiscoverRepo("../../testdata/repos/aws-cdk")
	if err != nil {
		t.Fatalf("DiscoverRepo error: %v", err)
	}
	assertDetected(t, ctx, "cdk")
	assertEvidencePath(t, ctx, "cdk.json")
}

func TestDiscoverRepoServerlessFixture(t *testing.T) {
	ctx, err := DiscoverRepo("../../testdata/repos/aws-serverless")
	if err != nil {
		t.Fatalf("DiscoverRepo error: %v", err)
	}
	assertDetected(t, ctx, "serverless")
	assertEvidencePath(t, ctx, "serverless.yml")
}

func assertDetected(t *testing.T, ctx RepoContext, want string) {
	t.Helper()
	for _, n := range ctx.DetectedTools {
		if n == want {
			return
		}
	}
	t.Fatalf("DetectedTools = %v, want to include %q", ctx.DetectedTools, want)
}

func assertEvidencePath(t *testing.T, ctx RepoContext, wantPath string) {
	t.Helper()
	for _, ev := range ctx.Evidence {
		if ev.Path == wantPath {
			return
		}
	}
	t.Fatalf("no Evidence with Path=%q in %+v", wantPath, ctx.Evidence)
}
```

- [ ] **Step 6: Run the integration test and verify failure**

```bash
go test ./internal/discovery -run TestDiscoverRepo -v
```

Expected: build failure — `terraformDetector` and the other four are not defined yet.

- [ ] **Step 7: Implement the Terraform detector**

Create `internal/discovery/terraform.go`:

```go
package discovery

import (
	"io/fs"
	"path/filepath"
	"strings"
)

type terraformDetector struct{}

func (t *terraformDetector) Name() string { return "terraform" }

func (t *terraformDetector) Matches(_, path string, _ fs.FileInfo) bool {
	ext := filepath.Ext(path)
	return ext == ".tf" || ext == ".tfvars"
}

func (t *terraformDetector) Extract(_, path string, content []byte) []Evidence {
	s := string(content)
	lower := strings.ToLower(s)
	if !strings.Contains(s, "aws_") && !strings.Contains(lower, `provider "aws"`) {
		return nil
	}
	return []Evidence{{
		Path:    path,
		Kind:    "aws-iac",
		Snippet: firstCodeLine(content),
	}}
}
```

Create `internal/discovery/terraform_test.go`:

```go
package discovery

import "testing"

func TestTerraformDetector(t *testing.T) {
	d := &terraformDetector{}

	cases := []struct {
		name    string
		path    string
		content string
		wantHit bool
	}{
		{"aws tf", "main.tf", `resource "aws_vpc" "main" {}`, true},
		{"non-aws tf", "main.tf", `resource "google_compute_instance" "x" {}`, false},
		{"non-tf ext", "main.yaml", `resource "aws_vpc" "main" {}`, false},
		{"tfvars aws", "vars.tfvars", "region = \"us-east-1\"\naws_role = \"arn:aws:iam::...\"", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := d.Matches("", tc.path, nil); got != (filepathExtIsTfOrTfvars(tc.path)) {
				t.Fatalf("Matches = %v", got)
			}
			ev := d.Extract("", tc.path, []byte(tc.content))
			if tc.wantHit && len(ev) == 0 {
				t.Fatalf("expected evidence, got none")
			}
			if !tc.wantHit && len(ev) > 0 {
				t.Fatalf("expected no evidence, got %+v", ev)
			}
		})
	}
}

func filepathExtIsTfOrTfvars(p string) bool {
	return len(p) > 3 && (p[len(p)-3:] == ".tf" || (len(p) > 7 && p[len(p)-7:] == ".tfvars"))
}
```

- [ ] **Step 8: Implement the CloudFormation detector**

Create `internal/discovery/cloudformation.go`:

```go
package discovery

import (
	"io/fs"
	"path/filepath"
	"strings"
)

type cloudformationDetector struct{}

func (c *cloudformationDetector) Name() string { return "cloudformation" }

func (c *cloudformationDetector) Matches(_, path string, _ fs.FileInfo) bool {
	switch filepath.Ext(path) {
	case ".yaml", ".yml", ".json":
		return true
	}
	return false
}

func (c *cloudformationDetector) Extract(_, path string, content []byte) []Evidence {
	for _, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "AWSTemplateFormatVersion") {
			return []Evidence{{
				Path:    path,
				Kind:    "aws-cloudformation",
				Snippet: firstCodeLine(content),
			}}
		}
	}
	return nil
}
```

Create `internal/discovery/cloudformation_test.go`:

```go
package discovery

import "testing"

func TestCloudFormationDetector(t *testing.T) {
	d := &cloudformationDetector{}

	yamlHit := "AWSTemplateFormatVersion: '2010-09-09'\nResources:\n  VPC:\n    Type: AWS::EC2::VPC"
	yamlMiss := "version: 2\njobs: {}"
	cases := []struct {
		name    string
		path    string
		content string
		wantHit bool
	}{
		{"cfn yaml", "template.yaml", yamlHit, true},
		{"cfn yml", "template.yml", yamlHit, true},
		{"non-cfn yaml", "ci.yaml", yamlMiss, false},
		{"non-yaml ext", "main.tf", yamlHit, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if !d.Matches("", tc.path, nil) && tc.wantHit {
				t.Fatalf("expected Matches=true")
			}
			ev := d.Extract("", tc.path, []byte(tc.content))
			if tc.wantHit && len(ev) == 0 {
				t.Fatalf("expected evidence, got none")
			}
			if !tc.wantHit && len(ev) > 0 {
				t.Fatalf("expected no evidence, got %+v", ev)
			}
		})
	}
}
```

- [ ] **Step 9: Implement the SAM detector**

Create `internal/discovery/sam.go`:

```go
package discovery

import (
	"io/fs"
	"path/filepath"
	"strings"
)

type samDetector struct{}

func (s *samDetector) Name() string { return "sam" }

func (s *samDetector) Matches(_, path string, _ fs.FileInfo) bool {
	base := filepath.Base(path)
	return base == "template.yaml" || base == "template.yml" || base == "template.json"
}

func (s *samDetector) Extract(_, path string, content []byte) []Evidence {
	if !strings.Contains(string(content), "Transform: AWS::Serverless") {
		return nil
	}
	return []Evidence{{
		Path:    path,
		Kind:    "aws-sam",
		Snippet: firstCodeLine(content),
	}}
}
```

Create `internal/discovery/sam_test.go`:

```go
package discovery

import "testing"

func TestSAMDetector(t *testing.T) {
	d := &samDetector{}

	samContent := "AWSTemplateFormatVersion: '2010-09-09'\nTransform: AWS::Serverless-2016-10-31\nResources:\n  Fn:\n    Type: AWS::Serverless::Function"
	cfnOnly := "AWSTemplateFormatVersion: '2010-09-09'\nResources: {}"

	if !d.Matches("", "template.yaml", nil) {
		t.Fatalf("expected Matches=true for template.yaml")
	}
	if d.Matches("", "stack.yaml", nil) {
		t.Fatalf("expected Matches=false for stack.yaml")
	}
	if ev := d.Extract("", "template.yaml", []byte(samContent)); len(ev) == 0 {
		t.Fatalf("expected SAM evidence")
	}
	if ev := d.Extract("", "template.yaml", []byte(cfnOnly)); len(ev) != 0 {
		t.Fatalf("expected no SAM evidence for plain CFN")
	}
}
```

- [ ] **Step 10: Implement the CDK detector**

Create `internal/discovery/cdk.go`:

```go
package discovery

import (
	"io/fs"
	"path/filepath"
)

type cdkDetector struct{}

func (c *cdkDetector) Name() string { return "cdk" }

func (c *cdkDetector) Matches(_, path string, _ fs.FileInfo) bool {
	return filepath.Base(path) == "cdk.json"
}

func (c *cdkDetector) Extract(_, path string, content []byte) []Evidence {
	return []Evidence{{
		Path:    path,
		Kind:    "aws-cdk",
		Snippet: firstCodeLine(content),
	}}
}
```

Create `internal/discovery/cdk_test.go`:

```go
package discovery

import "testing"

func TestCDKDetector(t *testing.T) {
	d := &cdkDetector{}
	if !d.Matches("", "cdk.json", nil) {
		t.Fatalf("expected Matches=true for cdk.json")
	}
	if d.Matches("", "package.json", nil) {
		t.Fatalf("expected Matches=false for package.json")
	}
	ev := d.Extract("", "cdk.json", []byte(`{"app":"npx ts-node bin/app.ts"}`))
	if len(ev) != 1 || ev[0].Kind != "aws-cdk" {
		t.Fatalf("Extract result = %+v", ev)
	}
}
```

- [ ] **Step 11: Implement the Serverless Framework detector**

Create `internal/discovery/serverless.go`:

```go
package discovery

import (
	"io/fs"
	"path/filepath"
	"regexp"
)

type serverlessDetector struct{}

var serverlessProviderRE = regexp.MustCompile(`(?m)^\s*provider:\s*aws\s*$|^\s*name:\s*aws\s*$`)

func (s *serverlessDetector) Name() string { return "serverless" }

func (s *serverlessDetector) Matches(_, path string, _ fs.FileInfo) bool {
	base := filepath.Base(path)
	return base == "serverless.yml" || base == "serverless.yaml"
}

func (s *serverlessDetector) Extract(_, path string, content []byte) []Evidence {
	if !serverlessProviderRE.Match(content) {
		return nil
	}
	return []Evidence{{
		Path:    path,
		Kind:    "aws-serverless",
		Snippet: firstCodeLine(content),
	}}
}
```

Create `internal/discovery/serverless_test.go`:

```go
package discovery

import "testing"

func TestServerlessDetector(t *testing.T) {
	d := &serverlessDetector{}

	flat := "service: hello\nprovider: aws\nfunctions: {}"
	nested := "service: hello\nprovider:\n  name: aws\n  runtime: nodejs20.x"
	googleProvider := "service: hello\nprovider:\n  name: google\n  runtime: nodejs20.x"

	if !d.Matches("", "serverless.yml", nil) {
		t.Fatalf("expected Matches=true")
	}
	if d.Matches("", "compose.yml", nil) {
		t.Fatalf("expected Matches=false")
	}
	if ev := d.Extract("", "serverless.yml", []byte(flat)); len(ev) == 0 {
		t.Fatalf("expected evidence for flat provider")
	}
	if ev := d.Extract("", "serverless.yml", []byte(nested)); len(ev) == 0 {
		t.Fatalf("expected evidence for nested provider")
	}
	if ev := d.Extract("", "serverless.yml", []byte(googleProvider)); len(ev) != 0 {
		t.Fatalf("expected no evidence for google provider")
	}
}
```

- [ ] **Step 12: Run all discovery tests and verify pass**

```bash
go test ./internal/discovery -v
```

Expected: PASS for every test in every file.

- [ ] **Step 13: Commit**

```bash
git add internal/discovery testdata/repos
git commit -m "feat: pluggable detector-based repo discovery"
```

---

### Task 4: Add Curated AWS Docs Evidence

**Files:**
- Create: `internal/docs/aws.go`
- Create: `internal/docs/aws_test.go`

- [ ] **Step 1: Write docs mapping tests**

Create `internal/docs/aws_test.go`:

```go
package docs

import "testing"

func TestAWSDocsForDomain(t *testing.T) {
	got := AWSDocsForDomain("networking")
	if len(got) == 0 {
		t.Fatalf("expected networking docs")
	}
	if got[0].Title != "Amazon VPC documentation" {
		t.Fatalf("first title = %q", got[0].Title)
	}
	if got[0].URL == "" {
		t.Fatalf("first URL was empty")
	}
	if unknown := AWSDocsForDomain("unknown"); len(unknown) != 0 {
		t.Fatalf("unknown docs length = %d, want 0", len(unknown))
	}
}
```

- [ ] **Step 2: Run the test and verify failure**

```bash
go test ./internal/docs -v
```

Expected: FAIL because `AWSDocsForDomain` is undefined.

- [ ] **Step 3: Implement AWS docs mapping**

Create `internal/docs/aws.go`:

```go
package docs

type DocLink struct {
	Title string
	URL   string
}

func AWSDocsForDomain(domain string) []DocLink {
	switch domain {
	case "networking":
		return []DocLink{
			{Title: "Amazon VPC documentation", URL: "https://docs.aws.amazon.com/vpc/"},
			{Title: "VPC security groups", URL: "https://docs.aws.amazon.com/vpc/latest/userguide/vpc-security-groups.html"},
			{Title: "Route tables", URL: "https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Route_Tables.html"},
		}
	case "compute":
		return []DocLink{
			{Title: "Amazon ECS documentation", URL: "https://docs.aws.amazon.com/ecs/"},
			{Title: "AWS Lambda documentation", URL: "https://docs.aws.amazon.com/lambda/"},
		}
	case "database":
		return []DocLink{
			{Title: "Amazon RDS documentation", URL: "https://docs.aws.amazon.com/rds/"},
			{Title: "Amazon DynamoDB documentation", URL: "https://docs.aws.amazon.com/dynamodb/"},
		}
	default:
		return nil
	}
}
```

- [ ] **Step 4: Run the test and verify pass**

```bash
go test ./internal/docs -v
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/docs/aws.go internal/docs/aws_test.go
git commit -m "feat: add AWS docs evidence map"
```

---

### Task 5: Add Clarification Gate

**Files:**
- Create: `internal/planner/clarify.go`
- Create: `internal/planner/clarify_test.go`

- [ ] **Step 1: Write clarification tests**

Create `internal/planner/clarify_test.go`:

```go
package planner

import "testing"

func TestClarificationQuestion(t *testing.T) {
	tests := []struct {
		name string
		req  Request
		want string
	}{
		{"known", Request{Raw: "make a vpc", PrimaryDomain: DomainNetworking}, ""},
		{"unknown", Request{Raw: "make infra better", PrimaryDomain: DomainUnknown}, "Which AWS area should this plan focus on first: networking, compute, or database?"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ClarificationQuestion(tt.req); got != tt.want {
				t.Fatalf("ClarificationQuestion() = %q, want %q", got, tt.want)
			}
		})
	}
}
```

- [ ] **Step 2: Run the test and verify failure**

```bash
go test ./internal/planner -run TestClarificationQuestion -v
```

Expected: FAIL because `ClarificationQuestion` is undefined.

- [ ] **Step 3: Implement clarification gate**

Create `internal/planner/clarify.go`:

```go
package planner

func ClarificationQuestion(req Request) string {
	if req.PrimaryDomain == DomainUnknown {
		return "Which AWS area should this plan focus on first: networking, compute, or database?"
	}
	return ""
}
```

- [ ] **Step 4: Run planner tests and verify pass**

```bash
go test ./internal/planner -v
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/planner/clarify.go internal/planner/clarify_test.go
git commit -m "feat: add minimal clarification gate"
```

---

### Task 6: Plan Model and Markdown Renderer with Frontmatter

**Files:**
- Create: `internal/planner/plan.go`
- Create: `internal/render/markdown.go`
- Create: `internal/render/markdown_test.go`

- [ ] **Step 1: Add the plan + metadata types**

Create `internal/planner/plan.go`:

```go
package planner

import "time"

type Confidence string

const (
	ConfidenceLow    Confidence = "low"
	ConfidenceMedium Confidence = "medium"
	ConfidenceHigh   Confidence = "high"
)

type CostHint struct {
	ResourceType       string  `yaml:"resource_type"        json:"resource_type"`
	MonthlyUSDEstimate float64 `yaml:"monthly_usd_estimate" json:"monthly_usd_estimate"`
	Notes              string  `yaml:"notes"                json:"notes"`
}

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

type Step struct {
	Title string
	Body  string
}

type Plan struct {
	Metadata        Metadata
	Title           string
	ProblemSummary  string
	RecommendedPath string
	Steps           []Step
	Verification    []string
}
```

- [ ] **Step 2: Write renderer tests**

Create `internal/render/markdown_test.go`:

```go
package render

import (
	"strings"
	"testing"
	"time"

	"github.com/ShubhanYenuganti/infra-planning-cli/internal/planner"
)

func sampleFullPlan() planner.Plan {
	return planner.Plan{
		Metadata: planner.Metadata{
			SchemaVersion: 1,
			GeneratedAt:   time.Date(2026, 5, 19, 1, 23, 45, 0, time.UTC),
			Request:       "provision private networking",
			Domain:        planner.DomainNetworking,
			Confidence:    planner.ConfidenceHigh,
			Evidence: []planner.Evidence{{
				Source: "repo", Kind: "aws-iac", Path: "main.tf",
				Snippet: `resource "aws_vpc" "main"`,
			}},
			Conflicts:     []planner.Conflict{},
			Prerequisites: []string{"terraform >= 1.5"},
			CostHints: []planner.CostHint{{
				ResourceType:       "aws_nat_gateway",
				MonthlyUSDEstimate: 32.40,
				Notes:              "estimate; verify against current AWS pricing",
			}},
			AlreadySatisfied: false,
		},
		Title:           "Provision Private Networking",
		ProblemSummary:  "Create private networking for a new AWS service.",
		RecommendedPath: "Use an existing Terraform VPC pattern.",
		Steps: []planner.Step{
			{Title: "Inspect existing Terraform", Body: "Review current VPC resources."},
		},
		Verification: []string{"terraform validate exits 0"},
	}
}

func TestMarkdownPlanIncludesFrontmatterAndSections(t *testing.T) {
	md := MarkdownPlan(sampleFullPlan())

	if !strings.HasPrefix(md, "---\n") {
		t.Fatalf("expected frontmatter at top, got:\n%s", md)
	}
	for _, want := range []string{
		"schema_version: 1",
		"domain: networking",
		"confidence: high",
		"cost_hints:",
		"already_satisfied: false",
		"# Provision Private Networking",
		"## Problem summary",
		"## Recommended path",
		"## Steps",
		"## Verification",
	} {
		if !strings.Contains(md, want) {
			t.Fatalf("markdown missing %q:\n%s", want, md)
		}
	}
}

func TestMarkdownBodyExcludesFrontmatter(t *testing.T) {
	body := MarkdownBody(sampleFullPlan())
	if strings.Contains(body, "schema_version") {
		t.Fatalf("body must not contain frontmatter:\n%s", body)
	}
	if !strings.HasPrefix(body, "# Provision Private Networking") {
		t.Fatalf("body must start with title, got:\n%s", body)
	}
}
```

- [ ] **Step 3: Run renderer tests and verify failure**

```bash
go test ./internal/render -v
```

Expected: FAIL because the package and functions are undefined.

- [ ] **Step 4: Implement the renderer**

Create `internal/render/markdown.go`:

```go
package render

import (
	"fmt"
	"strings"

	"github.com/ShubhanYenuganti/infra-planning-cli/internal/planner"
	"gopkg.in/yaml.v3"
)

func MarkdownPlan(p planner.Plan) string {
	var b strings.Builder
	b.WriteString("---\n")
	out, err := yaml.Marshal(p.Metadata)
	if err == nil {
		b.Write(out)
	}
	b.WriteString("---\n\n")
	b.WriteString(MarkdownBody(p))
	return b.String()
}

func MarkdownBody(p planner.Plan) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", p.Title)
	fmt.Fprintf(&b, "## Problem summary\n\n%s\n\n", p.ProblemSummary)
	fmt.Fprintf(&b, "## Recommended path\n\n%s\n\n", p.RecommendedPath)
	b.WriteString("## Steps\n\n")
	for i, s := range p.Steps {
		fmt.Fprintf(&b, "%d. **%s**\n\n   %s\n\n", i+1, s.Title, s.Body)
	}
	b.WriteString("## Verification\n\n")
	for _, item := range p.Verification {
		fmt.Fprintf(&b, "- %s\n", item)
	}
	return b.String()
}
```

- [ ] **Step 5: Run renderer tests and verify pass**

```bash
go test ./internal/render -v
```

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add internal/planner/plan.go internal/render/markdown.go internal/render/markdown_test.go
git commit -m "feat: render plans with YAML frontmatter, cost hints, and already_satisfied"
```

---

### Task 7: Plan Dispatch and Per-Domain Builders

**Files:**
- Create: `internal/planner/dispatch.go`
- Create: `internal/planner/dispatch_test.go`
- Create: `internal/planner/costs.go` (stub; replaced in Task 13)
- Create: `internal/planner/satisfaction.go` (stub; replaced in Task 14)
- Create: `internal/planner/conflicts.go` (stub; replaced in Task 15)

- [ ] **Step 1: Write dispatch + confidence tests**

Create `internal/planner/dispatch_test.go`:

```go
package planner

import (
	"strings"
	"testing"

	"github.com/ShubhanYenuganti/infra-planning-cli/internal/discovery"
)

func TestBuildPlanNetworkingHighConfidence(t *testing.T) {
	req := ClassifyRequest("provision a vpc")
	repo := discovery.RepoContext{
		Evidence:      []discovery.Evidence{{Path: "main.tf", Kind: "aws-iac", Snippet: `resource "aws_vpc"`}},
		DetectedTools: []string{"terraform"},
	}
	plan := BuildPlan(req, repo)
	if plan.Metadata.Domain != DomainNetworking {
		t.Fatalf("domain = %q", plan.Metadata.Domain)
	}
	if plan.Metadata.Confidence != ConfidenceHigh {
		t.Fatalf("confidence = %q, want high", plan.Metadata.Confidence)
	}
	if !strings.Contains(plan.Title, "Networking") {
		t.Fatalf("title = %q", plan.Title)
	}
	if len(plan.Steps) == 0 {
		t.Fatalf("no steps generated")
	}
}

func TestBuildPlanNetworkingMediumWithoutEvidence(t *testing.T) {
	req := ClassifyRequest("provision a vpc")
	plan := BuildPlan(req, discovery.RepoContext{})
	if plan.Metadata.Confidence != ConfidenceMedium {
		t.Fatalf("confidence = %q, want medium", plan.Metadata.Confidence)
	}
}

func TestBuildPlanComputeLow(t *testing.T) {
	req := ClassifyRequest("deploy a lambda worker")
	plan := BuildPlan(req, discovery.RepoContext{})
	if plan.Metadata.Domain != DomainCompute {
		t.Fatalf("domain = %q", plan.Metadata.Domain)
	}
	if plan.Metadata.Confidence != ConfidenceLow {
		t.Fatalf("confidence = %q, want low", plan.Metadata.Confidence)
	}
	if !strings.Contains(plan.Title, "Compute") {
		t.Fatalf("title = %q", plan.Title)
	}
}

func TestBuildPlanDatabaseLow(t *testing.T) {
	req := ClassifyRequest("create a postgres rds database")
	plan := BuildPlan(req, discovery.RepoContext{})
	if plan.Metadata.Domain != DomainDatabase {
		t.Fatalf("domain = %q", plan.Metadata.Domain)
	}
	if plan.Metadata.Confidence != ConfidenceLow {
		t.Fatalf("confidence = %q, want low", plan.Metadata.Confidence)
	}
}

func TestBuildPlanConflictsNonNilSlice(t *testing.T) {
	req := ClassifyRequest("provision a vpc")
	plan := BuildPlan(req, discovery.RepoContext{})
	if plan.Metadata.Conflicts == nil {
		t.Fatalf("conflicts must be non-nil empty slice for stable wire shape")
	}
}

func TestBuildPlanCostHintsNonNil(t *testing.T) {
	req := ClassifyRequest("provision a vpc")
	plan := BuildPlan(req, discovery.RepoContext{})
	if plan.Metadata.CostHints == nil {
		t.Fatalf("cost_hints must be non-nil empty slice for stable wire shape")
	}
}

func TestBuildPlanAlreadySatisfiedDefaultFalse(t *testing.T) {
	req := ClassifyRequest("provision a vpc")
	plan := BuildPlan(req, discovery.RepoContext{})
	if plan.Metadata.AlreadySatisfied {
		t.Fatalf("already_satisfied must be false when no matching snippet")
	}
}
```

- [ ] **Step 2: Run the test and verify failure**

```bash
go test ./internal/planner -run TestBuildPlan -v
```

Expected: FAIL because `BuildPlan` is undefined.

- [ ] **Step 3: Create stub helper files**

Create `internal/planner/costs.go` (stub; fully implemented in Task 13):

```go
package planner

// CostHintsForDomain returns seed cost estimates for a domain.
// Stub: replaced in Task 13.
func CostHintsForDomain(domain Domain) []CostHint { return []CostHint{} }
```

Create `internal/planner/satisfaction.go` (stub; fully implemented in Task 14):

```go
package planner

import "github.com/ShubhanYenuganti/infra-planning-cli/internal/discovery"

// DetectSatisfaction returns (true, evidencePath) when the repo already
// contains a resource matching the request. Stub: replaced in Task 14.
func DetectSatisfaction(req Request, repo discovery.RepoContext) (bool, string) {
	return false, ""
}
```

Create `internal/planner/conflicts.go` (stub; fully implemented in Task 15):

```go
package planner

import "github.com/ShubhanYenuganti/infra-planning-cli/internal/discovery"

// DetectConflicts returns conflicts when multiple IaC tools are present.
// Stub: replaced in Task 15.
func DetectConflicts(repo discovery.RepoContext) []Conflict { return []Conflict{} }
```

- [ ] **Step 4: Implement dispatch and per-domain builders**

Create `internal/planner/dispatch.go`:

```go
package planner

import (
	"fmt"
	"time"

	"github.com/ShubhanYenuganti/infra-planning-cli/internal/discovery"
	"github.com/ShubhanYenuganti/infra-planning-cli/internal/docs"
)

func BuildPlan(req Request, repo discovery.RepoContext) Plan {
	now := time.Now().UTC()

	conflicts := DetectConflicts(repo)
	satisfied, satisfiedPath := DetectSatisfaction(req, repo)

	var plan Plan
	switch req.PrimaryDomain {
	case DomainNetworking:
		plan = buildNetworkingPlan(req, repo)
	case DomainCompute:
		plan = buildComputePlan(req, repo)
	case DomainDatabase:
		plan = buildDatabasePlan(req, repo)
	default:
		plan = buildNetworkingPlan(req, repo) // safety net; clarification gate should have fired
	}

	if satisfied {
		plan.Steps = append([]Step{{
			Title: "Confirm whether existing resource already satisfies this request",
			Body:  fmt.Sprintf("Detected matching resource at `%s`. Verify with the human before generating new IaC.", satisfiedPath),
		}}, plan.Steps...)
	}

	plan.Metadata.SchemaVersion = 1
	plan.Metadata.GeneratedAt = now
	plan.Metadata.Request = req.Raw
	plan.Metadata.Domain = req.PrimaryDomain
	plan.Metadata.Confidence = computeConfidence(req.PrimaryDomain, repo.Evidence, conflicts)
	plan.Metadata.Conflicts = conflicts
	plan.Metadata.Evidence = combineEvidence(repo.Evidence, docs.AWSDocsForDomain(string(req.PrimaryDomain)))
	plan.Metadata.CostHints = CostHintsForDomain(req.PrimaryDomain)
	plan.Metadata.AlreadySatisfied = satisfied
	if plan.Metadata.Prerequisites == nil {
		plan.Metadata.Prerequisites = []string{}
	}
	return plan
}

func computeConfidence(domain Domain, repoEv []discovery.Evidence, conflicts []Conflict) Confidence {
	if len(conflicts) > 0 {
		if domain == DomainNetworking && len(repoEv) > 0 {
			return ConfidenceMedium
		}
		return ConfidenceLow
	}
	if domain == DomainCompute || domain == DomainDatabase {
		return ConfidenceLow
	}
	if domain == DomainNetworking {
		if len(repoEv) > 0 {
			return ConfidenceHigh
		}
		return ConfidenceMedium
	}
	return ConfidenceLow
}

func combineEvidence(repoEv []discovery.Evidence, docLinks []docs.DocLink) []Evidence {
	out := make([]Evidence, 0, len(repoEv)+len(docLinks))
	for _, ev := range repoEv {
		out = append(out, Evidence{
			Source:  "repo",
			Kind:    ev.Kind,
			Path:    ev.Path,
			Snippet: ev.Snippet,
		})
	}
	for _, link := range docLinks {
		out = append(out, Evidence{
			Source: "docs",
			Kind:   "official-docs",
			Title:  link.Title,
			URL:    link.URL,
		})
	}
	return out
}

func buildNetworkingPlan(req Request, repo discovery.RepoContext) Plan {
	evidenceLine := "No existing AWS IaC evidence was found. Start by adding a minimal Terraform networking module."
	if len(repo.Evidence) > 0 {
		ev := repo.Evidence[0]
		evidenceLine = fmt.Sprintf("Use existing repo evidence from `%s`: `%s`.", ev.Path, ev.Snippet)
	}
	return Plan{
		Title:           "Provision Private Networking",
		ProblemSummary:  fmt.Sprintf("Request: %s", req.Raw),
		RecommendedPath: "Prefer the repository's existing Terraform/AWS patterns, then add or extend VPC, private subnet, route table, and security group definitions in small reviewed changes.",
		Steps: []Step{
			{Title: "Inspect existing networking IaC", Body: evidenceLine},
			{Title: "Add or extend VPC resources", Body: "Create the smallest Terraform change that defines the required VPC CIDR, private subnets, route tables, and security groups."},
			{Title: "Review blast radius", Body: "Run `terraform plan` and verify that only expected networking resources are created or changed."},
		},
		Verification: []string{
			"`terraform fmt -check -recursive` exits 0",
			"`terraform validate` exits 0",
			"`terraform plan` shows only expected networking changes",
		},
		Metadata: Metadata{
			Prerequisites: []string{"terraform >= 1.5"},
		},
	}
}

func buildComputePlan(req Request, repo discovery.RepoContext) Plan {
	return Plan{
		Title:           "Deploy Compute Workload",
		ProblemSummary:  fmt.Sprintf("Request: %s", req.Raw),
		RecommendedPath: "Identify the appropriate compute primitive (EC2, ECS, EKS, or Lambda) from the request, then add the smallest viable resource definition alongside a least-privilege IAM role.",
		Steps: []Step{
			{Title: "Identify compute primitive", Body: "From the request, pick exactly one of EC2, ECS, EKS, or Lambda. If the request is ambiguous, stop and ask the human before generating IaC."},
			{Title: "Inspect existing compute IaC and IAM patterns", Body: "Search the repo for existing compute resources, task definitions, and IAM roles. Prefer extending established patterns over introducing a new one."},
			{Title: "Add the smallest viable resource definition with a least-privilege IAM role", Body: "Create one resource (function/service/task) and one IAM role. Grant only the permissions required for the smoke test."},
		},
		Verification: []string{
			"`terraform validate` exits 0",
			"`terraform plan` shows only the expected compute resources and one IAM role",
		},
		Metadata: Metadata{
			Prerequisites: []string{"terraform >= 1.5", "AWS account access (configured outside this CLI)"},
		},
	}
}

func buildDatabasePlan(req Request, repo discovery.RepoContext) Plan {
	return Plan{
		Title:           "Provision Managed Database",
		ProblemSummary:  fmt.Sprintf("Request: %s", req.Raw),
		RecommendedPath: "Identify the database engine (RDS Postgres, RDS MySQL, Aurora, or DynamoDB), then add the smallest viable instance plus parameter group and security group attachment.",
		Steps: []Step{
			{Title: "Identify engine", Body: "From the request, pick exactly one of RDS Postgres, RDS MySQL, Aurora, or DynamoDB. If ambiguous, stop and ask the human before generating IaC."},
			{Title: "Inspect existing database IaC, subnet groups, and secrets handling", Body: "Look for existing database resources, DB subnet groups, and how secrets/passwords are managed (Secrets Manager, SSM Parameter Store, etc.). Reuse patterns when present."},
			{Title: "Add the smallest viable instance + parameter group + security group attachment", Body: "Create one instance with minimal storage, a parameter group, and a security group restricted to the application's VPC CIDR."},
		},
		Verification: []string{
			"`terraform validate` exits 0",
			"`terraform plan` shows only the expected database resources and supporting groups",
		},
		Metadata: Metadata{
			Prerequisites: []string{"terraform >= 1.5"},
		},
	}
}
```

- [ ] **Step 5: Run dispatch tests and verify pass**

```bash
go test ./internal/planner -v
```

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add internal/planner/dispatch.go internal/planner/dispatch_test.go \
        internal/planner/costs.go internal/planner/satisfaction.go internal/planner/conflicts.go
git commit -m "feat: per-domain plan dispatch with confidence rule and stub quality helpers"
```

---

### Task 8: JSON Output Renderer

**Files:**
- Create: `internal/render/json.go`
- Create: `internal/render/json_test.go`

- [ ] **Step 1: Write JSON renderer tests**

Create `internal/render/json_test.go`:

```go
package render

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestJSONPlanContainsMetadataAndPlanMD(t *testing.T) {
	payload, err := JSONPlan(sampleFullPlan())
	if err != nil {
		t.Fatalf("JSONPlan error: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	for _, key := range []string{
		"schema_version", "domain", "confidence",
		"clarification_needed", "evidence", "conflicts", "prerequisites",
		"cost_hints", "already_satisfied", "plan_md",
	} {
		if _, ok := decoded[key]; !ok {
			t.Fatalf("missing key %q in JSON: %s", key, payload)
		}
	}

	planMD, _ := decoded["plan_md"].(string)
	if !strings.HasPrefix(planMD, "# Provision Private Networking") {
		t.Fatalf("plan_md must contain the body without frontmatter, got: %q", planMD)
	}
	if strings.Contains(planMD, "schema_version") {
		t.Fatalf("plan_md must not contain frontmatter, got: %q", planMD)
	}
}
```

- [ ] **Step 2: Run the test and verify failure**

```bash
go test ./internal/render -run TestJSONPlan -v
```

Expected: FAIL because `JSONPlan` is undefined.

- [ ] **Step 3: Implement JSON renderer**

Create `internal/render/json.go`:

```go
package render

import (
	"encoding/json"

	"github.com/ShubhanYenuganti/infra-planning-cli/internal/planner"
)

type jsonPayload struct {
	planner.Metadata
	PlanMD string `json:"plan_md"`
}

func JSONPlan(p planner.Plan) ([]byte, error) {
	payload := jsonPayload{
		Metadata: p.Metadata,
		PlanMD:   MarkdownBody(p),
	}
	return json.MarshalIndent(payload, "", "  ")
}
```

- [ ] **Step 4: Run the test and verify pass**

```bash
go test ./internal/render -v
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/render/json.go internal/render/json_test.go
git commit -m "feat: render plan as JSON payload"
```

---

### Task 9: CLI Orchestration with --json Mode

**Files:**
- Modify: `internal/cli/plan.go`
- Create: `internal/cli/plan_test.go`

- [ ] **Step 1: Write CLI integration tests**

Create `internal/cli/plan_test.go`:

```go
package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func runCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	cmd := NewRootCommand()
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stdout)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return stdout.String(), err
}

func makeTerraformRepo(t *testing.T) string {
	t.Helper()
	repo := filepath.Join(t.TempDir(), "repo")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "main.tf"),
		[]byte(`resource "aws_vpc" "main" { cidr_block = "10.0.0.0/16" }`), 0o644); err != nil {
		t.Fatal(err)
	}
	return repo
}

func TestPlanWritesMarkdownWithFrontmatter(t *testing.T) {
	repo := makeTerraformRepo(t)
	out := filepath.Join(repo, "docs", "infra-planning-cli", "2026-05-19-vpc.md")

	stdout, err := runCmd(t, "plan", "provision private networking", "--repo", repo, "--out", out, "--no-web")
	if err != nil {
		t.Fatalf("execute: %v\n%s", err, stdout)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	md := string(data)
	for _, want := range []string{
		"---\n",
		"schema_version: 1",
		"domain: networking",
		"confidence: high",
		"cost_hints:",
		"# Provision Private Networking",
		"## Problem summary",
		"## Recommended path",
		"## Steps",
		"## Verification",
		"terraform validate",
	} {
		if !strings.Contains(md, want) {
			t.Fatalf("markdown missing %q:\n%s", want, md)
		}
	}
}

func TestPlanJSONModeEmitsPayloadAndWritesNoFile(t *testing.T) {
	repo := makeTerraformRepo(t)

	stdout, err := runCmd(t, "plan", "provision private networking", "--repo", repo, "--json", "--no-web")
	if err != nil {
		t.Fatalf("execute: %v\n%s", err, stdout)
	}
	var decoded map[string]any
	if err := json.Unmarshal([]byte(stdout), &decoded); err != nil {
		t.Fatalf("not valid JSON: %v\n%s", err, stdout)
	}
	if decoded["domain"] != "networking" {
		t.Fatalf("domain = %v", decoded["domain"])
	}
	if decoded["confidence"] != "high" {
		t.Fatalf("confidence = %v", decoded["confidence"])
	}
	if _, ok := decoded["plan_md"]; !ok {
		t.Fatalf("plan_md missing")
	}
	if entries, _ := os.ReadDir(filepath.Join(repo, "docs")); len(entries) != 0 {
		t.Fatalf("--json must not write a file, found: %v", entries)
	}
}

func TestPlanClarificationDefaultMode(t *testing.T) {
	repo := makeTerraformRepo(t)

	stdout, err := runCmd(t, "plan", "make our platform better", "--repo", repo, "--no-web")
	if err != nil {
		t.Fatalf("execute: %v\n%s", err, stdout)
	}
	if !strings.Contains(stdout, "Which AWS area") {
		t.Fatalf("expected clarification text in stdout, got: %s", stdout)
	}
	if entries, _ := os.ReadDir(filepath.Join(repo, "docs")); len(entries) != 0 {
		t.Fatalf("clarification must not write a file, found: %v", entries)
	}
}

func TestPlanClarificationJSONMode(t *testing.T) {
	repo := makeTerraformRepo(t)

	stdout, err := runCmd(t, "plan", "make our platform better", "--repo", repo, "--json", "--no-web")
	if err != nil {
		t.Fatalf("execute: %v\n%s", err, stdout)
	}
	var decoded map[string]any
	if err := json.Unmarshal([]byte(stdout), &decoded); err != nil {
		t.Fatalf("not valid JSON: %v\n%s", err, stdout)
	}
	if decoded["clarification_needed"] != true {
		t.Fatalf("clarification_needed = %v", decoded["clarification_needed"])
	}
	if !strings.Contains(decoded["clarification_question"].(string), "Which AWS area") {
		t.Fatalf("clarification_question = %v", decoded["clarification_question"])
	}
	if decoded["plan_md"].(string) != "" {
		t.Fatalf("plan_md must be empty during clarification, got: %q", decoded["plan_md"])
	}
}

func TestPlanJSONAndOutAreMutuallyExclusive(t *testing.T) {
	repo := makeTerraformRepo(t)
	_, err := runCmd(t, "plan", "vpc", "--repo", repo, "--json", "--out", "/tmp/x.md")
	if err == nil {
		t.Fatalf("expected error when --json and --out are combined")
	}
}

func TestPlanNoWebFlagAccepted(t *testing.T) {
	repo := makeTerraformRepo(t)
	_, err := runCmd(t, "plan", "provision private networking", "--repo", repo, "--json", "--no-web")
	if err != nil {
		t.Fatalf("--no-web flag should be accepted: %v", err)
	}
}
```

- [ ] **Step 2: Run the tests and verify failure**

```bash
go test ./internal/cli -v
```

Expected: every test FAILs because the stub command does not yet orchestrate.

- [ ] **Step 3: Replace the stub with full orchestration**

Replace `internal/cli/plan.go` with:

```go
package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/ShubhanYenuganti/infra-planning-cli/internal/discovery"
	"github.com/ShubhanYenuganti/infra-planning-cli/internal/planner"
	"github.com/ShubhanYenuganti/infra-planning-cli/internal/render"
	"github.com/ShubhanYenuganti/infra-planning-cli/internal/web"
	"github.com/spf13/cobra"
)

type planOptions struct {
	repo  string
	out   string
	json  bool
	noWeb bool
}

func newPlanCommand() *cobra.Command {
	opts := &planOptions{}

	cmd := &cobra.Command{
		Use:   "plan REQUEST",
		Short: "Turn an infra request into a markdown execution plan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.json && opts.out != "" {
				return fmt.Errorf("--json and --out are mutually exclusive")
			}

			req := planner.ClassifyRequest(args[0])

			if question := planner.ClarificationQuestion(req); question != "" {
				return emitClarification(cmd, args[0], req, question, opts.json)
			}

			repoCtx, err := discovery.DiscoverRepo(opts.repo)
			if err != nil {
				return err
			}

			var fetcher *web.Fetcher
			if !opts.noWeb {
				fetcher = web.New()
			}

			plan := planner.BuildPlan(req, repoCtx, fetcher)

			if opts.json {
				payload, err := render.JSONPlan(plan)
				if err != nil {
					return err
				}
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), string(payload))
				return nil
			}

			outPath := opts.out
			if outPath == "" {
				outPath = defaultOutputPath(opts.repo, args[0], time.Now())
			}
			if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
				return err
			}
			if err := os.WriteFile(outPath, []byte(render.MarkdownPlan(plan)), 0o644); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "wrote plan: %s\n", outPath)
			return nil
		},
	}

	cmd.Flags().StringVar(&opts.repo, "repo", ".", "repository path to inspect")
	cmd.Flags().StringVar(&opts.out, "out", "", "output markdown path (default: docs/infra-planning-cli/YYYY-MM-DD-<slug>.md)")
	cmd.Flags().BoolVar(&opts.json, "json", false, "emit JSON payload to stdout instead of writing a file (mutually exclusive with --out)")
	cmd.Flags().BoolVar(&opts.noWeb, "no-web", false, "skip HTTP fetches for docs evidence")
	return cmd
}

func emitClarification(cmd *cobra.Command, raw string, req planner.Request, question string, asJSON bool) error {
	if !asJSON {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), question)
		return nil
	}
	meta := planner.Metadata{
		SchemaVersion:         1,
		GeneratedAt:           time.Now().UTC(),
		Request:               raw,
		Domain:                req.PrimaryDomain,
		Confidence:            planner.ConfidenceLow,
		ClarificationNeeded:   true,
		ClarificationQuestion: question,
		Evidence:              []planner.Evidence{},
		Conflicts:             []planner.Conflict{},
		Prerequisites:         []string{},
		CostHints:             []planner.CostHint{},
	}
	payload := struct {
		planner.Metadata
		PlanMD string `json:"plan_md"`
	}{Metadata: meta, PlanMD: ""}

	out, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), string(out))
	return nil
}

func defaultOutputPath(repo string, request string, now time.Time) string {
	return filepath.Join(repo, "docs", "infra-planning-cli", now.Format("2006-01-02")+"-"+slug(request)+".md")
}

func slug(value string) string {
	lower := strings.ToLower(value)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	s := strings.Trim(re.ReplaceAllString(lower, "-"), "-")
	if s == "" {
		return "plan"
	}
	if len(s) > 60 {
		s = strings.Trim(s[:60], "-")
	}
	return s
}
```

Note: this imports `internal/web` which does not exist yet — the package will be created in Task 12. The CLI tests above pass `--no-web` so that `fetcher` is `nil` and no HTTP calls are made during testing. **Do not run the CLI tests until Task 12 is complete.**

- [ ] **Step 4: Run all non-CLI tests to confirm no regressions**

```bash
go test ./internal/planner ./internal/render ./internal/discovery ./internal/docs -v
```

Expected: PASS for every package.

- [ ] **Step 5: Commit**

```bash
git add internal/cli/plan.go internal/cli/plan_test.go
git commit -m "feat: --json, --no-web, and per-domain CLI orchestration"
```

---

### Task 10: Update README

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Replace README with final content**

Replace `README.md` with:

```markdown
# infra-planning-cli

An agent-facing AWS infrastructure planning CLI for Claude Code, Codex, and similar coding agents.

Turns a vague infra request into a single, detailed execution plan — repo-aware, correctness-first, plan-only.

---

## What it does

1. Accepts a natural-language infra request.
2. Inspects the current repo for existing IaC patterns (Terraform, CloudFormation, SAM, CDK, Serverless Framework).
3. Fetches the relevant official AWS docs pages (unless `--no-web` is passed) and attaches HTTP status + timestamp to each evidence record.
4. Detects whether the requested resource already exists in the repo (already-satisfied check) and prepends a confirmation step.
5. Detects multi-IaC conflicts (e.g. Terraform + CDK in the same repo) and caps confidence accordingly.
6. Asks exactly one clarifying question only when blocked.
7. Writes a clean markdown execution plan (with YAML frontmatter) into the repo — or emits the same data as JSON to stdout.

It never executes Terraform, calls AWS APIs, or asks for credentials.

---

## Installation

```bash
git clone https://github.com/ShubhanYenuganti/infra-planning-cli
cd infra-planning-cli
go build ./cmd/infra-plan
```

Requires Go 1.24+.

---

## Usage

Default (writes a markdown file):

```bash
infra-plan plan "provision private networking for a new service" --repo /path/to/repo
```

The plan is written to:

```
/path/to/repo/docs/infra-planning-cli/YYYY-MM-DD-<slug>.md
```

JSON mode (no file written, stdout payload for agent pipelines):

```bash
infra-plan plan "provision private networking" --repo /path/to/repo --json
```

Skip HTTP fetches (offline / test mode):

```bash
infra-plan plan "provision private networking" --repo /path/to/repo --no-web
```

`--json` and `--out` are mutually exclusive.

---

## Output schema (frontmatter and JSON)

The markdown file ships with YAML frontmatter; `--json` emits the same shape plus a `plan_md` field carrying the rendered body.

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
  - source: repo
    kind: aws-iac
    path: main.tf
    snippet: 'resource "aws_vpc" "main"'
  - source: docs
    kind: official-docs
    title: "Amazon VPC documentation"
    url: "https://docs.aws.amazon.com/vpc/"
    status: 200
    fetched_at: 2026-05-19T01:23:46Z
conflicts: []
prerequisites:
  - "terraform >= 1.5"
cost_hints:
  - resource_type: aws_nat_gateway
    monthly_usd_estimate: 32.40
    notes: "estimate; verify against current AWS pricing"
  - resource_type: "aws_vpc_endpoint (Interface)"
    monthly_usd_estimate: 7.20
    notes: "estimate; verify against current AWS pricing"
already_satisfied: false
---
```

| Confidence | When |
|---|---|
| `high` | Networking request + matching repo evidence + no IaC conflicts |
| `medium` | Networking request without repo evidence, or any request with multi-IaC conflicts |
| `low` | Compute or database request (stub builders), or domain unknown |

---

## Discovery

Five IaC detectors run against the target repo. Each emits structured `evidence` records:

| Detector | Matches | Triggers on |
|---|---|---|
| `terraform` | `*.tf`, `*.tfvars` | `aws_` or `provider "aws"` |
| `cloudformation` | `*.yaml`, `*.yml`, `*.json` | `AWSTemplateFormatVersion` |
| `sam` | `template.yaml` / `.yml` / `.json` | `Transform: AWS::Serverless` |
| `cdk` | `cdk.json` | file presence |
| `serverless` | `serverless.yml` / `.yaml` | `provider: aws` |

Skipped directories: `.git`, `node_modules`, `vendor`, `dist`, `build`, `.terraform`, `cdk.out`. Files over 1 MiB and symlinks are skipped.

When more than one detector fires, `conflicts` is populated and confidence is capped at `medium`.

---

## Supported AWS domains

| Domain | Builder | Example requests |
|---|---|---|
| `networking` | Full hero plan | VPC, subnets, security groups, NAT gateway |
| `compute` | Minimal stub (low confidence) | EC2, ECS, EKS, Lambda |
| `database` | Minimal stub (low confidence) | RDS, DynamoDB, Postgres, MySQL |

If the request cannot be classified, the CLI prints *"Which AWS area should this plan focus on first: networking, compute, or database?"* and exits 0. Under `--json`, the same question comes back as `clarification_needed: true` with an empty `plan_md`.

---

## Output format

Each plan has four mandatory sections after the frontmatter:

```markdown
# <Title>

## Problem summary
## Recommended path
## Steps
## Verification
```

When `already_satisfied` is true, the first step is a confirmation prompt asking the human to verify the existing resource before generating new IaC.

---

## Project layout

```
cmd/infra-plan/          binary entrypoint
internal/
  cli/                   Cobra root + plan command (--json, --no-web)
  discovery/             pluggable detectors (terraform, cloudformation, sam, cdk, serverless)
  planner/               request, clarification, plan model, dispatch, cost hints, satisfaction, conflicts
  render/                markdown (with frontmatter) and JSON renderers
  docs/                  curated AWS docs URL map
  web/                   HTTP fetcher with allowlist and retry
testdata/repos/          fixture repos per detector (including aws-mixed for conflict tests)
docs/                    specs (01–04), design (06), this implementation plan (05)
```

---

## Running tests

```bash
go test ./...
```

---

## Design constraints

- **Plan-only.** Never provisions resources or executes cloud commands.
- **No credentials.** Never asks for AWS keys, tokens, or secrets.
- **Repo-first.** Prefers existing patterns in the target repo over generic advice.
- **Conservative.** Best possible plan over fastest possible plan.
- **Minimal interruption.** One clarifying question maximum, only when truly blocked.
- **AWS-first.** Five AWS-native IaC detectors today.
```

- [ ] **Step 2: Run docs smoke check**

```bash
grep -q "schema_version" README.md && grep -q -- "--json" README.md && grep -q "cost_hints" README.md
```

Expected: command exits 0.

- [ ] **Step 3: Commit**

```bash
git add README.md
git commit -m "docs: update README for --json, --no-web, cost hints, and conflict detection"
```

---

### Task 12: Web Evidence Fetcher

**Files:**
- Create: `internal/web/fetcher.go`
- Create: `internal/web/fetcher_test.go`
- Modify: `internal/planner/dispatch.go` (update `BuildPlan` and `combineEvidence` to accept `*web.Fetcher`)
- Modify: `internal/planner/dispatch_test.go` (pass `nil` as third argument)
- Run: `go test ./internal/cli -v` (finally runnable now that `internal/web` exists)

- [ ] **Step 1: Write fetcher tests**

Create `internal/web/fetcher_test.go`:

```go
package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetcherAllowlistRejectsNonAWS(t *testing.T) {
	f := New()
	_, _, err := f.Fetch("https://example.com/vpc")
	if err == nil {
		t.Fatalf("expected error for non-allowlist host")
	}
}

func TestFetcherCacheHitSkipsHTTP(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f := New()
	const cachedURL = "https://docs.aws.amazon.com/vpc/"
	// Pre-seed the cache to avoid hitting the real network.
	f.cache[cachedURL] = result{status: 200, fetchedAt: time.Now()}

	status, _, err := f.Fetch(cachedURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != 200 {
		t.Fatalf("expected 200, got %d", status)
	}
	if callCount != 0 {
		t.Fatalf("expected no HTTP calls for cache hit, got %d", callCount)
	}
}

func TestFetcherRetriesOn5xx(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f := &Fetcher{client: srv.Client(), cache: make(map[string]result)}
	status, _, err := f.fetchWithRetry(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected 200 after retry, got %d", status)
	}
	if callCount != 2 {
		t.Fatalf("expected 2 calls (1 fail + 1 retry), got %d", callCount)
	}
}

func TestFetcherSetsUserAgent(t *testing.T) {
	var gotUA string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f := &Fetcher{client: srv.Client(), cache: make(map[string]result)}
	_, _, err := f.fetchWithRetry(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotUA != "infra-planning-cli/0.1" {
		t.Fatalf("User-Agent = %q, want infra-planning-cli/0.1", gotUA)
	}
}
```

- [ ] **Step 2: Run fetcher tests and verify failure**

```bash
go test ./internal/web -v
```

Expected: FAIL because the package is undefined.

- [ ] **Step 3: Implement the fetcher**

Create `internal/web/fetcher.go`:

```go
package web

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const allowedHost = "docs.aws.amazon.com"

type result struct {
	status    int
	fetchedAt time.Time
}

// Fetcher makes HEAD/GET requests to docs.aws.amazon.com with a per-run
// in-memory cache, a 5-second timeout, and one retry on 5xx responses.
type Fetcher struct {
	client *http.Client
	cache  map[string]result
}

// New returns a ready-to-use Fetcher.
func New() *Fetcher {
	return &Fetcher{
		client: &http.Client{Timeout: 5 * time.Second},
		cache:  make(map[string]result),
	}
}

// Fetch returns (statusCode, fetchedAt, nil) for allowlisted URLs.
// Non-allowlist URLs return an error without making any HTTP call.
// Results are cached for the lifetime of the Fetcher.
func (f *Fetcher) Fetch(rawURL string) (int, time.Time, error) {
	u, err := url.Parse(rawURL)
	if err != nil || u.Host != allowedHost {
		return 0, time.Time{}, fmt.Errorf("url not in allowlist: %s", rawURL)
	}
	if cached, ok := f.cache[rawURL]; ok {
		return cached.status, cached.fetchedAt, nil
	}
	status, err := f.fetchWithRetry(rawURL)
	if err != nil {
		return 0, time.Time{}, err
	}
	r := result{status: status, fetchedAt: time.Now().UTC()}
	f.cache[rawURL] = r
	return r.status, r.fetchedAt, nil
}

func (f *Fetcher) fetchWithRetry(rawURL string) (int, error) {
	do := func() (int, error) {
		req, err := http.NewRequest(http.MethodGet, rawURL, nil)
		if err != nil {
			return 0, err
		}
		req.Header.Set("User-Agent", "infra-planning-cli/0.1")
		resp, err := f.client.Do(req)
		if err != nil {
			return 0, err
		}
		resp.Body.Close()
		return resp.StatusCode, nil
	}

	status, err := do()
	if err != nil {
		return 0, err
	}
	if status >= 500 {
		return do() // one retry on 5xx
	}
	return status, nil
}
```

- [ ] **Step 4: Run fetcher tests and verify pass**

```bash
go test ./internal/web -v
```

Expected: PASS.

- [ ] **Step 5: Update dispatch.go to accept a Fetcher**

Replace `internal/planner/dispatch.go` with:

```go
package planner

import (
	"fmt"
	"time"

	"github.com/ShubhanYenuganti/infra-planning-cli/internal/discovery"
	"github.com/ShubhanYenuganti/infra-planning-cli/internal/docs"
	"github.com/ShubhanYenuganti/infra-planning-cli/internal/web"
)

// BuildPlan constructs a domain-specific Plan from the classified request and
// discovered repo context. Pass a non-nil fetcher to enrich docs evidence with
// HTTP status and fetch timestamp; pass nil to skip all network calls.
func BuildPlan(req Request, repo discovery.RepoContext, fetcher *web.Fetcher) Plan {
	now := time.Now().UTC()

	conflicts := DetectConflicts(repo)
	satisfied, satisfiedPath := DetectSatisfaction(req, repo)

	var plan Plan
	switch req.PrimaryDomain {
	case DomainNetworking:
		plan = buildNetworkingPlan(req, repo)
	case DomainCompute:
		plan = buildComputePlan(req, repo)
	case DomainDatabase:
		plan = buildDatabasePlan(req, repo)
	default:
		plan = buildNetworkingPlan(req, repo)
	}

	if satisfied {
		plan.Steps = append([]Step{{
			Title: "Confirm whether existing resource already satisfies this request",
			Body:  fmt.Sprintf("Detected matching resource at `%s`. Verify with the human before generating new IaC.", satisfiedPath),
		}}, plan.Steps...)
	}

	plan.Metadata.SchemaVersion = 1
	plan.Metadata.GeneratedAt = now
	plan.Metadata.Request = req.Raw
	plan.Metadata.Domain = req.PrimaryDomain
	plan.Metadata.Confidence = computeConfidence(req.PrimaryDomain, repo.Evidence, conflicts)
	plan.Metadata.Conflicts = conflicts
	plan.Metadata.Evidence = combineEvidence(repo.Evidence, docs.AWSDocsForDomain(string(req.PrimaryDomain)), fetcher)
	plan.Metadata.CostHints = CostHintsForDomain(req.PrimaryDomain)
	plan.Metadata.AlreadySatisfied = satisfied
	if plan.Metadata.Prerequisites == nil {
		plan.Metadata.Prerequisites = []string{}
	}
	return plan
}

func computeConfidence(domain Domain, repoEv []discovery.Evidence, conflicts []Conflict) Confidence {
	if len(conflicts) > 0 {
		if domain == DomainNetworking && len(repoEv) > 0 {
			return ConfidenceMedium
		}
		return ConfidenceLow
	}
	if domain == DomainCompute || domain == DomainDatabase {
		return ConfidenceLow
	}
	if domain == DomainNetworking {
		if len(repoEv) > 0 {
			return ConfidenceHigh
		}
		return ConfidenceMedium
	}
	return ConfidenceLow
}

func combineEvidence(repoEv []discovery.Evidence, docLinks []docs.DocLink, fetcher *web.Fetcher) []Evidence {
	out := make([]Evidence, 0, len(repoEv)+len(docLinks))
	for _, ev := range repoEv {
		out = append(out, Evidence{
			Source:  "repo",
			Kind:    ev.Kind,
			Path:    ev.Path,
			Snippet: ev.Snippet,
		})
	}
	for _, link := range docLinks {
		ev := Evidence{
			Source: "docs",
			Kind:   "official-docs",
			Title:  link.Title,
			URL:    link.URL,
		}
		if fetcher != nil {
			if status, fetchedAt, err := fetcher.Fetch(link.URL); err == nil {
				ev.Status = status
				ev.FetchedAt = fetchedAt
			}
		}
		out = append(out, ev)
	}
	return out
}

func buildNetworkingPlan(req Request, repo discovery.RepoContext) Plan {
	evidenceLine := "No existing AWS IaC evidence was found. Start by adding a minimal Terraform networking module."
	if len(repo.Evidence) > 0 {
		ev := repo.Evidence[0]
		evidenceLine = fmt.Sprintf("Use existing repo evidence from `%s`: `%s`.", ev.Path, ev.Snippet)
	}
	return Plan{
		Title:           "Provision Private Networking",
		ProblemSummary:  fmt.Sprintf("Request: %s", req.Raw),
		RecommendedPath: "Prefer the repository's existing Terraform/AWS patterns, then add or extend VPC, private subnet, route table, and security group definitions in small reviewed changes.",
		Steps: []Step{
			{Title: "Inspect existing networking IaC", Body: evidenceLine},
			{Title: "Add or extend VPC resources", Body: "Create the smallest Terraform change that defines the required VPC CIDR, private subnets, route tables, and security groups."},
			{Title: "Review blast radius", Body: "Run `terraform plan` and verify that only expected networking resources are created or changed."},
		},
		Verification: []string{
			"`terraform fmt -check -recursive` exits 0",
			"`terraform validate` exits 0",
			"`terraform plan` shows only expected networking changes",
		},
		Metadata: Metadata{Prerequisites: []string{"terraform >= 1.5"}},
	}
}

func buildComputePlan(req Request, repo discovery.RepoContext) Plan {
	return Plan{
		Title:           "Deploy Compute Workload",
		ProblemSummary:  fmt.Sprintf("Request: %s", req.Raw),
		RecommendedPath: "Identify the appropriate compute primitive (EC2, ECS, EKS, or Lambda) from the request, then add the smallest viable resource definition alongside a least-privilege IAM role.",
		Steps: []Step{
			{Title: "Identify compute primitive", Body: "From the request, pick exactly one of EC2, ECS, EKS, or Lambda. If the request is ambiguous, stop and ask the human before generating IaC."},
			{Title: "Inspect existing compute IaC and IAM patterns", Body: "Search the repo for existing compute resources, task definitions, and IAM roles. Prefer extending established patterns over introducing a new one."},
			{Title: "Add the smallest viable resource definition with a least-privilege IAM role", Body: "Create one resource (function/service/task) and one IAM role. Grant only the permissions required for the smoke test."},
		},
		Verification: []string{
			"`terraform validate` exits 0",
			"`terraform plan` shows only the expected compute resources and one IAM role",
		},
		Metadata: Metadata{Prerequisites: []string{"terraform >= 1.5", "AWS account access (configured outside this CLI)"}},
	}
}

func buildDatabasePlan(req Request, repo discovery.RepoContext) Plan {
	return Plan{
		Title:           "Provision Managed Database",
		ProblemSummary:  fmt.Sprintf("Request: %s", req.Raw),
		RecommendedPath: "Identify the database engine (RDS Postgres, RDS MySQL, Aurora, or DynamoDB), then add the smallest viable instance plus parameter group and security group attachment.",
		Steps: []Step{
			{Title: "Identify engine", Body: "From the request, pick exactly one of RDS Postgres, RDS MySQL, Aurora, or DynamoDB. If ambiguous, stop and ask the human before generating IaC."},
			{Title: "Inspect existing database IaC, subnet groups, and secrets handling", Body: "Look for existing database resources, DB subnet groups, and how secrets/passwords are managed (Secrets Manager, SSM Parameter Store, etc.). Reuse patterns when present."},
			{Title: "Add the smallest viable instance + parameter group + security group attachment", Body: "Create one instance with minimal storage, a parameter group, and a security group restricted to the application's VPC CIDR."},
		},
		Verification: []string{
			"`terraform validate` exits 0",
			"`terraform plan` shows only the expected database resources and supporting groups",
		},
		Metadata: Metadata{Prerequisites: []string{"terraform >= 1.5"}},
	}
}
```

- [ ] **Step 6: Update dispatch_test.go to pass nil fetcher**

Replace `internal/planner/dispatch_test.go` with:

```go
package planner

import (
	"strings"
	"testing"

	"github.com/ShubhanYenuganti/infra-planning-cli/internal/discovery"
)

func TestBuildPlanNetworkingHighConfidence(t *testing.T) {
	req := ClassifyRequest("provision a vpc")
	repo := discovery.RepoContext{
		Evidence:      []discovery.Evidence{{Path: "main.tf", Kind: "aws-iac", Snippet: `resource "aws_vpc"`}},
		DetectedTools: []string{"terraform"},
	}
	plan := BuildPlan(req, repo, nil)
	if plan.Metadata.Domain != DomainNetworking {
		t.Fatalf("domain = %q", plan.Metadata.Domain)
	}
	if plan.Metadata.Confidence != ConfidenceHigh {
		t.Fatalf("confidence = %q, want high", plan.Metadata.Confidence)
	}
	if !strings.Contains(plan.Title, "Networking") {
		t.Fatalf("title = %q", plan.Title)
	}
	if len(plan.Steps) == 0 {
		t.Fatalf("no steps generated")
	}
}

func TestBuildPlanNetworkingMediumWithoutEvidence(t *testing.T) {
	req := ClassifyRequest("provision a vpc")
	plan := BuildPlan(req, discovery.RepoContext{}, nil)
	if plan.Metadata.Confidence != ConfidenceMedium {
		t.Fatalf("confidence = %q, want medium", plan.Metadata.Confidence)
	}
}

func TestBuildPlanComputeLow(t *testing.T) {
	req := ClassifyRequest("deploy a lambda worker")
	plan := BuildPlan(req, discovery.RepoContext{}, nil)
	if plan.Metadata.Domain != DomainCompute {
		t.Fatalf("domain = %q", plan.Metadata.Domain)
	}
	if plan.Metadata.Confidence != ConfidenceLow {
		t.Fatalf("confidence = %q, want low", plan.Metadata.Confidence)
	}
	if !strings.Contains(plan.Title, "Compute") {
		t.Fatalf("title = %q", plan.Title)
	}
}

func TestBuildPlanDatabaseLow(t *testing.T) {
	req := ClassifyRequest("create a postgres rds database")
	plan := BuildPlan(req, discovery.RepoContext{}, nil)
	if plan.Metadata.Domain != DomainDatabase {
		t.Fatalf("domain = %q", plan.Metadata.Domain)
	}
	if plan.Metadata.Confidence != ConfidenceLow {
		t.Fatalf("confidence = %q, want low", plan.Metadata.Confidence)
	}
}

func TestBuildPlanConflictsNonNilSlice(t *testing.T) {
	req := ClassifyRequest("provision a vpc")
	plan := BuildPlan(req, discovery.RepoContext{}, nil)
	if plan.Metadata.Conflicts == nil {
		t.Fatalf("conflicts must be non-nil empty slice for stable wire shape")
	}
}

func TestBuildPlanCostHintsNonNil(t *testing.T) {
	req := ClassifyRequest("provision a vpc")
	plan := BuildPlan(req, discovery.RepoContext{}, nil)
	if plan.Metadata.CostHints == nil {
		t.Fatalf("cost_hints must be non-nil empty slice for stable wire shape")
	}
}

func TestBuildPlanAlreadySatisfiedDefaultFalse(t *testing.T) {
	req := ClassifyRequest("provision a vpc")
	plan := BuildPlan(req, discovery.RepoContext{}, nil)
	if plan.Metadata.AlreadySatisfied {
		t.Fatalf("already_satisfied must be false when no matching snippet")
	}
}
```

- [ ] **Step 7: Run all tests**

```bash
go test ./...
```

Expected: PASS for every package including `internal/cli` (which previously could not compile).

- [ ] **Step 8: Commit**

```bash
git add internal/web/fetcher.go internal/web/fetcher_test.go \
        internal/planner/dispatch.go internal/planner/dispatch_test.go
git commit -m "feat: web evidence fetcher with allowlist, cache, and retry"
```

---

### Task 13: Implement CostHintsForDomain

**Files:**
- Create: `internal/planner/costs_test.go`
- Modify: `internal/planner/costs.go` (replace stub)

- [ ] **Step 1: Write cost-hints tests**

Create `internal/planner/costs_test.go`:

```go
package planner

import "testing"

func TestCostHintsForDomainNetworking(t *testing.T) {
	hints := CostHintsForDomain(DomainNetworking)
	if len(hints) != 2 {
		t.Fatalf("expected 2 networking hints, got %d", len(hints))
	}
	found := false
	for _, h := range hints {
		if h.ResourceType == "aws_nat_gateway" {
			found = true
			if h.MonthlyUSDEstimate != 32.40 {
				t.Fatalf("nat gateway estimate = %.2f, want 32.40", h.MonthlyUSDEstimate)
			}
		}
	}
	if !found {
		t.Fatalf("expected aws_nat_gateway hint in networking domain")
	}
}

func TestCostHintsForDomainCompute(t *testing.T) {
	hints := CostHintsForDomain(DomainCompute)
	if len(hints) != 2 {
		t.Fatalf("expected 2 compute hints, got %d", len(hints))
	}
}

func TestCostHintsForDomainDatabase(t *testing.T) {
	hints := CostHintsForDomain(DomainDatabase)
	if len(hints) != 2 {
		t.Fatalf("expected 2 database hints, got %d", len(hints))
	}
}

func TestCostHintsAllHaveNotes(t *testing.T) {
	for _, domain := range []Domain{DomainNetworking, DomainCompute, DomainDatabase} {
		for _, h := range CostHintsForDomain(domain) {
			if h.Notes == "" {
				t.Fatalf("domain %s: hint %q missing notes", domain, h.ResourceType)
			}
		}
	}
}
```

- [ ] **Step 2: Run and verify failure**

```bash
go test ./internal/planner -run TestCostHints -v
```

Expected: FAIL — `TestCostHintsForDomainNetworking` fails because the stub returns an empty slice (`len == 0`, not 2).

- [ ] **Step 3: Replace the stub**

Replace `internal/planner/costs.go` with:

```go
package planner

// CostHintsForDomain returns seed monthly cost estimates for the most common
// resources in a domain. All figures are estimates; users must verify against
// current AWS pricing before committing to a budget.
func CostHintsForDomain(domain Domain) []CostHint {
	notes := "estimate; verify against current AWS pricing"
	switch domain {
	case DomainNetworking:
		return []CostHint{
			{ResourceType: "aws_nat_gateway", MonthlyUSDEstimate: 32.40, Notes: notes},
			{ResourceType: "aws_vpc_endpoint (Interface)", MonthlyUSDEstimate: 7.20, Notes: notes},
		}
	case DomainCompute:
		return []CostHint{
			{ResourceType: "aws_lambda_function (1M req/mo)", MonthlyUSDEstimate: 0.20, Notes: notes},
			{ResourceType: "aws_ecs_fargate_task (0.25 vCPU 24/7)", MonthlyUSDEstimate: 8.90, Notes: notes},
		}
	case DomainDatabase:
		return []CostHint{
			{ResourceType: "aws_db_instance (db.t3.micro)", MonthlyUSDEstimate: 12.40, Notes: notes},
			{ResourceType: "aws_dynamodb_table (on-demand baseline)", MonthlyUSDEstimate: 0.00, Notes: notes},
		}
	default:
		return []CostHint{}
	}
}
```

- [ ] **Step 4: Run and verify pass**

```bash
go test ./internal/planner -v
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/planner/costs.go internal/planner/costs_test.go
git commit -m "feat: cost hints seed values per domain"
```

---

### Task 14: Implement DetectSatisfaction

**Files:**
- Create: `internal/planner/satisfaction_test.go`
- Modify: `internal/planner/satisfaction.go` (replace stub)

- [ ] **Step 1: Write satisfaction tests**

Create `internal/planner/satisfaction_test.go`:

```go
package planner

import (
	"strings"
	"testing"

	"github.com/ShubhanYenuganti/infra-planning-cli/internal/discovery"
)

func TestDetectSatisfactionNetworkingMatch(t *testing.T) {
	req := ClassifyRequest("provision a vpc")
	repo := discovery.RepoContext{
		Evidence: []discovery.Evidence{{
			Path:    "main.tf",
			Kind:    "aws-iac",
			Snippet: `resource "aws_vpc" "main" { cidr_block = "10.0.0.0/16" }`,
		}},
	}
	satisfied, path := DetectSatisfaction(req, repo)
	if !satisfied {
		t.Fatalf("expected satisfaction detected for aws_vpc snippet")
	}
	if path != "main.tf" {
		t.Fatalf("path = %q, want main.tf", path)
	}
}

func TestDetectSatisfactionNoMatch(t *testing.T) {
	req := ClassifyRequest("provision a vpc")
	repo := discovery.RepoContext{
		Evidence: []discovery.Evidence{{
			Path:    "main.tf",
			Kind:    "aws-iac",
			Snippet: `resource "aws_s3_bucket" "data" {}`,
		}},
	}
	satisfied, _ := DetectSatisfaction(req, repo)
	if satisfied {
		t.Fatalf("expected no satisfaction for unrelated snippet")
	}
}

func TestDetectSatisfactionComputeMatch(t *testing.T) {
	req := ClassifyRequest("deploy a lambda worker")
	repo := discovery.RepoContext{
		Evidence: []discovery.Evidence{{
			Path:    "serverless.yml",
			Kind:    "aws-iac",
			Snippet: `AWS::Serverless::Function`,
		}},
	}
	satisfied, path := DetectSatisfaction(req, repo)
	if !satisfied {
		t.Fatalf("expected satisfaction for Serverless::Function")
	}
	if path != "serverless.yml" {
		t.Fatalf("path = %q, want serverless.yml", path)
	}
}

func TestBuildPlanPrependsConfirmationStepWhenSatisfied(t *testing.T) {
	req := ClassifyRequest("provision a vpc")
	repo := discovery.RepoContext{
		Evidence: []discovery.Evidence{{
			Path:    "main.tf",
			Kind:    "aws-iac",
			Snippet: `resource "aws_vpc" "main"`,
		}},
		DetectedTools: []string{"terraform"},
	}
	plan := BuildPlan(req, repo, nil)
	if len(plan.Steps) == 0 {
		t.Fatalf("expected steps")
	}
	if !strings.Contains(plan.Steps[0].Title, "Confirm whether existing resource") {
		t.Fatalf("first step = %q, want confirmation step", plan.Steps[0].Title)
	}
	if !plan.Metadata.AlreadySatisfied {
		t.Fatalf("expected AlreadySatisfied = true")
	}
}
```

- [ ] **Step 2: Run and verify failure**

```bash
go test ./internal/planner -run TestDetectSatisfaction -v
go test ./internal/planner -run TestBuildPlanPrependsConfirmation -v
```

Expected: `TestDetectSatisfactionNetworkingMatch` and `TestDetectSatisfactionComputeMatch` FAIL — the stub returns `(false, "")`.

- [ ] **Step 3: Replace the stub**

Replace `internal/planner/satisfaction.go` with:

```go
package planner

import (
	"strings"

	"github.com/ShubhanYenuganti/infra-planning-cli/internal/discovery"
)

var satisfactionKeywords = map[Domain][]string{
	DomainNetworking: {"aws_vpc", "AWS::EC2::VPC"},
	DomainCompute: {
		"aws_lambda_function", "aws_ecs_service", "aws_instance",
		"AWS::Serverless::Function", "AWS::Lambda::Function",
	},
	DomainDatabase: {
		"aws_db_instance", "aws_rds_cluster", "aws_dynamodb_table",
		"AWS::RDS::DBInstance", "AWS::DynamoDB::Table",
	},
}

// DetectSatisfaction returns (true, evidencePath) when any existing repo
// snippet contains a keyword associated with the request domain.
func DetectSatisfaction(req Request, repo discovery.RepoContext) (bool, string) {
	keywords, ok := satisfactionKeywords[req.PrimaryDomain]
	if !ok {
		return false, ""
	}
	for _, ev := range repo.Evidence {
		for _, kw := range keywords {
			if strings.Contains(ev.Snippet, kw) {
				return true, ev.Path
			}
		}
	}
	return false, ""
}
```

- [ ] **Step 4: Run and verify pass**

```bash
go test ./internal/planner -v
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/planner/satisfaction.go internal/planner/satisfaction_test.go
git commit -m "feat: already-satisfied detection with confirmation step"
```

---

### Task 15: Implement DetectConflicts + Mixed-IaC Fixture

**Files:**
- Create: `testdata/repos/aws-mixed/main.tf`
- Create: `testdata/repos/aws-mixed/cdk.json`
- Create: `internal/planner/conflicts_test.go`
- Modify: `internal/planner/conflicts.go` (replace stub)

- [ ] **Step 1: Create the mixed-IaC fixture**

Create `testdata/repos/aws-mixed/main.tf`:

```hcl
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}
```

Create `testdata/repos/aws-mixed/cdk.json`:

```json
{
  "app": "npx ts-node --prefer-ts-exts bin/app.ts",
  "outputsFile": "outputs.json"
}
```

- [ ] **Step 2: Write conflict tests**

Create `internal/planner/conflicts_test.go`:

```go
package planner

import (
	"strings"
	"testing"

	"github.com/ShubhanYenuganti/infra-planning-cli/internal/discovery"
)

func TestDetectConflictsNoConflict(t *testing.T) {
	repo := discovery.RepoContext{DetectedTools: []string{"terraform"}}
	conflicts := DetectConflicts(repo)
	if len(conflicts) != 0 {
		t.Fatalf("expected no conflicts for single tool, got %+v", conflicts)
	}
}

func TestDetectConflictsEmptyTools(t *testing.T) {
	repo := discovery.RepoContext{}
	conflicts := DetectConflicts(repo)
	if len(conflicts) != 0 {
		t.Fatalf("expected no conflicts for empty tools, got %+v", conflicts)
	}
}

func TestDetectConflictsMultipleTools(t *testing.T) {
	repo := discovery.RepoContext{DetectedTools: []string{"terraform", "cdk"}}
	conflicts := DetectConflicts(repo)
	if len(conflicts) != 1 {
		t.Fatalf("expected 1 conflict for two tools, got %d", len(conflicts))
	}
	if !strings.Contains(conflicts[0].Description, "terraform") {
		t.Fatalf("description missing 'terraform': %q", conflicts[0].Description)
	}
	if !strings.Contains(conflicts[0].Description, "cdk") {
		t.Fatalf("description missing 'cdk': %q", conflicts[0].Description)
	}
	if conflicts[0].ResolutionQuestion == "" {
		t.Fatalf("conflict missing resolution question")
	}
	if len(conflicts[0].Options) != 3 {
		t.Fatalf("expected 3 options, got %d", len(conflicts[0].Options))
	}
}

func TestBuildPlanConfidenceCappedWithConflicts(t *testing.T) {
	req := ClassifyRequest("provision a vpc")
	repo := discovery.RepoContext{
		Evidence:      []discovery.Evidence{{Path: "main.tf", Kind: "aws-iac", Snippet: `resource "aws_vpc"`}},
		DetectedTools: []string{"terraform", "cdk"},
	}
	plan := BuildPlan(req, repo, nil)
	// Networking + evidence normally → high, but conflicts cap it to medium.
	if plan.Metadata.Confidence != ConfidenceMedium {
		t.Fatalf("confidence = %q, want medium when conflicts detected", plan.Metadata.Confidence)
	}
	if len(plan.Metadata.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict in metadata, got %d", len(plan.Metadata.Conflicts))
	}
}
```

- [ ] **Step 3: Run and verify failure**

```bash
go test ./internal/planner -run TestDetectConflicts -v
go test ./internal/planner -run TestBuildPlanConfidenceCapped -v
```

Expected: `TestDetectConflictsMultipleTools` and `TestBuildPlanConfidenceCappedWithConflicts` FAIL — the stub returns empty slice.

- [ ] **Step 4: Replace the stub**

Replace `internal/planner/conflicts.go` with:

```go
package planner

import (
	"sort"
	"strings"

	"github.com/ShubhanYenuganti/infra-planning-cli/internal/discovery"
)

// DetectConflicts returns one Conflict record when the repo contains evidence
// from more than one IaC tool. The confidence rule in computeConfidence then
// caps the plan confidence at medium.
func DetectConflicts(repo discovery.RepoContext) []Conflict {
	if len(repo.DetectedTools) <= 1 {
		return []Conflict{}
	}
	tools := make([]string, len(repo.DetectedTools))
	copy(tools, repo.DetectedTools)
	sort.Strings(tools)
	return []Conflict{{
		Description: "Multiple IaC tools detected: " + strings.Join(tools, ", "),
		Options: []string{
			"Consolidate to a single IaC tool before adding new resources.",
			"Ensure each tool manages a distinct resource scope with no overlap.",
			"Document the tool boundary explicitly and proceed with one tool for this change.",
		},
		ResolutionQuestion: "Which IaC tool should be used for this change?",
	}}
}
```

- [ ] **Step 5: Run all tests**

```bash
go test ./...
```

Expected: PASS.

- [ ] **Step 6: Verify the mixed fixture triggers conflict detection**

```bash
go test ./internal/discovery -run TestDiscoverRepo -v
```

Expected: all five single-tool fixtures PASS. (The mixed fixture is exercised via the conflict tests in `internal/planner`.)

- [ ] **Step 7: Commit**

```bash
git add internal/planner/conflicts.go internal/planner/conflicts_test.go \
        testdata/repos/aws-mixed
git commit -m "feat: multi-IaC conflict detection with confidence cap"
```

---

### Task 11: End-to-End CLI Tests

**Files:**
- Create: `tests/e2e/plan_test.go`
- Create: `testdata/repos/aws-satisfied/main.tf`

- [ ] **Step 1: Create the already-satisfied fixture**

Create `testdata/repos/aws-satisfied/main.tf`:

```hcl
terraform {
  required_providers {
    aws = { source = "hashicorp/aws", version = "~> 5.0" }
  }
}

provider "aws" {
  region = "us-east-1"
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
  tags = {
    Name = "main"
  }
}
```

- [ ] **Step 2: Write the E2E test file**

Create `tests/e2e/plan_test.go`:

```go
package e2e_test

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ShubhanYenuganti/infra-planning-cli/internal/planner"
)

var builtBin string

func TestMain(m *testing.M) {
	bin, err := buildBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "build failed: %v\n", err)
		os.Exit(1)
	}
	builtBin = bin
	os.Exit(m.Run())
}

func buildBinary() (string, error) {
	bin := filepath.Join(os.TempDir(), "infra-planning-e2e-bin")
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/infra-plan")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("go build: %v\n%s", err, out)
	}
	return bin, nil
}

// repoPath resolves a testdata fixture path relative to the module root.
// Tests in tests/e2e/ are two directories below the module root.
func repoPath(name string) string {
	return filepath.Join("..", "..", "testdata", "repos", name)
}

// runPlan runs "infra-plan plan --request <req> --repo <repo> [extra...]"
// and returns combined output. It fails the test if the command returns a
// non-zero exit code for any reason other than clarification-needed.
func runPlan(t *testing.T, repo, request string, extra ...string) []byte {
	t.Helper()
	args := []string{"plan", "--repo", repo, "--request", request}
	args = append(args, extra...)
	out, err := exec.Command(builtBin, args...).CombinedOutput()
	if err != nil {
		t.Fatalf("infra-plan plan failed: %v\noutput:\n%s", err, out)
	}
	return out
}

// TestE2ENetworkingMarkdown verifies a networking request against the
// aws-networking fixture produces YAML-frontmatter markdown with the
// correct domain and confidence fields.
func TestE2ENetworkingMarkdown(t *testing.T) {
	out := runPlan(t, repoPath("aws-networking"), "Add a private subnet")
	body := string(out)

	if !strings.HasPrefix(body, "---\n") {
		t.Fatalf("expected YAML frontmatter to start output, got:\n%s", body)
	}
	if !strings.Contains(body, "domain: networking") {
		t.Errorf("expected 'domain: networking' in frontmatter")
	}
	if !strings.Contains(body, "confidence:") {
		t.Errorf("expected 'confidence:' field in frontmatter")
	}
	if !strings.Contains(body, "## Steps") {
		t.Errorf("expected '## Steps' section in plan body")
	}
}

// TestE2EComputeJSONHasCostHints verifies --json output for a compute request
// contains cost_hints with at least two entries.
func TestE2EComputeJSONHasCostHints(t *testing.T) {
	out := runPlan(t, repoPath("aws-compute"), "Add a Lambda function", "--json")

	var plan planner.Plan
	if err := json.Unmarshal(out, &plan); err != nil {
		t.Fatalf("json.Unmarshal failed: %v\noutput was:\n%s", err, out)
	}
	if len(plan.Metadata.CostHints) < 2 {
		t.Errorf("expected >=2 cost hints for compute domain, got %d", len(plan.Metadata.CostHints))
	}
	// Spot-check that the first hint has a non-zero estimate.
	if plan.Metadata.CostHints[0].MonthlyUSDEstimate == 0 {
		t.Errorf("expected non-zero MonthlyUSDEstimate in first cost hint")
	}
}

// TestE2ENoWebFlagSuppressesFetch verifies that --no-web produces evidence
// records with zero FetchedAt timestamps (no HTTP calls were made).
func TestE2ENoWebFlagSuppressesFetch(t *testing.T) {
	out := runPlan(t, repoPath("aws-networking"), "Add a private subnet", "--no-web", "--json")

	var plan planner.Plan
	if err := json.Unmarshal(out, &plan); err != nil {
		t.Fatalf("json.Unmarshal failed: %v\noutput was:\n%s", err, out)
	}
	for _, ev := range plan.Metadata.Evidence {
		if !ev.FetchedAt.IsZero() {
			t.Errorf("--no-web set but evidence %q has FetchedAt=%v", ev.Source, ev.FetchedAt)
		}
	}
}

// TestE2EMixedIaCConflict verifies the aws-mixed fixture (Terraform + CDK)
// produces confidence=medium and exactly one conflict record.
func TestE2EMixedIaCConflict(t *testing.T) {
	out := runPlan(t, repoPath("aws-mixed"), "Add a VPC", "--json")

	var plan planner.Plan
	if err := json.Unmarshal(out, &plan); err != nil {
		t.Fatalf("json.Unmarshal failed: %v\noutput was:\n%s", err, out)
	}
	if plan.Metadata.Confidence != planner.ConfidenceMedium {
		t.Errorf("expected confidence=%q for mixed-IaC repo, got %q",
			planner.ConfidenceMedium, plan.Metadata.Confidence)
	}
	if len(plan.Metadata.Conflicts) != 1 {
		t.Errorf("expected exactly 1 conflict for mixed-IaC repo, got %d",
			len(plan.Metadata.Conflicts))
	}
	if !strings.Contains(plan.Metadata.Conflicts[0].Description, "terraform") {
		t.Errorf("expected conflict description to mention 'terraform', got: %s",
			plan.Metadata.Conflicts[0].Description)
	}
}

// TestE2EAlreadySatisfied verifies a repo that already contains the requested
// resource type produces already_satisfied=true and a confirmation step first.
func TestE2EAlreadySatisfied(t *testing.T) {
	out := runPlan(t, repoPath("aws-satisfied"), "Add a VPC", "--json")

	var plan planner.Plan
	if err := json.Unmarshal(out, &plan); err != nil {
		t.Fatalf("json.Unmarshal failed: %v\noutput was:\n%s", err, out)
	}
	if !plan.Metadata.AlreadySatisfied {
		t.Error("expected already_satisfied=true for aws-satisfied fixture")
	}
	if len(plan.Steps) == 0 {
		t.Fatal("expected at least one step in plan")
	}
	if !strings.Contains(plan.Steps[0].Title, "Confirm") {
		t.Errorf("expected first step title to contain 'Confirm', got: %q",
			plan.Steps[0].Title)
	}
}

// TestE2EJSONAndOutMutuallyExclusive verifies the CLI rejects --json and
// --out used together.
func TestE2EJSONAndOutMutuallyExclusive(t *testing.T) {
	cmd := exec.Command(builtBin, "plan",
		"--repo", repoPath("aws-networking"),
		"--request", "Add a subnet",
		"--json",
		"--out", "/tmp/test-plan.md",
	)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected non-zero exit when --json and --out both set, got output:\n%s", out)
	}
}
```

- [ ] **Step 3: Run tests — verify they fail (binary not yet built in test context)**

```bash
go test ./tests/e2e/... -v 2>&1 | head -30
```

Expected: `FAIL` — `buildBinary` succeeds but at least one sub-test fails because Tasks 12–15 have not been committed yet. If all tasks 1–15 are already complete, the tests should pass instead; skip to step 5.

- [ ] **Step 4: Verify the binary builds cleanly**

```bash
go build ./cmd/infra-plan
```

Expected: exits 0, binary produced in current directory.

- [ ] **Step 5: Run the full E2E suite**

```bash
go test ./tests/e2e/... -v -timeout 60s
```

Expected output (all PASS):

```
--- PASS: TestE2ENetworkingMarkdown
--- PASS: TestE2EComputeJSONHasCostHints
--- PASS: TestE2ENoWebFlagSuppressesFetch
--- PASS: TestE2EMixedIaCConflict
--- PASS: TestE2EAlreadySatisfied
--- PASS: TestE2EJSONAndOutMutuallyExclusive
PASS
ok      github.com/ShubhanYenuganti/infra-planning-cli/tests/e2e
```

- [ ] **Step 6: Run the full unit test suite**

```bash
go test ./...
```

Expected: PASS across all packages.

- [ ] **Step 7: Commit**

```bash
git add tests/e2e/plan_test.go testdata/repos/aws-satisfied
git commit -m "test: end-to-end CLI tests covering --no-web, cost hints, conflict, and already-satisfied"
```

---

## Self-Review

### 1. Spec Coverage

| Spec section | Task |
|---|---|
| §1 Natural-language request → markdown plan | Task 9 (CLI orchestration) |
| §2 Repo discovery (5 IaC detectors, pluggable) | Task 3 |
| §3 AWS docs evidence capture | Task 4 |
| §4 Clarification gate (one question, exit 0) | Task 5 |
| §5 Per-domain plan dispatch (networking hero, compute/db stubs) | Task 7 |
| §6 Frontmatter markdown output | Task 6 |
| §7 JSON output (`--json` flag) | Task 8 |
| §8 CLI contract (flags, exit codes, mutex check) | Task 9 |
| §9 Web evidence fetching (5s timeout, allowlist, 1 retry, `--no-web`) | Task 12 |
| §10.1 Cost hints (`CostHintsForDomain`, seed values) | Tasks 7 (stub), 13 |
| §10.2 Already-satisfied detection (`DetectSatisfaction`, confirmation step) | Tasks 7 (stub), 14 |
| §10.3 Multi-IaC conflict detection (`DetectConflicts`, confidence cap) | Tasks 7 (stub), 15 |
| README documentation | Task 10 |
| E2E tests | Task 11 |

**No gaps found.**

### 2. Placeholder Scan

No "TBD", "TODO", "implement later", "fill in details", "add appropriate error handling", or "similar to Task N" patterns present. All code blocks contain complete, runnable code.

### 3. Type Consistency

| Symbol | Defined in | Used correctly in |
|---|---|---|
| `Evidence.FetchedAt time.Time` | Task 6 (`plan.go`) | Task 12 (`fetcher.go`, `dispatch.go`), Task 11 (E2E test) |
| `Evidence.Status int` | Task 6 (`plan.go`) | Task 12 (`fetcher.go`, `dispatch.go`) |
| `CostHint` struct | Task 6 (`plan.go`) | Task 7 (stub), Task 13 (impl), Task 11 (E2E) |
| `Metadata.CostHints []CostHint` | Task 6 (`plan.go`) | Task 7, 12 (emitClarification stub), Task 13, Task 11 |
| `Metadata.AlreadySatisfied bool` | Task 6 (`plan.go`) | Task 7, Task 14, Task 11 |
| `Metadata.Conflicts []Conflict` | Task 6 (`plan.go`) | Task 7, Task 15, Task 11 |
| `Conflict` struct | Task 7 (`plan.go`) | Task 15, Task 11 |
| `BuildPlan(req, repo, fetcher)` 3-arg | Task 12 (`dispatch.go`) | Task 12 (`dispatch_test.go`, `cli/plan.go`) |
| `BuildPlan(req, repo)` 2-arg | Tasks 7–9 | Replaced entirely in Task 12 |
| `CostHintsForDomain(domain Domain)` | Task 7 stub, Task 13 impl | Task 7 (`dispatch.go`) |
| `DetectSatisfaction(req, repo)` | Task 7 stub, Task 14 impl | Task 7 (`dispatch.go`) |
| `DetectConflicts(repo)` | Task 7 stub, Task 15 impl | Task 7 (`dispatch.go`) |
| `web.Fetcher` / `web.New()` | Task 12 (`fetcher.go`) | Task 12 (`cli/plan.go`, `dispatch.go`) |
| `ConfidenceMedium` | Task 7 (`plan.go`) | Task 15, Task 11 (E2E) |
| `discovery.RepoContext.DetectedTools []string` | Task 3 | Task 15 (`conflicts.go`) |
| `planner.Plan.Steps []Step` | Task 6 | Task 11 (E2E `.Steps[0].Title`) |

**No mismatches found.**

---

## Execution Handoff

Plan complete and saved to `docs/05-implementation-plan.md`.

**Two execution options:**

**1. Subagent-Driven (recommended)** — dispatch a fresh subagent per task, review between tasks, fast iteration. Use `superpowers:subagent-driven-development`.

**2. Inline Execution** — execute tasks in this session using `superpowers:executing-plans`, batch execution with checkpoints.

Which approach?
