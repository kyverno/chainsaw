package validation

import (
	"testing"

	v1alpha1 "github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestValidateTestStep(t *testing.T) {
	validTestStepSpec := v1alpha1.TestStepSpec{
		Try: []v1alpha1.Operation{
			{
				Apply: &v1alpha1.Apply{
					FileRefOrResource: v1alpha1.FileRefOrResource{
						FileRef: v1alpha1.FileRef{
							File: "file",
						},
					},
				},
			},
		},
	}
	invalidTestStepSpec := v1alpha1.TestStepSpec{
		Try: []v1alpha1.Operation{
			{
				Apply: &v1alpha1.Apply{
					FileRefOrResource: v1alpha1.FileRefOrResource{
						FileRef: v1alpha1.FileRef{
							File: "file",
						},
					},
				},
				Assert: &v1alpha1.Assert{
					FileRefOrResource: v1alpha1.FileRefOrResource{
						FileRef: v1alpha1.FileRef{
							File: "file",
						},
					},
				},
			},
		},
	}
	tests := []struct {
		name      string
		input     *v1alpha1.TestStep
		expectErr bool
	}{{
		name: "Valid TestStepSpec",
		input: &v1alpha1.TestStep{
			Spec: validTestStepSpec,
		},
		expectErr: false,
	}, {
		name: "Invalid TestStepSpec",
		input: &v1alpha1.TestStep{
			Spec: invalidTestStepSpec,
		},
		expectErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateTestStep(tt.input)
			if tt.expectErr {
				assert.NotEmpty(t, errs)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}
