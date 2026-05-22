package smoke

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRecipeSmoke(t *testing.T) {
	root, err := RepoRoot()
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	cat, err := LoadCatalog(root)
	if err != nil {
		t.Fatalf("load catalog: %v", err)
	}

	// Build PATH with every CLI binary directory prepended.
	binDirs := make([]string, 0, len(cat.CLIs))
	for _, cli := range cat.CLIs {
		binDirs = append(binDirs, filepath.Join(root, cli.Path, "cmd", cli.Binary))
	}
	augmentedPath := strings.Join(binDirs, ":") + ":" + os.Getenv("PATH")

	script := filepath.Join(root, "scripts", "recipes", "doctor-all.sh")
	cmd := exec.Command("bash", script)
	env := os.Environ()
	for i, kv := range env {
		if strings.HasPrefix(kv, "PATH=") {
			env[i] = "PATH=" + augmentedPath
			break
		}
	}
	cmd.Env = env

	out, _ := cmd.Output() // Accept any exit code; assert structure of stdout.

	var report struct {
		CLIs []struct {
			CLI    string          `json:"cli"`
			Report json.RawMessage `json:"report"`
		} `json:"clis"`
	}
	if err := json.Unmarshal(out, &report); err != nil {
		t.Fatalf("recipe stdout not parseable JSON: %v\noutput: %s", err, out)
	}
	if got, want := len(report.CLIs), len(cat.CLIs); got != want {
		t.Fatalf("recipe report listed %d CLIs, expected %d", got, want)
	}
}
