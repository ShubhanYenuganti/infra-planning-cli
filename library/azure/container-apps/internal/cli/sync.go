package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newSyncCmd(_ *rootFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Sync container-apps resources to local SQLite cache",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("sync: not yet implemented (v1.1)")
		},
	}
}
