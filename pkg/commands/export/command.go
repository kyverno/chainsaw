package export

import (
	"github.com/kyverno/chainsaw/pkg/commands/export/schemas"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "export",
		Short:        "Export commands",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(
		schemas.Command(),
	)
	return cmd
}
