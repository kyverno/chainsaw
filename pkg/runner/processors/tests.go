package processors

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/names"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	opdelete "github.com/kyverno/chainsaw/pkg/runner/operations/delete"
	"github.com/kyverno/chainsaw/pkg/runner/summary"
	"github.com/kyverno/chainsaw/pkg/runner/timeout"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/kyverno/ext/output/color"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/utils/clock"
)

type TestsProcessor interface {
	Run(ctx context.Context)
	CreateTestProcessor(test discovery.Test) TestProcessor
}

func NewTestsProcessor(
	config v1alpha1.ConfigurationSpec,
	client client.Client,
	clock clock.PassiveClock,
	summary *summary.Summary,
	tests ...discovery.Test,
) TestsProcessor {
	return &testsProcessor{
		config:  config,
		client:  client,
		clock:   clock,
		summary: summary,
		tests:   tests,
	}
}

type testsProcessor struct {
	config  v1alpha1.ConfigurationSpec
	client  client.Client
	clock   clock.PassiveClock
	summary *summary.Summary
	tests   []discovery.Test
}

func (p *testsProcessor) Run(ctx context.Context) {
	t := testing.FromContext(ctx)
	var nspacer namespacer.Namespacer
	if p.config.Namespace != "" {
		namespace := client.Namespace(p.config.Namespace)
		if err := p.client.Get(ctx, client.ObjectKey(&namespace), namespace.DeepCopy()); err != nil {
			if !errors.IsNotFound(err) {
				// Get doesn't log
				logging.Log(ctx, logging.Get, color.BoldRed, err)
				t.FailNow()
			}
			t.Cleanup(func() {
				operation := operation{
					continueOnError: false,
					timeout:         timeout.DefaultCleanupTimeout,
					operation:       opdelete.New(p.client, namespace.DeepCopy()),
				}
				operation.execute(ctx)
			})
			if err := p.client.Create(ctx, namespace.DeepCopy()); err != nil {
				t.FailNow()
			}
		}
		nspacer = namespacer.New(p.client, p.config.Namespace)
	}
	for _, test := range p.tests {
		name, err := names.Test(p.config, test)
		if err != nil {
			logging.Log(ctx, logging.Internal, color.BoldRed, err)
			t.FailNow()
		}
		t.Run(name, func(t *testing.T) {
			t.Helper()
			processor := p.CreateTestProcessor(test)
			processor.Run(testing.IntoContext(ctx, t), nspacer)
		})
	}
}

func (p *testsProcessor) CreateTestProcessor(test discovery.Test) TestProcessor {
	return NewTestProcessor(p.config, p.client, p.clock, p.summary, test)
}
