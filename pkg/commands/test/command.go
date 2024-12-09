package test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha2"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/loaders/config"
	"github.com/kyverno/chainsaw/pkg/loaders/values"
	"github.com/kyverno/chainsaw/pkg/runner"
	"github.com/kyverno/chainsaw/pkg/runner/failer"
	flagutils "github.com/kyverno/chainsaw/pkg/utils/flag"
	fsutils "github.com/kyverno/chainsaw/pkg/utils/fs"
	restutils "github.com/kyverno/chainsaw/pkg/utils/rest"
	"github.com/kyverno/chainsaw/pkg/version"
	"github.com/kyverno/pkg/ext/output/color"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/clock"
	"k8s.io/utils/ptr"
)

type options struct {
	config                      string
	testFile                    string
	applyTimeout                metav1.Duration
	assertTimeout               metav1.Duration
	errorTimeout                metav1.Duration
	deleteTimeout               metav1.Duration
	cleanupTimeout              metav1.Duration
	execTimeout                 metav1.Duration
	testDirs                    []string
	skipDelete                  bool
	template                    bool
	defaultCompiler             string
	failFast                    bool
	parallel                    int
	repeatCount                 int
	reportFormat                string
	reportPath                  string
	reportName                  string
	namespace                   string
	deletionPropagationPolicy   string
	fullName                    bool
	excludeTestRegex            string
	includeTestRegex            string
	noColor                     bool
	kubeConfigOverrides         clientcmd.ConfigOverrides
	forceTerminationGracePeriod metav1.Duration
	delayBeforeCleanup          metav1.Duration
	selector                    []string
	noCluster                   bool
	pauseOnFailure              bool
	values                      []string
	clusters                    []string
	remarshal                   bool
	shardIndex                  int
	shardCount                  int
}

