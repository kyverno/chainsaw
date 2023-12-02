package validation

import (
	"path/filepath"
	"testing"

	v1alpha1 "github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateFileRef(t *testing.T) {
	tests := []struct {
		name      string
		input     v1alpha1.FileRef
		expectErr bool
		errMsg    string
	}{
		{
			name: "File field is empty",
			input: v1alpha1.FileRef{
				File: "",
			},
			expectErr: true,
			errMsg:    "a file reference must be specified",
		},
		{
			name: "File field is provided",
			input: v1alpha1.FileRef{
				File: filepath.Join("..", "..", "testdata", "validation", "example-file.yaml"),
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateFileRef(field.NewPath("testPath"), tt.input)
			if tt.expectErr {
				assert.NotEmpty(t, errs)
				assert.Contains(t, errs.ToAggregate().Error(), tt.errMsg)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}
