package create

import (
	"github.com/kyverno/chainsaw/pkg/commands/create/test"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "create",
		Short:        "Create Chainsaw resources",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(
		test.Command(),
	)
	return cmd
}
