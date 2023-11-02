package runner

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	runnerclient "github.com/kyverno/chainsaw/pkg/runner/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/names"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/summary"
	"github.com/kyverno/chainsaw/pkg/runner/testing"
	"github.com/kyverno/kyverno/ext/output/color"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/utils/clock"
)

type testsRunner struct {
	config  v1alpha1.ConfigurationSpec
	client  client.Client
	clock   clock.PassiveClock
	summary *summary.Summary
}

func (r *testsRunner) runTests(goctx context.Context, tests ...discovery.Test) {
	t := testing.FromContext(goctx)
	t.Helper()
	ctx := Context{
		clock:  r.clock,
		client: runnerclient.New(r.client),
	}
	if r.config.Namespace != "" {
		namespace := client.Namespace(r.config.Namespace)
		if err := ctx.client.Get(goctx, client.ObjectKey(&namespace), namespace.DeepCopy()); err != nil {
			if !errors.IsNotFound(err) {
				// Get doesn't log
				logging.Log(goctx, "GET   ", color.BoldRed, err)
				t.FailNow()
			}
			// TODO
			t.Cleanup(func() {
				// if err := operations.Delete(goctx, mainLogger, &namespace, c); err != nil {
				// 	t.FailNow()
				// }
			})
			if err := ctx.client.Create(goctx, namespace.DeepCopy()); err != nil {
				t.FailNow()
			}
		}
		ctx.namespacer = namespacer.New(ctx.client, r.config.Namespace)
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
		name, err := names.Test(r.config, test)
		if err != nil {
			logging.Log(goctx, "INTERN", color.BoldRed, err)
			t.FailNow()
		}
		t.Run(name, func(t *testing.T) {
			t.Helper()
			goctx := testing.IntoContext(goctx, t)
			if r.summary != nil {
				t.Cleanup(func() {
					if t.Skipped() {
						r.summary.IncSkipped()
					} else {
						if t.Failed() {
							r.summary.IncFailed()
						} else {
							r.summary.IncPassed()
						}
					}
				})
			}
			if test.Spec.Concurrent == nil || *test.Spec.Concurrent {
				t.Parallel()
			}
			if test.Spec.Skip != nil && *test.Spec.Skip {
				t.SkipNow()
			}
			beginLogger := logging.NewLogger(t, r.clock, test.Name, fmt.Sprintf("%-*s", size, "@begin"))
			// cleanLogger := logging.NewLogger(t, clock, test.Name, fmt.Sprintf("%-*s", size, "@clean"))
			ctx := ctx
			if test.Spec.Namespace != "" {
				namespace := client.Namespace(test.Spec.Namespace)
				if err := ctx.client.Get(logging.IntoContext(goctx, beginLogger), client.ObjectKey(&namespace), namespace.DeepCopy()); err != nil {
					if !errors.IsNotFound(err) {
						// Get doesn't log
						beginLogger.Log("GET   ", color.BoldRed, err)
						t.FailNow()
					}
					if err := ctx.client.Create(logging.IntoContext(goctx, beginLogger), namespace.DeepCopy()); err != nil {
						t.FailNow()
					}
					// TODO
					t.Cleanup(func() {
						// if err := operations.Delete(logging.IntoContext(goctx, cleanLogger), cleanLogger, &namespace, ctx.clientFactory(cleanLogger)); err != nil {
						// 	t.FailNow()
						// }
					})
				}
				ctx.namespacer = namespacer.New(ctx.client, test.Spec.Namespace)
			}
			if ctx.namespacer == nil {
				namespace := client.PetNamespace()
				if err := ctx.client.Create(logging.IntoContext(goctx, beginLogger), namespace.DeepCopy()); err != nil {
					t.FailNow()
				}
				t.Cleanup(func() {
					// TODO
					// if err := operations.Delete(context.Background(), cleanLogger, &namespace, ctx.clientFactory(cleanLogger)); err != nil {
					// 	t.FailNow()
					// }
				})
				ctx.namespacer = namespacer.New(r.client, namespace.Name)
			}
			runTest(goctx, ctx, r.config, test, size)
		})
	}
}
