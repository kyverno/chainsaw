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
			logger := logging.NewOperationLogger(t, clock)
			ctx := Context{
				clock: clock,
				clientFactory: func(logger logging.Logger) client.Client {
					return runnerclient.New(logger, c)
				},
			}
			if config.Namespace != "" {
				c := ctx.clientFactory(logger)
				namespace := client.Namespace(config.Namespace)
				if err := c.Get(context.Background(), client.ObjectKey(&namespace), namespace.DeepCopy()); err != nil {
					if !errors.IsNotFound(err) {
						logger.Log(err)
						t.FailNow()
					}
					t.Cleanup(func() {
						if err := operations.Delete(context.Background(), logger, &namespace, c); err != nil {
							logger.Log(err)
							t.FailNow()
						}
					})
					if err := c.Create(context.Background(), namespace.DeepCopy()); err != nil {
						logger.Log(err)
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
				name, err := testName(config, test)
				if err != nil {
					logger.Log(err)
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
					logger := logging.NewOperationLogger(t, clock, test.Name)
					ctx := ctx
					if ctx.namespacer == nil {
						namespace := client.PetNamespace()
						c := ctx.clientFactory(logger)
						if err := c.Create(context.Background(), namespace.DeepCopy()); err != nil {
							logger.Log(err)
							t.FailNow()
						}
						t.Cleanup(func() {
							if err := operations.Delete(context.Background(), logger, &namespace, c); err != nil {
								logger.Log(err)
								t.FailNow()
							}
						})
						ctx.namespacer = namespacer.New(c, namespace.Name)
					}
					runTest(t, ctx, config, test)
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

func runTest(t *testing.T, ctx Context, config v1alpha1.ConfigurationSpec, test discovery.Test) {
	t.Helper()
	for i, step := range test.Spec.Steps {
		name := step.Name
		if name == "" {
			name = fmt.Sprintf("step-%d", i+1)
		}
		logger := logging.NewOperationLogger(t, ctx.clock, test.Name, name)
		executeStep(t, logger, ctx, test.BasePath, config, test.Spec, step)
	}
}
