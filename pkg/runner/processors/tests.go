package processors

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha2"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/engine"
	"github.com/kyverno/chainsaw/pkg/engine/logging"
	"github.com/kyverno/chainsaw/pkg/engine/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/failer"
	"github.com/kyverno/chainsaw/pkg/runner/names"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/pkg/ext/output/color"
	"k8s.io/utils/clock"
)

func RunTests(ctx context.Context, clock clock.PassiveClock, nsOptions v1alpha2.NamespaceOptions, tc engine.Context, tests ...discovery.Test) {
	t := testing.FromContext(ctx)
	// setup cleaner
	cleaner := setupCleanup(ctx, tc)
	// setup namespace
	var nspacer namespacer.Namespacer
	if nsOptions.Name != "" {
		compilers := tc.Compilers()
		if nsOptions.Compiler != nil {
			compilers = compilers.WithDefaultCompiler(string(*nsOptions.Compiler))
		}
		namespaceData := namespaceData{
			cleaner:   cleaner,
			compilers: compilers,
			name:      nsOptions.Name,
			template:  nsOptions.Template,
		}
		nsTc, namespace, err := setupNamespace(ctx, tc, namespaceData)
		if err != nil {
			logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
			tc.IncFailed()
			failer.FailNow(ctx)
		}
		tc = nsTc
		if namespace != nil {
			nspacer = namespacer.New(namespace.GetName())
		}
	}
	// loop through tests
	for i := range tests {
		test := tests[i]
		name, err := names.Test(tc.FullName(), test)
		if err != nil {
			logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
			tc.IncFailed()
			failer.FailNow(ctx)
		}
		// compute test scenarios
		scenarios := []v1alpha1.Scenario{{}}
		if test.Test == nil {
			scenarios = nil
		} else if len(test.Test.Spec.Scenarios) != 0 {
			scenarios = test.Test.Spec.Scenarios
		}
		// loop through test scenarios
		for s := range scenarios {
			// run each test scenario in a separate T
			t.Run(name, func(t *testing.T) {
				t.Helper()
				tc, err := engine.WithBindings(ctx, tc, scenarios[s].Bindings...)
				if err != nil {
					logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
					tc.IncFailed()
					failer.FailNow(ctx)
				}
				ctx := testing.IntoContext(ctx, t)
				size := len("@chainsaw")
				for i, step := range test.Test.Spec.Steps {
					name := step.Name
					if name == "" {
						name = fmt.Sprintf("step-%d", i+1)
					}
					if size < len(name) {
						size = len(name)
					}
				}
				ctx = logging.IntoContext(ctx, logging.NewLogger(t, clock, test.Test.Name, fmt.Sprintf("%-*s", size, "@chainsaw")))
				info := TestInfo{
					Id:         i + 1,
					ScenarioId: s + 1,
					Metadata:   test.Test.ObjectMeta,
				}
				tc = tc.WithBinding(ctx, "test", info)
				t.Cleanup(func() {
					if t.Skipped() {
						tc.IncSkipped()
					} else {
						if t.Failed() {
							tc.IncFailed()
						} else {
							tc.IncPassed()
						}
					}
				})
				// TODO: move into each test processor
				if test.Test.Spec.Concurrent == nil || *test.Test.Spec.Concurrent {
					t.Parallel()
				}
				if test.Test.Spec.Skip != nil && *test.Test.Spec.Skip {
					t.SkipNow()
				}
				if test.Test.Spec.FailFast != nil {
					tc = tc.WithFailFast(ctx, *test.Test.Spec.FailFast)
				}
				if tc.FailFast() && tc.Failed() > 0 {
					t.SkipNow()
				}
				processor := createTestProcessor(clock, nsOptions, test, size)
				processor.Run(ctx, nspacer, tc)
			})
		}
	}
}

func createTestProcessor(clock clock.PassiveClock, nsOptions v1alpha2.NamespaceOptions, test discovery.Test, size int) TestProcessor {
	return NewTestProcessor(
		test,
		size,
		clock,
		nsOptions.Template,
		nsOptions.Compiler,
	)
}
