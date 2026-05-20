package checks

import (
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

type VersionIsSemver struct{}

func (*VersionIsSemver) Name() string { return "version-is-semver" }

var semver = regexp.MustCompile(`^v?\d+\.\d+\.\d+(?:-[0-9A-Za-z.-]+)?$`)

func (*VersionIsSemver) Run(t *testing.T, cli CLI) {
	cmd := exec.Command(cli.Binary, "--version")
	out, err := cmd.Output()
	if err != nil {
		t.Errorf("%s: --version failed: %v", cli.Name, err)
		return
	}
	got := strings.TrimSpace(string(out))
	if !semver.MatchString(got) {
		t.Errorf("%s: --version output %q is not SemVer", cli.Name, got)
	}
}
