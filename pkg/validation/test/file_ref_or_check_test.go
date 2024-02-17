package test

import (
	"path/filepath"
	"testing"

	v1alpha1 "github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateFileRefOrCheck(t *testing.T) {
	check := &v1alpha1.Check{
		Value: 42,
	}
	tests := []struct {
		name      string
		input     v1alpha1.FileRefOrCheck
		expectErr bool
		errMsg    string
	}{{
		name: "Both File and Resource are empty",
		input: v1alpha1.FileRefOrCheck{
			FileRef: v1alpha1.FileRef{
				File: "",
			},
			Check: nil,
		},
		expectErr: true,
		errMsg:    "a file reference or raw check must be specified",
	}, {
		name: "Both File and Resource are provided",
		input: v1alpha1.FileRefOrCheck{
			FileRef: v1alpha1.FileRef{
				File: "file",
			},
			Check: check,
		},
		expectErr: true,
		errMsg:    "a file reference or raw check must be specified (found both)",
	}, {
		name: "Only File is provided",
		input: v1alpha1.FileRefOrCheck{
			FileRef: v1alpha1.FileRef{
				File: filepath.Join("..", "..", "testdata", "validation", "example-file.yaml"),
			},
			Check: nil,
		},
		expectErr: false,
	}, {
		name: "Only Resource is provided",
		input: v1alpha1.FileRefOrCheck{
			FileRef: v1alpha1.FileRef{
				File: "",
			},
			Check: check,
		},
		expectErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateFileRefOrCheck(field.NewPath("testPath"), tt.input)
			if tt.expectErr {
				assert.NotEmpty(t, errs)
				assert.Contains(t, errs.ToAggregate().Error(), tt.errMsg)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}
