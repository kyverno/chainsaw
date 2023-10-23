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

func executeStep(t *testing.T, logger logging.Logger, ctx Context, basePath string, step v1alpha1.TestStepSpec) {
	t.Helper()
	c := ctx.clientFactory(t, logger)
	stepCtx := ctx.ctx
	if ctx.config.Timeout != nil {
		timeoutCtx, cancel := context.WithTimeout(stepCtx, ctx.config.Timeout.Duration)
		defer cancel()
		stepCtx = timeoutCtx
	}
	for _, delete := range step.Delete {
		var resource unstructured.Unstructured
		resource.SetAPIVersion(delete.APIVersion)
		resource.SetKind(delete.Kind)
		resource.SetName(delete.Name)
		resource.SetNamespace(delete.Namespace)
		resource.SetLabels(delete.Labels)
		if err := ctx.namespacer.Apply(&resource); err != nil {
			t.Fatal(err)
		}
		logging.ResourceOp(logger, "DELETE", client.ObjectKey(&resource), &resource)
		if err := operations.Delete(stepCtx, resource, c); err != nil {
			t.Fatal(err)
		}
	}
	for _, apply := range step.Apply {
		resources, err := resource.Load(filepath.Join(basePath, apply.File))
		if err != nil {
			t.Fatal(err)
		}
		for i := range resources {
			resource := &resources[i]
			if err := ctx.namespacer.Apply(resource); err != nil {
				t.Fatal(err)
			}
			logging.ResourceOp(logger, "APPLY", client.ObjectKey(resource), resource)
			err := operations.Apply(stepCtx, resource, c)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
	for _, assert := range step.Assert {
		resources, err := resource.Load(filepath.Join(basePath, assert.File))
		if err != nil {
			t.Fatal(err)
		}
		for i := range resources {
			resource := &resources[i]
			if err := ctx.namespacer.Apply(resource); err != nil {
				t.Fatal(err)
			}
			logging.ResourceOp(logger, "ASSERT", client.ObjectKey(resource), resource)
			if err := operations.Assert(stepCtx, resources[i], c); err != nil {
				t.Fatal(err)
			}
		}
	}
	for _, e := range step.Error {
		resources, err := resource.Load(filepath.Join(basePath, e.File))
		if err != nil {
			t.Fatal(err)
		}
		for i := range resources {
			resource := &resources[i]
			if err := ctx.namespacer.Apply(resource); err != nil {
				t.Fatal(err)
			}
			logging.ResourceOp(logger, "ERROR", client.ObjectKey(resource), resource)
			// Using the Error function to handle the error assertion
			err := operations.Error(stepCtx, resources[i], c)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}
