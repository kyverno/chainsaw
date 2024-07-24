package processors

import (
	"context"
	"time"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/cleanup/cleaner"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/kyverno/chainsaw/pkg/report"
	apibindings "github.com/kyverno/chainsaw/pkg/runner/bindings"
	"github.com/kyverno/chainsaw/pkg/runner/failer"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/mutate"
	"github.com/kyverno/chainsaw/pkg/runner/names"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/chainsaw/pkg/utils/kube"
	"github.com/kyverno/pkg/ext/output/color"
	"github.com/kyverno/pkg/ext/resource/convert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/utils/clock"
)

type TestsProcessor interface {
	Run(context.Context, model.TestContext, ...discovery.Test)
}

func NewTestsProcessor(config model.Configuration, clock clock.PassiveClock, report *report.Report) TestsProcessor {
	return &testsProcessor{
		config: config,
		clock:  clock,
		report: report,
	}
}

type testsProcessor struct {
	config model.Configuration
	clock  clock.PassiveClock
	report *report.Report
}

func (p *testsProcessor) Run(ctx context.Context, tc model.TestContext, tests ...discovery.Test) {
	t := testing.FromContext(ctx)
	tc, nspacer := p.setup(ctx, tc)
	for i := range tests {
		test := tests[i]
		name, err := names.Test(p.config.Discovery.FullName, test)
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
				ctx := testing.IntoContext(ctx, t)
				info := TestInfo{
					Id:         i + 1,
					ScenarioId: s + 1,
					Metadata:   test.Test.ObjectMeta,
				}
				tc := tc.WithBinding(ctx, "test", info)
				t.Cleanup(func() {
					if t.Skipped() {
						tc.IncSkipped()
					} else {
						if t.Failed() {
							tc.IncFailed()
						} else {
							tc.IncPassed()
						}
					}
				})
				if test.Test.Spec.Concurrent == nil || *test.Test.Spec.Concurrent {
					t.Parallel()
				}
				if test.Test.Spec.Skip != nil && *test.Test.Spec.Skip {
					t.SkipNow()
				}
				if p.config.Execution.FailFast {
					if tc.Failed() > 0 {
						t.SkipNow()
					}
				}
				processor := p.createTestProcessor(test)
				processor.Run(ctx, nspacer, tc, test)
			})
		}
	}
}

func buildNamespace(ctx context.Context, name string, template *v1alpha1.Any, bindings binding.Bindings) (*corev1.Namespace, error) {
	namespace := kube.Namespace(name)
	if template == nil {
		return &namespace, nil
	}
	if template.Value == nil {
		return &namespace, nil
	}
	object := kube.ToUnstructured(&namespace)
	bindings = apibindings.RegisterNamedBinding(ctx, bindings, "namespace", object.GetName())
	merged, err := mutate.Merge(ctx, object, bindings, *template)
	if err != nil {
		return nil, err
	}
	return convert.To[corev1.Namespace](merged)
}

func (p *testsProcessor) setup(ctx context.Context, tc model.TestContext) (model.TestContext, namespacer.Namespacer) {
	t := testing.FromContext(ctx)
	if p.report != nil {
		p.report.SetStartTime(time.Now())
		t.Cleanup(func() {
			p.report.SetEndTime(time.Now())
		})
	}
	cleaner := cleaner.New(tc.Timeouts().Cleanup, nil)
	t.Cleanup(func() {
		if !cleaner.Empty() {
			logging.Log(ctx, logging.Cleanup, logging.RunStatus, color.BoldFgCyan)
			defer func() {
				logging.Log(ctx, logging.Cleanup, logging.DoneStatus, color.BoldFgCyan)
			}()
			for _, err := range cleaner.Run(ctx) {
				logging.Log(ctx, logging.Cleanup, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
				failer.Fail(ctx)
			}
		}
	})
	var nspacer namespacer.Namespacer
	if p.config.Namespace.Name != "" {
		namespace, err := buildNamespace(ctx, p.config.Namespace.Name, p.config.Namespace.Template, tc.Bindings())
		if err != nil {
			logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
			failer.FailNow(ctx)
		}
		_, clusterClient, err := tc.CurrentClusterClient()
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
			} else if err := clusterClient.Create(ctx, namespace.DeepCopy()); err != nil {
				failer.FailNow(ctx)
			} else if tc.Cleanup() {
				cleaner.Add(clusterClient, namespace)
			}
		}
	}
	return tc, nspacer
}

func (p *testsProcessor) createTestProcessor(test discovery.Test) TestProcessor {
	var report *report.TestReport
	if p.report != nil {
		report = p.report.ForTest(&test)
	}
	return NewTestProcessor(p.config, p.clock, report)
}
