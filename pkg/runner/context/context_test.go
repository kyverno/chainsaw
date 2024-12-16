package context

import (
	"testing"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/engine/clusters"
	"github.com/kyverno/chainsaw/pkg/loaders/config"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/kyverno/kyverno-json/pkg/core/expression"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/utils/clock"
	tclock "k8s.io/utils/clock/testing"
	"k8s.io/utils/ptr"
)

func TestEmptyContext(t *testing.T) {
	clock := tclock.NewFakePassiveClock(time.Time{})
	tests := []struct {
		name string
		want TestContext
	}{{
		want: TestContext{
			Summary: &model.Summary{},
			Report: &model.Report{
				Name:      "chainsaw-report",
				StartTime: clock.Now(),
			},
			bindings:            apis.NewBindings(),
			clusters:            clusters.NewRegistry(nil),
			compilers:           apis.DefaultCompilers,
			deletionPropagation: metav1.DeletePropagationBackground,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EmptyContext(clock)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTestContext_Bindings(t *testing.T) {
	parent := EmptyContext(clock.RealClock{})
	child := parent.WithBinding("foo", "bar")
	{
		value, err := parent.Bindings().Get("$foo")
		assert.Error(t, err)
		assert.Nil(t, value)
	}
	{
		value, err := child.Bindings().Get("$foo")
		assert.NoError(t, err)
		assert.NotNil(t, value)
	}
}

func TestTestContext_Catch(t *testing.T) {
	parent := EmptyContext(clock.RealClock{})
	child := parent.WithCatch([]v1alpha1.CatchFinally{{}}...)
	{
		value := parent.Catch()
		assert.Nil(t, value)
	}
	{
		value := child.Catch()
		assert.NotNil(t, value)
	}
}

func TestTestContext_Cluster(t *testing.T) {
	config, err := config.DefaultConfiguration()
	assert.NoError(t, err)
	defaultCluster := &rest.Config{}
	parent, err := InitContext(config.Spec, defaultCluster, nil)
	assert.NoError(t, err)
	child := parent.
		WithCluster(clusters.DefaultClient, nil).
		WithCluster("foo", clusters.NewClusterFromKubeconfig("foo", "bar"))
	{
		d := parent.Cluster(clusters.DefaultClient)
		f := parent.Cluster("foo")
		assert.NotNil(t, d)
		assert.Nil(t, f)
	}
	{
		d := child.Cluster(clusters.DefaultClient)
		f := child.Cluster("foo")
		assert.Nil(t, d)
		assert.NotNil(t, f)
	}
}

func TestTestContext_Clusters(t *testing.T) {
	config, err := config.DefaultConfiguration()
	assert.NoError(t, err)
	defaultCluster := &rest.Config{}
	parent, err := InitContext(config.Spec, defaultCluster, nil)
	assert.NoError(t, err)
	child := parent.WithCluster("foo", clusters.NewClusterFromKubeconfig("foo", "bar"))
	{
		assert.NotNil(t, parent.Clusters())
		d := parent.Clusters().Lookup(clusters.DefaultClient)
		f := parent.Clusters().Lookup("foo")
		assert.NotNil(t, d)
		assert.Nil(t, f)
	}
	{
		assert.NotNil(t, child.Clusters())
		d := child.Clusters().Lookup(clusters.DefaultClient)
		f := child.Clusters().Lookup("foo")
		assert.NotNil(t, d)
		assert.NotNil(t, f)
	}
}

func TestTestContext_Compilers(t *testing.T) {
	parent := EmptyContext(clock.RealClock{})
	child := parent.WithDefaultCompiler(expression.CompilerCEL)
	{
		value := parent.Compilers()
		assert.Equal(t, value.Jp, value.Default)
	}
	{
		value := child.Compilers()
		assert.Equal(t, value.Cel, value.Default)
	}
}

func TestTestContext_CurrentCluster(t *testing.T) {
	config, err := config.DefaultConfiguration()
	assert.NoError(t, err)
	defaultCluster := &rest.Config{}
	parent, err := InitContext(config.Spec, defaultCluster, nil)
	assert.NoError(t, err)
	child := parent.WithCurrentCluster("foo")
	{
		assert.NotNil(t, parent.CurrentCluster())
	}
	{
		assert.Nil(t, child.CurrentCluster())
	}
}

func TestTestContext_CurrentClusterClient(t *testing.T) {
	config, err := config.DefaultConfiguration()
	assert.NoError(t, err)
	defaultCluster := &rest.Config{}
	parent, err := InitContext(config.Spec, defaultCluster, nil)
	assert.NoError(t, err)
	child := parent.WithCurrentCluster("foo")
	child2 := parent.WithDryRun(true)
	{
		config, client, err := parent.CurrentClusterClient()
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.NotNil(t, client)
	}
	{
		config, client, err := child.CurrentClusterClient()
		assert.NoError(t, err)
		assert.Nil(t, config)
		assert.Nil(t, client)
	}
	{
		config, client, err := child2.CurrentClusterClient()
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.NotNil(t, client)
	}
}

func TestTestContext_DelayBeforeCleanup(t *testing.T) {
	parent := EmptyContext(clock.RealClock{})
	child := parent.WithDelayBeforeCleanup(ptr.To(30 * time.Second))
	{
		value := parent.DelayBeforeCleanup()
		assert.Nil(t, value)
	}
	{
		value := child.DelayBeforeCleanup()
		assert.Equal(t, value, ptr.To(30*time.Second))
	}
}

func TestTestContext_DeletionPropagation(t *testing.T) {
	parent := EmptyContext(clock.RealClock{})
	child := parent.WithDeletionPropagation(metav1.DeletePropagationOrphan)
	{
		value := parent.DeletionPropagation()
		assert.Equal(t, metav1.DeletePropagationBackground, value)
	}
	{
		value := child.DeletionPropagation()
		assert.Equal(t, metav1.DeletePropagationOrphan, value)
	}
}

func TestTestContext_DryRun(t *testing.T) {
	parent := EmptyContext(clock.RealClock{})
	child := parent.WithDryRun(true)
	{
		value := parent.DryRun()
		assert.False(t, value)
	}
	{
		value := child.DryRun()
		assert.True(t, value)
	}
}

func TestTestContext_FailFast(t *testing.T) {
	parent := EmptyContext(clock.RealClock{})
	child := parent.WithFailFast(true)
	{
		value := parent.FailFast()
		assert.False(t, value)
	}
	{
		value := child.FailFast()
		assert.True(t, value)
	}
}

func TestTestContext_FullName(t *testing.T) {
	parent := EmptyContext(clock.RealClock{})
	child := parent.WithFullName(true)
	{
		value := parent.FullName()
		assert.False(t, value)
	}
	{
		value := child.FullName()
		assert.True(t, value)
	}
}

func TestTestContext_SkipDelete(t *testing.T) {
	parent := EmptyContext(clock.RealClock{})
	child := parent.WithSkipDelete(true)
	{
		value := parent.SkipDelete()
		assert.False(t, value)
	}
	{
		value := child.SkipDelete()
		assert.True(t, value)
	}
}

func TestTestContext_Templating(t *testing.T) {
	parent := EmptyContext(clock.RealClock{})
	child := parent.WithTemplating(true)
	{
		value := parent.Templating()
		assert.False(t, value)
	}
	{
		value := child.Templating()
		assert.True(t, value)
	}
}

func TestTestContext_TerminationGrace(t *testing.T) {
	parent := EmptyContext(clock.RealClock{})
	child := parent.WithTerminationGrace(ptr.To(30 * time.Second))
	{
		value := parent.TerminationGrace()
		assert.Nil(t, value)
	}
	{
		value := child.TerminationGrace()
		assert.Equal(t, ptr.To(30*time.Second), value)
	}
}

func TestTestContext_Timeouts(t *testing.T) {
	parent := EmptyContext(clock.RealClock{})
	child := parent.WithTimeouts(v1alpha1.Timeouts{
		Apply:   &metav1.Duration{Duration: 10 * time.Second},
		Assert:  &metav1.Duration{Duration: 20 * time.Second},
		Cleanup: &metav1.Duration{Duration: 30 * time.Second},
		Delete:  &metav1.Duration{Duration: 40 * time.Second},
		Error:   &metav1.Duration{Duration: 50 * time.Second},
		Exec:    &metav1.Duration{Duration: 60 * time.Second},
	})
	child2 := parent.WithTimeouts(v1alpha1.Timeouts{})
	{
		value := parent.Timeouts()
		assert.Equal(t, v1alpha1.DefaultTimeouts{}, value)
	}
	{
		value := child.Timeouts()
		assert.Equal(t, v1alpha1.DefaultTimeouts{
			Apply:   metav1.Duration{Duration: 10 * time.Second},
			Assert:  metav1.Duration{Duration: 20 * time.Second},
			Cleanup: metav1.Duration{Duration: 30 * time.Second},
			Delete:  metav1.Duration{Duration: 40 * time.Second},
			Error:   metav1.Duration{Duration: 50 * time.Second},
			Exec:    metav1.Duration{Duration: 60 * time.Second},
		}, value)
	}
	{
		value := child2.Timeouts()
		assert.Equal(t, v1alpha1.DefaultTimeouts{}, value)
	}
}
