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
	background := context.Background()
	if timeout := timeout(config, test, step); timeout != nil {
		return context.WithTimeout(background, *timeout)
	}
	return background, cancelNoOp
}

func cmdtimeoutCtx(cmd v1alpha1.Command, config v1alpha1.ConfigurationSpec, test v1alpha1.TestSpec, step v1alpha1.TestStepSpec) (context.Context, context.CancelFunc) {
	if cmd.Timeout != nil && cmd.Timeout.Abs() > 0 {
		return context.WithTimeout(context.Background(), cmd.Timeout.Duration)
	} else {
		return timeoutCtx(config, test, step)
	}
}
