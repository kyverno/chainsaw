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
	Run(ctx context.Context, test discovery.Test, step v1alpha1.TestStepSpec)
	CreateOperationProcessor(operation v1alpha1.Operation) OperationProcessor
}

func NewStepProcessor(config v1alpha1.ConfigurationSpec, client operations.OperationClient, namespacer namespacer.Namespacer, clock clock.PassiveClock) StepProcessor {
	return &stepProcessor{
		config:          config,
		operationClient: client,
		namespacer:      namespacer,
		clock:           clock,
	}
}

type stepProcessor struct {
	config          v1alpha1.ConfigurationSpec
	operationClient operations.OperationClient
	namespacer      namespacer.Namespacer
	clock           clock.PassiveClock
}

func (p *stepProcessor) Run(ctx context.Context, test discovery.Test, step v1alpha1.TestStepSpec) {
	t := testing.FromContext(ctx)
	logger := logging.FromContext(ctx)
	defer func() {
		t.Cleanup(func() {
			for _, handler := range step.Finally {
				if handler.PodLogs != nil {
					cmd, err := collect.PodLogs(handler.PodLogs)
					if err != nil {
						logger.Log("COLLEC", color.BoldRed, err)
						t.Fail()
					} else if err := p.operationClient.Command(ctx, nil, *cmd); err != nil {
						t.Fail()
					}
				}
				if handler.Events != nil {
					cmd, err := collect.Events(handler.Events)
					if err != nil {
						logger.Log("COLLEC", color.BoldRed, err)
						t.Fail()
					} else if err := p.operationClient.Command(ctx, nil, *cmd); err != nil {
						t.Fail()
					}
				}
				if handler.Command != nil {
					if err := p.operationClient.Command(ctx, nil, *handler.Command); err != nil {
						t.Fail()
					}
				}
				if handler.Script != nil {
					if err := p.operationClient.Script(ctx, nil, *handler.Script); err != nil {
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
					if handler.PodLogs != nil {
						cmd, err := collect.PodLogs(handler.PodLogs)
						if err != nil {
							logger.Log("COLLEC", color.BoldRed, err)
							t.Fail()
						} else if err := p.operationClient.Command(ctx, nil, *cmd); err != nil {
							t.Fail()
						}
					}
					if handler.Events != nil {
						cmd, err := collect.Events(handler.Events)
						if err != nil {
							logger.Log("COLLEC", color.BoldRed, err)
							t.Fail()
						} else if err := p.operationClient.Command(ctx, nil, *cmd); err != nil {
							t.Fail()
						}
					}
					if handler.Command != nil {
						if err := p.operationClient.Command(ctx, nil, *handler.Command); err != nil {
							t.Fail()
						}
					}
					if handler.Script != nil {
						if err := p.operationClient.Script(ctx, nil, *handler.Script); err != nil {
							t.Fail()
						}
					}
				}
			})
		}
	}()
	for _, operation := range step.Try {
		processor := p.CreateOperationProcessor(operation)
		processor.Run(ctx, test, step, operation)
	}
}

func (p *stepProcessor) CreateOperationProcessor(_ v1alpha1.Operation) OperationProcessor {
	return NewOperationProcessor(p.config, p.operationClient, p.namespacer, p.clock)
}
