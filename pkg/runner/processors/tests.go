package processors

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/kyverno/chainsaw/pkg/report"
	apibindings "github.com/kyverno/chainsaw/pkg/runner/bindings"
	"github.com/kyverno/chainsaw/pkg/runner/cleanup"
	"github.com/kyverno/chainsaw/pkg/runner/failer"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/names"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	opdelete "github.com/kyverno/chainsaw/pkg/runner/operations/delete"
	"github.com/kyverno/chainsaw/pkg/runner/summary"
	"github.com/kyverno/chainsaw/pkg/runner/timeout"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/chainsaw/pkg/utils/kube"
	"github.com/kyverno/pkg/ext/output/color"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/clock"
)

type TestsProcessor interface {
	Run(context.Context, model.GlobalContext, ...discovery.Test)
}

func NewTestsProcessor(
	clock clock.PassiveClock,
	summary *summary.Summary,
	report *report.Report,
) TestsProcessor {
	return &testsProcessor{
		clock:   clock,
		summary: summary,
		report:  report,
	}
}

type testsProcessor struct {
	clock   clock.PassiveClock
	summary *summary.Summary
	report  *report.Report
	// state
	shouldFailFast atomic.Bool
}

func (p *testsProcessor) Run(ctx context.Context, tc model.GlobalContext, tests ...discovery.Test) {
	t := testing.FromContext(ctx)
	if p.report != nil {
		p.report.SetStartTime(time.Now())
		t.Cleanup(func() {
			p.report.SetEndTime(time.Now())
		})
	}
	config := tc.Configuration()
	bindings := tc.Bindings()
	var nspacer namespacer.Namespacer
	namespace, err := tc.Namespace(ctx)
	if err != nil {
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		failer.FailNow(ctx)
	}
	if namespace != nil {
		clusterConfig, clusterClient, err := tc.Clusters().Resolve(false)
		if err != nil {
			logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
			failer.FailNow(ctx)
		}
		nspacer = namespacer.New(clusterClient, namespace.GetName())
		if err := clusterClient.Get(ctx, client.ObjectKey(namespace), namespace.DeepCopy()); err != nil {
			if !errors.IsNotFound(err) {
				// Get doesn't log
				logging.Log(ctx, logging.Get, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
				failer.FailNow(ctx)
			}
			if !cleanup.Skip(config.Cleanup.SkipDelete, nil, nil) {
				t.Cleanup(func() {
					operation := newOperation(
						model.OperationInfo{},
						false,
						timeout.Get(nil, config.Timeouts.CleanupDuration()),
						func(ctx context.Context, bindings binding.Bindings) (operations.Operation, binding.Bindings, error) {
							bindings = apibindings.RegisterClusterBindings(ctx, bindings, clusterConfig, clusterClient)
							return opdelete.New(clusterClient, kube.ToUnstructured(namespace), nspacer, false, metav1.DeletePropagationBackground), bindings, nil
						},
						nil,
					)
					operation.execute(ctx, bindings)
				})
			}
			if err := clusterClient.Create(ctx, namespace.DeepCopy()); err != nil {
				failer.FailNow(ctx)
			}
		}
	}
	for i := range tests {
		test := tests[i]
		name, err := names.Test(config, test)
		if err != nil {
			logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
			failer.FailNow(ctx)
		}
		var scenarios []discovery.Test
		if test.Test != nil {
			if len(test.Test.Spec.Scenarios) == 0 {
				scenarios = append(scenarios, test)
			} else {
				for s := range test.Test.Spec.Scenarios {
					scenario := test.Test.Spec.Scenarios[s]
					test := test
					test.Test = test.Test.DeepCopy()
					test.Test.Spec.Scenarios = nil
					bindings := scenario.Bindings
					bindings = append(bindings, test.Test.Spec.Bindings...)
					test.Test.Spec.Bindings = bindings
					scenarios = append(scenarios, test)
				}
			}
		}
		for s := range scenarios {
			test := scenarios[s]
			t.Run(name, func(t *testing.T) {
				t.Helper()
				t.Cleanup(func() {
					if t.Failed() {
						p.shouldFailFast.Store(true)
					}
				})
				processor := p.createTestProcessor(test)
				processor.Run(
					testing.IntoContext(ctx, t),
					tc.TestContext(ctx, test.Test, i, s),
					nspacer,
					test,
				)
			})
		}
	}
}

func (p *testsProcessor) createTestProcessor(test discovery.Test) TestProcessor {
	var report *report.TestReport
	if p.report != nil {
		report = p.report.ForTest(&test)
	}
	return NewTestProcessor(p.clock, p.summary, report, &p.shouldFailFast)
}
