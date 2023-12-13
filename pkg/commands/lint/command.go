package lint

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
)

func getTestSchema() string {
	return filepath.Join("..", "..", "..", ".schemas", "json", "test-chainsaw-v1alpha1.json")
}

func getConfigurationSchema() string {
	return filepath.Join("..", "..", "..", ".schemas", "json", "configuration-chainsaw-v1alpha1.json")
}

func Command() *cobra.Command {
	var fileFlag string
	var stdInFlag bool

	cmd := &cobra.Command{
		Use:       "lint [test|configuration]",
		Short:     "Lint a file or read from standard input",
		Long:      `Use chainsaw lint to lint a specific file or read from standard input for either test or configuration.`,
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		ValidArgs: []string{"test", "configuration"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if stdInFlag {
				input, err := io.ReadAll(os.Stdin)
				if err != nil {
					return err
				}
				return lintInput(input, args[0])
			} else if fileFlag != "" {
				input, err := os.ReadFile(fileFlag)
				if err != nil {
					return err
				}
				return lintInput(input, args[0])
			} else {
				return fmt.Errorf("either --file or --std-in must be specified")
			}
		},
	}

	cmd.Flags().StringVarP(&fileFlag, "file", "f", "", "Specify the file to lint")
	cmd.Flags().BoolVar(&stdInFlag, "std-in", false, "Read from standard input")

	return cmd
}

func lintInput(input []byte, schema string) error {
	fmt.Println("Processing input...")
	schemaLoader := gojsonschema.NewReferenceLoader(getScheme(schema))
	documentLoader := gojsonschema.NewBytesLoader(input)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		fmt.Println("The schema is not valid. See errors :")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
		return fmt.Errorf("document is not valid")
	}

	return nil
}

func getScheme(schema string) string {
	switch schema {
	case "configuration":
		return getConfigurationSchema()
	default:
		return getTestSchema()
	}
}
