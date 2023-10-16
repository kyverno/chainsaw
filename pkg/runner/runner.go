package runner

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/resource"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/rest"
)

func Run(cfg *rest.Config, config v1alpha1.ConfigurationSpec, tests ...discovery.Test) (int, error) {
	if len(tests) == 0 {
		return 0, nil
	}
	testing.Init()
	flag.Parse()
	// Set the verbose test flag to true since we are not using the regular go test CLI.
	if err := flag.Set("test.v", "true"); err != nil {
		return 0, err
	}
	if err := flag.Set("test.parallel", strconv.Itoa(config.Parallel)); err != nil {
		return 0, err
	}
	if err := flag.Set("test.timeout", config.Timeout.String()); err != nil {
		return 0, err
	}
	if err := flag.Set("test.failfast", fmt.Sprint(config.FailFast)); err != nil {
		return 0, err
	}
	if err := flag.Set("test.paniconexit0", "true"); err != nil {
		return 0, err
	}
	if err := flag.Set("test.fullpath", "false"); err != nil {
		return 0, err
	}
	// regex related flags
	var testDeps testDeps
	m := testing.MainStart(
		&testDeps,
		[]testing.InternalTest{{
			Name: "chainsaw",
			F: func(t *testing.T) {
				t.Helper()
				run(t, cfg, config, tests...)
			},
		}},
		nil,
		nil,
		nil,
	)
	return m.Run(), nil
}

func run(t *testing.T, cfg *rest.Config, config v1alpha1.ConfigurationSpec, tests ...discovery.Test) {
	t.Helper()
	c, err := client.New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	ctx := Context{
		client: c,
	}
	if config.Namespace != "" {
		namespace := client.Namespace(config.Namespace)
		if err := c.Get(context.Background(), client.ObjectKey(&namespace), nil); err != nil {
			if errors.IsNotFound(err) {
				if err := c.Create(context.Background(), namespace.DeepCopy()); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					t.Logf("cleanup namespace: %s", config.Namespace)
					if err := client.BlockingDelete(context.Background(), c, &namespace); err != nil {
						t.Fatal(err)
					}
				})
				ctx.namespacer = namespacer.New(ctx.client, config.Namespace)
			}
		}
	}
	for i := range tests {
		test := tests[i]
		t.Run(test.GetName(), func(t *testing.T) {
			t.Helper()
			runTest(t, ctx, test)
		})
	}
}

func runTest(t *testing.T, ctx Context, test discovery.Test) {
	t.Helper()
	t.Parallel()
	if ctx.namespacer == nil {
		namespace := client.PetNamespace()
		if err := ctx.client.Create(context.Background(), namespace.DeepCopy()); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			t.Logf("cleanup namespace: %s", namespace.Name)
			if err := client.BlockingDelete(context.Background(), ctx.client, &namespace); err != nil {
				t.Fatal(err)
			}
		})
		ctx.namespacer = namespacer.New(ctx.client, namespace.Name)
	}
	for i := range test.Spec.Steps {
		step := test.Spec.Steps[i]
		t.Logf("step-%d", i+1)
		executeStep(t, ctx, test.BasePath, step)
	}
}

func executeStep(t *testing.T, ctx Context, basePath string, step v1alpha1.TestStepSpec) {
	t.Helper()
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
			t.Logf("apply %s (%s/%s)", client.ObjectKey(resource), resource.GetAPIVersion(), resource.GetKind())
			cleanup, err := client.CreateOrUpdate(context.Background(), ctx.client, resource)
			if err != nil {
				t.Fatal(err)
			}
			if cleanup {
				t.Cleanup(func() {
					t.Logf("cleanup resource: %s (%s/%s)", client.ObjectKey(resource), resource.GetAPIVersion(), resource.GetKind())
					if err := client.BlockingDelete(context.Background(), ctx.client, resource); err != nil {
						t.Fatal(err)
					}
				})
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
			t.Logf("assert %s (%s/%s)", client.ObjectKey(resource), resource.GetAPIVersion(), resource.GetKind())
			err := client.Assert(context.Background(), resources[i], ctx.client)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}
