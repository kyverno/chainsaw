package processors

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/kyverno/chainsaw/pkg/report"
	apibindings "github.com/kyverno/chainsaw/pkg/runner/bindings"
	"github.com/kyverno/chainsaw/pkg/runner/cleanup"
	"github.com/kyverno/chainsaw/pkg/runner/clusters"
	"github.com/kyverno/chainsaw/pkg/runner/failer"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/mutate"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	opdelete "github.com/kyverno/chainsaw/pkg/runner/operations/delete"
	"github.com/kyverno/chainsaw/pkg/runner/summary"
	"github.com/kyverno/chainsaw/pkg/runner/timeout"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/pkg/ext/output/color"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/clock"
)

type TestProcessor interface {
	Run(context.Context, binding.Bindings, namespacer.Namespacer)
	CreateStepProcessor(namespacer.Namespacer, clusters.Registry, v1alpha1.TestStep) StepProcessor
}

func NewTestProcessor(
	config model.Configuration,
	clusters clusters.Registry,
	clock clock.PassiveClock,
	summary *summary.Summary,
	report *report.TestReport,
	test discovery.Test,
	shouldFailFast *atomic.Bool,
) TestProcessor {
	return &testProcessor{
		config:         config,
		clusters:       clusters,
		clock:          clock,
		summary:        summary,
		report:         report,
		test:           test,
		shouldFailFast: shouldFailFast,
		timeouts:       config.Timeouts.Combine(test.Test.Spec.Timeouts),
	}
}

type testProcessor struct {
	config         model.Configuration
	clusters       clusters.Registry
	clock          clock.PassiveClock
	summary        *summary.Summary
	report         *report.TestReport
	test           discovery.Test
	shouldFailFast *atomic.Bool
	timeouts       v1alpha1.Timeouts
}

