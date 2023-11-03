package processors

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/runner/collect"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/runner/testing"
	"github.com/kyverno/kyverno/ext/output/color"
	"k8s.io/utils/clock"
)

type StepProcessor interface {
	Run(ctx context.Context, nspacer namespacer.Namespacer, test discovery.Test, step v1alpha1.TestSpecStep)
	// TODO
	// CreateOperationProcessor(operation v1alpha1.Operation) OperationProcessor
}

func NewStepProcessor(config v1alpha1.ConfigurationSpec, client client.Client, clock clock.PassiveClock) StepProcessor {
	return &stepProcessor{
		config: config,
		client: client,
		clock:  clock,
	}
}

type stepProcessor struct {
	config v1alpha1.ConfigurationSpec
	client client.Client
	clock  clock.PassiveClock
}

func (r *stepProcessor) Run(ctx context.Context, nspacer namespacer.Namespacer, test discovery.Test, step v1alpha1.TestSpecStep) {
	t := testing.FromContext(ctx)
	logger := logging.FromContext(ctx)
	operationsClient := operations.NewClient(nspacer, r.client, r.config, test.Spec, step.Spec)
	defer func() {
		t.Cleanup(func() {
			for _, handler := range step.Spec.Finally {
				collectors, err := collect.Commands(handler.Collect)
				if err != nil {
					logger.Log("COLLEC", color.BoldRed, err)
					t.Fail()
				} else {
					for _, collector := range collectors {
						exec := v1alpha1.Exec{
							Command: collector,
						}
						if err := operationsClient.Exec(ctx, exec, true, nspacer.GetNamespace()); err != nil {
							t.Fail()
						}
					}
				}
				if handler.Exec != nil {
					if err := operationsClient.Exec(ctx, *handler.Exec, true, nspacer.GetNamespace()); err != nil {
						t.Fail()
					}
				}
			}
		})
	}()
	defer func() {
		if t.Failed() {
			t.Cleanup(func() {
				for _, handler := range step.Spec.Catch {
					collectors, err := collect.Commands(handler.Collect)
					if err != nil {
						logger.Log("COLLEC", color.BoldRed, err)
						t.Fail()
					} else {
						for _, collector := range collectors {
							exec := v1alpha1.Exec{
								Command: collector,
							}
							if err := operationsClient.Exec(ctx, exec, true, nspacer.GetNamespace()); err != nil {
								t.Fail()
							}
						}
					}
					if handler.Exec != nil {
						if err := operationsClient.Exec(ctx, *handler.Exec, true, nspacer.GetNamespace()); err != nil {
							t.Fail()
						}
					}
				}
			})
		}
	}()
	for _, operation := range step.Spec.Try {
		// TODO
		runner := NewOperationProcessor(r.config, operationsClient, r.clock)
		runner.Run(ctx, nspacer, test, step, operation)
	}
}

// TODO
// func (r *stepProcessor) CreateOperationProcessor(_ v1alpha1.Operation) OperationProcessor {
// 	return NewOperationProcessor(r.config, operationsClient, r.clock)
// }
