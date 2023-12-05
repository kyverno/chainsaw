package migrate

import (
	"github.com/kyverno/chainsaw/pkg/commands/migrate/kuttl"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "migrate",
		Short:        "Migrate resources to Chainsaw",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(
		kuttl.Command(),
	)
	return cmd
}
