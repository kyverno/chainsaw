package test

import (
	"path/filepath"
	"testing"

	v1alpha1 "github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateTestStep(t *testing.T) {
	validTestStepSpec := v1alpha1.TestStepSpec{
		Try: []v1alpha1.Operation{
			{
				Apply: &v1alpha1.Apply{
					ActionResourceRef: v1alpha1.ActionResourceRef{
						FileRef: v1alpha1.FileRef{
							File: filepath.Join("..", "..", "testdata", "validation", "example-file.yaml"),
						},
					},
				},
			},
		},
	}
	tests := []struct {
		name      string
		input     v1alpha1.TestStep
		expectErr bool
	}{{
		name: "no name",
		input: v1alpha1.TestStep{
			Name:         "",
			TestStepSpec: validTestStepSpec,
		},
		expectErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateTestStep(field.NewPath("testPath"), tt.input)
			if tt.expectErr {
				assert.NotEmpty(t, errs)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}
