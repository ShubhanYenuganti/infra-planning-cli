package cli

import (
	"fmt"
	"os"

	"github.com/ShubhanYenuganti/infra-planning-cli/internal/discovery"
	"github.com/ShubhanYenuganti/infra-planning-cli/internal/planner"
	"github.com/ShubhanYenuganti/infra-planning-cli/internal/render"
	"github.com/ShubhanYenuganti/infra-planning-cli/internal/web"
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
			repo, err := discovery.DiscoverRepo(opts.repo)
			if err != nil {
				return fmt.Errorf("discover repo: %w", err)
			}

			req := planner.ClassifyRequest(args[0])
			fetcher := web.New()
			plan := planner.BuildPlan(req, repo, fetcher)

			if plan.Metadata.ClarificationNeeded {
				fmt.Fprintln(cmd.ErrOrStderr(), plan.Metadata.ClarificationQuestion)
				return nil
			}

			if opts.json {
				b, err := render.JSONPlan(plan)
				if err != nil {
					return fmt.Errorf("render json: %w", err)
				}
				_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s\n", b)
				return err
			}

			md := render.MarkdownPlan(plan)

			if opts.out == "" {
				_, err = fmt.Fprint(cmd.OutOrStdout(), md)
				return err
			}

			return os.WriteFile(opts.out, []byte(md), 0644)
		},
	}

	cmd.Flags().StringVar(&opts.repo, "repo", ".", "repository path to inspect")
	cmd.Flags().StringVar(&opts.out, "out", "", "optional output markdown path")
	cmd.Flags().BoolVar(&opts.json, "json", false, "emit JSON payload to stdout instead of writing a file")
	return cmd
}
