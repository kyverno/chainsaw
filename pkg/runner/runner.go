package runner

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
	testing.Init()
	if err := flag.Set("test.v", "true"); err != nil {
		return 0, nil, err
	}
	if err := flag.Set("test.parallel", strconv.Itoa(config.Parallel)); err != nil {
		return 0, nil, err
	}
	if err := flag.Set("test.timeout", config.Timeout.Duration.String()); err != nil {
		return 0, nil, err
	}
	if err := flag.Set("test.failfast", fmt.Sprint(config.FailFast)); err != nil {
		return 0, nil, err
	}
	if err := flag.Set("test.paniconexit0", "true"); err != nil {
		return 0, nil, err
	}
	if err := flag.Set("test.fullpath", "false"); err != nil {
		return 0, nil, err
	}
	if err := flag.Set("test.count", "1"); err != nil {
		return 0, nil, err
	}
	if err := flag.Set("test.run", config.IncludeTestRegex); err != nil {
		return 0, nil, err
	}
	if err := flag.Set("test.skip", config.ExcludeTestRegex); err != nil {
		return 0, nil, err
	}
	flag.Parse()
	run := func(t *testing.T) {
		t.Helper()
		run(t, cfg, config, &summary, tests...)
	}
	internalTest := []testing.InternalTest{{
		Name: "chainsaw",
		F:    run,
	}}
	var testDeps testDeps
	m := testing.MainStart(&testDeps, internalTest, nil, nil, nil)
	return m.Run(), &summary, nil
}

func run(t *testing.T, cfg *rest.Config, config v1alpha1.ConfigurationSpec, summary *Summary, tests ...discovery.Test) {
	t.Helper()
	c, err := client.New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	ctx := Context{
		clientFactory: func(t *testing.T, logger logging.Logger) client.Client {
			t.Helper()
			return runnerclient.New(t, logger, c, !config.SkipDelete)
		},
	}
	if config.Namespace != "" {
		namespace := client.Namespace(config.Namespace)
		c := ctx.clientFactory(t, logging.NewTestLogger(t))
		if err := c.Get(context.Background(), client.ObjectKey(&namespace), namespace.DeepCopy()); err != nil {
			if errors.IsNotFound(err) {
				if err := c.Create(context.Background(), namespace.DeepCopy()); err != nil {
					t.Fatal(err)
				}
				ctx.namespacer = namespacer.New(c, config.Namespace)
			}
		}
	}
	var failed, passed, skipped atomic.Int32
	t.Cleanup(func() {
		summary.FailedTests = failed.Load()
		summary.PassedTests = passed.Load()
		summary.SkippedTests = skipped.Load()
	})
	for i := range tests {
		test := tests[i]
		name := test.GetName()
		if config.FullName {
			if cwd, err := os.Getwd(); err == nil {
				if abs, err := filepath.Abs(test.BasePath); err == nil {
					if rel, err := filepath.Rel(cwd, abs); err == nil {
						name = fmt.Sprintf("%s[%s]", rel, name)
					} else {
						t.Error(err)
					}
				} else {
					t.Error(err)
				}
			} else {
				t.Error(err)
			}
		}
		t.Run(name, func(t *testing.T) {
			t.Helper()
			runTest(t, ctx, config, test)
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
	}
}

func runTest(t *testing.T, ctx Context, config v1alpha1.ConfigurationSpec, test discovery.Test) {
	t.Helper()
	t.Parallel()
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
		t.Run(fmt.Sprintf("step-%d", i+1), func(t *testing.T) {
			stepCtx, cancel := context.WithTimeout(context.Background(), config.Timeout.Duration)
			defer cancel()
			done := make(chan bool)
			go func() {
				executeStep(t, logging.NewStepLogger(t, fmt.Sprintf("step-%d", i+1)), ctx, test.BasePath, step)
				done <- true
			}()
			select {
			case <-done:
			case <-stepCtx.Done():
				t.Fatalf("Step %d timed out", i+1)
			}
		})
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
