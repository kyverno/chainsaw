package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha2"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/engine/namespacer"
	"github.com/kyverno/chainsaw/pkg/logging"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
	"github.com/kyverno/chainsaw/pkg/runner/internal"
	"github.com/kyverno/chainsaw/pkg/runner/names"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/pkg/ext/output/color"
	"k8s.io/utils/clock"
)

type Runner interface {
	Run(context.Context, v1alpha2.NamespaceOptions, enginecontext.TestContext, ...discovery.Test) error
}

func New(clock clock.PassiveClock, onFailure func()) Runner {
	return &runner{
		clock:     clock,
		onFailure: onFailure,
	}
}

type runner struct {
	clock     clock.PassiveClock
	onFailure func()
	deps      *internal.TestDeps
}

func (r *runner) Run(ctx context.Context, nsOptions v1alpha2.NamespaceOptions, tc enginecontext.TestContext, tests ...discovery.Test) error {
	return r.run(ctx, nil, nsOptions, tc, tests...)
}

func (r *runner) run(ctx context.Context, m mainstart, nsOptions v1alpha2.NamespaceOptions, tc enginecontext.TestContext, tests ...discovery.Test) error {
	defer func() {
		tc.Report.EndTime = time.Now()
	}()
	// sanity check
	if len(tests) == 0 {
		return nil
	}
	internalTests := []testing.InternalTest{{
		Name: "chainsaw",
		F: func(t *testing.T) {
			t.Helper()
			t.Parallel()
			// configure golang context
			ctx = logging.WithSink(ctx, newSink(r.clock, t.Log))
			ctx = logging.WithLogger(ctx, logging.NewLogger(t.Name(), "@chainsaw"))
			// setup cleaner
			cleaner := setupCleanup(ctx, t, r.onFail, tc)
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
					t.Fail()
					tc.IncFailed()
					logging.Log(ctx, logging.Internal, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
					r.onFail()
					return
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
					t.Fail()
					tc.IncFailed()
					logging.Log(ctx, logging.Internal, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
					r.onFail()
				} else {
					testId := i + 1
					if len(test.Test.Spec.Scenarios) == 0 {
						t.Run(name, func(t *testing.T) {
							t.Helper()
							ctx = logging.WithSink(ctx, newSink(r.clock, t.Log))
							r.runTest(ctx, t, nsOptions, nspacer, tc, test, testId, 0)
						})
					} else {
						for s := range test.Test.Spec.Scenarios {
							scenarioId := s + 1
							t.Run(name, func(t *testing.T) {
								t.Helper()
								ctx = logging.WithSink(ctx, newSink(r.clock, t.Log))
								r.runTest(ctx, t, nsOptions, nspacer, tc, test, testId, scenarioId, test.Test.Spec.Scenarios[s].Bindings...)
							})
						}
					}
				}
			}
		},
	}}
	deps := r.deps
	if deps == nil {
		deps = &internal.TestDeps{}
	}
	if m == nil {
		m = testing.MainStart(deps, internalTests, nil, nil, nil)
	}
	// m.Run() returns:
	// - 0 if everything went well
	// - 1 if some of the tests failed
	// - 2 if running the tests was not possible
	// In our case, we consider an error only when running the tests was not possible.
	// For now, the case where some of the tests failed will be covered by the summary.
	if code := m.Run(); code > 1 {
		return fmt.Errorf("testing framework exited with non zero code %d", code)
	}
	return nil
}

func (r *runner) onFail() {
	if r.onFailure != nil {
		r.onFailure()
	}
}
