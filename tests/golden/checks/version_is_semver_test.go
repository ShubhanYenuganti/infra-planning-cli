package checks

import "testing"

func TestVersionIsSemver(t *testing.T) {
	check := &VersionIsSemver{}
	check.Run(t, CLI{Name: "mock", Binary: "testdata/mock-binaries/version-ok.sh"})
}
