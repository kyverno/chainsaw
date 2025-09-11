package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestSkipField(t *testing.T) {
	// Create a test resource
	testResource := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata": map[string]interface{}{
			"name": "test-configmap",
		},
		"data": map[string]interface{}{
			"key": "value",
		},
	}

	// Convert to unstructured
	unstructuredObj := &unstructured.Unstructured{
		Object: testResource,
	}

	// Test with skip field set to true
	step := TestStep{
		Name: "Test Step",
		TestStepSpec: TestStepSpec{
			Skip: stringPtr("true"),
			Try: []Operation{
				{
					Apply: &Apply{
						ActionResourceRef: ActionResourceRef{
							Resource: unstructuredObj,
						},
					},
				},
			},
		},
	}

	// Verify the skip field is properly set
	assert.NotNil(t, step.Skip)
	assert.Equal(t, "true", *step.Skip)

	// Test with skip field set to a template expression
	step = TestStep{
		Name: "Test Step with Template",
		TestStepSpec: TestStepSpec{
			Skip: stringPtr("{{ .SKIP_STEP }}"),
			Try: []Operation{
				{
					Apply: &Apply{
						ActionResourceRef: ActionResourceRef{
							Resource: unstructuredObj,
						},
					},
				},
			},
		},
	}

	// Verify the skip field is properly set
	assert.NotNil(t, step.Skip)
	assert.Equal(t, "{{ .SKIP_STEP }}", *step.Skip)
}

// Helper function to create a string pointer
func stringPtr(s string) *string {
	return &s
}
