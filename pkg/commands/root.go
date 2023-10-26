package commands

import (
	"github.com/kyverno/chainsaw/pkg/commands/docs"
	"github.com/kyverno/chainsaw/pkg/commands/test"
	"github.com/kyverno/chainsaw/pkg/commands/version"
	"github.com/spf13/cobra"
)

func RootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "chainsaw",
		Short:        "Stronger tool for e2e testing",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(
		docs.Command(),
		test.Command(),
		version.Command(),
	)
	return cmd
}
