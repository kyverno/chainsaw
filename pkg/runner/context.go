package runner

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type contextData struct {
	basePath            string
	catch               []v1alpha1.CatchFinally
	cluster             *string
	clusters            v1alpha1.Clusters
	delayBeforeCleanup  *metav1.Duration
	deletionPropagation *metav1.DeletionPropagation
	dryRun              *bool
	skipDelete          *bool
	templating          *bool
	terminationGrace    *metav1.Duration
	timeouts            *v1alpha1.Timeouts
}

func setupContext(tc enginecontext.TestContext, data contextData) (enginecontext.TestContext, error) {
	if len(data.catch) > 0 {
		tc = tc.WithCatch(data.catch...)
	}
	if data.dryRun != nil {
		tc = tc.WithDryRun(*data.dryRun)
	}
	if data.delayBeforeCleanup != nil {
		tc = tc.WithDelayBeforeCleanup(&data.delayBeforeCleanup.Duration)
	}
	if data.deletionPropagation != nil {
		tc = tc.WithDeletionPropagation(*data.deletionPropagation)
	}
	if data.skipDelete != nil {
		tc = tc.WithSkipDelete(*data.skipDelete)
	}
	if data.templating != nil {
		tc = tc.WithTemplating(*data.templating)
	}
	if data.terminationGrace != nil {
		tc = tc.WithTerminationGrace(&data.terminationGrace.Duration)
	}
	if data.timeouts != nil {
		tc = tc.WithTimeouts(*data.timeouts)
	}
	tc = enginecontext.WithClusters(tc, data.basePath, data.clusters)
	if data.cluster != nil {
		if _tc, err := enginecontext.WithCurrentCluster(tc, *data.cluster); err != nil {
			return tc, err
		} else {
			tc = _tc
		}
	}
	return tc, nil
}

func setupBindings(tc enginecontext.TestContext, bindings ...v1alpha1.Binding) (enginecontext.TestContext, error) {
	if _tc, err := enginecontext.WithBindings(tc, bindings...); err != nil {
		return tc, err
	} else {
		tc = _tc
	}
	return tc, nil
}

func setupContextAndBindings(tc enginecontext.TestContext, data contextData, bindings ...v1alpha1.Binding) (enginecontext.TestContext, error) {
	if tc, err := setupContext(tc, data); err != nil {
		return tc, err
	} else {
		return setupBindings(tc, bindings...)
	}
}
