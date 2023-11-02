package runner

import (
	"path/filepath"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/resource"
	"github.com/kyverno/chainsaw/pkg/runner/cleanup"
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
	operationsClient := operations.NewClient(
		logger,
		ctx.namespacer,
		c,
		config,
		test,
		step.Spec,
	)
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
							if err := operationsClient.Exec(exec, true, ctx.namespacer.GetNamespace()); err != nil {
								t.Fail()
							}
						}
					}
					if handler.Exec != nil {
						if err := operationsClient.Exec(*handler.Exec, true, ctx.namespacer.GetNamespace()); err != nil {
							t.Fail()
						}
					}
				}
			})
		}
	}()
	for _, operation := range step.Spec.Delete {
		var resource unstructured.Unstructured
		resource.SetAPIVersion(operation.APIVersion)
		resource.SetKind(operation.Kind)
		resource.SetName(operation.Name)
		resource.SetNamespace(operation.Namespace)
		resource.SetLabels(operation.Labels)
		if err := operationsClient.Delete(operation.Timeout, &resource); err != nil {
			fail(t, operation.ContinueOnError)
		}
	}
	var doCleanup operations.CleanupFunc
	if !cleanup.Skip(config.SkipDelete, test.SkipDelete, step.Spec.SkipDelete) {
		doCleanup = func(obj ctrlclient.Object, c client.Client) {
			t.Cleanup(func() {
				if err := operationsClient.Delete(nil, obj); err != nil {
					t.Fail()
				}
			})
		}
	}
	for _, operation := range step.Spec.Exec {
		if err := operationsClient.Exec(operation.Exec, !operation.SkipLogOutput, ctx.namespacer.GetNamespace()); err != nil {
			fail(t, operation.ContinueOnError)
		}
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
			if err := operationsClient.Apply(operation.Timeout, resource, shouldFail, doCleanup); err != nil {
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
		for _, resource := range resources {
			if err := operationsClient.Assert(operation.Timeout, resource); err != nil {
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
		for _, resource := range resources {
			if err := operationsClient.Error(operation.Timeout, resource); err != nil {
				fail(t, operation.ContinueOnError)
			}
		}
	}
}
