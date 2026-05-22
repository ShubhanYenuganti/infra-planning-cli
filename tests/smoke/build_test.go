package smoke

import (
	"os/exec"
	"path/filepath"
	"testing"
)

func TestBuildMatrix(t *testing.T) {
	root, err := RepoRoot()
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	cat, err := LoadCatalog(root)
	if err != nil {
		t.Fatalf("load catalog: %v", err)
	}
	for _, cli := range cat.CLIs {
		cli := cli
		t.Run(cli.Name, func(t *testing.T) {
			cliDir := filepath.Join(root, cli.Path)

			t.Run("go-build-all", func(t *testing.T) {
				cmd := exec.Command("go", "build", "./...")
				cmd.Dir = cliDir
				if out, err := cmd.CombinedOutput(); err != nil {
					t.Fatalf("go build ./... failed in %s:\n%s", cli.Path, out)
				}
			})

			t.Run("binary-builds", func(t *testing.T) {
				bin := BinaryPath(root, cli)
				cmd := exec.Command("go", "build", "-o", bin, "./cmd/"+cli.Binary+"/")
				cmd.Dir = cliDir
				if out, err := cmd.CombinedOutput(); err != nil {
					t.Fatalf("binary build failed for %s:\n%s", cli.Name, out)
				}
			})
		})
	}
}
