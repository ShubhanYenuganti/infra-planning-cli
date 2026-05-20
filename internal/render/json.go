package render

import (
	"encoding/json"

	"github.com/ShubhanYenuganti/infra-planning-cli/internal/planner"
)

// JSONPlan serializes p to indented JSON (2-space indent).
func JSONPlan(p planner.Plan) ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}
