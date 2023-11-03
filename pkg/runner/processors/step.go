package processors

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/runner/collect"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/kyverno/ext/output/color"
	"k8s.io/utils/clock"
)

type StepProcessor interface {
	Run(ctx context.Context, nspacer namespacer.Namespacer, test discovery.Test, step v1alpha1.TestStepSpec)
	CreateOperationProcessor(operation v1alpha1.Operation) OperationProcessor
}

func NewStepProcessor(config v1alpha1.ConfigurationSpec, client operations.Client, clock clock.PassiveClock) StepProcessor {
	return &stepProcessor{
		config: config,
		client: client,
		clock:  clock,
	}
}

type stepProcessor struct {
	config v1alpha1.ConfigurationSpec
	client operations.Client
	clock  clock.PassiveClock
}

func (p *stepProcessor) Run(ctx context.Context, nspacer namespacer.Namespacer, test discovery.Test, step v1alpha1.TestStepSpec) {
	t := testing.FromContext(ctx)
	logger := logging.FromContext(ctx)
	defer func() {
		t.Cleanup(func() {
			for _, handler := range step.Finally {
				collectors, err := collect.Commands(handler.Collect)
				if err != nil {
					logger.Log("COLLEC", color.BoldRed, err)
					t.Fail()
				} else {
					for _, collector := range collectors {
						exec := v1alpha1.Exec{
							Command: collector,
						}
						if err := p.client.Exec(ctx, exec, true, nspacer.GetNamespace()); err != nil {
							t.Fail()
						}
					}
				}
				if handler.Exec != nil {
					if err := p.client.Exec(ctx, *handler.Exec, true, nspacer.GetNamespace()); err != nil {
						t.Fail()
					}
				}
			}
		})
	}()
	defer func() {
		if t.Failed() {
			t.Cleanup(func() {
				for _, handler := range step.Catch {
					collectors, err := collect.Commands(handler.Collect)
					if err != nil {
						logger.Log("COLLEC", color.BoldRed, err)
						t.Fail()
					} else {
						for _, collector := range collectors {
							exec := v1alpha1.Exec{
								Command: collector,
							}
							if err := p.client.Exec(ctx, exec, true, nspacer.GetNamespace()); err != nil {
								t.Fail()
							}
						}
					}
					if handler.Exec != nil {
						if err := p.client.Exec(ctx, *handler.Exec, true, nspacer.GetNamespace()); err != nil {
							t.Fail()
						}
					}
				}
			})
		}
	}()
	for _, operation := range step.Try {
		processor := p.CreateOperationProcessor(operation)
		processor.Run(ctx, nspacer.GetNamespace(), test, step, operation)
	}
}

func (p *stepProcessor) CreateOperationProcessor(_ v1alpha1.Operation) OperationProcessor {
	return NewOperationProcessor(p.config, p.client, p.clock)
}
