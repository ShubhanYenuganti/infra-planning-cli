package checks

// Check is a single convention rule applied to a CLI directory.
type Check interface {
	// Name returns a short identifier (e.g., "pressfile-schema").
	Name() string
	// Run inspects the CLI rooted at cliDir and returns nil on pass,
	// or an error describing the violation.
	Run(cliDir string) error
}

// Registry returns all checks the linter knows about.
// Add new checks here as conventions evolve.
func Registry() []Check {
	return []Check{
		&PressfileCheck{},
		&SubcommandsCheck{},
		&FlagsCheck{},
		&SQLitePathCheck{},
		&AgentsMDCheck{},
	}
}
