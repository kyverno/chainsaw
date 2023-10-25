package runner

import (
	"path/filepath"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/resource"
	runnerclient "github.com/kyverno/chainsaw/pkg/runner/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func executeStep(
	t *testing.T,
	logger logging.Logger,
	client runnerclient.Client,
	ctx Context,
	basePath string,
	config v1alpha1.ConfigurationSpec,
	test v1alpha1.TestSpec,
	step v1alpha1.TestSpecStep,
) {
	t.Helper()
	stepCtx, cancel := timeoutCtx(config, test, step.Spec)
	defer cancel()
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
		if err := client.Delete(stepCtx, &resource); err != nil {
			logger.Log(err)
			t.FailNow()
		}
	}
	// var cleanup operations.CleanupFunc
	// if skip := skipDelete(config, test, step.Spec); skip == nil || !*skip {
	// 	cleanup = func(obj ctrlclient.Object, c client.Client) {
	// 		t.Cleanup(func() {
	// 			cleanupCtx, cancel := timeoutCtx(config, test, step.Spec)
	// 			defer cancel()
	// 			if err := client.Delete(cleanupCtx, logger, obj, c); err != nil {
	// 				logger.Log(err)
	// 				t.FailNow()
	// 			}
	// 		})
	// 	}
	// }
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
			if err := client.Apply(stepCtx, resource); err != nil {
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
			if err := client.Assert(stepCtx, resources[i]); err != nil {
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
			if err := client.Error(stepCtx, resources[i]); err != nil {
				logger.Log(err)
				t.FailNow()
			}
		}
	}
}
