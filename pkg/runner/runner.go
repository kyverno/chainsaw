package runner

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

type Summary struct {
	PassedTest *int
	FailedTest *int
}

func Run(cfg *rest.Config, config v1alpha1.ConfigurationSpec, summary *Summary, tests ...discovery.Test) (int, error) {
	if len(tests) == 0 {
		return 0, nil
	}
	testing.Init()
	if err := flag.Set("test.v", "true"); err != nil {
		return 0, err
	}
	if err := flag.Set("test.parallel", strconv.Itoa(config.Parallel)); err != nil {
		return 0, err
	}
	if err := flag.Set("test.timeout", config.Timeout.Duration.String()); err != nil {
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
	if err := flag.Set("test.count", "1"); err != nil {
		return 0, err
	}
	if err := flag.Set("test.run", config.IncludeTestRegex); err != nil {
		return 0, err
	}
	if err := flag.Set("test.skip", config.ExcludeTestRegex); err != nil {
		return 0, err
	}
	flag.Parse()
	run := func(t *testing.T) {
		t.Helper()
		run(t, cfg, config, summary, tests...)
	}
	internalTest := []testing.InternalTest{{
		Name: "chainsaw",
		F:    run,
	}}
	var testDeps testDeps
	m := testing.MainStart(&testDeps, internalTest, nil, nil, nil)
	return m.Run(), nil
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
		summary: summary,
	}
	if config.Namespace != "" {
		namespace := client.Namespace(config.Namespace)
		c := ctx.clientFactory(t, logging.NewTestLogger(t))
		if err := c.Get(context.Background(), client.ObjectKey(&namespace), nil); err != nil {
			if errors.IsNotFound(err) {
				if err := c.Create(context.Background(), namespace.DeepCopy()); err != nil {
					t.Fatal(err)
				}
				ctx.namespacer = namespacer.New(c, config.Namespace)
			}
		}
	}
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
			runTest(t, ctx, test)
		})
	}
}

func runTest(t *testing.T, ctx Context, test discovery.Test) {
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
		executeStep(t, logging.NewStepLogger(t, fmt.Sprintf("step-%d", i+1)), ctx, test.BasePath, step)
		if t.Failed() {
			*ctx.summary.FailedTest++
		} else {
			*ctx.summary.PassedTest++
		}
	}
}

func executeStep(t *testing.T, logger logging.Logger, ctx Context, basePath string, step v1alpha1.TestStepSpec) {
	t.Helper()
	c := ctx.clientFactory(t, logger)
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
}
