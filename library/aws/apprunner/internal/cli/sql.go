package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newSqlCmd(_ *rootFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "sql",
		Short: "Run a raw SQL query against the local apprunner cache",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("sql: not yet implemented (v1.1)")
		},
	}
}
