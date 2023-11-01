package commands

import (
	"github.com/kyverno/chainsaw/pkg/commands/docs"
	"github.com/kyverno/chainsaw/pkg/commands/kuttl"
	"github.com/kyverno/chainsaw/pkg/commands/root"
	"github.com/kyverno/chainsaw/pkg/commands/test"
	"github.com/kyverno/chainsaw/pkg/commands/version"
	"github.com/spf13/cobra"
)

func RootCommand() *cobra.Command {
	cmd := root.Command()
	cmd.AddCommand(
		docs.Command(),
		kuttl.Command(),
		test.Command(),
		version.Command(),
	)
	return cmd
}
