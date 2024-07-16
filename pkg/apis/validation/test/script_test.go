package test

import (
	"testing"

	v1alpha1 "github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateScript(t *testing.T) {
	// Test cases
	tests := []struct {
		name      string
		input     *v1alpha1.Script
		expectErr bool
		errMsg    string
	}{
		{
			name: "Content is empty",
			input: &v1alpha1.Script{
				Content: "",
			},
			expectErr: true,
			errMsg:    "content must be specified",
		},
		{
			name: "Content is provided",
			input: &v1alpha1.Script{
				Content: "echo Hello, World!",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateScript(field.NewPath("testPath"), tt.input)
			if tt.expectErr {
				assert.NotEmpty(t, errs)
				assert.Contains(t, errs.ToAggregate().Error(), tt.errMsg)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}
