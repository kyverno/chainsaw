package test

import (
	"fmt"
	"os"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type options struct {
	config             string
	timeout            metav1.Duration
	testDirs           []string
	skipDelete         bool
	stopOnFirstFailure bool
	parallel           int
	reportFormat       string
	reportName         string
	namespace          string
	suppress           []string
	fullName           bool
	skipTestRegex      string
}

func Command() *cobra.Command {
	var options options
	cmd := &cobra.Command{
		Use:   "test [flags]... [test directories]...",
		Short: "Stronger tool for e2e testing",
		RunE: func(cmd *cobra.Command, _ []string) error {
			out := cmd.OutOrStdout()
			var configuration v1alpha1.Configuration
			// if no config file was provided, give a chance to the default config name
			if options.config == "" {
				if _, err := os.Stat(config.DefaultFileName); err == nil {
					options.config = config.DefaultFileName
					fmt.Fprintf(out, "No configuration provided but found default file: %s\n", options.config)
				}
			}
			// try to load configuration file
			if options.config != "" {
				fmt.Fprintf(out, "Loading config (%s)...\n", options.config)
				config, err := config.Load(options.config)
				if err != nil {
					return err
				}
				configuration = *config
			} else {
				fmt.Fprintln(out, "Running without configuration")
			}
			// flags take precedence over configuration file
			flags := cmd.Flags()
			if isSet(flags, "duration") {
				configuration.Spec.Timeout = &options.timeout
			}
			if isSet(flags, "testDirs") {
				configuration.Spec.TestDirs = options.testDirs
			}
			if isSet(flags, "skipDelete") {
				configuration.Spec.SkipDelete = options.skipDelete
			}
			if isSet(flags, "stopOnFirstFailure") {
				configuration.Spec.StopOnFirstFailure = options.stopOnFirstFailure
			}
			if isSet(flags, "parallel") {
				configuration.Spec.Parallel = options.parallel
			}
			if isSet(flags, "reportFormat") {
				configuration.Spec.ReportFormat = v1alpha1.ReportFormatType(options.reportFormat)
			}
			if isSet(flags, "reportName") {
				configuration.Spec.ReportName = options.reportName
			}
			if isSet(flags, "namespace") {
				configuration.Spec.Namespace = options.namespace
			}
			if isSet(flags, "suppress") {
				configuration.Spec.Suppress = options.suppress
			}
			if isSet(flags, "fullName") {
				configuration.Spec.FullName = options.fullName
			}
			if isSet(flags, "skipTestRegex") {
				configuration.Spec.SkipTestRegex = options.skipTestRegex
			}
			// run tests
			fmt.Fprintln(out, "Running tests...")
			fmt.Fprintln(out, "- TODO")
			// done
			fmt.Fprintln(out, "Done.")
			return nil
		},
	}
	cmd.Flags().DurationVar(&options.timeout.Duration, "duration", 30*time.Second, "The duration to use as default for configuration.")
	cmd.Flags().StringVar(&options.config, "config", "", "Chainsaw configuration file.")
	cmd.Flags().StringSliceVar(&options.testDirs, "testDirs", []string{}, "Directories containing test cases to run.")
	cmd.Flags().BoolVar(&options.skipDelete, "skipDelete", false, "If set, do not delete the resources after running the tests.")
	cmd.Flags().BoolVar(&options.stopOnFirstFailure, "stopOnFirstFailure", false, "Stop the test upon encountering the first failure.")
	cmd.Flags().IntVar(&options.parallel, "parallel", 8, "The maximum number of tests to run at once.")
	cmd.Flags().StringVar(&options.reportFormat, "reportFormat", "", "Test report format (JSON|XML|nil).")
	cmd.Flags().StringVar(&options.reportName, "reportName", "chainsaw-report", "The name of the report to create.")
	cmd.Flags().StringVar(&options.namespace, "namespace", "", "Namespace to use for tests.")
	cmd.Flags().StringSliceVar(&options.suppress, "suppress", []string{}, "Logs to suppress.")
	cmd.Flags().BoolVar(&options.fullName, "fullName", false, "Use full test case folder path instead of folder name.")
	cmd.Flags().StringVar(&options.skipTestRegex, "skipTestRegex", "", "Regular expression to skip tests based on.")
	// TODO: panic ?
	if err := cmd.MarkFlagFilename("config"); err != nil {
		panic(err)
	}
	return cmd
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
