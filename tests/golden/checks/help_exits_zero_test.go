package checks

import "testing"

func TestHelpExitsZero_Pass(t *testing.T) {
	check := &HelpExitsZero{}
	check.Run(t, CLI{Name: "mock", Binary: "testdata/mock-binaries/help-ok.sh"})
}
