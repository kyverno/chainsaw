package processors

import (
	"context"
	"sync/atomic"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/report"
	"github.com/kyverno/chainsaw/pkg/runner/cleanup"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/mutate"
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
	Run(context.Context)
	CreateTestProcessor(discovery.Test, binding.Bindings) TestProcessor
}

func NewTestsProcessor(
	config v1alpha1.ConfigurationSpec,
	clusters clusters,
	clock clock.PassiveClock,
	summary *summary.Summary,
	testsReport *report.TestsReport,
	bindings binding.Bindings,
	tests ...discovery.Test,
) TestsProcessor {
	return &testsProcessor{
		config:      config,
		clusters:    clusters,
		clock:       clock,
		summary:     summary,
		testsReport: testsReport,
		bindings:    bindings,
		tests:       tests,
	}
}

type testsProcessor struct {
	config      v1alpha1.ConfigurationSpec
	clusters    clusters
	clock       clock.PassiveClock
	summary     *summary.Summary
	testsReport *report.TestsReport
	bindings    binding.Bindings
	tests       []discovery.Test
	// state
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
	bindings := p.bindings
	cluster := p.clusters.client(DefaultClient)
	if cluster != nil {
		if p.config.Namespace != "" {
			namespace := client.Namespace(p.config.Namespace)
			object := client.ToUnstructured(&namespace)
			bindings = p.bindings.Register("$namespace", binding.NewBinding(object.GetName()))
			if p.config.NamespaceTemplate != nil && p.config.NamespaceTemplate.Value != nil {
				template := v1alpha1.Any{
					Value: p.config.NamespaceTemplate.Value,
				}
				if merged, err := mutate.Merge(ctx, object, bindings, template); err != nil {
					t.FailNow()
				} else {
					object = merged
				}
				bindings = p.bindings.Register("$namespace", binding.NewBinding(object.GetName()))
			}
			nspacer = namespacer.New(cluster, object.GetName())
			if err := cluster.Get(ctx, client.ObjectKey(&object), object.DeepCopy()); err != nil {
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
							operation:       opdelete.New(cluster, object, nspacer, bindings, false),
						}
						operation.execute(ctx)
					})
				}
				if err := cluster.Create(ctx, object.DeepCopy()); err != nil {
					t.FailNow()
				}
			}
		}
	}
	for _, test := range p.tests {
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
			processor := p.CreateTestProcessor(test, bindings)
			processor.Run(testing.IntoContext(ctx, t), nspacer)
		})
	}
}

func (p *testsProcessor) CreateTestProcessor(test discovery.Test, bindings binding.Bindings) TestProcessor {
	testReport := report.NewTest(test.Name)
	if p.testsReport != nil {
		p.testsReport.AddTest(testReport)
	}
	return NewTestProcessor(p.config, p.clusters, p.clock, p.summary, testReport, test, &p.shouldFailFast, bindings)
}
