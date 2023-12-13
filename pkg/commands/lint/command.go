package lint

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
)

func Command() *cobra.Command {
	var fileFlag string
	var stdInFlag bool

	cmd := &cobra.Command{
		Use:   "lint",
		Short: "Lint a file or read from standard input",
		Long:  `Use chainsaw lint to lint a specific file or read from standard input.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if stdInFlag {
				input, err := io.ReadAll(os.Stdin)
				if err != nil {
					return err
				}
				return lintInput(input)
			} else if fileFlag != "" {
				input, err := os.ReadFile(fileFlag)
				if err != nil {
					return err
				}
				return lintInput(input)
			} else {
				return fmt.Errorf("either --file or --std-in must be specified")
			}
		},
	}

	cmd.Flags().StringVarP(&fileFlag, "file", "f", "", "Specify the file to lint")
	cmd.Flags().BoolVar(&stdInFlag, "std-in", false, "Read from standard input")

	return cmd
}

func lintInput(input []byte) error {
	fmt.Println("Processing input...")
	schemaLoader := gojsonschema.NewReferenceLoader("")
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
