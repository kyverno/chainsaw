package context

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ContextData struct {
	BasePath            string
	Catch               []v1alpha1.CatchFinally
	Cluster             *string
	Clusters            v1alpha1.Clusters
	DelayBeforeCleanup  *metav1.Duration
	DeletionPropagation *metav1.DeletionPropagation
	DryRun              *bool
	SkipDelete          *bool
	Templating          *bool
	TerminationGrace    *metav1.Duration
	Timeouts            *v1alpha1.Timeouts
}

func SetupContext(tc TestContext, data ContextData) (TestContext, error) {
	if len(data.Catch) > 0 {
		tc = tc.WithCatch(data.Catch...)
	}
	if data.DryRun != nil {
		tc = tc.WithDryRun(*data.DryRun)
	}
	if data.DelayBeforeCleanup != nil {
		tc = tc.WithDelayBeforeCleanup(&data.DelayBeforeCleanup.Duration)
	}
	if data.DeletionPropagation != nil {
		tc = tc.WithDeletionPropagation(*data.DeletionPropagation)
	}
	if data.SkipDelete != nil {
		tc = tc.WithSkipDelete(*data.SkipDelete)
	}
	if data.Templating != nil {
		tc = tc.WithTemplating(*data.Templating)
	}
	if data.TerminationGrace != nil {
		tc = tc.WithTerminationGrace(&data.TerminationGrace.Duration)
	}
	if data.Timeouts != nil {
		tc = tc.WithTimeouts(*data.Timeouts)
	}
	tc = WithClusters(tc, data.BasePath, data.Clusters)
	if data.Cluster != nil {
		if _tc, err := WithCurrentCluster(tc, *data.Cluster); err != nil {
			return tc, err
		} else {
			tc = _tc
		}
	}
	return tc, nil
}

func SetupBindings(tc TestContext, bindings ...v1alpha1.Binding) (TestContext, error) {
	if _tc, err := WithBindings(tc, bindings...); err != nil {
		return tc, err
	} else {
		tc = _tc
	}
	return tc, nil
}

func SetupContextAndBindings(tc TestContext, data ContextData, bindings ...v1alpha1.Binding) (TestContext, error) {
	if tc, err := SetupContext(tc, data); err != nil {
		return tc, err
	} else {
		return SetupBindings(tc, bindings...)
	}
}
