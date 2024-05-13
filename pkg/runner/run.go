package runner

import (
	"context"
	"fmt"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/report"
	apibindings "github.com/kyverno/chainsaw/pkg/runner/bindings"
	"github.com/kyverno/chainsaw/pkg/runner/clusters"
	"github.com/kyverno/chainsaw/pkg/runner/internal"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/processors"
	"github.com/kyverno/chainsaw/pkg/runner/summary"
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
	config v1alpha1.ConfigurationSpec,
	values map[string]any,
	tests ...discovery.Test,
) (*summary.Summary, error) {
	return run(ctx, cfg, clock, config, nil, values, tests...)
}

func run(
	ctx context.Context,
	cfg *rest.Config,
	clock clock.PassiveClock,
	config v1alpha1.ConfigurationSpec,
	m mainstart,
	values map[string]any,
	tests ...discovery.Test,
) (*summary.Summary, error) {
	var summary summary.Summary
	var testsReport *report.Report
	if config.ReportFormat != "" {
		testsReport = report.New(config.ReportName)
	}
	if len(tests) == 0 {
		return &summary, nil
	}
	if err := internal.SetupFlags(config); err != nil {
		return nil, err
	}
	bindings := binding.NewBindings()
	bindings = apibindings.RegisterNamedBinding(ctx, bindings, "values", values)
	registeredClusters := clusters.NewRegistry()
	if cfg != nil {
		cluster, err := clusters.NewClusterFromConfig(cfg)
		if err != nil {
			return nil, err
		}
		registeredClusters = registeredClusters.Register(clusters.DefaultClient, cluster)
	}
	registeredClusters = clusters.Register(registeredClusters, "", config.Clusters)
	internalTests := []testing.InternalTest{{
		Name: "chainsaw",
		F: func(t *testing.T) {
			t.Helper()
			t.Parallel()
			processor := processors.NewTestsProcessor(config, registeredClusters, clock, &summary, testsReport, tests...)
			ctx := testing.IntoContext(ctx, t)
			ctx = logging.IntoContext(ctx, logging.NewLogger(t, clock, t.Name(), "@main"))
			processor.Run(ctx, bindings)
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
		return &summary, fmt.Errorf("testing framework exited with non zero code %d", code)
	}
	if testsReport != nil && config.ReportFormat != "" {
		if err := testsReport.Save(config.ReportFormat, config.ReportPath, config.ReportName); err != nil {
			return &summary, err
		}
	}
	return &summary, nil
}
