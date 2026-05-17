# Infra Planning CLI Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build an AWS-first, agent-facing CLI that turns a vague infra request into one detailed, repo-specific markdown execution plan for Claude Code, Codex, or another coding agent to execute later.

**Architecture:** Use `printing-press-hermes` to scaffold the project, then implement a small Go CLI with a deterministic pipeline: repo discovery, AWS/domain classification, docs evidence capture, clarification gate, and plan rendering. The MVP is plan-only: it never provisions resources, never asks for credentials, and writes a single markdown file under `docs/infra-planning-cli/` in the target repo.

**Tech Stack:** Go 1.24+, Cobra for CLI commands, standard-library filesystem/process utilities, markdown templates, table-driven Go tests, and optional `hermes-press` research/generate/verify flow from `PrintingPressHermes`.

---

## File Structure

Target generated project root: `infra-planning-cli/`

- `go.mod` — module declaration for the standalone CLI.
- `cmd/infra-plan/main.go` — binary entrypoint.
- `internal/cli/root.go` — Cobra root command and flag wiring.
- `internal/cli/plan.go` — `plan` subcommand orchestration.
- `internal/discovery/context.go` — typed repo discovery model.
- `internal/discovery/discover.go` — repo file scanning and evidence extraction.
- `internal/discovery/discover_test.go` — discovery behavior tests.
- `internal/planner/request.go` — normalized request model and AWS domain classification.
- `internal/planner/request_test.go` — request classification tests.
- `internal/planner/clarify.go` — minimum-clarification gate.
- `internal/planner/clarify_test.go` — clarification gate tests.
- `internal/planner/plan.go` — execution plan data model.
- `internal/render/markdown.go` — markdown renderer.
- `internal/render/markdown_test.go` — output format tests.
- `internal/docs/aws.go` — curated AWS docs URL mapping for MVP domains.
- `internal/docs/aws_test.go` — AWS docs mapping tests.
- `testdata/repos/aws-networking/` — tiny fixture repo with Terraform networking examples.
- `docs/infra-planning-cli/` — generated markdown plan output directory in consuming repos.
- `README.md` — install and usage instructions.

## MVP CLI Contract

Command:

```bash
infra-plan plan "provision private networking for a new service" --repo .
```

Required behavior:

1. Inspect only the current repo or the path passed via `--repo`.
2. Classify the request into one or more AWS MVP domains: `networking`, `compute`, `database`, or `unknown`.
3. Prefer repo evidence over docs evidence.
4. Attach official AWS docs URLs as internal evidence, but keep the generated plan citation-free.
5. Ask exactly one minimal clarification question only when blocked.
6. Write one markdown file to `docs/infra-planning-cli/YYYY-MM-DD-<slug>.md`.
7. Never execute cloud commands or ask for credentials.

---

### Task 1: Scaffold the Go CLI Module

**Files:**
- Create: `infra-planning-cli/go.mod`
- Create: `infra-planning-cli/cmd/infra-plan/main.go`
- Create: `infra-planning-cli/internal/cli/root.go`
- Create: `infra-planning-cli/internal/cli/plan.go`

- [ ] **Step 1: Create the module file**

Create `infra-planning-cli/go.mod`:

```go
module github.com/ShubhanYenuganti/infra-planning-cli

go 1.24

require github.com/spf13/cobra v1.8.0
```

- [ ] **Step 2: Create the binary entrypoint**

Create `infra-planning-cli/cmd/infra-plan/main.go`:

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

Create `infra-planning-cli/internal/cli/root.go`:

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

Create `infra-planning-cli/internal/cli/plan.go`:

```go
package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

type planOptions struct {
	repo string
	out  string
}

func newPlanCommand() *cobra.Command {
	opts := &planOptions{}

	cmd := &cobra.Command{
		Use:   "plan REQUEST",
		Short: "Turn an infra request into a markdown execution plan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "request: %s\nrepo: %s\nout: %s\n", args[0], opts.repo, opts.out)
			return nil
		},
	}

	cmd.Flags().StringVar(&opts.repo, "repo", ".", "repository path to inspect")
	cmd.Flags().StringVar(&opts.out, "out", "", "optional output markdown path")
	return cmd
}
```

