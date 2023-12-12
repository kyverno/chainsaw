package docs

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/spf13/cobra"
)

//go:embed docs.tmpl
var docsTemplate string

type options struct {
	testFile   string
	readmeFile string
	testDirs   []string
}

func Command() *cobra.Command {
	var options options
	cmd := &cobra.Command{
		Use:          "docs",
		Short:        "Generate documentation for tests",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			tmpl, err := template.New("docs").Funcs(sprig.FuncMap()).Parse(docsTemplate)
			if err != nil {
				return err
			}
			tests, err := discovery.DiscoverTests(options.testFile, options.testDirs...)
			if err != nil {
				return err
			}
			for _, test := range tests {
				if test.Err == nil {
					file, err := os.Create(filepath.Join(test.BasePath, options.readmeFile))
					if err != nil {
						return err
					}
					defer file.Close()
					output := file
					if err := tmpl.Execute(output, test); err != nil {
						return err
					}
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "ERROR: failed to load test %s (%s)", test.BasePath, test.Err)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&options.testFile, "test-file", "chainsaw-test.yaml", "Name of the test file")
	cmd.Flags().StringVar(&options.readmeFile, "readme-file", "README.md", "Name of the generated docs file")
	cmd.Flags().StringArrayVar(&options.testDirs, "test-dir", []string{}, "Directories containing test cases to run")
	return cmd
}
