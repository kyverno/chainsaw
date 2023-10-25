package runner

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

func skipDelete(config v1alpha1.ConfigurationSpec, test v1alpha1.TestSpec, step v1alpha1.TestStepSpec) *bool {
	if step.SkipDelete != nil {
		return step.SkipDelete
	}
	if test.SkipDelete != nil {
		return test.SkipDelete
	}
	return &config.SkipDelete
}
