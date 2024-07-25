package processors

import (
	"context"
	"fmt"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/cleanup/cleaner"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/engine/logging"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/kyverno/chainsaw/pkg/report"
	"github.com/kyverno/chainsaw/pkg/runner/failer"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/pkg/ext/output/color"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/utils/clock"
)

type TestProcessor interface {
	Run(context.Context, namespacer.Namespacer, model.TestContext, discovery.Test)
}

func NewTestProcessor(config model.Configuration, clock clock.PassiveClock, report *report.TestReport) TestProcessor {
	return &testProcessor{
		config: config,
		clock:  clock,
		report: report,
	}
}

type testProcessor struct {
	config model.Configuration
	clock  clock.PassiveClock
	report *report.TestReport
}

func (p *testProcessor) Run(ctx context.Context, nspacer namespacer.Namespacer, tc model.TestContext, test discovery.Test) {
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
	size := len("@cleanup")
	for i, step := range test.Test.Spec.Steps {
		name := step.Name
		if name == "" {
			name = fmt.Sprintf("step-%d", i+1)
		}
		if size < len(name) {
			size = len(name)
		}
	}
	if test.Test.Spec.Timeouts != nil {
		tc = model.WithTimeouts(ctx, tc, *test.Test.Spec.Timeouts)
	}
	timeouts := tc.Timeouts()
	cleaner := cleaner.New(timeouts.Cleanup, nil)
	setupLogger := logging.NewLogger(t, p.clock, test.Test.Name, fmt.Sprintf("%-*s", size, "@setup"))
	cleanupLogger := logging.NewLogger(t, p.clock, test.Test.Name, fmt.Sprintf("%-*s", size, "@cleanup"))
	setupCtx := logging.IntoContext(ctx, setupLogger)
	cleanupCtx := logging.IntoContext(ctx, cleanupLogger)
	t.Cleanup(func() {
		if !cleaner.Empty() {
			logging.Log(cleanupCtx, logging.Cleanup, logging.RunStatus, color.BoldFgCyan)
			defer func() {
				logging.Log(cleanupCtx, logging.Cleanup, logging.DoneStatus, color.BoldFgCyan)
			}()
			for _, err := range cleaner.Run(cleanupCtx) {
				logging.Log(cleanupCtx, logging.Cleanup, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
				failer.Fail(cleanupCtx)
			}
		}
	})
	tc = model.WithClusters(ctx, tc, test.BasePath, test.Test.Spec.Clusters)
	_, clusterClient, _tc, err := model.WithCurrentCluster(ctx, tc, test.Test.Spec.Cluster)
	if err != nil {
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		failer.FailNow(ctx)
	}
	tc = _tc
	if test.Test.Spec.SkipDelete != nil {
		tc = tc.WithCleanup(ctx, !*test.Test.Spec.SkipDelete)
	}
	if clusterClient != nil {
		nsName := test.Test.Spec.Namespace
		if nspacer == nil && nsName == "" {
			nsName = fmt.Sprintf("chainsaw-%s", petname.Generate(2, "-"))
		}
		if nsName != "" {
			template := test.Test.Spec.NamespaceTemplate
			if template == nil || template.Value == nil {
				template = p.config.Namespace.Template
			}
			namespace, err := buildNamespace(ctx, nsName, template, tc.Bindings())
			if err != nil {
				logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
				failer.FailNow(ctx)
			}
			nspacer = namespacer.New(clusterClient, namespace.GetName())
			if err := clusterClient.Get(setupCtx, client.Key(namespace), namespace.DeepCopy()); err != nil {
				if !errors.IsNotFound(err) {
					// Get doesn't log
					setupLogger.Log(logging.Get, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
					failer.FailNow(ctx)
				} else if err := clusterClient.Create(logging.IntoContext(setupCtx, setupLogger), namespace.DeepCopy()); err != nil {
					failer.FailNow(ctx)
				} else if tc.Cleanup() {
					cleaner.Add(clusterClient, namespace)
				}
			}
			tc = tc.WithBinding(ctx, "namespace", namespace.GetName())
		}
	}
	if p.report != nil && nspacer != nil {
		p.report.SetNamespace(nspacer.GetNamespace())
	}
	if _tc, err := model.WithBindings(ctx, tc, test.Test.Spec.Bindings...); err != nil {
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		failer.FailNow(ctx)
	} else {
		tc = _tc
	}
	for i, step := range test.Test.Spec.Steps {
		name := step.Name
		if name == "" {
			name = fmt.Sprintf("step-%d", i+1)
		}
		ctx := logging.IntoContext(ctx, logging.NewLogger(t, p.clock, test.Test.Name, fmt.Sprintf("%-*s", size, name)))
		info := StepInfo{
			Id: i + 1,
		}
		tc := tc.WithBinding(ctx, "step", info)
		processor := p.createStepProcessor(nspacer, test, step)
		processor.Run(ctx, tc)
	}
}

func (p *testProcessor) createStepProcessor(nspacer namespacer.Namespacer, test discovery.Test, step v1alpha1.TestStep) StepProcessor {
	var report *report.StepReport
	if p.report != nil {
		report = p.report.ForStep(&step)
	}
	return NewStepProcessor(p.config, nspacer, p.clock, test, step, report)
}
