package processors

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/report"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/runner/summary"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/kyverno/ext/output/color"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/utils/clock"
)

type TestProcessor interface {
	Run(ctx context.Context, nspacer namespacer.Namespacer, test discovery.Test) report.TestReport
	CreateStepProcessor(client operations.OperationClient) StepProcessor
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

func (p *testProcessor) Run(ctx context.Context, nspacer namespacer.Namespacer, test discovery.Test) report.TestReport {
	t := testing.FromContext(ctx)
	// Create a TestReport
	testReport := report.TestReport{
		Name:      test.Name,
		StartTime: p.clock.Now(),
	}
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
	setupLogger := logging.NewLogger(t, p.clock, test.Name, fmt.Sprintf("%-*s", size, "@setup"))
	var namespace *corev1.Namespace
	if nspacer == nil || test.Spec.Namespace != "" {
		var ns corev1.Namespace
		if test.Spec.Namespace != "" {
			ns = client.Namespace(test.Spec.Namespace)
		} else {
			ns = client.PetNamespace()
		}
		namespace = &ns
	}
	if namespace != nil {
		nspacer = namespacer.New(p.client, namespace.Name)
		operationsClient := operations.NewOperationClient(nspacer, p.client, p.config, test.Spec, v1alpha1.Timeouts{})
		if err := p.client.Get(logging.IntoContext(ctx, setupLogger), client.ObjectKey(namespace), namespace.DeepCopy()); err != nil {
			if !errors.IsNotFound(err) {
				// Get doesn't log
				setupLogger.Log("GET   ", color.BoldRed, err)
				t.FailNow()
			}
			if err := p.client.Create(logging.IntoContext(ctx, setupLogger), namespace.DeepCopy()); err != nil {
				t.FailNow()
			}
			t.Cleanup(func() {
				if err := operationsClient.Delete(logging.IntoContext(ctx, setupLogger), nil, namespace); err != nil {
					t.FailNow()
				}
			})
		}
	}
	for i, step := range test.Spec.Steps {
		operationsClient := operations.NewOperationClient(nspacer, p.client, p.config, test.Spec, step.Spec.Timeouts)
		processor := p.CreateStepProcessor(operationsClient)
		name := step.Name
		if name == "" {
			name = fmt.Sprintf("step-%d", i+1)
		}
		stepReport := processor.Run(logging.IntoContext(ctx, logging.NewLogger(t, p.clock, test.Name, fmt.Sprintf("%-*s", size, name))), nspacer, test, step.Spec)
		testReport.Steps = append(testReport.Steps, stepReport)
	}
	testReport.EndTime = p.clock.Now()
	return testReport
}

func (p *testProcessor) CreateStepProcessor(client operations.OperationClient) StepProcessor {
	return NewStepProcessor(p.config, client, p.clock)
}
