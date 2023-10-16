package test

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/config"
	"github.com/kyverno/chainsaw/pkg/data"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/runner"
	flagutils "github.com/kyverno/chainsaw/pkg/utils/flag"
	restutils "github.com/kyverno/chainsaw/pkg/utils/rest"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

type options struct {
	config              string
	timeout             metav1.Duration
	testDirs            []string
	skipDelete          bool
	failFast            bool
	parallel            int
	reportFormat        string
	reportName          string
	namespace           string
	suppress            []string
	fullName            bool
	skipTestRegex       string
	kubeConfigOverrides clientcmd.ConfigOverrides
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
			if flagutils.IsSet(flags, "timeout") {
				configuration.Spec.Timeout = &options.timeout
			}
			if flagutils.IsSet(flags, "test-dir") {
				configuration.Spec.TestDirs = options.testDirs
			}
			if flagutils.IsSet(flags, "skip-delete") {
				configuration.Spec.SkipDelete = options.skipDelete
			}
			if flagutils.IsSet(flags, "stop-on-first-failure") {
				configuration.Spec.FailFast = options.failFast
			}
			if flagutils.IsSet(flags, "parallel") {
				configuration.Spec.Parallel = options.parallel
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
			if flagutils.IsSet(flags, "suppress") {
				configuration.Spec.Suppress = options.suppress
			}
			if flagutils.IsSet(flags, "full-name") {
				configuration.Spec.FullName = options.fullName
			}
			if flagutils.IsSet(flags, "skip-test-regex") {
				configuration.Spec.SkipTestRegex = options.skipTestRegex
			}
			fmt.Fprintf(out, "- Timeout %v\n", configuration.Spec.Timeout.Duration)
			fmt.Fprintf(out, "- TestDirs %v\n", configuration.Spec.TestDirs)
			fmt.Fprintf(out, "- SkipDelete %v\n", configuration.Spec.SkipDelete)
			fmt.Fprintf(out, "- FailFast %v\n", configuration.Spec.FailFast)
			fmt.Fprintf(out, "- Parallel %v\n", configuration.Spec.Parallel)
			fmt.Fprintf(out, "- ReportFormat '%v'\n", configuration.Spec.ReportFormat)
			fmt.Fprintf(out, "- ReportName '%v'\n", configuration.Spec.ReportName)
			fmt.Fprintf(out, "- Namespace '%v'\n", configuration.Spec.Namespace)
			fmt.Fprintf(out, "- Suppress %v\n", configuration.Spec.Suppress)
			fmt.Fprintf(out, "- FullName %v\n", configuration.Spec.FullName)
			fmt.Fprintf(out, "- SkipTestRegex '%v'\n", configuration.Spec.SkipTestRegex)
			// loading tests
			fmt.Fprintln(out, "Loading tests...")
			tests, err := discovery.DiscoverTests("chainsaw-test.yaml", configuration.Spec.TestDirs...)
			if err != nil {
				return err
			}
			for _, test := range tests {
				fmt.Fprintf(out, "- %s (%s)\n", test.Name, test.BasePath)
			}
			// run tests
			fmt.Fprintln(out, "Running tests...")
			cfg, err := restutils.Config(options.kubeConfigOverrides)
			if err != nil {
				return err
			}
			if _, err := runner.Run(cfg, configuration.Spec, tests...); err != nil {
				return err
			}
			// done
			fmt.Fprintln(out, "Done.")
			return nil
		},
	}
	cmd.Flags().DurationVar(&options.timeout.Duration, "timeout", 30*time.Second, "The timeout to use as default for configuration.")
	cmd.Flags().StringVar(&options.config, "config", "", "Chainsaw configuration file.")
	cmd.Flags().StringArrayVar(&options.testDirs, "test-dir", []string{}, "Directories containing test cases to run.")
	cmd.Flags().BoolVar(&options.skipDelete, "skip-delete", false, "If set, do not delete the resources after running the tests.")
	cmd.Flags().BoolVar(&options.failFast, "stop-on-first-failure", false, "Stop the test upon encountering the first failure.")
	cmd.Flags().IntVar(&options.parallel, "parallel", 8, "The maximum number of tests to run at once.")
	cmd.Flags().StringVar(&options.reportFormat, "report-format", "", "Test report format (JSON|XML|nil).")
	cmd.Flags().StringVar(&options.reportName, "report-name", "chainsaw-report", "The name of the report to create.")
	cmd.Flags().StringVar(&options.namespace, "namespace", "", "Namespace to use for tests.")
	cmd.Flags().StringArrayVar(&options.suppress, "suppress", []string{}, "Logs to suppress.")
	cmd.Flags().BoolVar(&options.fullName, "full-name", false, "Use full test case folder path instead of folder name.")
	cmd.Flags().StringVar(&options.skipTestRegex, "skip-test-regex", "", "Regular expression to skip tests based on.")
	clientcmd.BindOverrideFlags(&options.kubeConfigOverrides, cmd.Flags(), clientcmd.RecommendedConfigOverrideFlags("kube-"))
	if err := cmd.MarkFlagFilename("config"); err != nil {
		panic(err)
	}
	return cmd
}
