package runner

import (
	"context"
	"fmt"
	"path/filepath"
	"sync/atomic"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/resource"
	runnerclient "github.com/kyverno/chainsaw/pkg/runner/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/rest"
)

func Run(cfg *rest.Config, config v1alpha1.ConfigurationSpec, tests ...discovery.Test) (int, *Summary, error) {
	var summary Summary
	if len(tests) == 0 {
		return 0, &summary, nil
	}
	if err := setupFlags(config); err != nil {
		return 0, nil, err
	}
	var internalTests []testing.InternalTest
	var failed, passed, skipped atomic.Int32
	defer func() {
		summary.FailedTests = failed.Load()
		summary.PassedTests = passed.Load()
		summary.SkippedTests = skipped.Load()
	}()
	c, err := client.New(cfg)
	if err != nil {
		return 0, nil, err
	}
	var nspacer namespacer.Namespacer
	if config.Namespace != "" {
		nspacer = namespacer.New(c, config.Namespace)
		namespace := client.Namespace(config.Namespace)
		if err := c.Get(context.Background(), client.ObjectKey(&namespace), namespace.DeepCopy()); err != nil {
			if !errors.IsNotFound(err) {
				return 0, nil, err
			}
			if err := c.Create(context.Background(), namespace.DeepCopy()); err != nil {
				return 0, nil, err
			}
			defer func() {
				if err := client.BlockingDelete(context.Background(), c, &namespace); err != nil {
					panic(err)
				}
			}()
		}
	}
	for i := range tests {
		test := tests[i]
		name, err := testName(config, test)
		if err != nil {
			return 0, nil, err
		}
		internalTests = append(internalTests, testing.InternalTest{
			Name: name,
			F: func(t *testing.T) {
				t.Helper()
				t.Parallel()
				ctx := Context{
					clientFactory: func(t *testing.T, logger logging.Logger) client.Client {
						t.Helper()
						return runnerclient.New(t, logger, c, !config.SkipDelete)
					},
					namespacer: nspacer,
				}
				t.Cleanup(func() {
					if t.Skipped() {
						skipped.Add(1)
					} else {
						if t.Failed() {
							failed.Add(1)
						} else {
							passed.Add(1)
						}
					}
				})
				runTest(t, ctx, test)
			},
		})
	}
	var testDeps testDeps
	m := testing.MainStart(&testDeps, internalTests, nil, nil, nil)
	return m.Run(), &summary, nil
}

func runTest(t *testing.T, ctx Context, test discovery.Test) {
	t.Helper()
	if ctx.namespacer == nil {
		namespace := client.PetNamespace()
		c := ctx.clientFactory(t, logging.NewTestLogger(t))
		if err := c.Create(context.Background(), namespace.DeepCopy()); err != nil {
			t.Fatal(err)
		}
		ctx.namespacer = namespacer.New(c, namespace.Name)
	}
	for i := range test.Spec.Steps {
		step := test.Spec.Steps[i]
		executeStep(t, logging.NewStepLogger(t, fmt.Sprintf("step-%d", i+1)), ctx, test.BasePath, step)
	}
}

func executeStep(t *testing.T, logger logging.Logger, ctx Context, basePath string, step v1alpha1.TestStepSpec) {
	t.Helper()
	c := ctx.clientFactory(t, logger)
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
				if err := client.DeleteResource(context.TODO(), c, &currentItem); err != nil {
					t.Fatal(err)
				}
			}
		} else {
			resource := client.NewResource(delete.APIVersion, delete.Kind, delete.Name, delete.Namespace)
			t.Logf("=== DELETE %s/%s", delete.APIVersion, delete.Kind)
			if err := client.DeleteResource(context.TODO(), c, resource); err != nil {
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
			err := client.CreateOrUpdate(context.Background(), c, resource)
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
			if err := client.Assert(context.Background(), resources[i], c); err != nil {
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
			err := client.Error(context.Background(), resources[i], c)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}
