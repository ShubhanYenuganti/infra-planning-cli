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
