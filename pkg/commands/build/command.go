package build

import (
	"github.com/kyverno/chainsaw/pkg/commands/build/docs"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "build",
		Short:        "Build commands",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(
		docs.Command(),
	)
	return cmd
}
