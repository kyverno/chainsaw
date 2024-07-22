package renovate

import (
	"github.com/kyverno/chainsaw/pkg/commands/renovate/config"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "renovate",
		Short:        "Upgrade Chainsaw resources",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(
		config.Command(),
	)
	return cmd
}
