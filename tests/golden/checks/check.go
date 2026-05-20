package checks

import "testing"

// CLI describes one entry in catalog.yaml relevant to golden tests.
type CLI struct {
	Name   string
	Binary string // resolved path to the built binary
}

// Check is a single black-box behavioral assertion.
type Check interface {
	Name() string
	Run(t *testing.T, cli CLI)
}

func Registry() []Check {
	return []Check{
		&HelpExitsZero{},
		&UnknownCommandExits2{},
		&AutoJSONWhenPiped{},
		&VersionIsSemver{},
	}
}
