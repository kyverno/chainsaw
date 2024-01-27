package runner

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/report"
	runnerclient "github.com/kyverno/chainsaw/pkg/runner/client"
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

func Run(cfg *rest.Config, clock clock.PassiveClock, config v1alpha1.ConfigurationSpec, tests ...discovery.Test) (*summary.Summary, error) {
	return run(cfg, clock, config, nil, tests...)
}

func run(cfg *rest.Config, clock clock.PassiveClock, config v1alpha1.ConfigurationSpec, m mainstart, tests ...discovery.Test) (*summary.Summary, error) {
	var summary summary.Summary
	var testsReport *report.TestsReport
	if config.ReportFormat != "" {
		testsReport = report.NewTests(config.ReportName)
	}
	if len(tests) == 0 {
		return &summary, nil
	}
	if err := internal.SetupFlags(config); err != nil {
		return nil, err
	}
	var clusterClient client.Client
	if cfg != nil {
		client, err := client.New(cfg)
		if err != nil {
			return nil, err
		}
		clusterClient = runnerclient.New(client)
	}
	internalTests := []testing.InternalTest{{
		Name: "chainsaw",
		F: func(t *testing.T) {
			t.Helper()
			t.Parallel()
			processor := processors.NewTestsProcessor(config, clusterClient, clock, &summary, testsReport, tests...)
			ctx := testing.IntoContext(context.Background(), t)
			ctx = logging.IntoContext(ctx, logging.NewLogger(t, clock, t.Name(), "@main"))
			processor.Run(ctx)
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
		if err := testsReport.SaveReportBasedOnType(config.ReportFormat, config.ReportName); err != nil {
			return &summary, fmt.Errorf("failed to save test report: %v", err)
		}
	}
	return &summary, nil
}
