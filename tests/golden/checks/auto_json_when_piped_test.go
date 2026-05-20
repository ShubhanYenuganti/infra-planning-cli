package checks

import "testing"

func TestAutoJSONWhenPiped(t *testing.T) {
	check := &AutoJSONWhenPiped{}
	check.Run(t, CLI{Name: "mock", Binary: "testdata/mock-binaries/json-when-piped.sh"})
}
