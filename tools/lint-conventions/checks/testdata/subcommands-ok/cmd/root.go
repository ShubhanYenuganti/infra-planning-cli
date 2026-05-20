package cmd

// References the required subcommands: sync, search, sql.
var syncCmd, searchCmd, sqlCmd, rootCmd struct{}

func init() {
	_ = syncCmd
	_ = searchCmd
	_ = sqlCmd
	_ = rootCmd
}
