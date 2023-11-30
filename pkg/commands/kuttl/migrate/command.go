package migrate

import (
	"github.com/kyverno/chainsaw/pkg/commands/kuttl/migrate/config"
	"github.com/kyverno/chainsaw/pkg/commands/kuttl/migrate/tests"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "migrate",
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
