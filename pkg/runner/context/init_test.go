package context

import (
	"testing"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/loaders/config"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/utils/ptr"
)

func TestInitContext(t *testing.T) {
	config, err := config.DefaultConfiguration()
	assert.NoError(t, err)
	tests := []struct {
		name           string
		config         model.Configuration
		defaultCluster *rest.Config
		values         any
		wantErr        bool
		want           func(*testing.T, TestContext)
	}{{
		name:           "default config",
		config:         config.Spec,
		defaultCluster: nil,
		values:         nil,
		wantErr:        false,
		want: func(t *testing.T, tc TestContext) {
			t.Helper()
			assert.Equal(t, int32(0), tc.Passed())
			assert.Equal(t, int32(0), tc.Failed())
			assert.Equal(t, int32(0), tc.Skipped())
			assert.NotNil(t, tc.Bindings())
			assert.Nil(t, tc.Catch())
			assert.NotNil(t, tc.Compilers())
			assert.NotNil(t, tc.Clusters())
			assert.Nil(t, tc.CurrentCluster())
			assert.Nil(t, tc.DelayBeforeCleanup())
			assert.Equal(t, metav1.DeletePropagationBackground, tc.DeletionPropagation())
			assert.False(t, tc.DryRun())
			assert.False(t, tc.FailFast())
			assert.False(t, tc.FullName())
			assert.False(t, tc.SkipDelete())
			assert.True(t, tc.Templating())
			assert.Nil(t, tc.TerminationGrace())
			assert.Equal(t, config.Spec.Timeouts, tc.Timeouts())
		},
	}, {
		name: "with delay before cleanup",
		config: func(config model.Configuration) model.Configuration {
			config.Cleanup.DelayBeforeCleanup = ptr.To(metav1.Duration{Duration: 30 * time.Second})
			return config
		}(config.Spec),
		defaultCluster: nil,
		values:         nil,
		wantErr:        false,
		want: func(t *testing.T, tc TestContext) {
			t.Helper()
			assert.Equal(t, int32(0), tc.Passed())
			assert.Equal(t, int32(0), tc.Failed())
			assert.Equal(t, int32(0), tc.Skipped())
			assert.NotNil(t, tc.Bindings())
			assert.Nil(t, tc.Catch())
			assert.NotNil(t, tc.Compilers())
			assert.NotNil(t, tc.Clusters())
			assert.Nil(t, tc.CurrentCluster())
			assert.Equal(t, ptr.To(30*time.Second), tc.DelayBeforeCleanup())
			assert.Equal(t, metav1.DeletePropagationBackground, tc.DeletionPropagation())
			assert.False(t, tc.DryRun())
			assert.False(t, tc.FailFast())
			assert.False(t, tc.FullName())
			assert.False(t, tc.SkipDelete())
			assert.True(t, tc.Templating())
			assert.Nil(t, tc.TerminationGrace())
		},
	}, {
		name: "with compiler",
		config: func(config model.Configuration) model.Configuration {
			config.Templating.Compiler = ptr.To(v1alpha1.EngineJP)
			return config
		}(config.Spec),
		defaultCluster: nil,
		values:         nil,
		wantErr:        false,
		want: func(t *testing.T, tc TestContext) {
			t.Helper()
			assert.Equal(t, int32(0), tc.Passed())
			assert.Equal(t, int32(0), tc.Failed())
			assert.Equal(t, int32(0), tc.Skipped())
			assert.NotNil(t, tc.Bindings())
			assert.Nil(t, tc.Catch())
			assert.NotNil(t, tc.Compilers())
			assert.Equal(t, tc.Compilers().Default, tc.Compilers().Jp)
			assert.NotNil(t, tc.Clusters())
			assert.Nil(t, tc.CurrentCluster())
			assert.Nil(t, tc.DelayBeforeCleanup())
			assert.Equal(t, metav1.DeletePropagationBackground, tc.DeletionPropagation())
			assert.False(t, tc.DryRun())
			assert.False(t, tc.FailFast())
			assert.False(t, tc.FullName())
			assert.False(t, tc.SkipDelete())
			assert.True(t, tc.Templating())
			assert.Nil(t, tc.TerminationGrace())
		},
	}, {
		name: "with termination grace",
		config: func(config model.Configuration) model.Configuration {
			config.Execution.ForceTerminationGracePeriod = ptr.To(metav1.Duration{Duration: 30 * time.Second})
			return config
		}(config.Spec),
		defaultCluster: nil,
		values:         nil,
		wantErr:        false,
		want: func(t *testing.T, tc TestContext) {
			t.Helper()
			assert.Equal(t, int32(0), tc.Passed())
			assert.Equal(t, int32(0), tc.Failed())
			assert.Equal(t, int32(0), tc.Skipped())
			assert.NotNil(t, tc.Bindings())
			assert.Nil(t, tc.Catch())
			assert.NotNil(t, tc.Compilers())
			assert.NotNil(t, tc.Clusters())
			assert.Nil(t, tc.CurrentCluster())
			assert.Nil(t, tc.DelayBeforeCleanup())
			assert.Equal(t, metav1.DeletePropagationBackground, tc.DeletionPropagation())
			assert.False(t, tc.DryRun())
			assert.False(t, tc.FailFast())
			assert.False(t, tc.FullName())
			assert.False(t, tc.SkipDelete())
			assert.True(t, tc.Templating())
			assert.NotNil(t, tc.TerminationGrace())
			assert.Equal(t, ptr.To(30*time.Second), tc.TerminationGrace())
		},
	}, {
		name:           "with default cluster",
		config:         config.Spec,
		defaultCluster: &rest.Config{},
		values:         nil,
		wantErr:        false,
		want: func(t *testing.T, tc TestContext) {
			t.Helper()
			assert.Equal(t, int32(0), tc.Passed())
			assert.Equal(t, int32(0), tc.Failed())
			assert.Equal(t, int32(0), tc.Skipped())
			assert.NotNil(t, tc.Bindings())
			assert.Nil(t, tc.Catch())
			assert.NotNil(t, tc.Compilers())
			assert.NotNil(t, tc.Clusters())
			assert.NotNil(t, tc.CurrentCluster())
			assert.Nil(t, tc.DelayBeforeCleanup())
			assert.Equal(t, metav1.DeletePropagationBackground, tc.DeletionPropagation())
			assert.False(t, tc.DryRun())
			assert.False(t, tc.FailFast())
			assert.False(t, tc.FullName())
			assert.False(t, tc.SkipDelete())
			assert.True(t, tc.Templating())
			assert.Nil(t, tc.TerminationGrace())
		},
	}, {
		name: "with clusters",
		config: func(config model.Configuration) model.Configuration {
			config.Clusters = v1alpha1.Clusters{
				"foo": v1alpha1.Cluster{
					Kubeconfig: "foo",
					Context:    "bar",
				},
			}
			return config
		}(config.Spec),
		defaultCluster: nil,
		values:         nil,
		wantErr:        false,
		want: func(t *testing.T, tc TestContext) {
			t.Helper()
			assert.Equal(t, int32(0), tc.Passed())
			assert.Equal(t, int32(0), tc.Failed())
			assert.Equal(t, int32(0), tc.Skipped())
			assert.NotNil(t, tc.Bindings())
			assert.Nil(t, tc.Catch())
			assert.NotNil(t, tc.Compilers())
			assert.NotNil(t, tc.Clusters())
			assert.NotNil(t, tc.Clusters().Lookup("foo"))
			assert.Nil(t, tc.CurrentCluster())
			assert.Nil(t, tc.DelayBeforeCleanup())
			assert.Equal(t, metav1.DeletePropagationBackground, tc.DeletionPropagation())
			assert.False(t, tc.DryRun())
			assert.False(t, tc.FailFast())
			assert.False(t, tc.FullName())
			assert.False(t, tc.SkipDelete())
			assert.True(t, tc.Templating())
			assert.Nil(t, tc.TerminationGrace())
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InitContext(tt.config, tt.defaultCluster, tt.values)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if tt.want != nil {
				tt.want(t, got)
			}
		})
	}
}
