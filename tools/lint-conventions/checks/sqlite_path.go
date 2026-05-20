package checks

import (
	"fmt"
	"path/filepath"
	"strings"
)

type SQLitePathCheck struct{}

func (*SQLitePathCheck) Name() string { return "sqlite-path" }

func (*SQLitePathCheck) Run(cliDir string) error {
	var contents string
	for _, dir := range []string{"internal/store", "internal", "cmd"} {
		c, _ := concatGoSources(filepath.Join(cliDir, dir))
		contents += c
	}
	if contents == "" {
		return fmt.Errorf("no source files found in internal/store, internal, or cmd")
	}
	if !strings.Contains(contents, `.infra-press/`) {
		return fmt.Errorf(`SQLite path does not include "~/.infra-press/" prefix; ` +
			`stores must live at ~/.infra-press/<cli>.db`)
	}
	return nil
}
