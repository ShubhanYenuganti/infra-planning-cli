package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

type planOptions struct {
	repo string
	out  string
	json bool
}

func newPlanCommand() *cobra.Command {
	opts := &planOptions{}

	cmd := &cobra.Command{
		Use:   "plan REQUEST",
		Short: "Turn an infra request into a markdown execution plan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "request: %s\nrepo: %s\nout: %s\njson: %v\n", args[0], opts.repo, opts.out, opts.json)
			return nil
		},
	}

	cmd.Flags().StringVar(&opts.repo, "repo", ".", "repository path to inspect")
	cmd.Flags().StringVar(&opts.out, "out", "", "optional output markdown path")
	cmd.Flags().BoolVar(&opts.json, "json", false, "emit JSON payload to stdout instead of writing a file")
	return cmd
}
