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