- [ ] **Step 5: Download dependencies and verify the stub**

Run:

```bash
cd infra-planning-cli
go mod tidy
go run ./cmd/infra-plan plan "provision networking" --repo .
```

Expected output contains:

```text
request: provision networking
repo: .
out:
```

- [ ] **Step 6: Commit**

```bash
git add infra-planning-cli/go.mod infra-planning-cli/go.sum infra-planning-cli/cmd/infra-plan/main.go infra-planning-cli/internal/cli/root.go infra-planning-cli/internal/cli/plan.go
git commit -m "feat: scaffold infra planning CLI"
```

---

### Task 2: Add AWS Request Classification

**Files:**
- Create: `infra-planning-cli/internal/planner/request.go`
- Create: `infra-planning-cli/internal/planner/request_test.go`

- [ ] **Step 1: Write classification tests**

Create `infra-planning-cli/internal/planner/request_test.go`:

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

Run:

```bash
cd infra-planning-cli
go test ./internal/planner -run TestClassifyRequest -v
```

Expected: FAIL because `Domain` and `ClassifyRequest` are undefined.

- [ ] **Step 3: Implement minimal classification**

Create `infra-planning-cli/internal/planner/request.go`:

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

Run:

```bash
cd infra-planning-cli
go test ./internal/planner -run TestClassifyRequest -v
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add infra-planning-cli/internal/planner/request.go infra-planning-cli/internal/planner/request_test.go
git commit -m "feat: classify AWS infra requests"
```

---

### Task 3: Implement Repo Discovery

**Files:**
- Create: `infra-planning-cli/internal/discovery/context.go`
- Create: `infra-planning-cli/internal/discovery/discover.go`
- Create: `infra-planning-cli/internal/discovery/discover_test.go`
- Create: `infra-planning-cli/testdata/repos/aws-networking/main.tf`

- [ ] **Step 1: Create the fixture repo**

Create `infra-planning-cli/testdata/repos/aws-networking/main.tf`:

```hcl
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "private" {
  vpc_id     = aws_vpc.main.id
  cidr_block = "10.0.1.0/24"
}
```

- [ ] **Step 2: Write discovery tests**

Create `infra-planning-cli/internal/discovery/discover_test.go`:

```go
package discovery

import "testing"

func TestDiscoverRepoFindsTerraformAndAWSHints(t *testing.T) {
	ctx, err := DiscoverRepo("../../testdata/repos/aws-networking")
	if err != nil {
		t.Fatalf("DiscoverRepo returned error: %v", err)
	}

	if ctx.Root == "" {
		t.Fatalf("Root was empty")
	}
	if !ctx.HasTerraform {
		t.Fatalf("HasTerraform = false, want true")
	}
	if len(ctx.Evidence) == 0 {
		t.Fatalf("Evidence was empty")
	}
	if ctx.Evidence[0].Path != "main.tf" {
		t.Fatalf("first evidence path = %q, want main.tf", ctx.Evidence[0].Path)
	}
}
```

- [ ] **Step 3: Run the test and verify failure**

Run:

```bash
cd infra-planning-cli
go test ./internal/discovery -run TestDiscoverRepoFindsTerraformAndAWSHints -v
```

Expected: FAIL because `DiscoverRepo` is undefined.

- [ ] **Step 4: Add discovery types**

Create `infra-planning-cli/internal/discovery/context.go`:

```go
package discovery

type Evidence struct {
	Path    string
	Kind    string
	Snippet string
}

type RepoContext struct {
	Root         string
	HasTerraform bool
	HasAWS       bool
	Evidence     []Evidence
}
```

- [ ] **Step 5: Implement minimal repo discovery**

Create `infra-planning-cli/internal/discovery/discover.go`:

