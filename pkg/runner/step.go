package runner

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/resource"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
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
	// Delete the Objects before the test step is executed
	for _, delete := range step.Delete {
		// Use your dynamic listing logic if the name is not provided
		if delete.Name == "" {
			u, err := client.ListResourcesToDelete(c, delete)
			if err != nil {
				t.Fatal(err)
			}
			for _, item := range u.Items {
				currentItem := item
				t.Logf("=== DELETE %s/%s", delete.APIVersion, delete.Kind)
				if err := client.DeleteResource(stepCtx, c, &currentItem); err != nil {
					t.Fatal(err)
				}
			}
		} else {
			resource := client.NewResource(delete.APIVersion, delete.Kind, delete.Name, delete.Namespace)
			t.Logf("=== DELETE %s/%s", delete.APIVersion, delete.Kind)
			if err := client.DeleteResource(stepCtx, c, resource); err != nil {
				t.Fatal(err)
			}
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
			err := client.CreateOrUpdate(stepCtx, c, resource)
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
			if err := client.Assert(stepCtx, resources[i], c); err != nil {
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
			err := client.Error(stepCtx, resources[i], c)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}
