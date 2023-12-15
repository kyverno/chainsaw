package commands

import (
	"github.com/kyverno/chainsaw/pkg/commands/build"
	"github.com/kyverno/chainsaw/pkg/commands/create"
	"github.com/kyverno/chainsaw/pkg/commands/docs"
	"github.com/kyverno/chainsaw/pkg/commands/export"
	"github.com/kyverno/chainsaw/pkg/commands/migrate"
	"github.com/kyverno/chainsaw/pkg/commands/root"
	"github.com/kyverno/chainsaw/pkg/commands/test"
	"github.com/kyverno/chainsaw/pkg/commands/version"
	"github.com/spf13/cobra"
)

func RootCommand() *cobra.Command {
	cmd := root.Command()
	cmd.AddCommand(
		create.Command(),
		docs.Command(),
		export.Command(),
		build.Command(),
		migrate.Command(),
		test.Command(),
		version.Command(),
	)
	return cmd
}
