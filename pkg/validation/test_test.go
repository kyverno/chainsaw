package validation

import (
	"testing"

	v1alpha1 "github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestValidateTest(t *testing.T) {
	validTestSpec := v1alpha1.TestSpec{
		Steps: []v1alpha1.TestSpecStep{
			{
				Name: "step1",
				TestStepSpec: v1alpha1.TestStepSpec{
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
				},
			},
		},
	}
	invalidTestSpec := v1alpha1.TestSpec{
		Steps: []v1alpha1.TestSpecStep{
			{
				Name: "step1",
				TestStepSpec: v1alpha1.TestStepSpec{
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
				},
			},
		},
	}

	tests := []struct {
		name      string
		input     *v1alpha1.Test
		expectErr bool
	}{
		{
			name: "Valid TestSpec",
			input: &v1alpha1.Test{
				Spec: validTestSpec,
			},
			expectErr: false,
		},
		{
			name: "Invalid TestSpec",
			input: &v1alpha1.Test{
				Spec: invalidTestSpec,
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateTest(tt.input)
			if tt.expectErr {
				assert.NotEmpty(t, errs)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}
