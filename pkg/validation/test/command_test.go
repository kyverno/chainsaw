package test

import (
	"testing"

	v1alpha1 "github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateCommand(t *testing.T) {
	tests := []struct {
		name      string
		input     *v1alpha1.Command
		expectErr bool
		errMsg    string
	}{
		{
			name: "Entrypoint is empty",
			input: &v1alpha1.Command{
				Entrypoint: "",
			},
			expectErr: true,
			errMsg:    "entrypoint must be specified",
		},
		{
			name: "Entrypoint is provided",
			input: &v1alpha1.Command{
				Entrypoint: "echo Hello, World!",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateCommand(field.NewPath("testPath"), tt.input)
			if tt.expectErr {
				assert.NotEmpty(t, errs)
				assert.Contains(t, errs.ToAggregate().Error(), tt.errMsg)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}
