package runner

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/resource"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/kyverno/ext/output/color"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func fail(t *testing.T, continueOnError *bool) {
	t.Helper()
	if continueOnError != nil && *continueOnError {
		t.Fail()
	} else {
		t.FailNow()
	}
}

func executeStep(t *testing.T, logger logging.Logger, ctx Context, basePath string, config v1alpha1.ConfigurationSpec, test v1alpha1.TestSpec, step v1alpha1.TestSpecStep) {
	t.Helper()
	c := ctx.clientFactory(logger)
	defer func() {
		if t.Failed() {
			t.Cleanup(func() {
				for _, handler := range step.Spec.OnFailure {
					collectors, err := collect(handler.Collect)
					if err != nil {
						logger.Log("COLLEC", color.BoldRed, err)
						t.Fail()
					} else {
						for _, collector := range collectors {
							exec := v1alpha1.Exec{
								Command: collector,
							}
							if err := operations.Exec(context.Background(), logger, exec, true, ctx.namespacer.GetNamespace()); err != nil {
								t.Fail()
							}
						}
					}
					if handler.Exec != nil {
						if err := operations.Exec(context.Background(), logger, *handler.Exec, true, ctx.namespacer.GetNamespace()); err != nil {
							t.Fail()
						}
					}
				}
			})
		}
	}()
	for _, operation := range step.Spec.Delete {
		func() {
			var resource unstructured.Unstructured
			resource.SetAPIVersion(operation.APIVersion)
			resource.SetKind(operation.Kind)
			resource.SetName(operation.Name)
			resource.SetNamespace(operation.Namespace)
			resource.SetLabels(operation.Labels)
			if err := ctx.namespacer.Apply(&resource); err != nil {
				fail(t, operation.ContinueOnError)
			}
			operationCtx, cancel := timeoutCtx(DefaultDeleteTimeout, config.Timeouts.Delete, test.Timeouts.Delete, step.Spec.Timeouts.Delete, operation.Timeout)
			defer cancel()
			if err := operations.Delete(operationCtx, logger, &resource, c); err != nil {
				fail(t, operation.ContinueOnError)
			}
		}()
	}
	var cleanup operations.CleanupFunc
	if !skipDelete(config.SkipDelete, test.SkipDelete, step.Spec.SkipDelete) {
		cleanup = func(obj ctrlclient.Object, c client.Client) {
			t.Cleanup(func() {
				cleanupCtx, cancel := timeoutCtx(DefaultCleanupTimeout, config.Timeouts.Cleanup, test.Timeouts.Cleanup, step.Spec.Timeouts.Cleanup, nil)
				defer cancel()
				if err := operations.Delete(cleanupCtx, logger, obj, c); err != nil {
					t.Fail()
				}
			})
		}
	}
	for _, operation := range step.Spec.Exec {
		func() {
			operationCtx, cancel := timeoutCtx(DefaultExecTimeout, config.Timeouts.Exec, test.Timeouts.Exec, step.Spec.Timeouts.Exec, operation.Timeout)
			defer cancel()
			if err := operations.Exec(operationCtx, logger, operation.Exec, !operation.SkipLogOutput, ctx.namespacer.GetNamespace()); err != nil {
				fail(t, operation.ContinueOnError)
			}
		}()
	}
	for _, operation := range step.Spec.Apply {
		func() {
			resources, err := resource.Load(filepath.Join(basePath, operation.File))
			if err != nil {
				logger.Log("LOAD  ", color.BoldRed, err)
				fail(t, operation.ContinueOnError)
			}
			shouldFail := operation.ShouldFail != nil && *operation.ShouldFail
			for i := range resources {
				resource := &resources[i]
				if err := ctx.namespacer.Apply(resource); err != nil {
					logger.Log("LOAD  ", color.BoldRed, err)
					fail(t, operation.ContinueOnError)
				}
				operationCtx, cancel := timeoutCtx(DefaultApplyTimeout, config.Timeouts.Apply, test.Timeouts.Apply, step.Spec.Timeouts.Apply, operation.Timeout)
				defer cancel()
				if err := operations.Apply(operationCtx, logger, resource, c, shouldFail, cleanup); err != nil {
					fail(t, operation.ContinueOnError)
				}
			}
		}()
	}
	for _, operation := range step.Spec.Assert {
		func() {
			resources, err := resource.Load(filepath.Join(basePath, operation.File))
			if err != nil {
				logger.Log("LOAD  ", color.BoldRed, err)
				fail(t, operation.ContinueOnError)
			}
			for i := range resources {
				resource := &resources[i]
				if err := ctx.namespacer.Apply(resource); err != nil {
					logger.Log("LOAD  ", color.BoldRed, err)
					fail(t, operation.ContinueOnError)
				}
				operationCtx, cancel := timeoutCtx(DefaultAssertTimeout, config.Timeouts.Assert, test.Timeouts.Assert, step.Spec.Timeouts.Assert, operation.Timeout)
				defer cancel()
				if err := operations.Assert(operationCtx, logger, resources[i], c); err != nil {
					fail(t, operation.ContinueOnError)
				}
			}
		}()
	}
	for _, operation := range step.Spec.Error {
		func() {
			resources, err := resource.Load(filepath.Join(basePath, operation.File))
			if err != nil {
				logger.Log("LOAD  ", color.BoldRed, err)
				fail(t, operation.ContinueOnError)
			}
			for i := range resources {
				resource := &resources[i]
				if err := ctx.namespacer.Apply(resource); err != nil {
					logger.Log("LOAD  ", color.BoldRed, err)
					fail(t, operation.ContinueOnError)
				}
				operationCtx, cancel := timeoutCtx(DefaultErrorTimeout, config.Timeouts.Error, test.Timeouts.Error, step.Spec.Timeouts.Error, operation.Timeout)
				defer cancel()
				if err := operations.Error(operationCtx, logger, resources[i], c); err != nil {
					fail(t, operation.ContinueOnError)
				}
			}
		}()
	}
}
