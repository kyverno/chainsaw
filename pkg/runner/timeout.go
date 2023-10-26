package runner

import (
	"context"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

func timeout(config v1alpha1.ConfigurationSpec, test v1alpha1.TestSpec, step v1alpha1.TestStepSpec) *time.Duration {
	if step.Timeout != nil {
		return &step.Timeout.Duration
	}
	if test.Timeout != nil {
		return &test.Timeout.Duration
	}
	if config.Timeout != nil {
		return &config.Timeout.Duration
	}
	return nil
}

func cancelNoOp() {}

func timeoutCtx(config v1alpha1.ConfigurationSpec, test v1alpha1.TestSpec, step v1alpha1.TestStepSpec) (context.Context, context.CancelFunc) {
	ctx := context.Background()
	if timeout := timeout(config, test, step); timeout != nil {
		return context.WithTimeout(ctx, *timeout)
	}
	return ctx, cancelNoOp
}

func timeoutExecCtx(exec v1alpha1.Exec, config v1alpha1.ConfigurationSpec, test v1alpha1.TestSpec, step v1alpha1.TestStepSpec) (context.Context, context.CancelFunc) {
	if exec.Timeout != nil && exec.Timeout.Abs() > 0 {
		return context.WithTimeout(context.Background(), exec.Timeout.Duration)
	} else {
		return timeoutCtx(config, test, step)
	}
}
