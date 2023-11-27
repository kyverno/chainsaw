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

func Run(cfg *rest.Config, clock clock.PassiveClock, config v1alpha1.ConfigurationSpec, tests ...discovery.Test) (*summary.Summary, error) {
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
	client, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	client = runnerclient.New(client)
	internalTests := []testing.InternalTest{{
		Name: "chainsaw",
		F: func(t *testing.T) {
			t.Parallel()
			processor := processors.NewTestsProcessor(config, client, clock, &summary, testsReport, tests...)
			ctx := testing.IntoContext(context.Background(), t)
			ctx = logging.IntoContext(ctx, logging.NewLogger(t, clock, t.Name(), "@main"))
			processor.Run(ctx)
		},
	}}
	deps := &internal.TestDeps{}
	m := testing.MainStart(deps, internalTests, nil, nil, nil)
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
