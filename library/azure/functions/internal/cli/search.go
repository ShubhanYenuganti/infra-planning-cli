package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newSearchCmd(_ *rootFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "search",
		Short: "Full-text search over synced functions resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("search: not yet implemented (v1.1)")
		},
	}
}
