package checks

import "testing"

func TestAgentsMDCheck_Pass(t *testing.T) {
	check := &AgentsMDCheck{}
	if err := check.Run("testdata/agents-ok"); err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}

func TestAgentsMDCheck_Missing(t *testing.T) {
	check := &AgentsMDCheck{}
	if err := check.Run("testdata/agents-missing"); err == nil {
		t.Error("expected error when AGENTS.md missing")
	}
}

func TestAgentsMDCheck_TooShort(t *testing.T) {
	check := &AgentsMDCheck{}
	if err := check.Run("testdata/agents-empty"); err == nil {
		t.Error("expected error for trivial AGENTS.md")
	}
}
