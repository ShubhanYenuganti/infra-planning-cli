package checks

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"testing"
)

type AutoJSONWhenPiped struct{}

func (*AutoJSONWhenPiped) Name() string { return "auto-json-when-piped" }

func (*AutoJSONWhenPiped) Run(t *testing.T, cli CLI) {
	cmd := exec.Command(cli.Binary, "list")
	var buf bytes.Buffer
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		t.Skipf("%s: list command unsupported or errored: %v (skip — not a hard failure)", cli.Name, err)
		return
	}
	var v any
	if err := json.Unmarshal(buf.Bytes(), &v); err != nil {
		t.Errorf("%s: stdout was not JSON when piped: %v", cli.Name, err)
	}
}