func (p *testProcessor) Run(ctx context.Context, bindings binding.Bindings, nspacer namespacer.Namespacer) {
	if bindings == nil {
		bindings = binding.NewBindings()
	}
	t := testing.FromContext(ctx)
	if p.report != nil {
		p.report.SetStartTime(time.Now())
		t.Cleanup(func() {
			p.report.SetEndTime(time.Now())
			if t.Failed() {
				p.report.Fail()
			}
			if t.Skipped() {
				p.report.Skip()
			}
		})
	}
	size := len("@cleanup")
	for i, step := range p.test.Test.Spec.Steps {
		name := step.Name
		if name == "" {
			name = fmt.Sprintf("step-%d", i+1)
		}
		if size < len(name) {
			size = len(name)
		}
	}
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
	if p.test.Test.Spec.Concurrent == nil || *p.test.Test.Spec.Concurrent {
		t.Parallel()
	}
	if p.test.Test.Spec.Skip != nil && *p.test.Test.Spec.Skip {
		t.SkipNow()
	}
	if p.config.Execution.FailFast {
		if p.shouldFailFast.Load() {
			t.SkipNow()
		}
	}
	registeredClusters := clusters.Register(p.clusters, p.test.BasePath, p.test.Test.Spec.Clusters)
	clusterConfig, clusterClient, err := registeredClusters.Resolve(false, p.test.Test.Spec.Cluster)
	if err != nil {
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		failer.FailNow(ctx)
	}
	bindings = apibindings.RegisterClusterBindings(ctx, bindings, clusterConfig, clusterClient)
	setupLogger := logging.NewLogger(t, p.clock, p.test.Test.Name, fmt.Sprintf("%-*s", size, "@setup"))
	cleanupLogger := logging.NewLogger(t, p.clock, p.test.Test.Name, fmt.Sprintf("%-*s", size, "@cleanup"))
	var namespace *corev1.Namespace
	if clusterClient != nil {
		if nspacer == nil || p.test.Test.Spec.Namespace != "" {
			var ns corev1.Namespace
			if p.test.Test.Spec.Namespace != "" {
				ns = client.Namespace(p.test.Test.Spec.Namespace)
			} else {
				ns = client.PetNamespace()
			}
			namespace = &ns
		}
		if namespace != nil {
			object := client.ToUnstructured(namespace)
			bindings = apibindings.RegisterNamedBinding(ctx, bindings, "namespace", object.GetName())
			if p.test.Test.Spec.NamespaceTemplate != nil && p.test.Test.Spec.NamespaceTemplate.Value != nil {
				template := v1alpha1.Any{
					Value: p.test.Test.Spec.NamespaceTemplate.Value,
				}
				if merged, err := mutate.Merge(ctx, object, bindings, template); err != nil {
					failer.FailNow(ctx)
				} else {
					object = merged
				}
				bindings = apibindings.RegisterNamedBinding(ctx, bindings, "namespace", object.GetName())
			} else if p.config.Namespace.Template != nil && p.config.Namespace.Template.Value != nil {
				template := v1alpha1.Any{
					Value: p.config.Namespace.Template.Value,
				}
				if merged, err := mutate.Merge(ctx, object, bindings, template); err != nil {
					failer.FailNow(ctx)
				} else {
					object = merged
				}
				bindings = apibindings.RegisterNamedBinding(ctx, bindings, "namespace", object.GetName())
			}
			nspacer = namespacer.New(clusterClient, object.GetName())
			setupCtx := logging.IntoContext(ctx, setupLogger)
			cleanupCtx := logging.IntoContext(ctx, cleanupLogger)
			if err := clusterClient.Get(setupCtx, client.ObjectKey(&object), object.DeepCopy()); err != nil {
				if !errors.IsNotFound(err) {
					// Get doesn't log
					setupLogger.Log(logging.Get, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
					failer.FailNow(ctx)
				}
				if !cleanup.Skip(p.config.Cleanup.SkipDelete, p.test.Test.Spec.SkipDelete, nil) {
					t.Cleanup(func() {
						operation := newOperation(
							OperationInfo{},
							false,
							timeout.Get(nil, p.timeouts.CleanupDuration()),
							func(ctx context.Context, bindings binding.Bindings) (operations.Operation, binding.Bindings, error) {
								bindings = apibindings.RegisterClusterBindings(ctx, bindings, clusterConfig, clusterClient)
								return opdelete.New(clusterClient, object, nspacer, false, metav1.DeletePropagationBackground), bindings, nil
							},
							nil,
						)
						operation.execute(cleanupCtx, bindings)
					})
				}
				if err := clusterClient.Create(logging.IntoContext(setupCtx, setupLogger), object.DeepCopy()); err != nil {
					failer.FailNow(ctx)
				}
			}
		}
	}
	if p.report != nil && nspacer != nil {
		p.report.SetNamespace(nspacer.GetNamespace())
	}
	bindings, err = apibindings.RegisterBindings(ctx, bindings, p.test.Test.Spec.Bindings...)
	if err != nil {
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		failer.FailNow(ctx)
	}
	for i, step := range p.test.Test.Spec.Steps {
		processor := p.CreateStepProcessor(nspacer, registeredClusters, step)
		name := step.Name
		if name == "" {
			name = fmt.Sprintf("step-%d", i+1)
		}
		processor.Run(
			logging.IntoContext(ctx, logging.NewLogger(t, p.clock, p.test.Test.Name, fmt.Sprintf("%-*s", size, name))),
			apibindings.RegisterNamedBinding(ctx, bindings, "step", StepInfo{Id: i + 1}),
		)
	}
}

func (p *testProcessor) CreateStepProcessor(nspacer namespacer.Namespacer, clusters clusters.Registry, step v1alpha1.TestStep) StepProcessor {
	var report *report.StepReport
	if p.report != nil {
		report = p.report.ForStep(&step)
	}
	return NewStepProcessor(p.config, clusters, nspacer, p.clock, p.test, step, report)
}
