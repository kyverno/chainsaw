package kuttl

import (
	"github.com/kyverno/chainsaw/pkg/commands/kuttl/migrate"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "kuttl",
		Short:        "Work with KUTTL tests",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(
		migrate.Command(),
	)
	return cmd
}
