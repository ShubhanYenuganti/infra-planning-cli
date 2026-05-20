package planner

import "testing"

func TestCostHintsForDomainNetworking(t *testing.T) {
	hints := CostHintsForDomain(DomainNetworking)
	if len(hints) != 2 {
		t.Fatalf("expected 2 networking hints, got %d", len(hints))
	}
	found := false
	for _, h := range hints {
		if h.ResourceType == "aws_nat_gateway" {
			found = true
			if h.MonthlyUSDEstimate != 32.40 {
				t.Fatalf("nat gateway estimate = %.2f, want 32.40", h.MonthlyUSDEstimate)
			}
		}
	}
	if !found {
		t.Fatalf("expected aws_nat_gateway hint in networking domain")
	}
}

func TestCostHintsForDomainCompute(t *testing.T) {
	hints := CostHintsForDomain(DomainCompute)
	if len(hints) != 2 {
		t.Fatalf("expected 2 compute hints, got %d", len(hints))
	}
}

func TestCostHintsForDomainDatabase(t *testing.T) {
	hints := CostHintsForDomain(DomainDatabase)
	if len(hints) != 2 {
		t.Fatalf("expected 2 database hints, got %d", len(hints))
	}
}

func TestCostHintsAllHaveNotes(t *testing.T) {
	for _, domain := range []Domain{DomainNetworking, DomainCompute, DomainDatabase} {
		for _, h := range CostHintsForDomain(domain) {
			if h.Notes == "" {
				t.Fatalf("domain %s: hint %q missing notes", domain, h.ResourceType)
			}
		}
	}
}
