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
	"github.com/kyverno/chainsaw/pkg/report"
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
	report *report.TestReport,
	nsTemplate *v1alpha1.Any,
	delayBeforeCleanup *time.Duration,
	terminationGracePeriod *metav1.Duration,
	timeouts v1alpha1.DefaultTimeouts,
	deletionPropagationPolicy metav1.DeletionPropagation,
	templating bool,
	skipDelete bool,
	catch ...v1alpha1.CatchFinally,
) TestProcessor {
	if template := test.Test.Spec.NamespaceTemplate; template != nil && template.Value != nil {
		nsTemplate = template
	}
	if test.Test.Spec.DelayBeforeCleanup != nil {
		delayBeforeCleanup = &test.Test.Spec.DelayBeforeCleanup.Duration
	}
	if test.Test.Spec.ForceTerminationGracePeriod != nil {
		terminationGracePeriod = test.Test.Spec.ForceTerminationGracePeriod
	}
	if test.Test.Spec.Timeouts != nil {
		timeouts = withTimeouts(timeouts, *test.Test.Spec.Timeouts)
	}
	if test.Test.Spec.DeletionPropagationPolicy != nil {
		deletionPropagationPolicy = *test.Test.Spec.DeletionPropagationPolicy
	}
	if test.Test.Spec.Template != nil {
		templating = *test.Test.Spec.Template
	}
	if test.Test.Spec.SkipDelete != nil {
		skipDelete = *test.Test.Spec.SkipDelete
	}
	catch = append(catch, test.Test.Spec.Catch...)
	return &testProcessor{
		test:                      test,
		size:                      size,
		clock:                     clock,
		report:                    report,
		nsTemplate:                nsTemplate,
		delayBeforeCleanup:        delayBeforeCleanup,
		terminationGracePeriod:    terminationGracePeriod,
		timeouts:                  timeouts,
		deletionPropagationPolicy: deletionPropagationPolicy,
		templating:                templating,
		skipDelete:                skipDelete,
		catch:                     catch,
	}
}

type testProcessor struct {
	test                      discovery.Test
	size                      int
	clock                     clock.PassiveClock
	report                    *report.TestReport
	nsTemplate                *v1alpha1.Any
	delayBeforeCleanup        *time.Duration
	terminationGracePeriod    *metav1.Duration
	timeouts                  v1alpha1.DefaultTimeouts
	deletionPropagationPolicy metav1.DeletionPropagation
	templating                bool
	skipDelete                bool
	catch                     []v1alpha1.CatchFinally
}

func (p *testProcessor) Run(ctx context.Context, nspacer namespacer.Namespacer, tc engine.Context) {
	t := testing.FromContext(ctx)
	if p.report != nil {
		p.report.SetStartTime(time.Now())
		t.Cleanup(func() {
			p.report.SetEndTime(time.Now())
			if t.Failed() {
				p.report.Fail()
			}
			if t.Skipped() {
				p.report.Skip()
			}
		})
	}
	mainCleaner := cleaner.New(p.timeouts.Cleanup.Duration, nil, p.deletionPropagationPolicy)
	t.Cleanup(func() {
		if !mainCleaner.Empty() {
			logging.Log(ctx, logging.Cleanup, logging.RunStatus, color.BoldFgCyan)
			defer func() {
				logging.Log(ctx, logging.Cleanup, logging.DoneStatus, color.BoldFgCyan)
			}()
			for _, err := range mainCleaner.Run(ctx) {
				logging.Log(ctx, logging.Cleanup, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
				failer.Fail(ctx)
			}
		}
	})
	contextData := contextData{
		basePath: p.test.BasePath,
		clusters: p.test.Test.Spec.Clusters,
		cluster:  p.test.Test.Spec.Cluster,
		bindings: p.test.Test.Spec.Bindings,
	}
	nsName := p.test.Test.Spec.Namespace
	if nspacer == nil && nsName == "" {
		nsName = fmt.Sprintf("chainsaw-%s", petname.Generate(2, "-"))
	}
	if nsName != "" {
		var nsCleaner cleaner.CleanerCollector
		if !p.skipDelete {
			nsCleaner = mainCleaner
		}
		contextData.namespace = &namespaceData{
			name:     nsName,
			template: p.nsTemplate,
			cleaner:  nsCleaner,
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
	if p.report != nil && nspacer != nil {
		p.report.SetNamespace(nspacer.GetNamespace())
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
		processor := p.createStepProcessor(step)
		processor.Run(ctx, nspacer, tc)
	}
}

func (p *testProcessor) createStepProcessor(step v1alpha1.TestStep) StepProcessor {
	var report *report.StepReport
	if p.report != nil {
		report = p.report.ForStep(&step)
	}
	return NewStepProcessor(
		step,
		p.test.BasePath,
		report,
		p.delayBeforeCleanup,
		p.terminationGracePeriod,
		p.timeouts,
		p.deletionPropagationPolicy,
		p.templating,
		p.skipDelete,
		p.catch...,
	)
}
