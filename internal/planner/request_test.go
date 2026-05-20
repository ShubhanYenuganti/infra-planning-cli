package planner

import "testing"

func TestClassifyRequest(t *testing.T) {
	tests := []struct {
		name string
		text string
		want Domain
	}{
		{"vpc", "provision a vpc with private subnets", DomainNetworking},
		{"lambda", "deploy a lambda worker", DomainCompute},
		{"rds", "create a postgres rds database", DomainDatabase},
		{"unknown", "improve our platform", DomainUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassifyRequest(tt.text)
			if got.PrimaryDomain != tt.want {
				t.Fatalf("PrimaryDomain = %q, want %q", got.PrimaryDomain, tt.want)
			}
		})
	}
}
