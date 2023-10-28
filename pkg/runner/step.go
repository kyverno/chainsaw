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
				if step.Spec.OnFailure == nil {
					return
				}
				for _, failures := range step.Spec.OnFailure {
					if failures.Collect != nil {
						for _, collector := range collect(*failures.Collect) {
							cmd := collector
							exec := v1alpha1.Exec{
								Command: &cmd,
							}
							if err := operations.Exec(context.Background(), logger, exec, true, ctx.namespacer.GetNamespace()); err != nil {
								t.Fail()
							}
						}
					}
				}
			})
		}
	}()
	stepCtx, cancel := timeoutCtx(config, test, step.Spec)
	defer cancel()
	for _, operation := range step.Spec.Delete {
		var resource unstructured.Unstructured
		resource.SetAPIVersion(operation.APIVersion)
		resource.SetKind(operation.Kind)
		resource.SetName(operation.Name)
		resource.SetNamespace(operation.Namespace)
		resource.SetLabels(operation.Labels)
		if err := ctx.namespacer.Apply(&resource); err != nil {
			fail(t, operation.ContinueOnError)
		}
		if err := operations.Delete(stepCtx, logger, &resource, c); err != nil {
			fail(t, operation.ContinueOnError)
		}
	}
	var cleanup operations.CleanupFunc
	if skip := skipDelete(config, test, step.Spec); skip == nil || !*skip {
		cleanup = func(obj ctrlclient.Object, c client.Client) {
			t.Cleanup(func() {
				cleanupCtx, cancel := timeoutCtx(config, test, step.Spec)
				defer cancel()
				if err := operations.Delete(cleanupCtx, logger, obj, c); err != nil {
					t.Fail()
				}
			})
		}
	}
	for _, operation := range step.Spec.Exec {
		func() {
			cmdCtx, cancel := timeoutExecCtx(operation.Exec, config, test, step.Spec)
			defer cancel()
			if err := operations.Exec(cmdCtx, logger, operation.Exec, !operation.SkipLogOutput, ctx.namespacer.GetNamespace()); err != nil {
				fail(t, operation.ContinueOnError)
			}
		}()
	}
	for _, operation := range step.Spec.Apply {
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
			if err := operations.Apply(stepCtx, logger, resource, c, shouldFail, cleanup); err != nil {
				fail(t, operation.ContinueOnError)
			}
		}
	}
	for _, operation := range step.Spec.Assert {
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
			if err := operations.Assert(stepCtx, logger, resources[i], c); err != nil {
				fail(t, operation.ContinueOnError)
			}
		}
	}
	for _, operation := range step.Spec.Error {
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
			if err := operations.Error(stepCtx, logger, resources[i], c); err != nil {
				fail(t, operation.ContinueOnError)
			}
		}
	}
}
