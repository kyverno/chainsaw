package kuttl

import (
	"github.com/kyverno/chainsaw/pkg/commands/migrate/kuttl/config"
	"github.com/kyverno/chainsaw/pkg/commands/migrate/kuttl/tests"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "kuttl",
		Short:        "Migrate KUTTL resources to Chainsaw",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(
		config.Command(),
		tests.Command(),
	)
	return cmd
}
