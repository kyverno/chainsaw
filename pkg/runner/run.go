package runner

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/kyverno/chainsaw/pkg/report"
	"github.com/kyverno/chainsaw/pkg/runner/clusters"
	"github.com/kyverno/chainsaw/pkg/runner/internal"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/processors"
	"github.com/kyverno/chainsaw/pkg/testing"
	"k8s.io/client-go/rest"
	"k8s.io/utils/clock"
)

type mainstart interface {
	Run() int
}

func Run(
	ctx context.Context,
	cfg *rest.Config,
	clock clock.PassiveClock,
	config model.Configuration,
	values map[string]any,
	tests ...discovery.Test,
) (model.SummaryResult, error) {
	return run(ctx, cfg, clock, config, nil, values, tests...)
}

func run(
	ctx context.Context,
	cfg *rest.Config,
	clock clock.PassiveClock,
	config model.Configuration,
	m mainstart,
	values map[string]any,
	tests ...discovery.Test,
) (model.SummaryResult, error) {
	var testsReport *report.Report
	if config.Report != nil && config.Report.Format != "" {
		testsReport = report.New(config.Report.Name)
	}
	tc, err := setupTestContext(ctx, values, cfg, config)
	if err != nil {
		return nil, err
	}
	if len(tests) == 0 {
		return tc, nil
	}
	if err := internal.SetupFlags(config); err != nil {
		return nil, err
	}
	internalTests := []testing.InternalTest{{
		Name: "chainsaw",
		F: func(t *testing.T) {
			t.Helper()
			t.Parallel()
			ctx := testing.IntoContext(ctx, t)
			ctx = logging.IntoContext(ctx, logging.NewLogger(t, clock, t.Name(), "@main"))
			processor := processors.NewTestsProcessor(config, clock, testsReport)
			processor.Run(ctx, tc, tests...)
		},
	}}
	deps := &internal.TestDeps{}
	if m == nil {
		m = testing.MainStart(deps, internalTests, nil, nil, nil)
	}
	// m.Run() returns:
	// - 0 if everything went well
	// - 1 if some of the tests failed
	// - 2 if running the tests was not possible
	// In our case, we consider an error only when running the tests was not possible.
	// For now, the case where some of the tests failed will be covered by the summary.
	if code := m.Run(); code > 1 {
		return tc.Summary, fmt.Errorf("testing framework exited with non zero code %d", code)
	}
	if testsReport != nil && config.Report != nil && config.Report.Format != "" {
		if err := testsReport.Save(config.Report.Format, config.Report.Path, config.Report.Name); err != nil {
			return tc.Summary, err
		}
	}
	return tc.Summary, nil
}

func setupTestContext(ctx context.Context, values any, cluster *rest.Config, config model.Configuration) (model.TestContext, error) {
	tc := model.EmptyContext()
	tc = model.WithValues(ctx, tc, values)
	if cluster != nil {
		cluster, err := clusters.NewClusterFromConfig(cluster)
		if err != nil {
			return tc, err
		}
		tc = tc.WithCluster(ctx, clusters.DefaultClient, cluster)
		_tc, err := model.UseCluster(ctx, tc, clusters.DefaultClient)
		if err != nil {
			return tc, err
		}
		tc = _tc
	}
	tc = model.WithClusters(ctx, tc, "", config.Clusters)
	tc = model.WithDefaultTimeouts(ctx, tc, config.Timeouts)
	tc = tc.WithCleanup(ctx, !config.Cleanup.SkipDelete)
	return tc, nil
}
