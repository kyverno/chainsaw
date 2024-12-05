package processors

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/cleanup/cleaner"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/engine"
	"github.com/kyverno/chainsaw/pkg/engine/logging"
	"github.com/kyverno/chainsaw/pkg/engine/namespacer"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/kyverno/chainsaw/pkg/runner/failer"
	"github.com/kyverno/chainsaw/pkg/runner/names"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/pkg/ext/output/color"
	"k8s.io/utils/clock"
)

type TestsProcessor interface {
	Run(context.Context, engine.Context, ...discovery.Test)
}

func NewTestsProcessor(config model.Configuration, clock clock.PassiveClock) TestsProcessor {
	return &testsProcessor{
		config: config,
		clock:  clock,
	}
}

type testsProcessor struct {
	config model.Configuration
	clock  clock.PassiveClock
}

func (p *testsProcessor) Run(ctx context.Context, tc engine.Context, tests ...discovery.Test) {
	// 1. setup context
	t := testing.FromContext(ctx)
	mainCleaner := cleaner.New(p.config.Timeouts.Cleanup.Duration, nil, p.config.Deletion.Propagation)
	t.Cleanup(func() {
		if !mainCleaner.Empty() {
			logging.Log(ctx, logging.Cleanup, logging.BeginStatus, color.BoldFgCyan)
			defer func() {
				logging.Log(ctx, logging.Cleanup, logging.EndStatus, color.BoldFgCyan)
			}()
			for _, err := range mainCleaner.Run(ctx, nil) {
				logging.Log(ctx, logging.Cleanup, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
				failer.Fail(ctx)
			}
		}
	})
	contextData := contextData{
		basePath: "",
		clusters: p.config.Clusters,
	}
	if p.config.Namespace.Name != "" {
		var nsCleaner cleaner.CleanerCollector
		if !tc.SkipDelete() {
			nsCleaner = mainCleaner
		}
		compilers := tc.Compilers()
		if p.config.Namespace.Compiler != nil {
			compilers = compilers.WithDefaultCompiler(string(*p.config.Namespace.Compiler))
		}
		contextData.namespace = &namespaceData{
			name:      p.config.Namespace.Name,
			template:  p.config.Namespace.Template,
			compilers: compilers,
			cleaner:   nsCleaner,
		}
	}
	tc, namespace, err := setupContextData(ctx, tc, contextData)
	if err != nil {
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		tc.IncFailed()
		failer.FailNow(ctx)
	}
	var nspacer namespacer.Namespacer
	if namespace != nil {
		nspacer = namespacer.New(namespace.GetName())
	}
	// 2. loop through tests
	for i := range tests {
		test := tests[i]
		name, err := names.Test(p.config.Discovery.FullName, test)
		if err != nil {
			logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
			tc.IncFailed()
			failer.FailNow(ctx)
		}
		// 3. compute test scenarios
		scenarios := applyScenarios(test)
		// 4. loop through test scenarios
		for s := range scenarios {
			test := scenarios[s]
			// 5. run each test scenario in a separate T
			t.Run(name, func(t *testing.T) {
				t.Helper()
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
				ctx = logging.IntoContext(ctx, logging.NewLogger(t, p.clock, test.Test.Name, fmt.Sprintf("%-*s", size, "@chainsaw")))
				info := TestInfo{
					Id:         i + 1,
					ScenarioId: s + 1,
					Metadata:   test.Test.ObjectMeta,
				}
				tc := tc.WithBinding(ctx, "test", info)
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
				if test.Test.Spec.Concurrent == nil || *test.Test.Spec.Concurrent {
					t.Parallel()
				}
				if test.Test.Spec.Skip != nil && *test.Test.Spec.Skip {
					t.SkipNow()
				}
				failFast := p.config.Execution.FailFast
				if test.Test.Spec.FailFast != nil {
					failFast = *test.Test.Spec.FailFast
				}
				if failFast {
					if tc.Failed() > 0 {
						t.SkipNow()
					}
				}
				processor := p.createTestProcessor(test, size)
				processor.Run(ctx, nspacer, tc)
			})
		}
	}
}

func (p *testsProcessor) createTestProcessor(test discovery.Test, size int) TestProcessor {
	return NewTestProcessor(
		test,
		size,
		p.clock,
		p.config.Namespace.Template,
		p.config.Namespace.Compiler,
		p.config.Timeouts,
		p.config.Error.Catch...,
	)
}
