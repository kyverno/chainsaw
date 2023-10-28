package runner

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	runnerclient "github.com/kyverno/chainsaw/pkg/runner/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/kyverno/ext/output/color"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/rest"
	"k8s.io/utils/clock"
)

func Run(cfg *rest.Config, clock clock.PassiveClock, config v1alpha1.ConfigurationSpec, tests ...discovery.Test) (*Summary, error) {
	var summary Summary
	if len(tests) == 0 {
		return &summary, nil
	}
	if err := setupFlags(config); err != nil {
		return nil, err
	}
	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	var failed, passed, skipped atomic.Int32
	defer func() {
		summary.FailedTests = failed.Load()
		summary.PassedTests = passed.Load()
		summary.SkippedTests = skipped.Load()
	}()
	internalTests := []testing.InternalTest{{
		Name: "chainsaw",
		F: func(t *testing.T) {
			t.Helper()
			mainLogger := logging.NewLogger(t, clock, t.Name(), "@init")
			ctx := Context{
				clock: clock,
				clientFactory: func(logger logging.Logger) client.Client {
					return runnerclient.New(logger, c)
				},
			}
			if config.Namespace != "" {
				c := ctx.clientFactory(mainLogger)
				namespace := client.Namespace(config.Namespace)
				if err := c.Get(context.Background(), client.ObjectKey(&namespace), namespace.DeepCopy()); err != nil {
					if !errors.IsNotFound(err) {
						// Get doesn't log
						mainLogger.Log("GET   ", color.BoldRed, err)
						t.FailNow()
					}
					t.Cleanup(func() {
						if err := operations.Delete(context.Background(), mainLogger, &namespace, c); err != nil {
							t.FailNow()
						}
					})
					if err := c.Create(context.Background(), namespace.DeepCopy()); err != nil {
						t.FailNow()
					}
				}
				ctx.namespacer = namespacer.New(c, config.Namespace)
			}
			// TODO: shall we precompute subtest names ?
			// t.Cleanup(func() {
			// 	if t.Skipped() {
			// 		skipped.Add(1)
			// 	} else {
			// 		if t.Failed() {
			// 			failed.Add(1)
			// 		} else {
			// 			passed.Add(1)
			// 		}
			// 	}
			// })
			for i := range tests {
				test := tests[i]
				size := 0
				for i, step := range test.Spec.Steps {
					name := step.Name
					if name == "" {
						name = fmt.Sprintf("step-%d", i+1)
					}
					if size < len(name) {
						size = len(name)
					}
				}
				name, err := testName(config, test)
				if err != nil {
					mainLogger.Log("INTERN", color.BoldRed, err)
					t.FailNow()
				}
				t.Run(name, func(t *testing.T) {
					t.Helper()
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
					if test.Spec.Concurrent == nil || *test.Spec.Concurrent {
						t.Parallel()
					}
					if test.Spec.Skip != nil && *test.Spec.Skip {
						t.SkipNow()
					}
					beginLogger := logging.NewLogger(t, clock, test.Name, fmt.Sprintf("%-*s", size, "@begin"))
					cleanLogger := logging.NewLogger(t, clock, test.Name, fmt.Sprintf("%-*s", size, "@clean"))
					ctx := ctx
					if test.Spec.Namespace != "" {
						namespace := client.Namespace(test.Spec.Namespace)
						c := ctx.clientFactory(beginLogger)
						if err := c.Get(context.Background(), client.ObjectKey(&namespace), namespace.DeepCopy()); err != nil {
							if !errors.IsNotFound(err) {
								// Get doesn't log
								beginLogger.Log("GET   ", color.BoldRed, err)
								t.FailNow()
							}
							if err := c.Create(context.Background(), namespace.DeepCopy()); err != nil {
								t.FailNow()
							}
							t.Cleanup(func() {
								if err := operations.Delete(context.Background(), cleanLogger, &namespace, ctx.clientFactory(cleanLogger)); err != nil {
									t.FailNow()
								}
							})
						}
						ctx.namespacer = namespacer.New(c, test.Spec.Namespace)
					}
					if ctx.namespacer == nil {
						namespace := client.PetNamespace()
						if err := ctx.clientFactory(beginLogger).Create(context.Background(), namespace.DeepCopy()); err != nil {
							t.FailNow()
						}
						t.Cleanup(func() {
							if err := operations.Delete(context.Background(), cleanLogger, &namespace, ctx.clientFactory(cleanLogger)); err != nil {
								t.FailNow()
							}
						})
						ctx.namespacer = namespacer.New(c, namespace.Name)
					}
					runTest(t, ctx, config, test, size)
				})
			}
		},
	}}
	m := testing.MainStart(&testDeps{}, internalTests, nil, nil, nil)
	if code := m.Run(); code > 1 {
		return &summary, fmt.Errorf("testing framework exited with non zero code %d", code)
	}
	return &summary, nil
}

func runTest(t *testing.T, ctx Context, config v1alpha1.ConfigurationSpec, test discovery.Test, size int) {
	t.Helper()
	for i, step := range test.Spec.Steps {
		name := step.Name
		if name == "" {
			name = fmt.Sprintf("step-%d", i+1)
		}
		logger := logging.NewLogger(t, ctx.clock, test.Name, fmt.Sprintf("%-*s", size, name))
		executeStep(t, logger, ctx, test.BasePath, config, test.Spec, step)
	}
}
