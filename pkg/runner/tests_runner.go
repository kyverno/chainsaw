package runner

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
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

func (r *testsRunner) runTests(ctx context.Context, tests ...discovery.Test) {
	t := testing.FromContext(ctx)
	t.Helper()
	var nspacer namespacer.Namespacer
	if r.config.Namespace != "" {
		namespace := client.Namespace(r.config.Namespace)
		if err := r.client.Get(ctx, client.ObjectKey(&namespace), namespace.DeepCopy()); err != nil {
			if !errors.IsNotFound(err) {
				// Get doesn't log
				logging.Log(ctx, "GET   ", color.BoldRed, err)
				t.FailNow()
			}
			t.Cleanup(func() {
				// TODO: wait
				if err := r.client.Delete(ctx, &namespace); err != nil {
					t.FailNow()
				}
			})
			if err := r.client.Create(ctx, namespace.DeepCopy()); err != nil {
				t.FailNow()
			}
		}
		nspacer = namespacer.New(r.client, r.config.Namespace)
	}
	for _, test := range tests {
		name, err := names.Test(r.config, test)
		if err != nil {
			logging.Log(ctx, "INTERN", color.BoldRed, err)
			t.FailNow()
		}
		t.Run(name, func(t *testing.T) {
			t.Helper()
			goctx := testing.IntoContext(ctx, t)
			runner := testRunner{
				config:  r.config,
				client:  r.client,
				clock:   r.clock,
				summary: r.summary,
			}
			runner.runTest(goctx, nspacer, test)
		})
	}
}
