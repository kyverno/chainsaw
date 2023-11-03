package runner

import (
	"context"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/resource"
	"github.com/kyverno/chainsaw/pkg/runner/cleanup"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/runner/testing"
	"github.com/kyverno/kyverno/ext/output/color"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/utils/clock"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type operationRunner struct {
	config v1alpha1.ConfigurationSpec
	client operations.Client
	clock  clock.PassiveClock
}

func (r *operationRunner) executeOperation(goctx context.Context, ctx Context, test discovery.Test, step v1alpha1.TestSpecStep, operation v1alpha1.Operation) {
	fail := func(t *testing.T, continueOnError *bool) {
		t.Helper()
		if continueOnError != nil && *continueOnError {
			t.Fail()
		} else {
			t.FailNow()
		}
	}
	t := testing.FromContext(goctx)
	t.Helper()
	// Handle Delete
	if operation.Delete != nil {
		for _, operation := range operation.Delete {
			var resource unstructured.Unstructured
			resource.SetAPIVersion(operation.APIVersion)
			resource.SetKind(operation.Kind)
			resource.SetName(operation.Name)
			resource.SetNamespace(operation.Namespace)
			resource.SetLabels(operation.Labels)
			if err := r.client.Delete(goctx, operation.Timeout, &resource); err != nil {
				fail(t, operation.ContinueOnError)
			}
		}
	}

	// Handle Exec
	if operation.Exec != nil {
		for _, operation := range operation.Exec {
			if err := r.client.Exec(goctx, operation.Exec, !operation.SkipLogOutput, ctx.namespacer.GetNamespace()); err != nil {
				fail(t, operation.ContinueOnError)
			}
		}
	}

	var doCleanup operations.CleanupFunc
	if !cleanup.Skip(r.config.SkipDelete, test.Spec.SkipDelete, step.Spec.SkipDelete) {
		doCleanup = func(obj ctrlclient.Object, c client.Client) {
			t.Cleanup(func() {
				if err := r.client.Delete(goctx, nil, obj); err != nil {
					t.Fail()
				}
			})
		}
	}
	// Handle Apply
	if operation.Apply != nil {
		for _, operation := range operation.Apply {
			resources, err := resource.Load(filepath.Join(test.BasePath, operation.File))
			if err != nil {
				logging.FromContext(goctx).Log("LOAD  ", color.BoldRed, err)
				fail(t, operation.ContinueOnError)
			}
			shouldFail := operation.ShouldFail != nil && *operation.ShouldFail
			for i := range resources {
				resource := &resources[i]
				if err := r.client.Apply(goctx, operation.Timeout, resource, shouldFail, doCleanup); err != nil {
					fail(t, operation.ContinueOnError)
				}
			}
		}
	}

	// Handle Assert
	if operation.Assert != nil {
		for _, operation := range operation.Assert {
			resources, err := resource.Load(filepath.Join(test.BasePath, operation.File))
			if err != nil {
				logging.FromContext(goctx).Log("LOAD  ", color.BoldRed, err)
				fail(t, operation.ContinueOnError)
			}
			for _, resource := range resources {
				if err := r.client.Assert(goctx, operation.Timeout, resource); err != nil {
					fail(t, operation.ContinueOnError)
				}
			}
		}
	}

	// Handle Error
	if operation.Error != nil {
		for _, operation := range operation.Error {
			resources, err := resource.Load(filepath.Join(test.BasePath, operation.File))
			if err != nil {
				logging.FromContext(goctx).Log("LOAD  ", color.BoldRed, err)
				fail(t, operation.ContinueOnError)
			}
			for _, resource := range resources {
				if err := r.client.Error(goctx, operation.Timeout, resource); err != nil {
					fail(t, operation.ContinueOnError)
				}
			}
		}
	}
}
