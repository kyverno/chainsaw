package timeout

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGet(t *testing.T) {
	operation := metav1.Duration{Duration: 4 * time.Minute}
	fallback := 10 * time.Second
	tests := []struct {
		name      string
		fallback  time.Duration
		operation *metav1.Duration
		want      time.Duration
	}{{
		name:      "fallback",
		fallback:  fallback,
		operation: nil,
		want:      fallback,
	}, {
		name:      "operation",
		fallback:  fallback,
		operation: &operation,
		want:      operation.Duration,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Get(tt.operation, tt.fallback)
			assert.Equal(t, tt.want, got)
		})
	}
}
