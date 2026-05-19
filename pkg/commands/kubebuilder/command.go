package kubebuilder

import (
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "kubebuilder",
		Short:        "Kubebuilder plugin for Chainsaw",
		Long:         `A kubebuilder plugin that scaffolds Chainsaw e2e tests for kubebuilder-based operators.`,
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(
		scaffoldCommand(),
	)
	return cmd
}
