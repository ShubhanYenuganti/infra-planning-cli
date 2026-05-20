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
