package planner

import "testing"

func TestClarificationQuestion(t *testing.T) {
	tests := []struct {
		name string
		req  Request
		want string
	}{
		{"known", Request{Raw: "make a vpc", PrimaryDomain: DomainNetworking}, ""},
		{"unknown", Request{Raw: "make infra better", PrimaryDomain: DomainUnknown}, "Which AWS area should this plan focus on first: networking, compute, or database?"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ClarificationQuestion(tt.req); got != tt.want {
				t.Fatalf("ClarificationQuestion() = %q, want %q", got, tt.want)
			}
		})
	}
}
