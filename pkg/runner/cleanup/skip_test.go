package cleanup

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
)

func TestSkip(t *testing.T) {
	tests := []struct {
		name   string
		config bool
		test   *bool
		step   *bool
		want   bool
	}{{
		name:   "from config",
		config: true,
		test:   nil,
		step:   nil,
		want:   true,
	}, {
		name:   "from test",
		config: false,
		test:   ptr.To(true),
		step:   nil,
		want:   true,
	}, {
		name:   "from test",
		config: true,
		test:   ptr.To(false),
		step:   nil,
		want:   false,
	}, {
		name:   "from step",
		config: false,
		test:   ptr.To(false),
		step:   ptr.To(true),
		want:   true,
	}, {
		name:   "from step",
		config: true,
		test:   ptr.To(true),
		step:   ptr.To(false),
		want:   false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Skip(tt.config, tt.test, tt.step)
			assert.Equal(t, tt.want, got)
		})
	}
}
