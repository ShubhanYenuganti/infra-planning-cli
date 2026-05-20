package checks

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

type AgentsMDCheck struct{}

func (*AgentsMDCheck) Name() string { return "agents-md-non-trivial" }

func (*AgentsMDCheck) Run(cliDir string) error {
	path := filepath.Join(cliDir, "AGENTS.md")
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("AGENTS.md missing or unreadable: %w", err)
	}
	defer f.Close()
	lines := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines++
	}
	if lines < 8 {
		return fmt.Errorf("AGENTS.md too short (%d lines, need >=8)", lines)
	}
	return nil
}
