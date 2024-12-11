package runner

import (
	"context"
	"fmt"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha2"
	"github.com/kyverno/chainsaw/pkg/cleanup/cleaner"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/engine/namespacer"
	"github.com/kyverno/chainsaw/pkg/logging"
	"github.com/kyverno/chainsaw/pkg/model"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/pkg/ext/output/color"
)

func (r *runner) runTest(
	ctx context.Context,
	t testing.TTest,
	nsOptions v1alpha2.NamespaceOptions,
	nspacer namespacer.Namespacer,
	tc enginecontext.TestContext,
	test discovery.Test,
	testId int,
	scenarioId int,
	bindings ...v1alpha1.Binding,
) {
	// configure golang context
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
	ctx = logging.WithLogger(ctx, logging.NewLogger(test.Test.Name, fmt.Sprintf("%-*s", size, "@chainsaw")))
	// setup summary
	t.Cleanup(func() {
		if t.Skipped() {
			tc.IncSkipped()
		} else if t.Failed() {
			tc.IncFailed()
		} else {
			tc.IncPassed()
		}
	})
	// setup concurrency
	if test.Test.Spec.Concurrent == nil || *test.Test.Spec.Concurrent {
		t.Parallel()
	}
	// setup reporting
	report := &model.TestReport{
		BasePath:   test.BasePath,
		Name:       test.Test.Name,
		Concurrent: test.Test.Spec.Concurrent,
		StartTime:  time.Now(),
	}
	stepReport := &model.StepReport{
		Name:      "main",
		StartTime: time.Now(),
	}
	t.Cleanup(func() {
		report.EndTime = time.Now()
		if t.Skipped() {
			report.Skipped = true
		}
		tc.Report.Add(report)
	})
	// setup context
	tc = tc.WithBinding(ctx, "test", TestInfo{
		Id:         testId,
		ScenarioId: scenarioId,
		Metadata:   test.Test.ObjectMeta,
	})
	tc, err := enginecontext.WithBindings(ctx, tc, bindings...)
	if err != nil {
		t.Fail()
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
		r.onFail()
		return
	}
	contextData := contextData{
		basePath:            test.BasePath,
		catch:               test.Test.Spec.Catch,
		cluster:             test.Test.Spec.Cluster,
		clusters:            test.Test.Spec.Clusters,
		delayBeforeCleanup:  test.Test.Spec.DelayBeforeCleanup,
		deletionPropagation: test.Test.Spec.DeletionPropagationPolicy,
		skipDelete:          test.Test.Spec.SkipDelete,
		templating:          test.Test.Spec.Template,
		terminationGrace:    test.Test.Spec.ForceTerminationGracePeriod,
		timeouts:            test.Test.Spec.Timeouts,
	}
	tc, err = setupContext(ctx, tc, contextData)
	if err != nil {
		t.Fail()
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
		r.onFail()
		return
	}
	// skip checks
	if test.Test.Spec.Skip != nil && *test.Test.Spec.Skip {
		t.Skip()
		return
	}
	if tc.FailFast() && tc.Failed() > 0 {
		t.Skip()
		return
	}
	// setup cleaner
	mainCleaner := cleaner.New(tc.Timeouts().Cleanup.Duration, nil, tc.DeletionPropagation())
	t.Cleanup(func() {
		if !mainCleaner.Empty() {
			logging.Log(ctx, logging.Cleanup, logging.BeginStatus, nil, color.BoldFgCyan)
			defer func() {
				logging.Log(ctx, logging.Cleanup, logging.EndStatus, nil, color.BoldFgCyan)
			}()
			stepReport := &model.StepReport{
				Name:      fmt.Sprintf("cleanup (%s)", stepReport.Name),
				StartTime: time.Now(),
			}
			defer func() {
				stepReport.EndTime = time.Now()
				report.Add(stepReport)
			}()
			for _, err := range mainCleaner.Run(ctx, stepReport) {
				t.Fail()
				logging.Log(ctx, logging.Cleanup, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
				r.onFail()
			}
		}
	})
	// setup namespace
	// TODO: should be part of setupContext ?
	if test.Test.Spec.Compiler != nil {
		tc = tc.WithDefaultCompiler(string(*test.Test.Spec.Compiler))
	}
	nsName := test.Test.Spec.Namespace
	if nspacer == nil && nsName == "" {
		nsName = fmt.Sprintf("chainsaw-%s", petname.Generate(2, "-"))
	}
	if nsName != "" {
		var nsCleaner cleaner.CleanerCollector
		if !tc.SkipDelete() {
			nsCleaner = mainCleaner
		}
		// TODO: this may not use the right default compiler if the template is coming from the config
		// but the default compiler is specified at the test level
		if template := test.Test.Spec.NamespaceTemplate; template != nil && template.Value() != nil {
			nsOptions.Template = template
			nsOptions.Compiler = test.Test.Spec.NamespaceTemplateCompiler
		}
		compilers := tc.Compilers()
		if nsOptions.Compiler != nil {
			compilers = compilers.WithDefaultCompiler(string(*nsOptions.Compiler))
		}
		namespaceData := namespaceData{
			cleaner:   nsCleaner,
			compilers: compilers,
			name:      nsName,
			template:  nsOptions.Template,
		}
		nsTc, namespace, err := setupNamespace(ctx, tc, namespaceData)
		if err != nil {
			t.Fail()
			logging.Log(ctx, logging.Internal, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
			r.onFail()
			return
		}
		tc = nsTc
		if namespace != nil {
			nspacer = namespacer.New(namespace.GetName())
		}
	}
	if nspacer != nil {
		report.Namespace = nspacer.GetNamespace()
	}
	// setup bindings
	tc, err = setupBindings(ctx, tc, test.Test.Spec.Bindings...)
	if err != nil {
		t.Fail()
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
		r.onFail()
		return
	}
	// run steps
	for i, step := range test.Test.Spec.Steps {
		name := step.Name
		if name == "" {
			name = fmt.Sprintf("step-%d", i+1)
		}
		ctx := logging.WithLogger(ctx, logging.NewLogger(test.Test.Name, fmt.Sprintf("%-*s", size, name)))
		info := StepInfo{
			Id: i + 1,
		}
		tc := tc.WithBinding(ctx, "step", info)
		if stop := r.runStep(ctx, t, test.BasePath, nspacer, tc, step, report); stop {
			return
		}
	}
}
