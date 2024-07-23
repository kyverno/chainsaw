package processors

import (
	"context"
	"fmt"
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
	"github.com/kyverno/chainsaw/pkg/runner/timeout"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/chainsaw/pkg/utils/kube"
	"github.com/kyverno/pkg/ext/output/color"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/clock"
)

type TestProcessor interface {
	Run(context.Context, namespacer.Namespacer, model.TestContext, discovery.Test)
}

func NewTestProcessor(clock clock.PassiveClock, report *report.TestReport) TestProcessor {
	return &testProcessor{
		clock:  clock,
		report: report,
	}
}

type testProcessor struct {
	clock  clock.PassiveClock
	report *report.TestReport
}

func (p *testProcessor) Run(ctx context.Context, nspacer namespacer.Namespacer, tc model.TestContext, test discovery.Test) {
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
	for i, step := range test.Test.Spec.Steps {
		name := step.Name
		if name == "" {
			name = fmt.Sprintf("step-%d", i+1)
		}
		if size < len(name) {
			size = len(name)
		}
	}
	config := tc.Configuration()
	timeouts := config.Timeouts.Combine(test.Test.Spec.Timeouts)
	registeredClusters := clusters.Register(tc.Clusters(), test.BasePath, test.Test.Spec.Clusters)
	clusterConfig, clusterClient, err := registeredClusters.Resolve(false, test.Test.Spec.Cluster)
	if err != nil {
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		failer.FailNow(ctx)
	}
	bindings := apibindings.RegisterClusterBindings(ctx, tc.Bindings(), clusterConfig, clusterClient)
	setupLogger := logging.NewLogger(t, p.clock, test.Test.Name, fmt.Sprintf("%-*s", size, "@setup"))
	cleanupLogger := logging.NewLogger(t, p.clock, test.Test.Name, fmt.Sprintf("%-*s", size, "@cleanup"))
	var namespace *corev1.Namespace
	if clusterClient != nil {
		if nspacer == nil || test.Test.Spec.Namespace != "" {
			var ns corev1.Namespace
			if test.Test.Spec.Namespace != "" {
				ns = kube.Namespace(test.Test.Spec.Namespace)
			} else {
				ns = kube.PetNamespace()
			}
			namespace = &ns
		}
		if namespace != nil {
			object := kube.ToUnstructured(namespace)
			bindings = apibindings.RegisterNamedBinding(ctx, bindings, "namespace", object.GetName())
			if test.Test.Spec.NamespaceTemplate != nil && test.Test.Spec.NamespaceTemplate.Value != nil {
				template := v1alpha1.Any{
					Value: test.Test.Spec.NamespaceTemplate.Value,
				}
				if merged, err := mutate.Merge(ctx, object, bindings, template); err != nil {
					failer.FailNow(ctx)
				} else {
					object = merged
				}
				bindings = apibindings.RegisterNamedBinding(ctx, bindings, "namespace", object.GetName())
			} else if config.Namespace.Template != nil && config.Namespace.Template.Value != nil {
				template := v1alpha1.Any{
					Value: config.Namespace.Template.Value,
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
				if !cleanup.Skip(config.Cleanup.SkipDelete, test.Test.Spec.SkipDelete, nil) {
					t.Cleanup(func() {
						operation := newOperation(
							OperationInfo{},
							false,
							timeout.Get(nil, timeouts.CleanupDuration()),
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
	bindings, err = apibindings.RegisterBindings(ctx, bindings, test.Test.Spec.Bindings...)
	if err != nil {
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		failer.FailNow(ctx)
	}
	for i, step := range test.Test.Spec.Steps {
		processor := p.createStepProcessor(nspacer, tc, test, step)
		name := step.Name
		if name == "" {
			name = fmt.Sprintf("step-%d", i+1)
		}
		processor.Run(
			logging.IntoContext(ctx, logging.NewLogger(t, p.clock, test.Test.Name, fmt.Sprintf("%-*s", size, name))),
			apibindings.RegisterNamedBinding(ctx, bindings, "step", StepInfo{Id: i + 1}),
		)
	}
}

func (p *testProcessor) createStepProcessor(nspacer namespacer.Namespacer, tc model.TestContext, test discovery.Test, step v1alpha1.TestStep) StepProcessor {
	var report *report.StepReport
	if p.report != nil {
		report = p.report.ForStep(&step)
	}
	return NewStepProcessor(tc.Configuration(), tc.Clusters(), nspacer, p.clock, test, step, report)
}
