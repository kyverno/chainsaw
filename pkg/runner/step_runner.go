package runner

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/runner/collect"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/runner/testing"
	"github.com/kyverno/kyverno/ext/output/color"
	"k8s.io/utils/clock"
)

type stepRunner struct {
	config v1alpha1.ConfigurationSpec
	client client.Client
	clock  clock.PassiveClock
}

func (r *stepRunner) runStep(goctx context.Context, ctx Context, test discovery.Test, step v1alpha1.TestSpecStep) {
	t := testing.FromContext(goctx)
	t.Helper()
	logger := logging.FromContext(goctx)
	operationsClient := operations.NewClient(
		ctx.namespacer,
		ctx.client,
		r.config,
		test.Spec,
		step.Spec,
	)
	defer func() {
		if t.Failed() {
			t.Cleanup(func() {
				for _, handler := range step.Spec.OnFailure {
					collectors, err := collect.Commands(handler.Collect)
					if err != nil {
						logger.Log("COLLEC", color.BoldRed, err)
						t.Fail()
					} else {
						for _, collector := range collectors {
							exec := v1alpha1.Exec{
								Command: collector,
							}
							if err := operationsClient.Exec(goctx, exec, true, ctx.namespacer.GetNamespace()); err != nil {
								t.Fail()
							}
						}
					}
					if handler.Exec != nil {
						if err := operationsClient.Exec(goctx, *handler.Exec, true, ctx.namespacer.GetNamespace()); err != nil {
							t.Fail()
						}
					}
				}
			})
		}
	}()
	for _, operation := range step.Spec.Operations {
		runner := operationRunner{
			config: r.config,
			client: operationsClient,
			clock:  r.clock,
		}
		runner.executeOperation(goctx, ctx, test, step, operation)
	}
}
