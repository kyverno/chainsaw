package test

import (
	"path/filepath"
	"testing"

	v1alpha1 "github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateFileRefOrResource(t *testing.T) {
	pod := &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata": map[string]any{
				"name": "example-pod",
			},
			"spec": map[string]any{
				"containers": []any{
					map[string]any{
						"name":  "nginx",
						"image": "nginx:latest",
					},
				},
			},
		},
	}
	tests := []struct {
		name      string
		input     v1alpha1.ActionResourceRef
		expectErr bool
		errMsg    string
	}{{
		name: "Both File and Resource are empty",
		input: v1alpha1.ActionResourceRef{
			FileRef: v1alpha1.FileRef{
				File: "",
			},
			Resource: nil,
		},
		expectErr: true,
		errMsg:    "a file reference or raw resource must be specified",
	}, {
		name: "Both File and Resource are provided",
		input: v1alpha1.ActionResourceRef{
			FileRef: v1alpha1.FileRef{
				File: "file",
			},
			Resource: pod,
		},
		expectErr: true,
		errMsg:    "a file reference or raw resource must be specified (found both)",
	}, {
		name: "Only File is provided",
		input: v1alpha1.ActionResourceRef{
			FileRef: v1alpha1.FileRef{
				File: filepath.Join("..", "..", "testdata", "validation", "example-file.yaml"),
			},
			Resource: nil,
		},
		expectErr: false,
	}, {
		name: "Only Resource is provided",
		input: v1alpha1.ActionResourceRef{
			FileRef: v1alpha1.FileRef{
				File: "",
			},
			Resource: pod,
		},
		expectErr: false,
	}}
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
