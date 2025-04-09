package runner

import (
	"context"
	"fmt"
	"testing"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha2"
	"github.com/kyverno/chainsaw/pkg/cleanup/cleaner"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/engine/namespacer"
	"github.com/kyverno/chainsaw/pkg/logging"
	"github.com/kyverno/chainsaw/pkg/model"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
	"github.com/kyverno/chainsaw/pkg/runner/internal"
	"github.com/kyverno/chainsaw/pkg/runner/names"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/pkg/ext/output/color"
	"go.uber.org/multierr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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
			fail := func(t *testing.T, err error) bool {
				t.Helper()
				if err != nil {
					t.Fail()
					return true
				}
				return false
			}
			// setup logger sink
			ctx = logging.WithSink(ctx, newSink(r.clock, t.Log))
			// setup logger
			ctx = logging.WithLogger(ctx, logging.NewLogger(t.Name(), "@chainsaw"))
			// setup cleanup
			cleanup := cleaner.New(tc.Timeouts().Cleanup.Duration, nil, tc.DeletionPropagation())
			t.Cleanup(func() {
				fail(t, r.cleanup(ctx, tc, cleanup))
			})
			// setup namespace
			tc, err := r.setupNamespace(ctx, nsOptions, tc, cleanup)
			if fail(t, err) {
				tc.IncFailed()
				return
			}
			// loop through tests
			for i := range tests {
				test := tests[i]
				name, err := names.Test(test, tc.FullName())
				if err != nil {
					t.Fail()
					tc.IncFailed()
					logging.Log(ctx, logging.Internal, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
					r.onFail()
				} else {
					// setup logger
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
					ctx := logging.WithLogger(ctx, logging.NewLogger(test.Test.Name, fmt.Sprintf("%-*s", size, "@chainsaw")))
					// helper to run test
					runTest := func(ctx context.Context, t *testing.T, testId int, scenarioId int, tc enginecontext.TestContext, bindings ...v1alpha1.Binding) {
						t.Helper()
						// setup logger sink
						ctx = logging.WithSink(ctx, newSink(r.clock, t.Log))
						// setup concurrency
						if test.Test.Spec.Concurrent == nil || *test.Test.Spec.Concurrent {
							t.Parallel()
						}
						// setup summary
						defer func() {
							if t.Skipped() {
								tc.IncSkipped()
							} else if t.Failed() {
								tc.IncFailed()
							} else {
								tc.IncPassed()
							}
						}()
						// setup reporting
						report := &model.TestReport{
							BasePath:   test.BasePath,
							Name:       test.Test.Name,
							Concurrent: test.Test.Spec.Concurrent,
							StartTime:  time.Now(),
						}
						defer func() {
							report.EndTime = time.Now()
							report.Skipped = t.Skipped()
							tc.Report.Add(report)
						}()
						// skip check
						if test.Test.Spec.Skip != nil && *test.Test.Spec.Skip {
							t.Skip()
							return
						}
						// setup context
						tc, err := r.setupTestContext(ctx, testId, scenarioId, tc, test, bindings...)
						if fail(t, err) {
							return
						}
						// fail fast check
						if tc.FailFast() && tc.Failed() > 0 {
							t.Skip()
							return
						}
						// setup cleaner
						cleanup := cleaner.New(tc.Timeouts().Cleanup.Duration, nil, tc.DeletionPropagation())
						defer func() {
							fail(t, r.testCleanup(ctx, tc, cleanup, report))
						}()
						// setup namespace
						// TODO: should be part of setupContext ?
						if test.Test.Spec.Compiler != nil {
							tc = tc.WithDefaultCompiler(string(*test.Test.Spec.Compiler))
						}
						nsOptions := nsOptions
						nsOptions.Name = test.Test.Spec.Namespace
						if tc.Namespacer() == nil && nsOptions.Name == "" {
							nsOptions.Name = fmt.Sprintf("chainsaw-%s", petname.Generate(2, "-"))
						}
						if template := test.Test.Spec.NamespaceTemplate; template != nil && template.Value() != nil {
							nsOptions.Template = template
							nsOptions.Compiler = test.Test.Spec.NamespaceTemplateCompiler
						}
						tc, err = r.setupNamespace(ctx, nsOptions, tc, cleanup)
						if fail(t, err) {
							return
						}
						// setup bindings
						tc, err = enginecontext.SetupBindings(tc, test.Test.Spec.Bindings...)
						if err != nil {
							t.Fail()
							logging.Log(ctx, logging.Internal, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
							r.onFail()
							return
						}
						// loop through steps
						for i, step := range test.Test.Spec.Steps {
							ctx := logging.WithLogger(ctx, logging.NewLogger(test.Test.Name, fmt.Sprintf("%-*s", size, names.Step(step, i))))
							info := StepInfo{
								Id: i + 1,
							}
							tc := tc.WithBinding("step", info)
							if stop := r.runStep(ctx, t.Cleanup, t.Fail, t.Failed, tc, step, report); stop {
								return
							}
						}
					}
					// run test scenarios
					testId := i + 1
					if len(test.Test.Spec.Scenarios) == 0 {
						t.Run(name, func(t *testing.T) {
							t.Helper()
							runTest(ctx, t, testId, 0, tc)
						})
					} else {
						for s := range test.Test.Spec.Scenarios {
							scenarioId := s + 1
							t.Run(name, func(t *testing.T) {
								t.Helper()
								runTest(ctx, t, testId, scenarioId, tc, test.Test.Spec.Scenarios[s].Bindings...)
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

func (r *runner) runStep(
	ctx context.Context,
	cleanup func(func()),
	fail func(),
	failed func() bool,
	tc enginecontext.TestContext,
	step v1alpha1.TestStep,
	testReport *model.TestReport,
) bool {
	report := &model.StepReport{
		Name:      step.Name,
		StartTime: time.Now(),
	}
	defer func() {
		report.EndTime = time.Now()
		testReport.Add(report)
	}()
	if step.Compiler != nil {
		tc = tc.WithDefaultCompiler(string(*step.Compiler))
	}
	contextData := enginecontext.ContextData{
		Catch:               step.Catch,
		Cluster:             step.Cluster,
		Clusters:            step.Clusters,
		DeletionPropagation: step.DeletionPropagationPolicy,
		SkipDelete:          step.SkipDelete,
		Templating:          step.Template,
		Timeouts:            step.Timeouts,
	}
	tc, err := enginecontext.SetupContextAndBindings(tc, contextData, step.Bindings...)
	if err != nil {
		fail()
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
		r.onFail()
		return true
	}
	cleaner := cleaner.New(tc.Timeouts().Cleanup.Duration, tc.DelayBeforeCleanup(), tc.DeletionPropagation())
	cleanup(func() {
		if !cleaner.Empty() || len(step.Cleanup) != 0 {
			report := &model.StepReport{
				Name:      fmt.Sprintf("cleanup (%s)", report.Name),
				StartTime: time.Now(),
			}
			defer func() {
				report.EndTime = time.Now()
				testReport.Add(report)
			}()
			logging.Log(ctx, logging.Cleanup, logging.BeginStatus, nil, color.BoldFgCyan)
			defer func() {
				logging.Log(ctx, logging.Cleanup, logging.EndStatus, nil, color.BoldFgCyan)
			}()
			if !cleaner.Empty() {
				if errs := cleaner.Run(ctx, report); len(errs) != 0 {
					fail()
					for _, err := range errs {
						logging.Log(ctx, logging.Cleanup, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
					}
					r.onFail()
				}
			}
			for i, operation := range step.Cleanup {
				if operation.Compiler != nil {
					tc = tc.WithDefaultCompiler(string(*operation.Compiler))
				}
				outputsTc, err := r.runCatch(ctx, tc, operation, i)
				if err != nil {
					fail()
				}
				tc = outputsTc
			}
		}
	})
	if len(step.Finally) != 0 {
		defer func() {
			logging.Log(ctx, logging.Finally, logging.BeginStatus, nil, color.BoldFgCyan)
			defer func() {
				logging.Log(ctx, logging.Finally, logging.EndStatus, nil, color.BoldFgCyan)
			}()
			for i, operation := range step.Finally {
				if operation.Compiler != nil {
					tc = tc.WithDefaultCompiler(string(*operation.Compiler))
				}
				outputsTc, err := r.runCatch(ctx, tc, operation, i)
				if err != nil {
					fail()
				}
				tc = outputsTc
			}
		}()
	}
	if catch := tc.Catch(); len(catch) != 0 {
		defer func() {
			if failed() {
				logging.Log(ctx, logging.Catch, logging.BeginStatus, nil, color.BoldFgCyan)
				defer func() {
					logging.Log(ctx, logging.Catch, logging.EndStatus, nil, color.BoldFgCyan)
				}()
				for i, operation := range catch {
					if operation.Compiler != nil {
						tc = tc.WithDefaultCompiler(string(*operation.Compiler))
					}
					outputsTc, err := r.runCatch(ctx, tc, operation, i)
					if err != nil {
						fail()
					}
					tc = outputsTc
				}
			}
		}()
	}
	logging.Log(ctx, logging.Try, logging.BeginStatus, nil, color.BoldFgCyan)
	defer func() {
		logging.Log(ctx, logging.Try, logging.EndStatus, nil, color.BoldFgCyan)
	}()
	for i, operation := range step.Try {
		continueOnError, outputsTc, err := r.runOperation(ctx, tc, operation, i, cleaner, report)
		if err != nil {
			fail()
			if !continueOnError {
				return true
			}
		}
		tc = outputsTc
	}
	return false
}

func (r *runner) runOperation(
	ctx context.Context,
	tc enginecontext.TestContext,
	operation v1alpha1.Operation,
	operationId int,
	cleaner cleaner.Cleaner,
	report *model.StepReport,
) (bool, enginecontext.TestContext, error) {
	if operation.Compiler != nil {
		tc = tc.WithDefaultCompiler(string(*operation.Compiler))
	}
	opType, actions, err := operations.TryOperation(ctx, tc, operation, cleaner)
	if err != nil {
		logging.Log(ctx, logging.Try, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
		r.onFail()
		return false, tc, err
	}
	continueOnError := operation.ContinueOnError != nil && *operation.ContinueOnError
	var errs []error
	for i, action := range actions {
		outputs, err := r.runAction(ctx, action, opType, operationId, i, tc, report)
		if err != nil {
			errs = append(errs, err)
			r.onFail()
		}
		for k, v := range outputs {
			tc = tc.WithBinding(k, v)
		}
	}
	return continueOnError, tc, multierr.Combine(errs...)
}

func (r *runner) runCatch(
	ctx context.Context,
	tc enginecontext.TestContext,
	operation v1alpha1.CatchFinally,
	operationId int,
) (enginecontext.TestContext, error) {
	if operation.Compiler != nil {
		tc = tc.WithDefaultCompiler(string(*operation.Compiler))
	}
	actions, err := operations.CatchOperation(ctx, tc, operation)
	if err != nil {
		logging.Log(ctx, logging.Try, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
		r.onFail()
		return tc, err
	}
	var errs []error
	for i, action := range actions {
		outputs, err := r.runAction(ctx, action, "", operationId, i, tc, nil)
		if err != nil {
			errs = append(errs, err)
			r.onFail()
		}
		for k, v := range outputs {
			tc = tc.WithBinding(k, v)
		}
	}
	return tc, multierr.Combine(errs...)
}

func (*runner) runAction(
	ctx context.Context,
	action operations.Operation,
	opType model.OperationType,
	operationId int,
	actionId int,
	tc enginecontext.TestContext,
	stepReport *model.StepReport,
) (outputs map[string]any, err error) {
	tc = tc.WithBinding("operation", OperationInfo{
		Id:         operationId + 1,
		ResourceId: actionId + 1,
	})
	if stepReport != nil {
		report := &model.OperationReport{
			Type:      opType,
			StartTime: time.Now(),
		}
		defer func() {
			report.EndTime = time.Now()
			report.Err = err
			stepReport.Add(report)
		}()
	}

	ctx = enginecontext.TestContextToCtx(ctx, &tc)

	return action.Execute(ctx, tc)
}

func (r *runner) onFail() {
	if r.onFailure != nil {
		r.onFailure()
	}
}

func (r *runner) cleanup(ctx context.Context, tc enginecontext.TestContext, cleaner cleaner.Cleaner) error {
	if tc.SkipDelete() {
		logging.Log(ctx, logging.Cleanup, logging.SkippedStatus, nil, color.BoldYellow)
		return nil
	}
	var errs []error
	if !cleaner.Empty() {
		logging.Log(ctx, logging.Cleanup, logging.BeginStatus, nil, color.BoldFgCyan)
		defer func() {
			logging.Log(ctx, logging.Cleanup, logging.EndStatus, nil, color.BoldFgCyan)
		}()
		for _, err := range cleaner.Run(ctx, nil) {
			errs = append(errs, err)
			logging.Log(ctx, logging.Cleanup, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
			r.onFail()
		}
	}
	return multierr.Combine(errs...)
}

func (r *runner) setupNamespace(ctx context.Context, nsOptions v1alpha2.NamespaceOptions, tc enginecontext.TestContext, cleanup cleaner.Cleaner) (enginecontext.TestContext, error) {
	if nsOptions.Name != "" {
		compilers := tc.Compilers()
		if nsOptions.Compiler != nil {
			compilers = compilers.WithDefaultCompiler(string(*nsOptions.Compiler))
		}
		var ns *corev1.Namespace
		if namespace, err := buildNamespace(ctx, compilers, nsOptions.Name, nsOptions.Template, tc); err != nil {
			logging.Log(ctx, logging.Internal, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
			r.onFail()
			return tc, err
		} else if _, clusterClient, err := tc.CurrentClusterClient(); err != nil {
			logging.Log(ctx, logging.Internal, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
			r.onFail()
			return tc, err
		} else if clusterClient != nil {
			if err := clusterClient.Get(ctx, client.Key(namespace), namespace.DeepCopy()); err != nil {
				if !errors.IsNotFound(err) {
					logging.Log(ctx, logging.Internal, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
					r.onFail()
					return tc, err
				} else if err := clusterClient.Create(ctx, namespace.DeepCopy()); err != nil {
					logging.Log(ctx, logging.Internal, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
					r.onFail()
					return tc, err
				} else if cleanup != nil {
					cleanup.Add(clusterClient, namespace)
				}
			}
			ns = namespace
		}
		if ns != nil {
			tc = tc.WithBinding("namespace", ns.GetName())
			tc = tc.WithNamespacer(namespacer.New(ns.GetName()))
			return tc, nil
		}
	}
	return tc, nil
}

func (r *runner) setupTestContext(ctx context.Context, testId int, scenarioId int, tc enginecontext.TestContext, test discovery.Test, bindings ...v1alpha1.Binding) (enginecontext.TestContext, error) {
	tc = tc.WithBinding("test", TestInfo{
		Id:         testId,
		ScenarioId: scenarioId,
		Metadata:   test.Test.ObjectMeta,
	})
	tc, err := enginecontext.WithBindings(tc, bindings...)
	if err != nil {
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
		r.onFail()
		return tc, err
	}
	contextData := enginecontext.ContextData{
		BasePath:            &test.BasePath,
		Catch:               test.Test.Spec.Catch,
		Cluster:             test.Test.Spec.Cluster,
		Clusters:            test.Test.Spec.Clusters,
		DelayBeforeCleanup:  test.Test.Spec.DelayBeforeCleanup,
		DeletionPropagation: test.Test.Spec.DeletionPropagationPolicy,
		SkipDelete:          test.Test.Spec.SkipDelete,
		Templating:          test.Test.Spec.Template,
		TerminationGrace:    test.Test.Spec.ForceTerminationGracePeriod,
		Timeouts:            test.Test.Spec.Timeouts,
	}
	tc, err = enginecontext.SetupContext(tc, contextData)
	if err != nil {
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
		r.onFail()
	}
	return tc, err
}

func (r *runner) testCleanup(ctx context.Context, tc enginecontext.TestContext, cleaner cleaner.Cleaner, report *model.TestReport) error {
	if tc.SkipDelete() {
		logging.Log(ctx, logging.Cleanup, logging.SkippedStatus, nil, color.BoldYellow)
		return nil
	}
	var errs []error
	if !cleaner.Empty() {
		logging.Log(ctx, logging.Cleanup, logging.BeginStatus, nil, color.BoldFgCyan)
		defer func() {
			logging.Log(ctx, logging.Cleanup, logging.EndStatus, nil, color.BoldFgCyan)
		}()
		stepReport := &model.StepReport{
			Name:      "@cleanup",
			StartTime: time.Now(),
		}
		defer func() {
			stepReport.EndTime = time.Now()
			report.Add(stepReport)
		}()
		for _, err := range cleaner.Run(ctx, stepReport) {
			errs = append(errs, err)
			logging.Log(ctx, logging.Cleanup, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
			r.onFail()
		}
	}
	return multierr.Combine(errs...)
}
