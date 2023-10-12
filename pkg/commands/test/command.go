package test

import (
	"fmt"
	"os"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/config"
	"github.com/kyverno/chainsaw/pkg/runner"
	flagutils "github.com/kyverno/chainsaw/pkg/utils/flag"
	"github.com/spf13/cobra"
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
		Use:          "test [flags]... [test directories]...",
		Short:        "Stronger tool for e2e testing",
		SilenceUsage: true,
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
			if flagutils.IsSet(flags, "duration") {
				configuration.Spec.Timeout = &options.timeout
			}
			if flagutils.IsSet(flags, "testDirs") {
				configuration.Spec.TestDirs = options.testDirs
			}
			if flagutils.IsSet(flags, "skipDelete") {
				configuration.Spec.SkipDelete = options.skipDelete
			}
			if flagutils.IsSet(flags, "stopOnFirstFailure") {
				configuration.Spec.StopOnFirstFailure = options.stopOnFirstFailure
			}
			if flagutils.IsSet(flags, "parallel") {
				configuration.Spec.Parallel = options.parallel
			}
			if flagutils.IsSet(flags, "reportFormat") {
				configuration.Spec.ReportFormat = v1alpha1.ReportFormatType(options.reportFormat)
			}
			if flagutils.IsSet(flags, "reportName") {
				configuration.Spec.ReportName = options.reportName
			}
			if flagutils.IsSet(flags, "namespace") {
				configuration.Spec.Namespace = options.namespace
			}
			if flagutils.IsSet(flags, "suppress") {
				configuration.Spec.Suppress = options.suppress
			}
			if flagutils.IsSet(flags, "fullName") {
				configuration.Spec.FullName = options.fullName
			}
			if flagutils.IsSet(flags, "skipTestRegex") {
				configuration.Spec.SkipTestRegex = options.skipTestRegex
			}
			// loading tests
			fmt.Fprintln(out, "Loading tests...")
			fmt.Fprintln(out, "- TODO")
			// TODO: load tests
			test := v1alpha1.Test{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "chainsaw.kyverno.io/v1alpha1",
					Kind:       "Test",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-1",
				},
				Spec: v1alpha1.TestSpec{
					Steps: []v1alpha1.TestStepSpec{{
						Apply: []v1alpha1.Apply{{
							File: "foo.yaml",
						}},
					}, {
						Assert: []v1alpha1.Assert{{
							File: "bar.yaml",
						}},
					}},
				},
			}
			// run tests
			fmt.Fprintln(out, "Running tests...")
			if _, err := runner.Run(test); err != nil {
				return err
			}
			// done
			fmt.Fprintln(out, "Done.")
			return nil
		},
	}
	cmd.Flags().DurationVar(&options.timeout.Duration, "timeout", 30*time.Second, "The timeout to use as default for configuration.")
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
