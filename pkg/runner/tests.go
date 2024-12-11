package runner

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha2"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/engine/namespacer"
	"github.com/kyverno/chainsaw/pkg/logging"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
	"github.com/kyverno/chainsaw/pkg/runner/names"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/pkg/ext/output/color"
)

func (r *runner) runTests(ctx context.Context, t testing.TTest, nsOptions v1alpha2.NamespaceOptions, tc enginecontext.TestContext, tests ...discovery.Test) {
	// configure golang context
	ctx = logging.WithSink(ctx, newSink(r.clock, t))
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
					ctx = logging.WithSink(ctx, newSink(r.clock, t))
					r.runTest(ctx, t, nsOptions, nspacer, tc, test, testId, 0)
				})
			} else {
				for s := range test.Test.Spec.Scenarios {
					scenarioId := s + 1
					t.Run(name, func(t *testing.T) {
						t.Helper()
						ctx = logging.WithSink(ctx, newSink(r.clock, t))
						r.runTest(ctx, t, nsOptions, nspacer, tc, test, testId, scenarioId, test.Test.Spec.Scenarios[s].Bindings...)
					})
				}
			}
		}
	}
}
