package checks

import "testing"

func TestUnknownCommandExits2(t *testing.T) {
	check := &UnknownCommandExits2{}
	check.Run(t, CLI{Name: "mock", Binary: "testdata/mock-binaries/help-ok.sh"})
}