func Command() *cobra.Command {
	var options options
	cmd := &cobra.Command{
		Use:          "test [flags]... [test directories]...",
		Short:        "Run tests",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			color.Init(options.noColor, true)
			clock := clock.RealClock{}
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Version: %s\n", version.Version())
			var configuration v1alpha2.Configuration
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
				config, err := config.DefaultConfiguration()
				if err != nil {
					return err
				}
				configuration = *config
			}
			// flags take precedence over configuration file
			flags := cmd.Flags()
			if flagutils.IsSet(flags, "test-file") {
				configuration.Spec.Discovery.TestFile = options.testFile
			}
			if flagutils.IsSet(flags, "apply-timeout") {
				configuration.Spec.Timeouts.Apply = options.applyTimeout
			}
			if flagutils.IsSet(flags, "assert-timeout") {
				configuration.Spec.Timeouts.Assert = options.assertTimeout
			}
			if flagutils.IsSet(flags, "error-timeout") {
				configuration.Spec.Timeouts.Error = options.errorTimeout
			}
			if flagutils.IsSet(flags, "delete-timeout") {
				configuration.Spec.Timeouts.Delete = options.deleteTimeout
			}
			if flagutils.IsSet(flags, "cleanup-timeout") {
				configuration.Spec.Timeouts.Cleanup = options.cleanupTimeout
			}
			if flagutils.IsSet(flags, "exec-timeout") {
				configuration.Spec.Timeouts.Exec = options.execTimeout
			}
			if flagutils.IsSet(flags, "skip-delete") {
				configuration.Spec.Cleanup.SkipDelete = options.skipDelete
			}
			if flagutils.IsSet(flags, "template") {
				configuration.Spec.Templating.Enabled = options.template
			}
			if flagutils.IsSet(flags, "default-compiler") {
				configuration.Spec.Templating.Compiler = ptr.To(v1alpha1.Compiler(options.defaultCompiler))
			}
			if flagutils.IsSet(flags, "fail-fast") {
				configuration.Spec.Execution.FailFast = options.failFast
			}
			if flagutils.IsSet(flags, "parallel") {
				configuration.Spec.Execution.Parallel = &options.parallel
			}
			if flagutils.IsSet(flags, "repeat-count") {
				configuration.Spec.Execution.RepeatCount = &options.repeatCount
			}
			if flagutils.IsSet(flags, "report-format") {
				if configuration.Spec.Report == nil {
					configuration.Spec.Report = &v1alpha2.ReportOptions{
						Format: v1alpha2.JSONFormat,
						Name:   "chainsaw-report",
					}
				}
				configuration.Spec.Report.Format = v1alpha2.ReportFormatType(options.reportFormat)
			}
			if flagutils.IsSet(flags, "report-path") {
				if configuration.Spec.Report == nil {
					configuration.Spec.Report = &v1alpha2.ReportOptions{
						Format: v1alpha2.JSONFormat,
						Name:   "chainsaw-report",
					}
				}
				configuration.Spec.Report.Path = options.reportPath
			}
			if flagutils.IsSet(flags, "report-name") {
				if configuration.Spec.Report == nil {
					configuration.Spec.Report = &v1alpha2.ReportOptions{
						Format: v1alpha2.JSONFormat,
						Name:   "chainsaw-report",
					}
				}
				configuration.Spec.Report.Name = options.reportName
			}
			if flagutils.IsSet(flags, "namespace") {
				configuration.Spec.Namespace.Name = options.namespace
			}
			if flagutils.IsSet(flags, "deletion-propagation-policy") {
				configuration.Spec.Deletion.Propagation = metav1.DeletionPropagation(options.deletionPropagationPolicy)
			}
			if flagutils.IsSet(flags, "full-name") {
				configuration.Spec.Discovery.FullName = options.fullName
			}
			if flagutils.IsSet(flags, "include-test-regex") {
				configuration.Spec.Discovery.IncludeTestRegex = options.includeTestRegex
			}
			if flagutils.IsSet(flags, "exclude-test-regex") {
				configuration.Spec.Discovery.ExcludeTestRegex = options.excludeTestRegex
			}
			if flagutils.IsSet(flags, "force-termination-grace-period") {
				configuration.Spec.Execution.ForceTerminationGracePeriod = &options.forceTerminationGracePeriod
			}
			if flagutils.IsSet(flags, "cleanup-delay") {
				configuration.Spec.Cleanup.DelayBeforeCleanup = &options.delayBeforeCleanup
			}
			if flagutils.IsSet(flags, "cluster") {
				for _, cluster := range options.clusters {
					parts1 := strings.Split(cluster, "=")
					if len(parts1) != 2 {
						return fmt.Errorf("failed to decode cluster argument %s", cluster)
					}
					name := parts1[0]
					parts2 := strings.Split(parts1[1], ":")
					var c v1alpha1.Cluster
					if len(parts2) == 1 {
						c = v1alpha1.Cluster{
							Kubeconfig: parts2[0],
						}
					} else if len(parts2) == 2 {
						c = v1alpha1.Cluster{
							Kubeconfig: parts2[0],
							Context:    parts2[1],
						}
					} else {
						return fmt.Errorf("failed to decode cluster argument %s", cluster)
					}
					if configuration.Spec.Clusters == nil {
						configuration.Spec.Clusters = map[string]v1alpha1.Cluster{}
					}
					configuration.Spec.Clusters[name] = c
				}
			}
			options.testDirs = append(options.testDirs, args...)
			if len(options.testDirs) == 0 {
				options.testDirs = append(options.testDirs, ".")
			}
			// if pause on failure is set, force non concurrency
			if options.pauseOnFailure {
				configuration.Spec.Execution.Parallel = ptr.To(1)
			}
			fmt.Fprintf(out, "- Using test file: %s\n", configuration.Spec.Discovery.TestFile)
			fmt.Fprintf(out, "- TestDirs %v\n", options.testDirs)
			fmt.Fprintf(out, "- SkipDelete %v\n", configuration.Spec.Cleanup.SkipDelete)
			fmt.Fprintf(out, "- FailFast %v\n", configuration.Spec.Execution.FailFast)
			if configuration.Spec.Report != nil {
				fmt.Fprintf(out, "- ReportFormat '%v'\n", configuration.Spec.Report.Format)
				fmt.Fprintf(out, "- ReportName '%v'\n", configuration.Spec.Report.Name)
				if configuration.Spec.Report.Path != "" {
					fmt.Fprintf(out, "- ReportPath '%v'\n", configuration.Spec.Report.Path)
				}
			}
			fmt.Fprintf(out, "- Namespace '%v'\n", configuration.Spec.Namespace.Name)
			fmt.Fprintf(out, "- FullName %v\n", configuration.Spec.Discovery.FullName)
			fmt.Fprintf(out, "- IncludeTestRegex '%v'\n", configuration.Spec.Discovery.IncludeTestRegex)
			fmt.Fprintf(out, "- ExcludeTestRegex '%v'\n", configuration.Spec.Discovery.ExcludeTestRegex)
			fmt.Fprintf(out, "- ApplyTimeout %v\n", configuration.Spec.Timeouts.Apply.Duration)
			fmt.Fprintf(out, "- AssertTimeout %v\n", configuration.Spec.Timeouts.Assert.Duration)
			fmt.Fprintf(out, "- CleanupTimeout %v\n", configuration.Spec.Timeouts.Cleanup.Duration)
			fmt.Fprintf(out, "- DeleteTimeout %v\n", configuration.Spec.Timeouts.Delete.Duration)
			fmt.Fprintf(out, "- ErrorTimeout %v\n", configuration.Spec.Timeouts.Error.Duration)
			fmt.Fprintf(out, "- ExecTimeout %v\n", configuration.Spec.Timeouts.Exec.Duration)
			fmt.Fprintf(out, "- DeletionPropagationPolicy %v\n", configuration.Spec.Deletion.Propagation)
			if configuration.Spec.Execution.Parallel != nil && *configuration.Spec.Execution.Parallel > 0 {
				fmt.Fprintf(out, "- Parallel %d\n", *configuration.Spec.Execution.Parallel)
			}
			if configuration.Spec.Execution.RepeatCount != nil {
				fmt.Fprintf(out, "- RepeatCount %v\n", *configuration.Spec.Execution.RepeatCount)
			}
			if configuration.Spec.Execution.ForceTerminationGracePeriod != nil {
				fmt.Fprintf(out, "- ForceTerminationGracePeriod %v\n", configuration.Spec.Execution.ForceTerminationGracePeriod.Duration)
			}
			if configuration.Spec.Cleanup.DelayBeforeCleanup != nil {
				fmt.Fprintf(out, "- DelayBeforeCleanup %v\n", configuration.Spec.Cleanup.DelayBeforeCleanup.Duration)
			}
			if len(options.selector) != 0 {
				fmt.Fprintf(out, "- Selector %v\n", options.selector)
			}
			if len(options.values) != 0 {
				fmt.Fprintf(out, "- Values %v\n", options.values)
			}
			fmt.Fprintf(out, "- Template %v\n", configuration.Spec.Templating.Enabled)
			if configuration.Spec.Templating.Compiler != nil {
				fmt.Fprintf(out, "- Default compiler %v\n", *configuration.Spec.Templating.Compiler)
			}
			if len(configuration.Spec.Clusters) != 0 {
				fmt.Fprintf(out, "- Clusters %v\n", configuration.Spec.Clusters)
			}
			if options.remarshal {
				fmt.Fprintf(out, "- Remarshal %v\n", options.remarshal)
			}
			fmt.Fprintf(out, "- NoCluster %v\n", options.noCluster)
			fmt.Fprintf(out, "- PauseOnFailure %v\n", options.pauseOnFailure)
			if options.shardCount > 0 {
				fmt.Fprintf(out, "- Shard %v / %v\n", options.shardIndex, options.shardCount)
			}
			// load tests
			fmt.Fprintln(out, "Loading tests...")
			if err := fsutils.CheckFolders(options.testDirs...); err != nil {
				return err
			}
			var selector labels.Selector
			if len(options.selector) != 0 {
				parsed, err := labels.Parse(strings.Join(options.selector, ","))
				if err != nil {
					return err
				}
				selector = parsed
			}
			tests, err := discovery.DiscoverTests(configuration.Spec.Discovery.TestFile, selector, options.remarshal, options.testDirs...)
			if err != nil {
				return err
			}
			// TODO: we may want to find a sort key here ?
			if options.shardCount > 0 && options.shardIndex < options.shardCount {
				shardLen := float64(len(tests)) / float64(options.shardCount)
				shardStart := int(shardLen * float64(options.shardIndex))
				shardEnd := int(shardLen * float64(options.shardIndex+1))
				if options.shardIndex == options.shardCount-1 {
					shardEnd = len(tests)
				}
				tests = tests[shardStart:shardEnd]
			}
			var testToRun []discovery.Test
			for _, test := range tests {
				if test.Err != nil {
					fmt.Fprintf(out, "- %s (%s) - (%s)\n", test.Test.Name, test.BasePath, test.Err)
				} else {
					fmt.Fprintf(out, "- %s (%s)\n", test.Test.Name, test.BasePath)
					testToRun = append(testToRun, test)
				}
			}
			// load values
			fmt.Fprintln(out, "Loading values...")
			values, err := values.Load(options.values...)
			if err != nil {
				return err
			}
			// run tests
			fmt.Fprintln(out, "Running tests...")
			var restConfig *rest.Config
			if !options.noCluster {
				cfg, err := restutils.DefaultConfig(options.kubeConfigOverrides)
				if err != nil {
					return err
				}
				restConfig = cfg
			}
			runner := runner.New(clock, failer.New(options.pauseOnFailure))
			summary, err := runner.Run(context.Background(), restConfig, configuration.Spec, values, testToRun...)
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
	config, err := config.DefaultConfiguration()
	if err != nil {
		panic(err)
	}
	// config
	cmd.Flags().StringVar(&options.config, "config", "", "Chainsaw configuration file")
	cmd.Flags().StringSliceVar(&options.testDirs, "test-dir", nil, "Directories containing test cases to run")
	clientcmd.BindOverrideFlags(&options.kubeConfigOverrides, cmd.Flags(), clientcmd.RecommendedConfigOverrideFlags("kube-"))
	// timeouts options
	cmd.Flags().DurationVar(&options.applyTimeout.Duration, "apply-timeout", config.Spec.Timeouts.Apply.Duration, "The apply timeout to use as default for configuration")
	cmd.Flags().DurationVar(&options.assertTimeout.Duration, "assert-timeout", config.Spec.Timeouts.Assert.Duration, "The assert timeout to use as default for configuration")
	cmd.Flags().DurationVar(&options.cleanupTimeout.Duration, "cleanup-timeout", config.Spec.Timeouts.Cleanup.Duration, "The cleanup timeout to use as default for configuration")
	cmd.Flags().DurationVar(&options.deleteTimeout.Duration, "delete-timeout", config.Spec.Timeouts.Delete.Duration, "The delete timeout to use as default for configuration")
	cmd.Flags().DurationVar(&options.errorTimeout.Duration, "error-timeout", config.Spec.Timeouts.Error.Duration, "The error timeout to use as default for configuration")
	cmd.Flags().DurationVar(&options.execTimeout.Duration, "exec-timeout", config.Spec.Timeouts.Exec.Duration, "The exec timeout to use as default for configuration")
	// discovery options
	cmd.Flags().StringVar(&options.testFile, "test-file", "chainsaw-test", "Name of the test file")
	cmd.Flags().BoolVar(&options.fullName, "full-name", false, "Use full test case folder path instead of folder name")
	cmd.Flags().StringVar(&options.includeTestRegex, "include-test-regex", "", "Regular expression to include tests")
	cmd.Flags().StringVar(&options.excludeTestRegex, "exclude-test-regex", "", "Regular expression to exclude tests")
	// execution options
	cmd.Flags().BoolVar(&options.failFast, "fail-fast", false, "Stop the test upon encountering the first failure")
	cmd.Flags().IntVar(&options.parallel, "parallel", 0, "The maximum number of tests to run at once")
	cmd.Flags().IntVar(&options.repeatCount, "repeat-count", 1, "Number of times to repeat each test")
	cmd.Flags().DurationVar(&options.forceTerminationGracePeriod.Duration, "force-termination-grace-period", 0, "If specified, overrides termination grace periods in applicable resources")
	// namespace options
	cmd.Flags().StringVar(&options.namespace, "namespace", "", "Namespace to use for tests")
	// templating options
	cmd.Flags().BoolVar(&options.template, "template", config.Spec.Templating.Enabled, "If set, resources will be considered for templating")
	cmd.Flags().StringVar(&options.defaultCompiler, "default-compiler", "", "If set, configures the default compiler (jp or cel)")
	// cleanup options
	cmd.Flags().BoolVar(&options.skipDelete, "skip-delete", false, "If set, do not delete the resources after running the tests")
	cmd.Flags().DurationVar(&options.delayBeforeCleanup.Duration, "cleanup-delay", 0, "Adds a delay between the time a test ends and the time cleanup starts")
	// deletion options
	cmd.Flags().StringVar(&options.deletionPropagationPolicy, "deletion-propagation-policy", "Background", "The deletion propagation policy (Foreground|Background|Orphan)")
	// error options
	// reporting options
	cmd.Flags().StringVar(&options.reportFormat, "report-format", "", "Test report format (JSON|XML|JUNIT-TEST|JUNIT-STEP|JUNIT-OPERATION)")
	cmd.Flags().StringVar(&options.reportName, "report-name", "chainsaw-report", "The name of the report to create")
	cmd.Flags().StringVar(&options.reportPath, "report-path", "", "The path of the report to create")
	// multi-cluster options
	cmd.Flags().StringSliceVar(&options.clusters, "cluster", nil, "Register cluster (format <cluster name>=<kubeconfig path>:[context name])")
	// pause options
	cmd.Flags().BoolVar(&options.pauseOnFailure, "pause-on-failure", false, "Pause test execution failure (implies no concurrency)")
	// no cluster options
	cmd.Flags().BoolVar(&options.noCluster, "no-cluster", false, "Runs without cluster")
	// label selectors
	cmd.Flags().StringSliceVar(&options.selector, "selector", nil, "Selector (label query) to filter on")
	// external values
	cmd.Flags().StringSliceVar(&options.values, "values", nil, "Values passed to the tests")
	// sharding
	cmd.Flags().IntVar(&options.shardIndex, "shard-index", 0, "Current shard index (if `--shard-count` > 0)")
	cmd.Flags().IntVar(&options.shardCount, "shard-count", 0, "Number of shards")
	// others
	cmd.Flags().BoolVar(&options.noColor, "no-color", false, "Removes output colors")
	cmd.Flags().BoolVar(&options.remarshal, "remarshal", false, "Remarshals tests yaml to apply anchors before parsing")
	if err := cmd.MarkFlagFilename("config"); err != nil {
		panic(err)
	}
	return cmd
}
