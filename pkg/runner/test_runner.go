package runner

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/summary"
	"github.com/kyverno/chainsaw/pkg/runner/testing"
	"github.com/kyverno/kyverno/ext/output/color"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/utils/clock"
)

type testRunner struct {
	config  v1alpha1.ConfigurationSpec
	client  client.Client
	clock   clock.PassiveClock
	summary *summary.Summary
}

func (r *testRunner) runTest(goctx context.Context, nspacer namespacer.Namespacer, test discovery.Test) {
	t := testing.FromContext(goctx)
	t.Helper()
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
	cleanLogger := logging.NewLogger(t, r.clock, test.Name, fmt.Sprintf("%-*s", size, "@clean"))
	if test.Spec.Namespace != "" {
		namespace := client.Namespace(test.Spec.Namespace)
		if err := r.client.Get(logging.IntoContext(goctx, beginLogger), client.ObjectKey(&namespace), namespace.DeepCopy()); err != nil {
			if !errors.IsNotFound(err) {
				// Get doesn't log
				beginLogger.Log("GET   ", color.BoldRed, err)
				t.FailNow()
			}
			if err := r.client.Create(logging.IntoContext(goctx, beginLogger), namespace.DeepCopy()); err != nil {
				t.FailNow()
			}
			t.Cleanup(func() {
				// TODO: wait
				if err := r.client.Delete(logging.IntoContext(goctx, cleanLogger), &namespace); err != nil {
					t.FailNow()
				}
			})
		}
		nspacer = namespacer.New(r.client, test.Spec.Namespace)
	}
	if nspacer == nil {
		namespace := client.PetNamespace()
		if err := r.client.Create(logging.IntoContext(goctx, beginLogger), namespace.DeepCopy()); err != nil {
			t.FailNow()
		}
		t.Cleanup(func() {
			// TODO: wait
			if err := r.client.Delete(logging.IntoContext(goctx, cleanLogger), &namespace); err != nil {
				t.FailNow()
			}
		})
		nspacer = namespacer.New(r.client, namespace.Name)
	}
	for i, step := range test.Spec.Steps {
		name := step.Name
		if name == "" {
			name = fmt.Sprintf("step-%d", i+1)
		}
		goctx := logging.IntoContext(goctx, logging.NewLogger(t, r.clock, test.Name, fmt.Sprintf("%-*s", size, name)))
		runner := stepRunner{
			config: r.config,
			client: r.client,
			clock:  r.clock,
		}
		runner.runStep(goctx, nspacer, test, step)
	}
}
