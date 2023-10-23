package runner

import (
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
