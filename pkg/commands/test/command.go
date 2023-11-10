package test

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/config"
	"github.com/kyverno/chainsaw/pkg/data"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/runner"
	"github.com/kyverno/chainsaw/pkg/runner/timeout"
	flagutils "github.com/kyverno/chainsaw/pkg/utils/flag"
	restutils "github.com/kyverno/chainsaw/pkg/utils/rest"
	"github.com/kyverno/kyverno/ext/output/color"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/clock"
)

type options struct {
	config              string
	testFile            string
	applyTimeout        metav1.Duration
	assertTimeout       metav1.Duration
	errorTimeout        metav1.Duration
	deleteTimeout       metav1.Duration
	cleanupTimeout      metav1.Duration
	execTimeout         metav1.Duration
	testDirs            []string
	skipDelete          bool
	failFast            bool
	parallel            int
	repeatCount         int
	reportFormat        string
	reportName          string
	namespace           string
	fullName            bool
	excludeTestRegex    string
	includeTestRegex    string
	noColor             bool
	kubeConfigOverrides clientcmd.ConfigOverrides
}

func Command() *cobra.Command {
	var options options
	cmd := &cobra.Command{
		Use:          "test [flags]... [test directories]...",
		Short:        "Run tests",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			color.Init(options.noColor, true)
			clock := clock.RealClock{}
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
				fmt.Fprintln(out, "Loading default configuration...")
				bytes, err := fs.ReadFile(data.Config(), filepath.Join("config", "default.yaml"))
				if err != nil {
					return err
				}
				config, err := config.LoadBytes(bytes)
				if err != nil {
					return err
				}
				configuration = *config
			}
			// flags take precedence over configuration file
			flags := cmd.Flags()
			if flagutils.IsSet(flags, "test-file") {
				configuration.Spec.TestFile = options.testFile
			}
			if flagutils.IsSet(flags, "apply-timeout") {
				configuration.Spec.Timeouts.Apply = &options.applyTimeout
			}
			if flagutils.IsSet(flags, "assert-timeout") {
				configuration.Spec.Timeouts.Assert = &options.assertTimeout
			}
			if flagutils.IsSet(flags, "error-timeout") {
				configuration.Spec.Timeouts.Error = &options.errorTimeout
			}
			if flagutils.IsSet(flags, "delete-timeout") {
				configuration.Spec.Timeouts.Delete = &options.deleteTimeout
			}
			if flagutils.IsSet(flags, "cleanup-timeout") {
				configuration.Spec.Timeouts.Cleanup = &options.cleanupTimeout
			}
			if flagutils.IsSet(flags, "exec-timeout") {
				configuration.Spec.Timeouts.Exec = &options.execTimeout
			}
			if flagutils.IsSet(flags, "test-dir") {
				configuration.Spec.TestDirs = options.testDirs
			}
			if flagutils.IsSet(flags, "skip-delete") {
				configuration.Spec.SkipDelete = options.skipDelete
			}
			if flagutils.IsSet(flags, "fail-fast") {
				configuration.Spec.FailFast = options.failFast
			}
			if flagutils.IsSet(flags, "parallel") {
				configuration.Spec.Parallel = &options.parallel
			}
			if flagutils.IsSet(flags, "repeat-count") {
				configuration.Spec.RepeatCount = &options.repeatCount
			}
			if flagutils.IsSet(flags, "report-format") {
				configuration.Spec.ReportFormat = v1alpha1.ReportFormatType(options.reportFormat)
			}
			if flagutils.IsSet(flags, "report-name") {
				configuration.Spec.ReportName = options.reportName
			}
			if flagutils.IsSet(flags, "namespace") {
				configuration.Spec.Namespace = options.namespace
			}
			if flagutils.IsSet(flags, "full-name") {
				configuration.Spec.FullName = options.fullName
			}
			if flagutils.IsSet(flags, "include-test-regex") {
				configuration.Spec.IncludeTestRegex = options.includeTestRegex
			}
			if flagutils.IsSet(flags, "exclude-test-regex") {
				configuration.Spec.ExcludeTestRegex = options.excludeTestRegex
			}
			fmt.Fprintf(out, "- Using test file: %s\n", configuration.Spec.TestFile)
			fmt.Fprintf(out, "- TestDirs %v\n", configuration.Spec.TestDirs)
			fmt.Fprintf(out, "- SkipDelete %v\n", configuration.Spec.SkipDelete)
			fmt.Fprintf(out, "- FailFast %v\n", configuration.Spec.FailFast)
			fmt.Fprintf(out, "- ReportFormat '%v'\n", configuration.Spec.ReportFormat)
			fmt.Fprintf(out, "- ReportName '%v'\n", configuration.Spec.ReportName)
			fmt.Fprintf(out, "- Namespace '%v'\n", configuration.Spec.Namespace)
			fmt.Fprintf(out, "- FullName %v\n", configuration.Spec.FullName)
			fmt.Fprintf(out, "- IncludeTestRegex '%v'\n", configuration.Spec.IncludeTestRegex)
			fmt.Fprintf(out, "- ExcludeTestRegex '%v'\n", configuration.Spec.ExcludeTestRegex)
			if configuration.Spec.Parallel != nil && *configuration.Spec.Parallel > 0 {
				fmt.Fprintf(out, "- Parallel %d\n", *configuration.Spec.Parallel)
			}
			if configuration.Spec.RepeatCount != nil {
				fmt.Fprintf(out, "- RepeatCount %v\n", *configuration.Spec.RepeatCount)
			}
			if configuration.Spec.Timeouts.Apply != nil {
				fmt.Fprintf(out, "- ApplyTimeout %v\n", configuration.Spec.Timeouts.Apply.Duration)
			}
			if configuration.Spec.Timeouts.Assert != nil {
				fmt.Fprintf(out, "- AssertTimeout %v\n", configuration.Spec.Timeouts.Assert.Duration)
			}
			if configuration.Spec.Timeouts.Error != nil {
				fmt.Fprintf(out, "- ErrorTimeout %v\n", configuration.Spec.Timeouts.Error.Duration)
			}
			if configuration.Spec.Timeouts.Delete != nil {
				fmt.Fprintf(out, "- DeleteTimeout %v\n", configuration.Spec.Timeouts.Delete.Duration)
			}
			if configuration.Spec.Timeouts.Cleanup != nil {
				fmt.Fprintf(out, "- CleanupTimeout %v\n", configuration.Spec.Timeouts.Cleanup.Duration)
			}
			if configuration.Spec.Timeouts.Exec != nil {
				fmt.Fprintf(out, "- ExecTimeout %v\n", configuration.Spec.Timeouts.Exec.Duration)
			}
			// loading tests
			fmt.Fprintln(out, "Loading tests...")
			tests, err := discovery.DiscoverTests(configuration.Spec.TestFile, configuration.Spec.TestDirs...)
			if err != nil {
				return err
			}
			var testToRun []discovery.Test
			for _, test := range tests {
				if test.Err != nil {
					fmt.Fprintf(out, "- %s (%s) - (%s)\n", test.Name, test.BasePath, test.Err)
				} else {
					fmt.Fprintf(out, "- %s (%s)\n", test.Name, test.BasePath)
					testToRun = append(testToRun, test)
				}
			}
			// run tests
			fmt.Fprintln(out, "Running tests...")
			cfg, err := restutils.Config(options.kubeConfigOverrides)
			if err != nil {
				return err
			}
			summary, err := runner.Run(cfg, clock, configuration.Spec, testToRun...)
			if summary != nil {
				fmt.Fprintln(out, "Tests Summary...")
				fmt.Fprintln(out, "- Passed  tests", summary.Passed())
				fmt.Fprintln(out, "- Failed  tests", summary.Failed())
				fmt.Fprintln(out, "- Skipped tests", summary.Skipped())
			}
			if err != nil {
				fmt.Fprintln(out, "Done with error.")
			} else if summary != nil && summary.Failed() > 0 {
				fmt.Fprintln(out, "Done with failures.")
				err = errors.New("some tests failed")
			} else {
				fmt.Fprintln(out, "Done.")
			}
			return err
		},
	}
	cmd.Flags().StringVar(&options.testFile, "test-file", "chainsaw-test.yaml", "Name of the test file.")
	cmd.Flags().DurationVar(&options.applyTimeout.Duration, "apply-timeout", timeout.DefaultApplyTimeout, "The apply timeout to use as default for configuration.")
	cmd.Flags().DurationVar(&options.assertTimeout.Duration, "assert-timeout", timeout.DefaultAssertTimeout, "The assert timeout to use as default for configuration.")
	cmd.Flags().DurationVar(&options.errorTimeout.Duration, "error-timeout", timeout.DefaultErrorTimeout, "The error timeout to use as default for configuration.")
	cmd.Flags().DurationVar(&options.deleteTimeout.Duration, "delete-timeout", timeout.DefaultDeleteTimeout, "The delete timeout to use as default for configuration.")
	cmd.Flags().DurationVar(&options.cleanupTimeout.Duration, "cleanup-timeout", timeout.DefaultCleanupTimeout, "The cleanup timeout to use as default for configuration.")
	cmd.Flags().DurationVar(&options.execTimeout.Duration, "exec-timeout", timeout.DefaultExecTimeout, "The exec timeout to use as default for configuration.")
	cmd.Flags().StringVar(&options.config, "config", "", "Chainsaw configuration file.")
	cmd.Flags().StringArrayVar(&options.testDirs, "test-dir", []string{}, "Directories containing test cases to run.")
	cmd.Flags().BoolVar(&options.skipDelete, "skip-delete", false, "If set, do not delete the resources after running the tests.")
	cmd.Flags().BoolVar(&options.failFast, "fail-fast", false, "Stop the test upon encountering the first failure.")
	cmd.Flags().IntVar(&options.parallel, "parallel", 0, "The maximum number of tests to run at once.")
	cmd.Flags().IntVar(&options.repeatCount, "repeat-count", 1, "Number of times to repeat each test.")
	cmd.Flags().StringVar(&options.reportFormat, "report-format", "", "Test report format (JSON|XML|nil).")
	cmd.Flags().StringVar(&options.reportName, "report-name", "chainsaw-report", "The name of the report to create.")
	cmd.Flags().StringVar(&options.namespace, "namespace", "", "Namespace to use for tests.")
	cmd.Flags().BoolVar(&options.fullName, "full-name", false, "Use full test case folder path instead of folder name.")
	cmd.Flags().StringVar(&options.includeTestRegex, "include-test-regex", "", "Regular expression to include tests.")
	cmd.Flags().StringVar(&options.excludeTestRegex, "exclude-test-regex", "", "Regular expression to exclude tests.")
	cmd.Flags().BoolVar(&options.noColor, "no-color", false, "Removes output colors.")
	clientcmd.BindOverrideFlags(&options.kubeConfigOverrides, cmd.Flags(), clientcmd.RecommendedConfigOverrideFlags("kube-"))
	if err := cmd.MarkFlagFilename("config"); err != nil {
		panic(err)
	}
	return cmd
}
