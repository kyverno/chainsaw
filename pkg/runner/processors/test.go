package processors

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	opdelete "github.com/kyverno/chainsaw/pkg/runner/operations/delete"
	"github.com/kyverno/chainsaw/pkg/runner/summary"
	"github.com/kyverno/chainsaw/pkg/runner/timeout"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/kyverno/ext/output/color"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/utils/clock"
)

type TestProcessor interface {
	Run(ctx context.Context, nspacer namespacer.Namespacer)
	CreateStepProcessor(nspacer namespacer.Namespacer, step v1alpha1.TestSpecStep) StepProcessor
}

func NewTestProcessor(
	config v1alpha1.ConfigurationSpec,
	client client.Client,
	clock clock.PassiveClock,
	summary *summary.Summary,
	test discovery.Test,
) TestProcessor {
	return &testProcessor{
		config:  config,
		client:  client,
		clock:   clock,
		summary: summary,
		test:    test,
	}
}

type testProcessor struct {
	config  v1alpha1.ConfigurationSpec
	client  client.Client
	clock   clock.PassiveClock
	summary *summary.Summary
	test    discovery.Test
}

func (p *testProcessor) Run(ctx context.Context, nspacer namespacer.Namespacer) {
	t := testing.FromContext(ctx)
	size := 0
	for i, step := range p.test.Spec.Steps {
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
	if p.test.Spec.Concurrent == nil || *p.test.Spec.Concurrent {
		t.Parallel()
	}
	if p.test.Spec.Skip != nil && *p.test.Spec.Skip {
		t.SkipNow()
	}
	setupLogger := logging.NewLogger(t, p.clock, p.test.Name, fmt.Sprintf("%-*s", size, "@setup"))
	var namespace *corev1.Namespace
	if nspacer == nil || p.test.Spec.Namespace != "" {
		var ns corev1.Namespace
		if p.test.Spec.Namespace != "" {
			ns = client.Namespace(p.test.Spec.Namespace)
		} else {
			ns = client.PetNamespace()
		}
		namespace = &ns
	}
	if namespace != nil {
		nspacer = namespacer.New(p.client, namespace.Name)
		ctx := logging.IntoContext(ctx, setupLogger)
		if err := p.client.Get(ctx, client.ObjectKey(namespace), namespace.DeepCopy()); err != nil {
			if !errors.IsNotFound(err) {
				// Get doesn't log
				setupLogger.Log(logging.Get, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
				t.FailNow()
			}
			t.Cleanup(func() {
				operation := operation{
					continueOnError: false,
					timeout:         timeout.DefaultCleanupTimeout,
					operation:       opdelete.New(p.client, namespace),
				}
				operation.execute(ctx)
			})
			if err := p.client.Create(logging.IntoContext(ctx, setupLogger), namespace.DeepCopy()); err != nil {
				t.FailNow()
			}
		}
	}
	for i, step := range p.test.Spec.Steps {
		processor := p.CreateStepProcessor(nspacer, step)
		name := step.Name
		if name == "" {
			name = fmt.Sprintf("step-%d", i+1)
		}
		processor.Run(logging.IntoContext(ctx, logging.NewLogger(t, p.clock, p.test.Name, fmt.Sprintf("%-*s", size, name))))
	}
}

func (p *testProcessor) CreateStepProcessor(nspacer namespacer.Namespacer, step v1alpha1.TestSpecStep) StepProcessor {
	return NewStepProcessor(p.config, p.client, nspacer, p.clock, p.test, step)
}
