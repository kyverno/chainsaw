package runner

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	defaultApplyTimeout   = 5 * time.Second
	defaultAssertTimeout  = 30 * time.Second
	defaultErrorTimeout   = 30 * time.Second
	defaultDeleteTimeout  = 15 * time.Second
	defaultCleanupTimeout = 30 * time.Second
	defaultExecTimeout    = 5 * time.Second
)

func timeout(fallback time.Duration, config *metav1.Duration, test *metav1.Duration, step *metav1.Duration, operation *metav1.Duration) time.Duration {
	if operation != nil {
		return operation.Duration
	}
	if step != nil {
		return step.Duration
	}
	if test != nil {
		return test.Duration
	}
	if config != nil {
		return config.Duration
	}
	return fallback
}

func timeoutCtx(fallback time.Duration, config *metav1.Duration, test *metav1.Duration, step *metav1.Duration, operation *metav1.Duration) (context.Context, context.CancelFunc) {
	timeout := timeout(fallback, config, test, step, operation)
	return context.WithTimeout(context.Background(), timeout)
}
