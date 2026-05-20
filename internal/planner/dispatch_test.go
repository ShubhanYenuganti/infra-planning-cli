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

func TestBuildPlanClarificationNeeded(t *testing.T) {
	req := Request{PrimaryDomain: DomainUnknown, Raw: "do something"}
	plan := BuildPlan(req, discovery.RepoContext{})
	if !plan.Metadata.ClarificationNeeded {
		t.Fatalf("clarification_needed must be true for DomainUnknown")
	}
	if plan.Metadata.ClarificationQuestion == "" {
		t.Fatalf("clarification_question must be non-empty for DomainUnknown")
	}
}
