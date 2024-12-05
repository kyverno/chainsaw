package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/engine"
	"github.com/kyverno/chainsaw/pkg/engine/clusters"
	enginecontext "github.com/kyverno/chainsaw/pkg/engine/context"
	"github.com/kyverno/chainsaw/pkg/engine/logging"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/kyverno/chainsaw/pkg/report"
	"github.com/kyverno/chainsaw/pkg/runner/internal"
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
			ctx = logging.IntoContext(ctx, logging.NewLogger(t, clock, t.Name(), "@chainsaw"))
			processor := processors.NewTestsProcessor(config, clock)
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
	if config.Report != nil && config.Report.Format != "" {
		tc.Report.EndTime = time.Now()
		if err := report.Save(tc.Report, config.Report.Format, config.Report.Path, config.Report.Name); err != nil {
			return tc.Summary, err
		}
	}
	return tc.Summary, nil
}

func setupTestContext(ctx context.Context, values any, cluster *rest.Config, config model.Configuration) (engine.Context, error) {
	tc := enginecontext.EmptyContext()
	// cleanup options
	tc = tc.WithSkipDelete(ctx, config.Cleanup.SkipDelete)
	// templating options
	tc = tc.WithTemplating(ctx, config.Templating.Enabled)
	if config.Templating.Compiler != nil {
		tc = tc.WithDefaultCompiler(string(*config.Templating.Compiler))
	}
	// deletion options
	tc = tc.WithDeletionPropagation(ctx, config.Deletion.Propagation)
	// values
	tc = engine.WithValues(ctx, tc, values)
	// default cluster
	if cluster != nil {
		cluster, err := clusters.NewClusterFromConfig(cluster)
		if err != nil {
			return tc, err
		}
		tc = tc.WithCluster(ctx, clusters.DefaultClient, cluster)
		return engine.WithCurrentCluster(ctx, tc, clusters.DefaultClient)
	}
	return tc, nil
}
