package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/engine/clusters"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/kyverno/chainsaw/pkg/report"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
	"github.com/kyverno/chainsaw/pkg/runner/internal"
	"github.com/kyverno/chainsaw/pkg/testing"
	"k8s.io/client-go/rest"
	"k8s.io/utils/clock"
)

type Runner interface {
	Run(context.Context, model.Configuration, enginecontext.TestContext, ...discovery.Test) (model.SummaryResult, error)
}

func New(clock clock.PassiveClock, onFailure func()) Runner {
	return &runner{
		clock:     clock,
		onFailure: onFailure,
	}
}

type runner struct {
	clock     clock.PassiveClock
	onFailure func()
	deps      *internal.TestDeps
}

func (r *runner) Run(ctx context.Context, config model.Configuration, tc enginecontext.TestContext, tests ...discovery.Test) (model.SummaryResult, error) {
	return r.run(ctx, nil, config, tc, tests...)
}

func (r *runner) run(ctx context.Context, m mainstart, config model.Configuration, tc enginecontext.TestContext, tests ...discovery.Test) (model.SummaryResult, error) {
	// sanity check
	if len(tests) == 0 {
		return nil, nil
	}
	internalTests := []testing.InternalTest{{
		Name: "chainsaw",
		F: func(t *testing.T) {
			t.Helper()
			t.Parallel()
			// run tests
			r.runTests(ctx, t, config.Namespace, tc, tests...)
		},
	}}
	deps := r.deps
	if deps == nil {
		deps = &internal.TestDeps{}
	}
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
		return nil, fmt.Errorf("testing framework exited with non zero code %d", code)
	}
	tc.Report.EndTime = time.Now()
	// TODO: move to the caller
	if config.Report != nil && config.Report.Format != "" {
		if err := report.Save(tc.Report, config.Report.Format, config.Report.Path, config.Report.Name); err != nil {
			return tc, err
		}
	}
	return tc, nil
}

func (r *runner) onFail() {
	if r.onFailure != nil {
		r.onFailure()
	}
}

func InitContext(config model.Configuration, defaultCluster *rest.Config, values any) (enginecontext.TestContext, error) {
	tc := enginecontext.EmptyContext()
	// cleanup options
	tc = tc.WithSkipDelete(config.Cleanup.SkipDelete)
	if config.Cleanup.DelayBeforeCleanup != nil {
		tc = tc.WithDelayBeforeCleanup(&config.Cleanup.DelayBeforeCleanup.Duration)
	}
	// templating options
	tc = tc.WithTemplating(config.Templating.Enabled)
	if config.Templating.Compiler != nil {
		tc = tc.WithDefaultCompiler(string(*config.Templating.Compiler))
	}
	// discovery options
	tc = tc.WithFullName(config.Discovery.FullName)
	// execution options
	tc = tc.WithFailFast(config.Execution.FailFast)
	if config.Execution.ForceTerminationGracePeriod != nil {
		tc = tc.WithTerminationGrace(&config.Execution.ForceTerminationGracePeriod.Duration)
	}
	// deletion options
	tc = tc.WithDeletionPropagation(config.Deletion.Propagation)
	// error options
	tc = tc.WithCatch(config.Error.Catch...)
	// timeouts
	tc = tc.WithTimeouts(v1alpha1.Timeouts{
		Apply:   &config.Timeouts.Apply,
		Assert:  &config.Timeouts.Assert,
		Cleanup: &config.Timeouts.Cleanup,
		Delete:  &config.Timeouts.Delete,
		Error:   &config.Timeouts.Error,
		Exec:    &config.Timeouts.Exec,
	})
	// values
	tc = enginecontext.WithValues(tc, values)
	// clusters
	tc = enginecontext.WithClusters(tc, "", config.Clusters)
	if defaultCluster != nil {
		cluster, err := clusters.NewClusterFromConfig(defaultCluster)
		if err != nil {
			return tc, err
		}
		tc = tc.WithCluster(clusters.DefaultClient, cluster)
		return enginecontext.WithCurrentCluster(tc, clusters.DefaultClient)
	}
	return tc, nil
}
