package flags

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha2"
	"github.com/kyverno/chainsaw/pkg/loaders/config"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
)

func Test_getFlags(t *testing.T) {
	tests := []struct {
		name   string
		config model.Configuration
		want   map[string]string
	}{{
		name:   "default",
		config: model.Configuration{},
		want: map[string]string{
			"test.v":            "true",
			"test.paniconexit0": "true",
			"test.fullpath":     "false",
			"test.run":          "",
			"test.skip":         "",
		},
	}, {
		name: "include",
		config: model.Configuration{
			Discovery: v1alpha2.DiscoveryOptions{
				IncludeTestRegex: "^.*$",
			},
		},
		want: map[string]string{
			"test.v":            "true",
			"test.paniconexit0": "true",
			"test.fullpath":     "false",
			"test.run":          "^.*$",
			"test.skip":         "",
		},
	}, {
		name: "exclude",
		config: model.Configuration{
			Discovery: v1alpha2.DiscoveryOptions{
				ExcludeTestRegex: "^.*$",
			},
		},
		want: map[string]string{
			"test.v":            "true",
			"test.paniconexit0": "true",
			"test.fullpath":     "false",
			"test.run":          "",
			"test.skip":         "^.*$",
		},
	}, {
		name: "parallel",
		config: model.Configuration{
			Execution: v1alpha2.ExecutionOptions{
				Parallel: ptr.To(10),
			},
		},
		want: map[string]string{
			"test.v":            "true",
			"test.paniconexit0": "true",
			"test.fullpath":     "false",
			"test.run":          "",
			"test.skip":         "",
			"test.parallel":     "10",
		},
	}, {
		name: "repeat count",
		config: model.Configuration{
			Execution: v1alpha2.ExecutionOptions{
				RepeatCount: ptr.To(10),
			},
		},
		want: map[string]string{
			"test.v":            "true",
			"test.paniconexit0": "true",
			"test.fullpath":     "false",
			"test.run":          "",
			"test.skip":         "",
			"test.count":        "10",
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getFlags(tt.config)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSetupFlags(t *testing.T) {
	config, err := config.DefaultConfiguration()
	assert.NoError(t, err)
	tests := []struct {
		name    string
		config  model.Configuration
		wantErr bool
	}{{
		name:   "empty",
		config: model.Configuration{},
	}, {
		name:   "default",
		config: config.Spec,
	}, {
		name: "error",
		config: model.Configuration{
			Execution: v1alpha2.ExecutionOptions{
				RepeatCount: ptr.To(-1),
			},
		},
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetupFlags(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
