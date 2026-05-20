package cmd

// References each required flag as its own quoted literal (mimics
// cobra's Flags().Bool("json", ...), Flags().String("dry-run", ...) etc.)
const (
	flagJSON       = "json"
	flagDryRun     = "dry-run"
	flagSelect     = "select"
	flagDataSource = "data-source"
	flagCompact    = "compact"
)

var _ = flagJSON
var _ = flagDryRun
var _ = flagSelect
var _ = flagDataSource
var _ = flagCompact
