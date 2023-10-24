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
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func executeStep(t *testing.T, logger logging.Logger, ctx Context, basePath string, config v1alpha1.ConfigurationSpec, test v1alpha1.TestSpec, step v1alpha1.TestSpecStep) {
	t.Helper()
	c := ctx.clientFactory(t, logger)
	stepCtx := context.Background()
	if timeout := timeout(config, test, step.Spec); timeout != nil {
		timeoutCtx, cancel := context.WithTimeout(stepCtx, *timeout)
		defer cancel()
		stepCtx = timeoutCtx
	}
	for _, delete := range step.Spec.Delete {
		var resource unstructured.Unstructured
		resource.SetAPIVersion(delete.APIVersion)
		resource.SetKind(delete.Kind)
		resource.SetName(delete.Name)
		resource.SetNamespace(delete.Namespace)
		resource.SetLabels(delete.Labels)
		if err := ctx.namespacer.Apply(&resource); err != nil {
			logger.Log(err)
			t.FailNow()
		}
		logging.ResourceOp(logger, "DELETE", client.ObjectKey(&resource), &resource)
		if err := operations.Delete(stepCtx, resource, c); err != nil {
			logger.Log(err)
			t.FailNow()
		}
	}
	for _, apply := range step.Spec.Apply {
		resources, err := resource.Load(filepath.Join(basePath, apply.File))
		if err != nil {
			logger.Log(err)
			t.FailNow()
		}
		for i := range resources {
			resource := &resources[i]
			if err := ctx.namespacer.Apply(resource); err != nil {
				logger.Log(err)
				t.FailNow()
			}
			logging.ResourceOp(logger, "APPLY", client.ObjectKey(resource), resource)
			err := operations.Apply(stepCtx, resource, c)
			if err != nil {
				logger.Log(err)
				t.FailNow()
			}
		}
	}
	for _, assert := range step.Spec.Assert {
		resources, err := resource.Load(filepath.Join(basePath, assert.File))
		if err != nil {
			logger.Log(err)
			t.FailNow()
		}
		for i := range resources {
			resource := &resources[i]
			if err := ctx.namespacer.Apply(resource); err != nil {
				logger.Log(err)
				t.FailNow()
			}
			logging.ResourceOp(logger, "ASSERT", client.ObjectKey(resource), resource)
			if err := operations.Assert(stepCtx, resources[i], c); err != nil {
				logger.Log(err)
				t.FailNow()
			}
		}
	}
	for _, e := range step.Spec.Error {
		resources, err := resource.Load(filepath.Join(basePath, e.File))
		if err != nil {
			logger.Log(err)
			t.FailNow()
		}
		for i := range resources {
			resource := &resources[i]
			if err := ctx.namespacer.Apply(resource); err != nil {
				logger.Log(err)
				t.FailNow()
			}
			logging.ResourceOp(logger, "ERROR", client.ObjectKey(resource), resource)
			// Using the Error function to handle the error assertion
			err := operations.Error(stepCtx, resources[i], c)
			if err != nil {
				logger.Log(err)
				t.FailNow()
			}
		}
	}
}