```go
package discovery

import (
	"os"
	"path/filepath"
	"strings"
)

func DiscoverRepo(root string) (RepoContext, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return RepoContext{}, err
	}

	ctx := RepoContext{Root: abs}
	err = filepath.WalkDir(abs, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			name := entry.Name()
			if name == ".git" || name == "node_modules" || name == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}

		rel, err := filepath.Rel(abs, path)
		if err != nil {
			return err
		}
		if filepath.Ext(path) != ".tf" && filepath.Ext(path) != ".tfvars" && filepath.Base(path) != "serverless.yml" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		content := string(data)
		if strings.HasSuffix(path, ".tf") || strings.HasSuffix(path, ".tfvars") {
			ctx.HasTerraform = true
		}
		if strings.Contains(content, "aws_") || strings.Contains(strings.ToLower(content), "provider \"aws\"") {
			ctx.HasAWS = true
			ctx.Evidence = append(ctx.Evidence, Evidence{Path: rel, Kind: "aws-iac", Snippet: firstNonEmptyLine(content)})
		}
		return nil
	})
	if err != nil {
		return RepoContext{}, err
	}

	return ctx, nil
}

func firstNonEmptyLine(content string) string {
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}
```

- [ ] **Step 6: Run the test and verify pass**

Run:

```bash
cd infra-planning-cli
go test ./internal/discovery -run TestDiscoverRepoFindsTerraformAndAWSHints -v
```

Expected: PASS.

- [ ] **Step 7: Commit**

```bash
git add infra-planning-cli/internal/discovery infra-planning-cli/testdata/repos/aws-networking/main.tf
git commit -m "feat: discover repo infrastructure context"
```

---

### Task 4: Add Curated AWS Docs Evidence

**Files:**
- Create: `infra-planning-cli/internal/docs/aws.go`
- Create: `infra-planning-cli/internal/docs/aws_test.go`

- [ ] **Step 1: Write docs mapping tests**

Create `infra-planning-cli/internal/docs/aws_test.go`:

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

	unknown := AWSDocsForDomain("unknown")
	if len(unknown) != 0 {
		t.Fatalf("unknown docs length = %d, want 0", len(unknown))
	}
}
```

- [ ] **Step 2: Run the test and verify failure**

Run:

```bash
cd infra-planning-cli
go test ./internal/docs -run TestAWSDocsForDomain -v
```

Expected: FAIL because `AWSDocsForDomain` is undefined.

- [ ] **Step 3: Implement AWS docs mapping**

Create `infra-planning-cli/internal/docs/aws.go`:

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

Run:

```bash
cd infra-planning-cli
go test ./internal/docs -run TestAWSDocsForDomain -v
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add infra-planning-cli/internal/docs/aws.go infra-planning-cli/internal/docs/aws_test.go
git commit -m "feat: add AWS docs evidence map"
```

---

### Task 5: Add Clarification Gate

**Files:**
- Create: `infra-planning-cli/internal/planner/clarify.go`
- Create: `infra-planning-cli/internal/planner/clarify_test.go`

- [ ] **Step 1: Write clarification tests**

Create `infra-planning-cli/internal/planner/clarify_test.go`:

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
			got := ClarificationQuestion(tt.req)
			if got != tt.want {
				t.Fatalf("ClarificationQuestion() = %q, want %q", got, tt.want)
			}
		})
	}
}
```

- [ ] **Step 2: Run the test and verify failure**

Run:

```bash
cd infra-planning-cli
go test ./internal/planner -run TestClarificationQuestion -v
```

Expected: FAIL because `ClarificationQuestion` is undefined.

- [ ] **Step 3: Implement clarification gate**

Create `infra-planning-cli/internal/planner/clarify.go`:

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

Run:

```bash
cd infra-planning-cli
go test ./internal/planner -v
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add infra-planning-cli/internal/planner/clarify.go infra-planning-cli/internal/planner/clarify_test.go
git commit -m "feat: add minimal clarification gate"
```

---

### Task 6: Add Plan Model and Markdown Renderer

**Files:**
- Create: `infra-planning-cli/internal/planner/plan.go`
- Create: `infra-planning-cli/internal/render/markdown.go`
- Create: `infra-planning-cli/internal/render/markdown_test.go`

- [ ] **Step 1: Write renderer tests**

Create `infra-planning-cli/internal/render/markdown_test.go`:

