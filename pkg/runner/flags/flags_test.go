package flags

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha2"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
)

func TestGetFlags(t *testing.T) {
	tests := []struct {
		name   string
		config v1alpha2.ConfigurationSpec
		want   map[string]string
	}{{
		name:   "default",
		config: v1alpha2.ConfigurationSpec{},
		want: map[string]string{
			"test.v":            "true",
			"test.paniconexit0": "true",
			"test.fullpath":     "false",
			"test.run":          "",
			"test.skip":         "",
		},
	}, {
		name: "include",
		config: v1alpha2.ConfigurationSpec{
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
		config: v1alpha2.ConfigurationSpec{
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
		config: v1alpha2.ConfigurationSpec{
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
		config: v1alpha2.ConfigurationSpec{
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
			got := GetFlags(tt.config)
			assert.Equal(t, tt.want, got)
		})
	}
}
