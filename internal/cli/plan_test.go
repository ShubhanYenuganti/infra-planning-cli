package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// runPlan executes the plan command with the given args and returns captured stdout.
func runPlan(t *testing.T, args ...string) string {
	t.Helper()
	cmd := newPlanCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("plan command failed: %v", err)
	}
	return buf.String()
}

func TestPlanMarkdownOutput(t *testing.T) {
	tmpDir := t.TempDir()
	out := runPlan(t, "--repo", tmpDir, "add a VPC with private subnets")
	if !strings.Contains(out, "---") {
		t.Errorf("expected YAML frontmatter delimiter in output, got:\n%s", out)
	}
	if !strings.Contains(out, "# ") {
		t.Errorf("expected markdown heading in output, got:\n%s", out)
	}
}

func TestPlanJSONOutput(t *testing.T) {
	tmpDir := t.TempDir()
	out := runPlan(t, "--repo", tmpDir, "--json", "add a VPC with private subnets")
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &m); err != nil {
		t.Fatalf("expected valid JSON output, got error: %v\noutput:\n%s", err, out)
	}
	if _, ok := m["metadata"]; !ok {
		t.Errorf("expected 'metadata' key in JSON output")
	}
}

func TestPlanWriteFile(t *testing.T) {
	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, "plan.md")
	// Run without capturing stdout since output goes to file.
	cmd := newPlanCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"--repo", tmpDir, "--out", outFile, "deploy an ECS service"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("plan command failed: %v", err)
	}
	// stdout should be empty
	if buf.Len() != 0 {
		t.Errorf("expected no stdout when --out is set, got: %s", buf.String())
	}
	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("expected output file to exist: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty output file")
	}
	if !strings.Contains(string(data), "---") {
		t.Errorf("expected YAML frontmatter in file, got:\n%s", string(data))
	}
}

func TestPlanClarificationNeeded(t *testing.T) {
	tmpDir := t.TempDir()
	// A vague request should trigger a clarification question.
	// We capture stderr to check for the clarification question.
	cmd := newPlanCommand()
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"--repo", tmpDir, "do something"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("plan command failed: %v", err)
	}
	// A vague request maps to DomainUnknown → ClarificationNeeded=true.
	// The question should appear on stderr and stdout should be empty.
	if stderr.Len() == 0 {
		t.Error("expected clarification question on stderr, got nothing")
	}
	if stdout.Len() != 0 {
		t.Errorf("expected no plan output on stdout when clarification needed, got: %s", stdout.String())
	}
}

func TestPlanInvalidRepo(t *testing.T) {
	cmd := newPlanCommand()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"--repo", "/nonexistent/path/that/does/not/exist", "add a VPC"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for nonexistent repo path, got nil")
	}
}

func TestPlanNonEmptyOutput(t *testing.T) {
	tmpDir := t.TempDir()
	// Ensure at least one known domain produces non-empty output.
	out := runPlan(t, "--repo", tmpDir, "provision RDS postgres database")
	if strings.TrimSpace(out) == "" {
		t.Error("expected non-empty output for database request")
	}
}