```go
package render

import (
	"strings"
	"testing"

	"github.com/ShubhanYenuganti/infra-planning-cli/internal/planner"
)

func TestMarkdownPlanIncludesMandatorySections(t *testing.T) {
	plan := planner.Plan{
		Title:           "Provision Private Networking",
		ProblemSummary:  "Create private networking for a new AWS service.",
		RecommendedPath: "Use an existing Terraform VPC pattern and add private subnets first.",
		Steps: []planner.Step{
			{Title: "Inspect existing Terraform", Body: "Review current VPC and subnet resources before editing."},
		},
		Verification: []string{"terraform validate exits 0"},
	}

	md := MarkdownPlan(plan)
	for _, section := range []string{"# Provision Private Networking", "## Problem summary", "## Recommended path", "## Steps", "## Verification"} {
		if !strings.Contains(md, section) {
			t.Fatalf("markdown missing section %q:\n%s", section, md)
		}
	}
}
```

- [ ] **Step 2: Run the test and verify failure**

Run:

```bash
cd infra-planning-cli
go test ./internal/render -run TestMarkdownPlanIncludesMandatorySections -v
```

Expected: FAIL because `planner.Plan` and `MarkdownPlan` are undefined.

- [ ] **Step 3: Add plan model**

Create `infra-planning-cli/internal/planner/plan.go`:

```go
package planner

type Plan struct {
	Title           string
	ProblemSummary  string
	RecommendedPath string
	Steps           []Step
	Verification    []string
}

type Step struct {
	Title string
	Body  string
}
```

- [ ] **Step 4: Implement markdown renderer**

Create `infra-planning-cli/internal/render/markdown.go`:

```go
package render

import (
	"fmt"
	"strings"

	"github.com/ShubhanYenuganti/infra-planning-cli/internal/planner"
)

func MarkdownPlan(plan planner.Plan) string {
	var b strings.Builder
	_, _ = fmt.Fprintf(&b, "# %s\n\n", plan.Title)
	_, _ = fmt.Fprintf(&b, "## Problem summary\n\n%s\n\n", plan.ProblemSummary)
	_, _ = fmt.Fprintf(&b, "## Recommended path\n\n%s\n\n", plan.RecommendedPath)
	b.WriteString("## Steps\n\n")
	for i, step := range plan.Steps {
		_, _ = fmt.Fprintf(&b, "%d. **%s**\n\n   %s\n\n", i+1, step.Title, step.Body)
	}
	b.WriteString("## Verification\n\n")
	for _, item := range plan.Verification {
		_, _ = fmt.Fprintf(&b, "- %s\n", item)
	}
	return b.String()
}
```

- [ ] **Step 5: Run renderer tests and verify pass**

Run:

```bash
cd infra-planning-cli
go test ./internal/render -v
```

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add infra-planning-cli/internal/planner/plan.go infra-planning-cli/internal/render
git commit -m "feat: render markdown execution plans"
```

---

### Task 7: Generate the MVP Networking Plan

**Files:**
- Modify: `infra-planning-cli/internal/cli/plan.go`
- Create: `infra-planning-cli/internal/cli/plan_test.go`

- [ ] **Step 1: Write CLI integration test**

Create `infra-planning-cli/internal/cli/plan_test.go`:

```go
package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPlanCommandWritesMarkdownFile(t *testing.T) {
	tmp := t.TempDir()
	repo := filepath.Join(tmp, "repo")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "main.tf"), []byte(`resource "aws_vpc" "main" { cidr_block = "10.0.0.0/16" }`), 0o644); err != nil {
		t.Fatal(err)
	}

	out := filepath.Join(repo, "docs", "infra-planning-cli", "2026-05-15-provision-networking.md")
	cmd := NewRootCommand()
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stdout)
	cmd.SetArgs([]string{"plan", "provision private networking", "--repo", repo, "--out", out})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, stdout.String())
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("expected output file: %v", err)
	}
	md := string(data)
	for _, want := range []string{"## Problem summary", "## Recommended path", "## Steps", "## Verification", "terraform validate"} {
		if !strings.Contains(md, want) {
			t.Fatalf("markdown missing %q:\n%s", want, md)
		}
	}
}
```

- [ ] **Step 2: Run the test and verify failure**

Run:

```bash
cd infra-planning-cli
go test ./internal/cli -run TestPlanCommandWritesMarkdownFile -v
```

Expected: FAIL because the current command only prints a stub.

- [ ] **Step 3: Replace stub command with orchestration**

Replace `infra-planning-cli/internal/cli/plan.go` with:

```go
package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/ShubhanYenuganti/infra-planning-cli/internal/discovery"
	"github.com/ShubhanYenuganti/infra-planning-cli/internal/planner"
	"github.com/ShubhanYenuganti/infra-planning-cli/internal/render"
	"github.com/spf13/cobra"
)

