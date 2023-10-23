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
	var internalTests []testing.InternalTest
	var failed, passed, skipped atomic.Int32
	defer func() {
		summary.FailedTests = failed.Load()
		summary.PassedTests = passed.Load()
		summary.SkippedTests = skipped.Load()
	}()
	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	var nspacer namespacer.Namespacer
	if config.Namespace != "" {
		nspacer = namespacer.New(c, config.Namespace)
		namespace := client.Namespace(config.Namespace)
		if err := c.Get(context.Background(), client.ObjectKey(&namespace), namespace.DeepCopy()); err != nil {
			if !errors.IsNotFound(err) {
				return nil, err
			}
			if err := c.Create(context.Background(), namespace.DeepCopy()); err != nil {
				return nil, err
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
			return nil, err
		}
		internalTests = append(internalTests, testing.InternalTest{
			Name: name,
			F: func(t *testing.T) {
				t.Helper()
				t.Parallel()
				ctx := Context{
					clock:      clock,
					config:     config,
					ctx:        context.Background(),
					namespacer: nspacer,
					clientFactory: func(t *testing.T, logger logging.Logger) client.Client {
						t.Helper()
						return runnerclient.New(t, logger, c, !config.SkipDelete)
					},
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
	if code := m.Run(); code > 1 {
		return &summary, fmt.Errorf("testing framework exited with non zero code %d", code)
	}
	return &summary, nil
}

func runTest(t *testing.T, ctx Context, test discovery.Test) {
	t.Helper()
	if ctx.namespacer == nil {
		namespace := client.PetNamespace()
		logger := logging.NewTestLogger(t, ctx.clock)
		c := ctx.clientFactory(t, logger)
		if err := c.Create(context.Background(), namespace.DeepCopy()); err != nil {
			logger.Log(err)
			t.FailNow()
		}
		ctx.namespacer = namespacer.New(c, namespace.Name)
	}
	for i := range test.Spec.Steps {
		step := test.Spec.Steps[i]
		executeStep(t, logging.NewStepLogger(t, ctx.clock, fmt.Sprintf("step-%d", i+1)), ctx, test.BasePath, step)
	}
}
