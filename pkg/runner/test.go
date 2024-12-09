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
	"github.com/kyverno/chainsaw/pkg/engine"
	"github.com/kyverno/chainsaw/pkg/engine/logging"
	"github.com/kyverno/chainsaw/pkg/engine/namespacer"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/kyverno/chainsaw/pkg/runner/processors"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/pkg/ext/output/color"
)

func (r *runner) runTest(
	ctx context.Context,
	t testing.TTest,
	nsOptions v1alpha2.NamespaceOptions,
	nspacer namespacer.Namespacer,
	tc engine.Context,
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
	ctx = logging.IntoContext(ctx, logging.NewLogger(t, r.clock, test.Test.Name, fmt.Sprintf("%-*s", size, "@chainsaw")))
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
	tc = tc.WithBinding(ctx, "test", processors.TestInfo{
		Id:         testId,
		ScenarioId: scenarioId,
		Metadata:   test.Test.ObjectMeta,
	})
	tc, err := engine.WithBindings(ctx, tc, bindings...)
	if err != nil {
		t.Fail()
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		r.failer.Fail()
		return
	}
	contextData := processors.ContextData{
		BasePath:            test.BasePath,
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
	tc, err = processors.SetupContext(ctx, tc, contextData)
	if err != nil {
		t.Fail()
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		r.failer.Fail()
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
			logging.Log(ctx, logging.Cleanup, logging.BeginStatus, color.BoldFgCyan)
			defer func() {
				logging.Log(ctx, logging.Cleanup, logging.EndStatus, color.BoldFgCyan)
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
				logging.Log(ctx, logging.Cleanup, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
				r.failer.Fail()
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
		namespaceData := processors.NamespaceData{
			Cleaner:   nsCleaner,
			Compilers: compilers,
			Name:      nsName,
			Template:  nsOptions.Template,
		}
		nsTc, namespace, err := processors.SetupNamespace(ctx, tc, namespaceData)
		if err != nil {
			t.Fail()
			logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
			r.failer.Fail()
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
	tc, err = processors.SetupBindings(ctx, tc, test.Test.Spec.Bindings...)
	if err != nil {
		t.Fail()
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		r.failer.Fail()
		return
	}
	// run steps
	for i, step := range test.Test.Spec.Steps {
		name := step.Name
		if name == "" {
			name = fmt.Sprintf("step-%d", i+1)
		}
		ctx := logging.IntoContext(ctx, logging.NewLogger(t, r.clock, test.Test.Name, fmt.Sprintf("%-*s", size, name)))
		info := processors.StepInfo{
			Id: i + 1,
		}
		tc := tc.WithBinding(ctx, "step", info)
		processor := processors.NewStepProcessor(step, report, test.BasePath)
		if stop := processor.Run(ctx, t, r.failer, nspacer, tc); stop {
			return
		}
	}
}
