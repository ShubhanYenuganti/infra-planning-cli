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
