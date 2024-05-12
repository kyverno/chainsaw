package processors

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/report"
	apibindings "github.com/kyverno/chainsaw/pkg/runner/bindings"
	"github.com/kyverno/chainsaw/pkg/runner/cleanup"
	"github.com/kyverno/chainsaw/pkg/runner/clusters"
	"github.com/kyverno/chainsaw/pkg/runner/failer"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/mutate"
	"github.com/kyverno/chainsaw/pkg/runner/names"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	opdelete "github.com/kyverno/chainsaw/pkg/runner/operations/delete"
	"github.com/kyverno/chainsaw/pkg/runner/summary"
	"github.com/kyverno/chainsaw/pkg/runner/timeout"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/pkg/ext/output/color"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/utils/clock"
)

type TestsProcessor interface {
	Run(context.Context, binding.Bindings)
	CreateTestProcessor(discovery.Test) TestProcessor
}

func NewTestsProcessor(
	config v1alpha1.ConfigurationSpec,
	clusters clusters.Registry,
	clock clock.PassiveClock,
	summary *summary.Summary,
	report *report.Report,
	tests ...discovery.Test,
) TestsProcessor {
	return &testsProcessor{
		config:   config,
		clusters: clusters,
		clock:    clock,
		summary:  summary,
		report:   report,
		tests:    tests,
	}
}

type testsProcessor struct {
	config   v1alpha1.ConfigurationSpec
	clusters clusters.Registry
	clock    clock.PassiveClock
	summary  *summary.Summary
	report   *report.Report
	tests    []discovery.Test
	// state
	shouldFailFast atomic.Bool
}

func (p *testsProcessor) Run(ctx context.Context, bindings binding.Bindings) {
	if bindings == nil {
		bindings = binding.NewBindings()
	}
	t := testing.FromContext(ctx)
	if p.report != nil {
		p.report.SetStartTime(time.Now())
		t.Cleanup(func() {
			p.report.SetEndTime(time.Now())
		})
	}
	var nspacer namespacer.Namespacer
	cluster := p.clusters.Resolve()
	bindings = apibindings.RegisterClusterBindings(ctx, bindings, cluster.Config(), cluster.Client())
	if cluster != nil {
		if p.config.Namespace != "" {
			namespace := client.Namespace(p.config.Namespace)
			object := client.ToUnstructured(&namespace)
			bindings = apibindings.RegisterNamedBinding(ctx, bindings, "namespace", object.GetName())
			if p.config.NamespaceTemplate != nil && p.config.NamespaceTemplate.Value != nil {
				template := v1alpha1.Any{
					Value: p.config.NamespaceTemplate.Value,
				}
				if merged, err := mutate.Merge(ctx, object, bindings, template); err != nil {
					failer.FailNow(ctx)
				} else {
					object = merged
				}
				bindings = apibindings.RegisterNamedBinding(ctx, bindings, "namespace", object.GetName())
			}
			nspacer = namespacer.New(cluster.Client(), object.GetName())
			if err := cluster.Client().Get(ctx, client.ObjectKey(&object), object.DeepCopy()); err != nil {
				if !errors.IsNotFound(err) {
					// Get doesn't log
					logging.Log(ctx, logging.Get, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
					failer.FailNow(ctx)
				}
				if !cleanup.Skip(p.config.SkipDelete, nil, nil) {
					t.Cleanup(func() {
						operation := newLazyOperation(
							cluster,
							OperationInfo{},
							false,
							timeout.Get(nil, p.config.Timeouts.CleanupDuration()),
							func(_ context.Context, _ binding.Bindings) (operations.Operation, error) {
								return opdelete.New(cluster.Client(), object, nspacer, false), nil
							},
							nil,
						)
						operation.execute(ctx, bindings)
					})
				}
				if err := cluster.Client().Create(ctx, object.DeepCopy()); err != nil {
					failer.FailNow(ctx)
				}
			}
		}
	}
	bindings, err := apibindings.RegisterBindings(ctx, bindings)
	if err != nil {
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		failer.FailNow(ctx)
	}
	for i, test := range p.tests {
		name, err := names.Test(p.config, test)
		if err != nil {
			logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
			failer.FailNow(ctx)
		}
		t.Run(name, func(t *testing.T) {
			t.Helper()
			t.Cleanup(func() {
				if t.Failed() {
					p.shouldFailFast.Store(true)
				}
			})
			processor := p.CreateTestProcessor(test)
			processor.Run(
				testing.IntoContext(ctx, t),
				apibindings.RegisterNamedBinding(ctx, bindings, "test", TestInfo{Id: i + 1}),
				nspacer,
			)
		})
	}
}

func (p *testsProcessor) CreateTestProcessor(test discovery.Test) TestProcessor {
	var report *report.TestReport
	if p.report != nil {
		report = p.report.ForTest(&test)
	}
	return NewTestProcessor(p.config, p.clusters, p.clock, p.summary, report, test, &p.shouldFailFast)
}
