package checks

import (
	"os/exec"
	"testing"
)

type HelpExitsZero struct{}

func (*HelpExitsZero) Name() string { return "help-exits-zero" }

func (*HelpExitsZero) Run(t *testing.T, cli CLI) {
	cmd := exec.Command(cli.Binary, "--help")
	if err := cmd.Run(); err != nil {
		t.Errorf("%s --help exited non-zero: %v", cli.Name, err)
	}
}
