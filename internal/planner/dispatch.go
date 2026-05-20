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
