package validation

import (
	"testing"

	v1alpha1 "github.com/kyverno/chainsaw/pkg/apis/v1alpha1"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateFileRefOrResource(t *testing.T) {
	pod := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata": map[string]interface{}{
				"name": "example-pod",
			},
			"spec": map[string]interface{}{
				"containers": []interface{}{
					map[string]interface{}{
						"name":  "nginx",
						"image": "nginx:latest",
					},
				},
			},
		},
	}
	tests := []struct {
		name      string
		input     v1alpha1.FileRefOrResource
		expectErr bool
		errMsg    string
	}{
		{
			name: "Both File and Resource are empty",
			input: v1alpha1.FileRefOrResource{
				FileRef: v1alpha1.FileRef{
					File: "",
				},
				Resource: nil,
			},
			expectErr: true,
			errMsg:    "a file reference or raw resource must be specified",
		},
		{
			name: "Both File and Resource are provided",
			input: v1alpha1.FileRefOrResource{
				FileRef: v1alpha1.FileRef{
					File: "file",
				},
				Resource: pod,
			},
			expectErr: true,
			errMsg:    "a file reference or raw resource must be specified (found both)",
		},
		{
			name: "Only File is provided",
			input: v1alpha1.FileRefOrResource{
				FileRef: v1alpha1.FileRef{
					File: "file",
				},
				Resource: nil,
			},
			expectErr: false,
		},
		{
			name: "Only Resource is provided",
			input: v1alpha1.FileRefOrResource{
				FileRef: v1alpha1.FileRef{
					File: "",
				},
				Resource: pod, // Replace SomeResourceType with the actual type
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateFileRefOrResource(field.NewPath("testPath"), tt.input)
			if tt.expectErr {
				assert.NotEmpty(t, errs)
				assert.Contains(t, errs.ToAggregate().Error(), tt.errMsg)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}
