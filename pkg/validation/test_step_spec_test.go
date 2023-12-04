package validation

import (
	"path/filepath"
	"testing"

	v1alpha1 "github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateTestStepSpec(t *testing.T) {
	validTry := []v1alpha1.Operation{{
		Apply: &v1alpha1.Apply{
			FileRefOrResource: v1alpha1.FileRefOrResource{
				FileRef: v1alpha1.FileRef{
					File: filepath.Join("..", "..", "testdata", "validation", "example-file.yaml"),
				},
			},
		},
	}}
	invalidTry := []v1alpha1.Operation{
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
	}
	validCatch := []v1alpha1.Catch{
		{
			Command: &v1alpha1.Command{
				Entrypoint: "echo",
				Args:       []string{"Hello, World!"},
			},
		},
	}
	invalidCatch := []v1alpha1.Catch{
		{
			Script: &v1alpha1.Script{
				Content: "echo Hello, World!",
			},
			Command: &v1alpha1.Command{
				Entrypoint: "echo",
				Args:       []string{"Hello, World!"},
			},
		},
	}

	validFinally := []v1alpha1.Finally{
		{
			Command: &v1alpha1.Command{
				Entrypoint: "echo",
				Args:       []string{"Hello, World!"},
			},
		},
	}
	invalidFinally := []v1alpha1.Finally{
		{
			Script: &v1alpha1.Script{
				Content: "echo Hello, World!",
			},
			Command: &v1alpha1.Command{
				Entrypoint: "echo",
				Args:       []string{"Hello, World!"},
			},
		},
	}

	tests := []struct {
		name      string
		input     v1alpha1.TestStepSpec
		expectErr bool
	}{
		{
			name: "Valid TestStepSpec",
			input: v1alpha1.TestStepSpec{
				Try:     validTry,
				Catch:   validCatch,
				Finally: validFinally,
			},
			expectErr: false,
		},
		{
			name: "Invalid Try in TestStepSpec",
			input: v1alpha1.TestStepSpec{
				Try:     invalidTry,
				Catch:   validCatch,
				Finally: validFinally,
			},
			expectErr: true,
		},
		{
			name: "Invalid Catch in TestStepSpec",
			input: v1alpha1.TestStepSpec{
				Try:     validTry,
				Catch:   invalidCatch,
				Finally: validFinally,
			},
			expectErr: true,
		},
		{
			name: "Invalid Finally in TestStepSpec",
			input: v1alpha1.TestStepSpec{
				Try:     validTry,
				Catch:   validCatch,
				Finally: invalidFinally,
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateTestStepSpec(field.NewPath("testPath"), tt.input)
			if tt.expectErr {
				assert.NotEmpty(t, errs)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}