type planOptions struct {
	repo string
	out  string
}

func newPlanCommand() *cobra.Command {
	opts := &planOptions{}

	cmd := &cobra.Command{
		Use:   "plan REQUEST",
		Short: "Turn an infra request into a markdown execution plan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := planner.ClassifyRequest(args[0])
			if question := planner.ClarificationQuestion(request); question != "" {
				return fmt.Errorf(question)
			}

			repoCtx, err := discovery.DiscoverRepo(opts.repo)
			if err != nil {
				return err
			}

			plan := networkingPlan(args[0], repoCtx)
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
	cmd.Flags().StringVar(&opts.out, "out", "", "optional output markdown path")
	return cmd
}

func networkingPlan(raw string, repoCtx discovery.RepoContext) planner.Plan {
	evidence := "No existing AWS IaC evidence was found. Start by adding a minimal Terraform networking module."
	if len(repoCtx.Evidence) > 0 {
		evidence = fmt.Sprintf("Use existing repo evidence from `%s`: `%s`.", repoCtx.Evidence[0].Path, repoCtx.Evidence[0].Snippet)
	}

	return planner.Plan{
		Title:           "Provision Private Networking",
		ProblemSummary:  fmt.Sprintf("Request: %s", raw),
		RecommendedPath: "Prefer the repository's existing Terraform/AWS patterns, then add or extend VPC, private subnet, route table, and security group definitions in small reviewed changes.",
		Steps: []planner.Step{
			{Title: "Inspect existing networking IaC", Body: evidence},
			{Title: "Add or extend VPC resources", Body: "Create the smallest Terraform change that defines the required VPC CIDR, private subnets, route tables, and security groups."},
			{Title: "Review blast radius", Body: "Run `terraform plan` and verify that only expected networking resources are created or changed."},
		},
		Verification: []string{
			"`terraform fmt -check -recursive` exits 0",
			"`terraform validate` exits 0",
			"`terraform plan` shows only expected networking changes",
		},
	}
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

- [ ] **Step 4: Run CLI tests and verify pass**

Run:

```bash
cd infra-planning-cli
go test ./internal/cli -v
```

Expected: PASS.

- [ ] **Step 5: Run all tests**

Run:

```bash
cd infra-planning-cli
go test ./...
```

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add infra-planning-cli/internal/cli/plan.go infra-planning-cli/internal/cli/plan_test.go
git commit -m "feat: generate networking execution plan"
```

---

### Task 8: Add README and Printing Press Handoff

**Files:**
- Create: `infra-planning-cli/README.md`
- Modify: `infra-planning-cli/internal/cli/plan.go`

- [ ] **Step 1: Create README usage docs**

Create `infra-planning-cli/README.md`:

```markdown
# infra-planning-cli

Agent-facing AWS infrastructure planning CLI for Claude Code, Codex, and similar coding agents.

## MVP behavior

`infra-plan` turns a vague infra request into a single markdown execution plan. It is plan-only: it does not execute Terraform, call AWS APIs, or ask for credentials.

## Usage

```bash
go run ./cmd/infra-plan plan "provision private networking for a new service" --repo /path/to/repo
```

The generated plan is written to:

```text
/path/to/repo/docs/infra-planning-cli/YYYY-MM-DD-provision-private-networking-for-a-new-service.md
```

## Discovery priority

1. Repo code paths and IaC patterns
2. Repo config
3. Official AWS docs evidence
4. Community examples in a future version

## MVP domains

- networking
- compute
- database

If the request is too vague to classify, the CLI returns one minimal clarification question.

## Printing Press handoff

This project is intended to be built and iterated with the `printing-press-hermes` workflow:

```bash
hermes-press research AWSInfraPlanning --spec ./docs/infra-planning-cli/01-position-core-concept.md --out ./runs/aws-infra-planning --json
hermes-press generate aws-infra-planning --run-dir ./runs/aws-infra-planning --json
hermes-press verify aws-infra-planning --work-dir ./runs/aws-infra-planning/working/aws-infra-planning --json --strict
```

The four planning documents plus `05-implementation-plan.md` are the project spec bundle.
```

- [ ] **Step 2: Run docs smoke check**

Run:

```bash
cd infra-planning-cli
test -f README.md && grep -q "Printing Press handoff" README.md
```

Expected: command exits 0.

- [ ] **Step 3: Run all tests**

Run:

```bash
cd infra-planning-cli
go test ./...
```

Expected: PASS.

- [ ] **Step 4: Commit**

```bash
git add infra-planning-cli/README.md
git commit -m "docs: document infra planning CLI MVP"
```

---

### Task 9: Verify End-to-End Demo

**Files:**
- No source file changes expected unless the verification fails.

- [ ] **Step 1: Build the CLI**

Run:

```bash
cd infra-planning-cli
go build ./cmd/infra-plan
```

Expected: command exits 0 and creates `infra-plan`.

- [ ] **Step 2: Run the demo against fixture repo**

Run:

```bash
cd infra-planning-cli
./infra-plan plan "provision private networking" --repo ./testdata/repos/aws-networking --out ./testdata/repos/aws-networking/docs/infra-planning-cli/2026-05-15-provision-networking.md
```

Expected output:

```text
wrote plan: ./testdata/repos/aws-networking/docs/infra-planning-cli/2026-05-15-provision-networking.md
```

- [ ] **Step 3: Inspect generated markdown**

Run:

```bash
cd infra-planning-cli
grep -E "^(# Provision Private Networking|## Problem summary|## Recommended path|## Steps|## Verification)" ./testdata/repos/aws-networking/docs/infra-planning-cli/2026-05-15-provision-networking.md
```

Expected output contains all five headings.

- [ ] **Step 4: Run final test suite**

Run:

```bash
cd infra-planning-cli
go test ./...
```

Expected: PASS.

- [ ] **Step 5: Commit final demo artifact only if intentionally keeping it**

If the generated fixture plan should remain as a golden file:

```bash
git add infra-planning-cli/testdata/repos/aws-networking/docs/infra-planning-cli/2026-05-15-provision-networking.md
git commit -m "test: add golden networking plan output"
```

If the generated fixture plan is only a smoke-test artifact:

```bash
rm infra-planning-cli/testdata/repos/aws-networking/docs/infra-planning-cli/2026-05-15-provision-networking.md
git status --short
```

Expected: no uncommitted generated fixture plan.

---

## Self-Review

**Spec coverage:**
- `01-position-core-concept.md`: implemented through the plan-only CLI, repo-first discovery, minimal clarification gate, and markdown artifact output.
- `02-discovery-workflow.md`: implemented through repo discovery first, AWS docs evidence mapping second, and conservative clarification for unknown domains.
- `03-output-format.md`: implemented through `docs/infra-planning-cli/YYYY-MM-DD-<slug>.md` and mandatory markdown sections.
- `04-mvp-scope.md`: implemented as AWS-first, networking demo first, with compute/database classification included but execution focused on networking.

**Placeholder scan:**
- No `TBD`, unresolved `TODO`, or unspecified implementation steps are required to complete the MVP.

**Type consistency:**
- `planner.Domain`, `planner.Request`, `planner.Plan`, `planner.Step`, `discovery.RepoContext`, and `render.MarkdownPlan` are defined before use.

## Execution Handoff

Plan complete. Recommended execution mode: use `superpowers:subagent-driven-development` and dispatch one fresh subagent per task, with review after each task. If running inline instead, use `superpowers:executing-plans` and commit after every task.
