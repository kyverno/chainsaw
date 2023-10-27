package runner

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/report"
	runnerclient "github.com/kyverno/chainsaw/pkg/runner/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/rest"
	"k8s.io/utils/clock"
)

func Run(cfg *rest.Config, clock clock.PassiveClock, config v1alpha1.ConfigurationSpec, tests ...discovery.Test) (*Summary, *report.Report, error) {
	var summary Summary
	var rep report.Report
	if len(tests) == 0 {
		return &summary, &rep, nil
	}
	if err := setupFlags(config); err != nil {
		return nil, &rep, err
	}
	c, err := client.New(cfg)
	if err != nil {
		return nil, &rep, err
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
			mainLogger := logging.NewLogger(t, clock)
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
						mainLogger.Log(err)
						t.FailNow()
					}
					t.Cleanup(func() {
						if err := operations.Delete(context.Background(), mainLogger, &namespace, c); err != nil {
							mainLogger.Log(err)
							t.FailNow()
						}
					})
					if err := c.Create(context.Background(), namespace.DeepCopy()); err != nil {
						mainLogger.Log(err)
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
					mainLogger.Log(err)
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
					testLogger := logging.NewLogger(t, clock, test.Name)
					beginLogger := testLogger.WithName("@begin")
					cleanLogger := testLogger.WithName("@clean")

					ctx := ctx
					if test.Spec.Namespace != "" {
						namespace := client.Namespace(test.Spec.Namespace)
						c := ctx.clientFactory(beginLogger)
						if err := c.Get(context.Background(), client.ObjectKey(&namespace), namespace.DeepCopy()); err != nil {
							if !errors.IsNotFound(err) {
								beginLogger.Log(err)
								t.FailNow()
							}
							if err := c.Create(context.Background(), namespace.DeepCopy()); err != nil {
								beginLogger.Log(err)
								t.FailNow()
							}
							t.Cleanup(func() {
								if err := operations.Delete(context.Background(), cleanLogger, &namespace, c); err != nil {
									cleanLogger.Log(err)
									t.FailNow()
								}
							})
						}
						ctx.namespacer = namespacer.New(c, test.Spec.Namespace)
					}
					if ctx.namespacer == nil {
						namespace := client.PetNamespace()
						c := ctx.clientFactory(beginLogger)
						if err := c.Create(context.Background(), namespace.DeepCopy()); err != nil {
							beginLogger.Log(err)
							t.FailNow()
						}
						t.Cleanup(func() {
							if err := operations.Delete(context.Background(), cleanLogger, &namespace, c); err != nil {
								cleanLogger.Log(err)
								t.FailNow()
							}
						})
						ctx.namespacer = namespacer.New(c, namespace.Name)
					}
					// Initialize a new TestReport for this test
					testReport := &report.TestReport{
						Name:      test.Name,
						StartTime: time.Now(),
					}
					runTest(t, testLogger, ctx, config, test, testReport)
					rep.Tests = append(rep.Tests, *testReport)
				})
			}
		},
	}}
	m := testing.MainStart(&testDeps{}, internalTests, nil, nil, nil)
	if code := m.Run(); code > 1 {
		return &summary, &rep, fmt.Errorf("testing framework exited with non zero code %d", code)
	}
	return &summary, &rep, nil
}

func runTest(t *testing.T, logger logging.Logger, ctx Context, config v1alpha1.ConfigurationSpec, test discovery.Test, testReport *report.TestReport) {
	t.Helper()
	for i, step := range test.Spec.Steps {
		name := step.Name
		if name == "" {
			name = fmt.Sprintf("step-%d", i+1)
		}
		stepReport := executeStep(t, logger.WithName(name), ctx, test.BasePath, config, test.Spec, step)
		testReport.Steps = append(testReport.Steps, *stepReport)
	}
	testReport.EndTime = time.Now()
}
