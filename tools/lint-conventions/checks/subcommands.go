package checks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SubcommandsCheck struct{}

func (*SubcommandsCheck) Name() string { return "subcommands-present" }

func (*SubcommandsCheck) Run(cliDir string) error {
	required := []string{"sync", "search", "sql"}
	cmdDir := filepath.Join(cliDir, "cmd")
	contents, err := concatGoSources(cmdDir)
	if err != nil {
		alt := filepath.Join(cliDir, "internal", "cmd")
		contents, err = concatGoSources(alt)
		if err != nil {
			return fmt.Errorf("no cmd directory found")
		}
	}
	for _, sub := range required {
		if !strings.Contains(contents, sub+"Cmd") &&
			!strings.Contains(contents, `Use: "`+sub+`"`) &&
			!strings.Contains(contents, `Use:   "`+sub+`"`) {
			return fmt.Errorf("missing required subcommand: %s", sub)
		}
	}
	return nil
}

func concatGoSources(dir string) (string, error) {
	var b strings.Builder
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		b.Write(data)
		b.WriteString("\n")
		return nil
	})
	if err != nil {
		return "", err
	}
	if b.Len() == 0 {
		return "", fmt.Errorf("no .go files in %s", dir)
	}
	return b.String(), nil
}
