package checks

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type PressfileCheck struct{}

func (*PressfileCheck) Name() string { return "pressfile-schema" }

type pressfileSchema struct {
	CLI          string `yaml:"cli"`
	Cloud        string `yaml:"cloud"`
	Service      string `yaml:"service"`
	PressVersion string `yaml:"press_version"`
	PressCommand string `yaml:"press_command"`
	SpecURL      string `yaml:"spec_url"`
	SpecVersion  string `yaml:"spec_version"`
	SpecEtag     string `yaml:"spec_etag"`
	GeneratedAt  string `yaml:"generated_at"`
	Patches      []struct {
		Path   string `yaml:"path"`
		Reason string `yaml:"reason"`
	} `yaml:"patches"`
}

func (*PressfileCheck) Run(cliDir string) error {
	path := filepath.Join(cliDir, "pressfile.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read pressfile: %w", err)
	}
	var p pressfileSchema
	if err := yaml.Unmarshal(data, &p); err != nil {
		return fmt.Errorf("parse pressfile: %w", err)
	}
	required := map[string]string{
		"cli":           p.CLI,
		"cloud":         p.Cloud,
		"service":       p.Service,
		"press_version": p.PressVersion,
		"spec_url":      p.SpecURL,
		"generated_at":  p.GeneratedAt,
	}
	for field, value := range required {
		if value == "" {
			return fmt.Errorf("missing required field: %s", field)
		}
	}
	return nil
}
