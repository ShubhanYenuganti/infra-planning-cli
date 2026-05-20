package e2e_test

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ShubhanYenuganti/infra-planning-cli/internal/planner"
)

var builtBin string

func TestMain(m *testing.M) {
	bin, err := buildBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "build failed: %v\n", err)
		os.Exit(1)
	}
	builtBin = bin
	os.Exit(m.Run())
}

func buildBinary() (string, error) {
	bin := filepath.Join(os.TempDir(), "infra-planning-e2e-bin")
	// Resolve module root: tests/e2e/ is two directories below the module root.
	moduleRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		return "", fmt.Errorf("resolve module root: %v", err)
	}
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/infra-plan")
	cmd.Dir = moduleRoot
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("go build: %v\n%s", err, out)
	}
	return bin, nil
}

// repoPath resolves a testdata fixture path relative to the module root.
// tests/e2e/ is two directories below the module root.
func repoPath(name string) string {
	return filepath.Join("..", "..", "testdata", "repos", name)
}

// runPlan runs "infra-plan plan <request> --repo <repo> [extra...]"
// and returns stdout. Fails the test on non-zero exit.
func runPlan(t *testing.T, repo, request string, extra ...string) []byte {
	t.Helper()
	args := []string{"plan", request, "--repo", repo}
	args = append(args, extra...)
	cmd := exec.Command(builtBin, args...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("infra-plan plan failed: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return []byte(stdout.String())
}

func TestE2ENetworkingMarkdown(t *testing.T) {
	out := runPlan(t, repoPath("aws-networking"), "Add a private subnet")
	body := string(out)

	if !strings.HasPrefix(body, "---\n") {
		t.Fatalf("expected YAML frontmatter to start output, got:\n%s", body)
	}
	if !strings.Contains(body, "domain: networking") {
		t.Errorf("expected 'domain: networking' in frontmatter")
	}
	if !strings.Contains(body, "confidence:") {
		t.Errorf("expected 'confidence:' field in frontmatter")
	}
	if !strings.Contains(body, "## Steps") {
		t.Errorf("expected '## Steps' section in plan body")
	}
}

func TestE2EComputeJSONHasCostHints(t *testing.T) {
	out := runPlan(t, repoPath("aws-compute"), "Add a Lambda function", "--json")

	var plan planner.Plan
	if err := json.Unmarshal(out, &plan); err != nil {
		t.Fatalf("json.Unmarshal failed: %v\noutput was:\n%s", err, out)
	}
	if len(plan.Metadata.CostHints) < 2 {
		t.Errorf("expected >=2 cost hints for compute domain, got %d", len(plan.Metadata.CostHints))
	}
}

func TestE2ENoWebFlagSuppressesFetch(t *testing.T) {
	out := runPlan(t, repoPath("aws-networking"), "Add a private subnet", "--no-web", "--json")

	var plan planner.Plan
	if err := json.Unmarshal(out, &plan); err != nil {
		t.Fatalf("json.Unmarshal failed: %v\noutput was:\n%s", err, out)
	}
	for _, ev := range plan.Metadata.Evidence {
		if !ev.FetchedAt.IsZero() {
			t.Errorf("--no-web set but evidence %q has FetchedAt=%v", ev.Source, ev.FetchedAt)
		}
	}
}

func TestE2EMixedIaCConflict(t *testing.T) {
	out := runPlan(t, repoPath("aws-mixed"), "Add a VPC", "--json")

	var plan planner.Plan
	if err := json.Unmarshal(out, &plan); err != nil {
		t.Fatalf("json.Unmarshal failed: %v\noutput was:\n%s", err, out)
	}
	if plan.Metadata.Confidence != planner.ConfidenceMedium {
		t.Errorf("expected confidence=%q for mixed-IaC repo, got %q",
			planner.ConfidenceMedium, plan.Metadata.Confidence)
	}
	if len(plan.Metadata.Conflicts) != 1 {
		t.Errorf("expected exactly 1 conflict for mixed-IaC repo, got %d",
			len(plan.Metadata.Conflicts))
	}
	if !strings.Contains(plan.Metadata.Conflicts[0].Description, "terraform") {
		t.Errorf("expected conflict description to mention 'terraform', got: %s",
			plan.Metadata.Conflicts[0].Description)
	}
}

func TestE2EAlreadySatisfied(t *testing.T) {
	out := runPlan(t, repoPath("aws-satisfied"), "Add a VPC", "--json")

	var plan planner.Plan
	if err := json.Unmarshal(out, &plan); err != nil {
		t.Fatalf("json.Unmarshal failed: %v\noutput was:\n%s", err, out)
	}
	if !plan.Metadata.AlreadySatisfied {
		t.Error("expected already_satisfied=true for aws-satisfied fixture")
	}
	if len(plan.Steps) == 0 {
		t.Fatal("expected at least one step in plan")
	}
	if !strings.Contains(plan.Steps[0].Title, "Confirm") {
		t.Errorf("expected first step title to contain 'Confirm', got: %q",
			plan.Steps[0].Title)
	}
}

func TestE2EJSONAndOutMutuallyExclusive(t *testing.T) {
	cmd := exec.Command(builtBin, "plan",
		"Add a subnet",
		"--repo", repoPath("aws-networking"),
		"--json",
		"--out", "/tmp/test-plan.md",
	)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected non-zero exit when --json and --out both set, got output:\n%s", out)
	}
}
