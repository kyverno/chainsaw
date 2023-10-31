package runner

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_timeout(t *testing.T) {
	config := metav1.Duration{Duration: 1 * time.Minute}
	test := metav1.Duration{Duration: 2 * time.Minute}
	step := metav1.Duration{Duration: 3 * time.Minute}
	operation := metav1.Duration{Duration: 4 * time.Minute}
	fallback := 10 * time.Second
	tests := []struct {
		name      string
		fallback  time.Duration
		config    *metav1.Duration
		test      *metav1.Duration
		step      *metav1.Duration
		operation *metav1.Duration
		want      time.Duration
	}{{
		name:      "none",
		fallback:  fallback,
		config:    nil,
		test:      nil,
		step:      nil,
		operation: nil,
		want:      fallback,
	}, {
		name:      "from config",
		fallback:  fallback,
		config:    &config,
		test:      nil,
		step:      nil,
		operation: nil,
		want:      config.Duration,
	}, {
		name:      "from test",
		fallback:  fallback,
		config:    &config,
		test:      &test,
		step:      nil,
		operation: nil,
		want:      test.Duration,
	}, {
		name:      "from step",
		fallback:  fallback,
		config:    &config,
		test:      &test,
		step:      &step,
		operation: nil,
		want:      step.Duration,
	}, {
		name:      "from operation",
		fallback:  fallback,
		config:    &config,
		test:      &test,
		step:      &step,
		operation: &operation,
		want:      operation.Duration,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := timeout(tt.fallback, tt.config, tt.test, tt.step, tt.operation)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_timeoutCtx(t *testing.T) {
	config := metav1.Duration{Duration: 1 * time.Minute}
	test := metav1.Duration{Duration: 2 * time.Minute}
	step := metav1.Duration{Duration: 3 * time.Minute}
	operation := metav1.Duration{Duration: 4 * time.Minute}
	fallback := 10 * time.Second
	tests := []struct {
		name      string
		fallback  time.Duration
		config    *metav1.Duration
		test      *metav1.Duration
		step      *metav1.Duration
		operation *metav1.Duration
	}{{
		name:      "none",
		fallback:  fallback,
		config:    nil,
		test:      nil,
		step:      nil,
		operation: nil,
	}, {
		name:      "from config",
		fallback:  fallback,
		config:    &config,
		test:      nil,
		step:      nil,
		operation: nil,
	}, {
		name:      "from test",
		fallback:  fallback,
		config:    &config,
		test:      &test,
		step:      nil,
		operation: nil,
	}, {
		name:      "from step",
		fallback:  fallback,
		config:    &config,
		test:      &test,
		step:      &step,
		operation: nil,
	}, {
		name:      "from operation",
		fallback:  fallback,
		config:    &config,
		test:      &test,
		step:      &step,
		operation: &operation,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, cancel := timeoutCtx(tt.fallback, tt.config, tt.test, tt.step, tt.operation)
			defer cancel()
			assert.NotNil(t, got)
			assert.NotNil(t, cancel)
			_, ok := got.Deadline()
			assert.Equal(t, true, ok)
		})
	}
}
