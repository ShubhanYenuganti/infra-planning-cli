package cli

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "infra-plan",
		Short: "Generate agent-ready AWS infrastructure execution plans",
	}
	cmd.AddCommand(newPlanCommand())
	return cmd
}
