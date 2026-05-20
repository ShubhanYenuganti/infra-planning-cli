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
