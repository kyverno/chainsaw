package processors

import (
	"context"
	"sync/atomic"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/report"
	"github.com/kyverno/chainsaw/pkg/runner/cleanup"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/names"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	opdelete "github.com/kyverno/chainsaw/pkg/runner/operations/delete"
	"github.com/kyverno/chainsaw/pkg/runner/summary"
	"github.com/kyverno/chainsaw/pkg/runner/timeout"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/kyverno/ext/output/color"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/clientcmd"
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
	testsReport *report.TestsReport,
	tests ...discovery.Test,
) TestsProcessor {
	return &testsProcessor{
		config:      config,
		client:      client,
		clock:       clock,
		summary:     summary,
		testsReport: testsReport,
		tests:       tests,
	}
}

type testsProcessor struct {
	config         v1alpha1.ConfigurationSpec
	client         client.Client
	clock          clock.PassiveClock
	summary        *summary.Summary
	testsReport    *report.TestsReport
	tests          []discovery.Test
	shouldFailFast atomic.Bool
}

func (p *testsProcessor) Run(ctx context.Context) {
	t := testing.FromContext(ctx)
	t.Cleanup(func() {
		if p.testsReport != nil {
			p.testsReport.Close()
		}
	})
	var nspacer namespacer.Namespacer
	if p.config.Namespace != "" {
		namespace := client.Namespace(p.config.Namespace)
		nspacer = namespacer.New(p.client, p.config.Namespace)
		if err := p.client.Get(ctx, client.ObjectKey(&namespace), namespace.DeepCopy()); err != nil {
			if !errors.IsNotFound(err) {
				// Get doesn't log
				logging.Log(ctx, logging.Get, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
				t.FailNow()
			}
			if !cleanup.Skip(p.config.SkipDelete, nil, nil) {
				t.Cleanup(func() {
					operation := operation{
						continueOnError: false,
						timeout:         timeout.Get(nil, p.config.Timeouts.CleanupDuration()),
						operation:       opdelete.New(p.client, client.ToUnstructured(namespace.DeepCopy()), nspacer),
					}
					operation.execute(ctx)
				})
			}
			if err := p.client.Create(ctx, namespace.DeepCopy()); err != nil {
				t.FailNow()
			}
		}
	}
	for _, test := range p.tests {
		if test.Spec.Kubeconfig != nil && test.Spec.Kubeconfig.File != "" {
			config, err := clientcmd.BuildConfigFromFlags("", test.Spec.Kubeconfig.File)
			if err != nil {
				logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
				t.FailNow()
			}
			kClient, err := client.New(config)
			if err != nil {
				logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
				t.FailNow()
			}
			p.client = kClient
		}
		name, err := names.Test(p.config, test)
		if err != nil {
			logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
			t.FailNow()
		}
		t.Run(name, func(t *testing.T) {
			t.Helper()
			t.Cleanup(func() {
				if t.Failed() {
					p.shouldFailFast.Store(true)
				}
			})
			processor := p.CreateTestProcessor(test)
			processor.Run(testing.IntoContext(ctx, t), nspacer)
		})
	}
}

func (p *testsProcessor) CreateTestProcessor(test discovery.Test) TestProcessor {
	testReport := report.NewTest(test.Name)
	if p.testsReport != nil {
		p.testsReport.AddTest(testReport)
	}
	return NewTestProcessor(p.config, p.client, p.clock, p.summary, testReport, test, &p.shouldFailFast)
}
