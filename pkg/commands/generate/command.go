package generate

import (
	"github.com/kyverno/chainsaw/pkg/commands/generate/docs"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "generate",
		Short:        "Generate commands",
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
