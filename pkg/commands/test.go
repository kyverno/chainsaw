package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	v1alpha1 "github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	utils "github.com/kyverno/chainsaw/pkg/utils"
)

func chainsawTestCMD() *cobra.Command {
	configPath := ""
	duration := 30

	options := v1alpha1.Configuration{}

	testCMD := &cobra.Command{
		Use:   "test [flags]... [test directories]...",
		Short: "Stronger tool for e2e testing",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.Flags()

			if configPath == "" {
				if _, err := os.Stat("chainsaw-test.yaml"); err == nil {
					configPath = "chainsaw-test.yaml"
				} else {
					log.Println("running without a 'chainsaw-test.yaml' configuration")
				}
			}

			if configPath != "" {
				objects, err := utils.LoadYAMLFromFile(configPath)
				if err != nil {
					return err
				}

				for _, obj := range objects {
					kind := obj.GetObjectKind().GroupVersionKind().Kind

					if kind == "Configuration" {
						switch ts := obj.(type) {
						case *v1alpha1.Configuration:
							options = *ts
						case *unstructured.Unstructured:
							log.Println(fmt.Errorf("bad configuration in file %q", configPath))
						}
					} else {
						log.Println(fmt.Errorf("unknown object type: %s", kind))
					}
				}
			}

			if isSet(flags, "duraiton") {
				options.Spec.Duration = duration
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Running tests...")
		},
	}

	testCMD.Flags().IntVar(&duration, "duration", 30, "The duration to use as default for configuration.")

	return testCMD
}

// isSet returns true if a flag is set on the command line.
func isSet(flagSet *pflag.FlagSet, name string) bool {
	found := false

	flagSet.Visit(func(flag *pflag.Flag) {
		if flag.Name == name {
			found = true
		}
	})

	return found
}
