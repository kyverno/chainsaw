package version

import (
	"fmt"

	"github.com/kyverno/chainsaw/pkg/version"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:          "version",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "Version: %s\n", version.Version())
			fmt.Fprintf(cmd.OutOrStdout(), "Time: %s\n", version.Time())
			fmt.Fprintf(cmd.OutOrStdout(), "Git commit ID: %s\n", version.Hash())
			return nil
		},
	}
}
