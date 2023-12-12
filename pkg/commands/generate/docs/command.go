package docs

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/spf13/cobra"
)

//go:embed docs.tmpl
var docsTemplate string

//go:embed catalog.tmpl
var catalogTemplate string

var (
	docsTmpl    = template.Must(template.New("docs").Funcs(sprig.FuncMap()).Parse(docsTemplate))
	catalogTmpl = template.Must(template.New("catalog").Funcs(sprig.FuncMap()).Parse(catalogTemplate))
)

type options struct {
	testFile   string
	readmeFile string
	catalog    string
	testDirs   []string
}

func Command() *cobra.Command {
	var options options
	cmd := &cobra.Command{
		Use:          "docs",
		Short:        "Generate tests documentation",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			tests, err := discovery.DiscoverTests(options.testFile, options.testDirs...)
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			if err := generateDocs(out, options.readmeFile, tests...); err != nil {
				return err
			}
			if options.catalog != "" {
				if err := generateCatalog(out, options.readmeFile, options.catalog, tests...); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&options.testFile, "test-file", "chainsaw-test.yaml", "Name of the test file")
	cmd.Flags().StringVar(&options.readmeFile, "readme-file", "README.md", "Name of the generated docs file")
	cmd.Flags().StringVar(&options.catalog, "catalog", "", "Path to the generated test catalog file")
	cmd.Flags().StringArrayVar(&options.testDirs, "test-dir", []string{}, "Directories containing test cases to run")
	return cmd
}

func generateDocs(out io.Writer, fileName string, tests ...discovery.Test) error {
	for _, test := range tests {
		if test.Err == nil {
			file, err := os.Create(filepath.Join(test.BasePath, fileName))
			if err != nil {
				return err
			}
			defer file.Close()
			output := file
			if err := docsTmpl.Execute(output, test); err != nil {
				return err
			}
		} else {
			fmt.Fprintf(out, "ERROR: failed to load test %s (%s)", test.BasePath, test.Err)
		}
	}
	return nil
}

func generateCatalog(out io.Writer, readme string, catalog string, tests ...discovery.Test) error {
	file, err := os.Create(catalog)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := catalogTmpl.Execute(file, map[string]any{
		"Readme": readme,
		"Tests":  tests,
	}); err != nil {
		return err
	}
	return nil
}
