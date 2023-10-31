package runner

import (
	"testing"

	"k8s.io/utils/ptr"
)

func Test_skipDelete(t *testing.T) {
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
			if got := skipDelete(tt.config, tt.test, tt.step); got != tt.want {
				t.Errorf("skipDelete() = %v, want %v", got, tt.want)
			}
		})
	}
}
