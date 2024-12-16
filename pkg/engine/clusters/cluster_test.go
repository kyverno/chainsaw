package clusters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/rest"
)

func TestNewClusterFromConfig(t *testing.T) {
	tests := []struct {
		name   string
		config *rest.Config
		want   Cluster
	}{{
		name: "nil",
		want: &fromConfig{
			config: nil,
		},
	}, {
		name:   "nil",
		config: &rest.Config{},
		want: &fromConfig{
			config: &rest.Config{},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewClusterFromConfig(tt.config)
			assert.Equal(t, tt.want, got)
			config, err := got.Config()
			assert.NoError(t, err)
			assert.Same(t, tt.config, config)
		})
	}
}

func TestNewClusterFromKubeconfig(t *testing.T) {
	tests := []struct {
		name       string
		kubeconfig string
		context    string
	}{{
		name: "none",
	}, {
		name:       "foo",
		kubeconfig: "foo",
		context:    "bar",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewClusterFromKubeconfig(tt.kubeconfig, tt.context)
			assert.NotNil(t, got)
			config, err := got.Config()
			assert.Error(t, err)
			assert.Nil(t, config)
		})
	}
}
