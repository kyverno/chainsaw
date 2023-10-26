package runner

import (
	"path/filepath"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/resource"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
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
			logger.Log(err)
			fail(t, operation.ContinueOnError)
		}
		if err := operations.Delete(stepCtx, logger, &resource, c); err != nil {
			logger.Log(err)
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
					logger.Log(err)
					t.Fail()
				}
			})
		}
	}
	for _, operation := range step.Spec.Apply {
		resources, err := resource.Load(filepath.Join(basePath, operation.File))
		if err != nil {
			logger.Log(err)
			fail(t, operation.ContinueOnError)
		}
		for i := range resources {
			resource := &resources[i]
			if err := ctx.namespacer.Apply(resource); err != nil {
				logger.Log(err)
				fail(t, operation.ContinueOnError)
			}
			if err := operations.Apply(stepCtx, logger, resource, c, cleanup); err != nil {
				logger.Log(err)
				fail(t, operation.ContinueOnError)
			}
		}
	}
	for _, operation := range step.Spec.Assert {
		resources, err := resource.Load(filepath.Join(basePath, operation.File))
		if err != nil {
			logger.Log(err)
			fail(t, operation.ContinueOnError)
		}
		for i := range resources {
			resource := &resources[i]
			if err := ctx.namespacer.Apply(resource); err != nil {
				logger.Log(err)
				fail(t, operation.ContinueOnError)
			}
			if err := operations.Assert(stepCtx, logger, resources[i], c); err != nil {
				logger.Log(err)
				fail(t, operation.ContinueOnError)
			}
		}
	}
	for _, operation := range step.Spec.Error {
		resources, err := resource.Load(filepath.Join(basePath, operation.File))
		if err != nil {
			logger.Log(err)
			fail(t, operation.ContinueOnError)
		}
		for i := range resources {
			resource := &resources[i]
			if err := ctx.namespacer.Apply(resource); err != nil {
				logger.Log(err)
				fail(t, operation.ContinueOnError)
			}
			if err := operations.Error(stepCtx, logger, resources[i], c); err != nil {
				logger.Log(err)
				fail(t, operation.ContinueOnError)
			}
		}
	}
}
