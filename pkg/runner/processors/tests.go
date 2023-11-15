package processors

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/report"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/names"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/summary"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/kyverno/ext/output/color"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/utils/clock"
)

type TestsProcessor interface {
	Run(ctx context.Context, tests ...discovery.Test)
	CreateTestProcessor(test discovery.Test) TestProcessor
}

func NewTestsProcessor(config v1alpha1.ConfigurationSpec, client client.Client, clock clock.PassiveClock, summary *summary.Summary, report *report.Report) TestsProcessor {
	return &testsProcessor{
		config:  config,
		client:  client,
		clock:   clock,
		summary: summary,
		report:  report,
	}
}

type testsProcessor struct {
	config  v1alpha1.ConfigurationSpec
	client  client.Client
	clock   clock.PassiveClock
	summary *summary.Summary
	report  *report.Report
}

func (p *testsProcessor) Run(ctx context.Context, tests ...discovery.Test) {
	t := testing.FromContext(ctx)
	var nspacer namespacer.Namespacer
	if p.config.Namespace != "" {
		namespace := client.Namespace(p.config.Namespace)
		if err := p.client.Get(ctx, client.ObjectKey(&namespace), namespace.DeepCopy()); err != nil {
			if !errors.IsNotFound(err) {
				// Get doesn't log
				logging.Log(ctx, "GET   ", color.BoldRed, err)
				t.FailNow()
			}
			t.Cleanup(func() {
				// TODO: wait
				if err := p.client.Delete(ctx, &namespace); err != nil {
					t.FailNow()
				}
			})
			if err := p.client.Create(ctx, namespace.DeepCopy()); err != nil {
				t.FailNow()
			}
		}
		nspacer = namespacer.New(p.client, p.config.Namespace)
	}
	for _, test := range tests {
		name, err := names.Test(p.config, test)
		if err != nil {
			logging.Log(ctx, "INTERN", color.BoldRed, err)
			t.FailNow()
		}
		t.Run(name, func(t *testing.T) {
			t.Helper()
			processor := p.CreateTestProcessor(test)
			processor.Run(testing.IntoContext(ctx, t), nspacer, test)
		})
	}
}

func (p *testsProcessor) CreateTestProcessor(_ discovery.Test) TestProcessor {
	return NewTestProcessor(p.config, p.client, p.clock, p.summary, p.report)
}
