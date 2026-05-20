package checks

import (
	"os/exec"
	"testing"
)

type UnknownCommandExits2 struct{}

func (*UnknownCommandExits2) Name() string { return "unknown-command-exits-2" }

func (*UnknownCommandExits2) Run(t *testing.T, cli CLI) {
	cmd := exec.Command(cli.Binary, "this-command-does-not-exist")
	err := cmd.Run()
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Errorf("%s: expected ExitError, got %v", cli.Name, err)
		return
	}
	if exitErr.ExitCode() != 2 {
		t.Errorf("%s: expected exit code 2 (user error), got %d", cli.Name, exitErr.ExitCode())
	}
}
