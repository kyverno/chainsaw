package processors

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/summary"
	"github.com/kyverno/chainsaw/pkg/runner/testing"
	"github.com/kyverno/kyverno/ext/output/color"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/utils/clock"
)

type TestProcessor interface {
	Run(ctx context.Context, nspacer namespacer.Namespacer, test discovery.Test)
	CreateStepProcessor(step v1alpha1.TestSpecStep) StepProcessor
}

func NewTestProcessor(config v1alpha1.ConfigurationSpec, client client.Client, clock clock.PassiveClock, summary *summary.Summary) TestProcessor {
	return &testProcessor{
		config:  config,
		client:  client,
		clock:   clock,
		summary: summary,
	}
}

type testProcessor struct {
	config  v1alpha1.ConfigurationSpec
	client  client.Client
	clock   clock.PassiveClock
	summary *summary.Summary
}

func (p *testProcessor) Run(ctx context.Context, nspacer namespacer.Namespacer, test discovery.Test) {
	t := testing.FromContext(ctx)
	size := 0
	for i, step := range test.Spec.Steps {
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
	if test.Spec.Concurrent == nil || *test.Spec.Concurrent {
		t.Parallel()
	}
	if test.Spec.Skip != nil && *test.Spec.Skip {
		t.SkipNow()
	}
	beginLogger := logging.NewLogger(t, p.clock, test.Name, fmt.Sprintf("%-*s", size, "@begin"))
	cleanLogger := logging.NewLogger(t, p.clock, test.Name, fmt.Sprintf("%-*s", size, "@clean"))
	if test.Spec.Namespace != "" {
		namespace := client.Namespace(test.Spec.Namespace)
		if err := p.client.Get(logging.IntoContext(ctx, beginLogger), client.ObjectKey(&namespace), namespace.DeepCopy()); err != nil {
			if !errors.IsNotFound(err) {
				// Get doesn't log
				beginLogger.Log("GET   ", color.BoldRed, err)
				t.FailNow()
			}
			if err := p.client.Create(logging.IntoContext(ctx, beginLogger), namespace.DeepCopy()); err != nil {
				t.FailNow()
			}
			t.Cleanup(func() {
				// TODO: wait
				if err := p.client.Delete(logging.IntoContext(ctx, cleanLogger), &namespace); err != nil {
					t.FailNow()
				}
			})
		}
		nspacer = namespacer.New(p.client, test.Spec.Namespace)
	}
	if nspacer == nil {
		namespace := client.PetNamespace()
		if err := p.client.Create(logging.IntoContext(ctx, beginLogger), namespace.DeepCopy()); err != nil {
			t.FailNow()
		}
		t.Cleanup(func() {
			// TODO: wait
			if err := p.client.Delete(logging.IntoContext(ctx, cleanLogger), &namespace); err != nil {
				t.FailNow()
			}
		})
		nspacer = namespacer.New(p.client, namespace.Name)
	}
	for i, step := range test.Spec.Steps {
		processor := p.CreateStepProcessor(step)
		name := step.Name
		if name == "" {
			name = fmt.Sprintf("step-%d", i+1)
		}
		processor.Run(logging.IntoContext(ctx, logging.NewLogger(t, p.clock, test.Name, fmt.Sprintf("%-*s", size, name))), nspacer, test, step)
	}
}

func (p *testProcessor) CreateStepProcessor(_ v1alpha1.TestSpecStep) StepProcessor {
	return NewStepProcessor(p.config, p.client, p.clock)
}
