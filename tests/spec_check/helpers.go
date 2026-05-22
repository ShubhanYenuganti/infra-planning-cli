package spec_check

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type CLI struct {
	Name    string `yaml:"name"`
	Cloud   string `yaml:"cloud"`
	Service string `yaml:"service"`
	Path    string `yaml:"path"`
	Binary  string `yaml:"binary"`
	Status  string `yaml:"status"`
	SpecURL string `yaml:"spec_url"`
}

type Catalog struct {
	Version int   `yaml:"version"`
	CLIs    []CLI `yaml:"clis"`
}

type Pressfile struct {
	SpecSHA256   string `yaml:"spec_sha256"`
	VendoredSpec string `yaml:"vendored_spec"`
	SpecURL      string `yaml:"spec_url"`
}

// RepoRoot walks up from the working directory looking for catalog.yaml.
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

func LoadPressfile(repoRoot, cliPath string) (*Pressfile, error) {
	b, err := os.ReadFile(filepath.Join(repoRoot, cliPath, "pressfile.yaml"))
	if err != nil {
		return nil, err
	}
	var p Pressfile
	if err := yaml.Unmarshal(b, &p); err != nil {
		return nil, err
	}
	return &p, nil
}
