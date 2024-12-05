package processors

import (
	"context"
	"fmt"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/cleanup/cleaner"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/engine"
	"github.com/kyverno/chainsaw/pkg/engine/logging"
	"github.com/kyverno/chainsaw/pkg/engine/namespacer"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/kyverno/chainsaw/pkg/runner/failer"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/pkg/ext/output/color"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/clock"
)

type TestProcessor interface {
	Run(context.Context, namespacer.Namespacer, engine.Context)
}

func NewTestProcessor(
	test discovery.Test,
	size int,
	clock clock.PassiveClock,
	nsTemplate *v1alpha1.Projection,
	nsTemplateCompiler *v1alpha1.Compiler,
	terminationGracePeriod *metav1.Duration,
	timeouts v1alpha1.DefaultTimeouts,
	catch ...v1alpha1.CatchFinally,
) TestProcessor {
	if template := test.Test.Spec.NamespaceTemplate; template != nil && template.Value() != nil {
		nsTemplate = template
		nsTemplateCompiler = test.Test.Spec.NamespaceTemplateCompiler
	}
	if test.Test.Spec.ForceTerminationGracePeriod != nil {
		terminationGracePeriod = test.Test.Spec.ForceTerminationGracePeriod
	}
	if test.Test.Spec.Timeouts != nil {
		timeouts = withTimeouts(timeouts, *test.Test.Spec.Timeouts)
	}
	catch = append(catch, test.Test.Spec.Catch...)
	return &testProcessor{
		test:                   test,
		size:                   size,
		clock:                  clock,
		nsTemplate:             nsTemplate,
		nsTemplateCompiler:     nsTemplateCompiler,
		terminationGracePeriod: terminationGracePeriod,
		timeouts:               timeouts,
		catch:                  catch,
	}
}

type testProcessor struct {
	test                   discovery.Test
	size                   int
	clock                  clock.PassiveClock
	nsTemplate             *v1alpha1.Projection
	nsTemplateCompiler     *v1alpha1.Compiler
	terminationGracePeriod *metav1.Duration
	timeouts               v1alpha1.DefaultTimeouts
	catch                  []v1alpha1.CatchFinally
}

func (p *testProcessor) Run(ctx context.Context, nspacer namespacer.Namespacer, tc engine.Context) {
	t := testing.FromContext(ctx)
	report := &model.TestReport{
		BasePath:   p.test.BasePath,
		Name:       p.test.Test.Name,
		Concurrent: p.test.Test.Spec.Concurrent,
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
	// TODO: not the expected deletion propagation
	mainCleaner := cleaner.New(p.timeouts.Cleanup.Duration, nil, tc.DeletionPropagation())
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
				logging.Log(ctx, logging.Cleanup, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
				failer.Fail(ctx)
			}
		}
	})
	if p.test.Test.Spec.Compiler != nil {
		tc = tc.WithDefaultCompiler(string(*p.test.Test.Spec.Compiler))
	}
	contextData := contextData{
		basePath:            p.test.BasePath,
		bindings:            p.test.Test.Spec.Bindings,
		cluster:             p.test.Test.Spec.Cluster,
		clusters:            p.test.Test.Spec.Clusters,
		deletionPropagation: p.test.Test.Spec.DeletionPropagationPolicy,
		delayBeforeCleanup:  p.test.Test.Spec.DelayBeforeCleanup,
		skipDelete:          p.test.Test.Spec.SkipDelete,
		templating:          p.test.Test.Spec.Template,
	}
	nsName := p.test.Test.Spec.Namespace
	if nspacer == nil && nsName == "" {
		nsName = fmt.Sprintf("chainsaw-%s", petname.Generate(2, "-"))
	}
	if nsName != "" {
		var nsCleaner cleaner.CleanerCollector
		if !tc.SkipDelete() {
			nsCleaner = mainCleaner
		}
		// TODO this may not use the right default compiler if the template is coming from the config
		// but the default compiler is specified at the test level
		compilers := tc.Compilers()
		if p.nsTemplateCompiler != nil {
			compilers = compilers.WithDefaultCompiler(string(*p.nsTemplateCompiler))
		}
		contextData.namespace = &namespaceData{
			name:      nsName,
			template:  p.nsTemplate,
			compilers: compilers,
			cleaner:   nsCleaner,
		}
	}
	tc, namespace, err := setupContextData(ctx, tc, contextData)
	if err != nil {
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		failer.FailNow(ctx)
	}
	if namespace != nil {
		nspacer = namespacer.New(namespace.GetName())
	}
	if nspacer != nil {
		report.Namespace = nspacer.GetNamespace()
	}
	for i, step := range p.test.Test.Spec.Steps {
		name := step.Name
		if name == "" {
			name = fmt.Sprintf("step-%d", i+1)
		}
		ctx := logging.IntoContext(ctx, logging.NewLogger(t, p.clock, p.test.Test.Name, fmt.Sprintf("%-*s", p.size, name)))
		info := StepInfo{
			Id: i + 1,
		}
		tc := tc.WithBinding(ctx, "step", info)
		processor := p.createStepProcessor(step, report)
		processor.Run(ctx, nspacer, tc)
	}
}

func (p *testProcessor) createStepProcessor(step v1alpha1.TestStep, report *model.TestReport) StepProcessor {
	return NewStepProcessor(
		step,
		report,
		p.test.BasePath,
		p.terminationGracePeriod,
		p.timeouts,
		p.catch...,
	)
}
