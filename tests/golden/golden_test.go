package golden_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ShubhanYenuganti/infra-press/tests/golden/checks"
	"gopkg.in/yaml.v3"
)

type catalogEntry struct {
	Name   string `yaml:"name"`
	Path   string `yaml:"path"`
	Binary string `yaml:"binary"`
}

type catalog struct {
	CLIs []catalogEntry `yaml:"clis"`
}

func loadCatalog(t *testing.T) catalog {
	root := repoRoot(t)
	data, err := os.ReadFile(filepath.Join(root, "catalog.yaml"))
	if err != nil {
		t.Fatalf("read catalog: %v", err)
	}
	var c catalog
	if err := yaml.Unmarshal(data, &c); err != nil {
		t.Fatalf("parse catalog: %v", err)
	}
	return c
}

func repoRoot(t *testing.T) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	return filepath.Clean(filepath.Join(wd, "..", ".."))
}

func TestGoldenAcrossLibrary(t *testing.T) {
	cat := loadCatalog(t)
	if len(cat.CLIs) == 0 {
		t.Skip("catalog.yaml has no CLIs yet; nothing to check")
	}
	root := repoRoot(t)
	for _, entry := range cat.CLIs {
		binary := filepath.Join(root, entry.Path, "cmd", entry.Binary, entry.Binary)
		cli := checks.CLI{Name: entry.Name, Binary: binary}
		for _, check := range checks.Registry() {
			t.Run(entry.Name+"/"+check.Name(), func(t *testing.T) {
				check.Run(t, cli)
			})
		}
	}
}
