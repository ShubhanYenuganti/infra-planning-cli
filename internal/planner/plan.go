package planner

import "time"

type Confidence string

const (
	ConfidenceLow    Confidence = "low"
	ConfidenceMedium Confidence = "medium"
	ConfidenceHigh   Confidence = "high"
)

type CostHint struct {
	ResourceType       string  `yaml:"resource_type"        json:"resource_type"`
	MonthlyUSDEstimate float64 `yaml:"monthly_usd_estimate" json:"monthly_usd_estimate"`
	Notes              string  `yaml:"notes"                json:"notes"`
}

type Evidence struct {
	Source    string    `yaml:"source"               json:"source"`
	Kind      string    `yaml:"kind"                 json:"kind"`
	Path      string    `yaml:"path,omitempty"       json:"path,omitempty"`
	Snippet   string    `yaml:"snippet,omitempty"    json:"snippet,omitempty"`
	Title     string    `yaml:"title,omitempty"      json:"title,omitempty"`
	URL       string    `yaml:"url,omitempty"        json:"url,omitempty"`
	FetchedAt time.Time `yaml:"fetched_at,omitempty" json:"fetched_at,omitempty"`
	Status    int       `yaml:"status,omitempty"     json:"status,omitempty"`
}

type Conflict struct {
	Description        string   `yaml:"description"         json:"description"`
	Options            []string `yaml:"options,omitempty"   json:"options,omitempty"`
	ResolutionQuestion string   `yaml:"resolution_question" json:"resolution_question"`
}

type Metadata struct {
	SchemaVersion         int        `yaml:"schema_version"         json:"schema_version"`
	GeneratedAt           time.Time  `yaml:"generated_at"           json:"generated_at"`
	Request               string     `yaml:"request"                json:"request"`
	Domain                Domain     `yaml:"domain"                 json:"domain"`
	Confidence            Confidence `yaml:"confidence"             json:"confidence"`
	ClarificationNeeded   bool       `yaml:"clarification_needed"   json:"clarification_needed"`
	ClarificationQuestion string     `yaml:"clarification_question" json:"clarification_question"`
	Evidence              []Evidence `yaml:"evidence"               json:"evidence"`
	Conflicts             []Conflict `yaml:"conflicts"              json:"conflicts"`
	Prerequisites         []string   `yaml:"prerequisites"          json:"prerequisites"`
	CostHints             []CostHint `yaml:"cost_hints"             json:"cost_hints"`
	AlreadySatisfied      bool       `yaml:"already_satisfied"      json:"already_satisfied"`
}

type Step struct {
	Title string
	Body  string
}

type Plan struct {
	Metadata        Metadata
	Title           string
	ProblemSummary  string
	RecommendedPath string
	Steps           []Step
	Verification    []string
}
