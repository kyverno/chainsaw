package lint

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/validation"
	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func Command() *cobra.Command {
	var fileFlag string
	cmd := &cobra.Command{
		Use:       "lint [test|configuration]",
		Short:     "Lint a file or read from standard input",
		Long:      "Use chainsaw lint to lint a specific file or read from standard input for either test or configuration.",
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		ValidArgs: []string{"test", "configuration"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if fileFlag == "-" {
				input, err := io.ReadAll(cmd.InOrStdin())
				if err != nil {
					return err
				}
				return lintInput(input, args[0], "", cmd.OutOrStdout())
			} else if fileFlag != "" {
				input, err := os.ReadFile(fileFlag)
				if err != nil {
					return err
				}
				return lintInput(input, args[0], filepath.Ext(fileFlag), cmd.OutOrStdout())
			} else {
				return fmt.Errorf("no file or standard input specified")
			}
		},
	}
	cmd.Flags().StringVarP(&fileFlag, "file", "f", "", "Specify the file to lint or '-' for standard input")
	return cmd
}

func lintInput(input []byte, schema string, format string, writer io.Writer) error {
	fmt.Fprintln(writer, "Processing input...")
	if err := lintSchema(input, schema, format, writer); err != nil {
		return err
	}
	if schema == "test" {
		if err := lintBusinessLogic(input, writer); err != nil {
			return err
		}
	}
	fmt.Fprintln(writer, "The document is valid")
	return nil
}

func lintSchema(input []byte, schema string, format string, writer io.Writer) error {
	processor, err := getProcessor(format, input)
	if err != nil {
		return err
	}
	jsonInput, err := processor.ToJSON(input)
	if err != nil {
		return err
	}
	goschema, err := getScheme(schema)
	if err != nil {
		return err
	}
	schemaLoader := gojsonschema.NewBytesLoader(goschema)
	documentLoader := gojsonschema.NewBytesLoader(jsonInput)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}
	if !result.Valid() {
		fmt.Fprintln(writer, "The schema is not valid. See errors:")
		for _, desc := range result.Errors() {
			fmt.Fprintf(writer, "- %s\n", desc)
		}
		return fmt.Errorf("document is not valid")
	}
	return nil
}

func lintBusinessLogic(input []byte, writer io.Writer) error {
	test := &v1alpha1.Test{}
	if err := yaml.UnmarshalStrict(input, test); err != nil {
		return err
	}
	return validation.ValidateTest(test).ToAggregate()
}
