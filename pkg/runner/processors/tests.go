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
	"github.com/kyverno/chainsaw/pkg/runner/failer"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/names"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	opdelete "github.com/kyverno/chainsaw/pkg/runner/operations/delete"
	"github.com/kyverno/chainsaw/pkg/runner/summary"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/chainsaw/pkg/utils/kube"
	"github.com/kyverno/pkg/ext/output/color"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/clock"
	"k8s.io/utils/ptr"
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
}

func (p *testsProcessor) Run(ctx context.Context, gc model.GlobalContext, tests ...discovery.Test) {
	t := testing.FromContext(ctx)
	if p.report != nil {
		p.report.SetStartTime(time.Now())
		t.Cleanup(func() {
			p.report.SetEndTime(time.Now())
		})
	}
	var nspacer namespacer.Namespacer
	namespace, err := gc.Namespace(ctx)
	if err != nil {
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		failer.FailNow(ctx)
	}
	if namespace != nil {
		_, clusterClient, err := gc.Clusters().Resolve(false)
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
			if gc.Cleanup() {
				t.Cleanup(func() {
					operation := newOperation(
						model.OperationInfo{},
						false,
						ptr.To(gc.Timeouts().CleanupDuration()),
						func(ctx context.Context, bindings binding.Bindings) (operations.Operation, binding.Bindings, error) {
							return opdelete.New(clusterClient, kube.ToUnstructured(namespace), nspacer, false, metav1.DeletePropagationBackground), bindings, nil
						},
						nil,
					)
					operation.execute(ctx, nil)
				})
			}
			if err := clusterClient.Create(ctx, namespace.DeepCopy()); err != nil {
				failer.FailNow(ctx)
			}
		}
	}
	hasFailures := atomic.Bool{}
	for i := range tests {
		test := tests[i]
		name, err := names.Test(gc.FullName(), test)
		if err != nil {
			logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
			failer.FailNow(ctx)
		}
		scenarios := applyScenarios(test)
		for s := range scenarios {
			test := scenarios[s]
			t.Run(name, func(t *testing.T) {
				t.Helper()
				ctx := testing.IntoContext(ctx, t)
				tc := gc.TestContext(ctx, test.Test, i, s)
				p.registerUpdateSummary(t)
				p.registerFailuresTracker(t, &hasFailures)
				if test.Test.Spec.Concurrent == nil || *test.Test.Spec.Concurrent {
					t.Parallel()
				}
				if test.Test.Spec.Skip != nil && *test.Test.Spec.Skip {
					t.SkipNow()
				}
				if gc.FailFast() && hasFailures.Load() {
					t.SkipNow()
				}
				processor := p.createTestProcessor(test)
				processor.Run(ctx, tc, nspacer, test)
			})
		}
	}
}

func (p *testsProcessor) createTestProcessor(test discovery.Test) TestProcessor {
	var report *report.TestReport
	if p.report != nil {
		report = p.report.ForTest(&test)
	}
	return NewTestProcessor(p.clock, report)
}

func (p *testsProcessor) registerUpdateSummary(t *testing.T) {
	t.Helper()
	if p.summary != nil {
		t.Cleanup(func() {
			if t.Skipped() {
				p.summary.IncSkipped()
			} else {
				if t.Failed() {
					p.summary.IncFailed()
				} else {
					p.summary.IncPassed()
				}
			}
		})
	}
}

func (p *testsProcessor) registerFailuresTracker(t *testing.T, flag *atomic.Bool) {
	t.Helper()
	t.Cleanup(func() {
		if t.Failed() {
			flag.Store(true)
		}
	})
}

func applyScenarios(test discovery.Test) []discovery.Test {
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
	return scenarios
}
