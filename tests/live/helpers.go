//go:build live

package live

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type LiveSmoke struct {
	ListSubcommand string   `yaml:"list_subcommand"`
	Args           []string `yaml:"args"`
}

type CLI struct {
	Name      string    `yaml:"name"`
	Cloud     string    `yaml:"cloud"`
	Path      string    `yaml:"path"`
	Binary    string    `yaml:"binary"`
	Status    string    `yaml:"status"`
	LiveSmoke LiveSmoke `yaml:"live_smoke"`
}

type Catalog struct {
	Version int   `yaml:"version"`
	CLIs    []CLI `yaml:"clis"`
}

func RepoRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "catalog.yaml")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

func LoadCatalog(repoRoot string) (*Catalog, error) {
	b, err := os.ReadFile(filepath.Join(repoRoot, "catalog.yaml"))
	if err != nil {
		return nil, err
	}
	var c Catalog
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func BinaryPath(repoRoot string, cli CLI) string {
	return filepath.Join(repoRoot, cli.Path, "cmd", cli.Binary, cli.Binary)
}
