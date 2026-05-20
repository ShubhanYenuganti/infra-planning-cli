package render

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/ShubhanYenuganti/infra-planning-cli/internal/planner"
)

func TestJSONPlanContainsTitle(t *testing.T) {
	p := sampleFullPlan()
	got, err := JSONPlan(p)
	if err != nil {
		t.Fatalf("JSONPlan returned error: %v", err)
	}
	if !strings.Contains(string(got), "Provision Private Networking") {
		t.Fatalf("JSON output missing title:\n%s", got)
	}
}

func TestJSONPlanIsValidJSON(t *testing.T) {
	p := sampleFullPlan()
	got, err := JSONPlan(p)
	if err != nil {
		t.Fatalf("JSONPlan returned error: %v", err)
	}
	var roundtrip planner.Plan
	if err := json.Unmarshal(got, &roundtrip); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, got)
	}
	if roundtrip.Title != p.Title {
		t.Fatalf("round-trip title mismatch: got %q, want %q", roundtrip.Title, p.Title)
	}
}

func TestJSONPlanGeneratedAtRFC3339(t *testing.T) {
	p := sampleFullPlan()
	got, err := JSONPlan(p)
	if err != nil {
		t.Fatalf("JSONPlan returned error: %v", err)
	}
	// RFC3339 representation of 2026-05-19T01:23:45Z
	want := "2026-05-19T01:23:45Z"
	if !strings.Contains(string(got), want) {
		t.Fatalf("JSON output missing RFC3339 time %q:\n%s", want, got)
	}
}

func TestJSONPlanIndented(t *testing.T) {
	p := planner.Plan{
		Metadata: planner.Metadata{
			SchemaVersion: 1,
			GeneratedAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		Title: "Simple Plan",
	}
	got, err := JSONPlan(p)
	if err != nil {
		t.Fatalf("JSONPlan returned error: %v", err)
	}
	// Indented JSON must have newlines and 2-space indentation
	if !strings.Contains(string(got), "\n  ") {
		t.Fatalf("JSON output does not appear to be indented:\n%s", got)
	}
}
