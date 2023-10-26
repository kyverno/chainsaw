package docs

import (
	"log"

	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	var options options
	cmd := &cobra.Command{
		Use:          "docs",
		Short:        "Generate reference documentation",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			root := cmd.Root()
			if err := options.validate(root); err != nil {
				return err
			}
			return options.execute(root)
		},
	}
	cmd.Flags().StringVarP(&options.path, "output", "o", ".", "Output path")
	cmd.Flags().BoolVar(&options.website, "website", false, "Website version")
	cmd.Flags().BoolVar(&options.autogenTag, "autogenTag", true, "Determines if the generated docs should contain a timestamp")
	if err := cmd.MarkFlagDirname("output"); err != nil {
		log.Println("WARNING", err)
	}
	if err := cmd.MarkFlagRequired("output"); err != nil {
		log.Println("WARNING", err)
	}
	return cmd
}
