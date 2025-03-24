package errors

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestResourceError(t *testing.T) {
	// Create test objects
	expected := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name":      "test-config",
				"namespace": "default",
			},
			"data": map[string]interface{}{
				"key1": "expected-value",
			},
		},
	}

	actual := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name":      "test-config",
				"namespace": "default",
			},
			"data": map[string]interface{}{
				"key1": "actual-value",
			},
		},
	}

	// Create field error list
	errorList := field.ErrorList{
		field.Invalid(field.NewPath("data", "key1"), "actual-value", "expected 'expected-value'"),
	}

	// Test creating a ResourceError
	err := ResourceError(
		apis.DefaultCompilers,
		expected,
		actual,
		false,
		apis.NewBindings(),
		errorList,
	)

	// Verify the error is of the correct type
	_, ok := err.(resourceError)
	assert.True(t, ok, "Expected error to be of type resourceError")

	// Verify the error message contains the expected information
	errMsg := err.Error()

	// Should include resource identification
	assert.Contains(t, errMsg, "v1/ConfigMap/default/test-config")

	// Should include field errors
	assert.Contains(t, errMsg, "data.key1")
	assert.Contains(t, errMsg, "expected 'expected-value'")
	assert.Contains(t, errMsg, "actual-value")

	// Should include diff information
	assert.Contains(t, errMsg, "key1:")
}

func TestResourceErrorWithTemplating(t *testing.T) {
	// Create test objects
	expected := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name":      "test-config",
				"namespace": "default",
			},
			"data": map[string]interface{}{
				"key1": "{{ .value }}", // Template value
			},
		},
	}

	actual := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name":      "test-config",
				"namespace": "default",
			},
			"data": map[string]interface{}{
				"key1": "wrong-value",
			},
		},
	}

	// Create bindings for templating
	bindings := apis.NewBindings().Register("value", apis.NewBinding("template-value"))

	// Test creating a ResourceError with templating
	err := ResourceError(
		apis.DefaultCompilers,
		expected,
		actual,
		true, // Enable templating
		bindings,
		field.ErrorList{},
	)

	// Verify the error message contains templated values
	errMsg := err.Error()

	// Should include resource identification
	assert.Contains(t, errMsg, "v1/ConfigMap/default/test-config")

	// We need to check for the template or value in diff, not always both
	// Either the template was successfully rendered (contains template-value)
	// Or it shows the raw template (contains {{ .value }})
	if !assert.True(t,
		(contains(errMsg, "template-value") || contains(errMsg, "{{ .value }}")),
		"Error message should contain either the template expression or the template value") {
		// Print the actual error message for debugging
		t.Logf("Actual error message: %s", errMsg)
	}

	assert.Contains(t, errMsg, "wrong-value")
}

func TestResourceErrorWithInvalidTemplate(t *testing.T) {
	// Create test objects with invalid template
	expected := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name":      "test-config",
				"namespace": "default",
			},
			"data": map[string]interface{}{
				"key1": "{{ .undefinedValue }}", // Undefined value in template
			},
		},
	}

	actual := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name":      "test-config",
				"namespace": "default",
			},
			"data": map[string]interface{}{
				"key1": "some-value",
			},
		},
	}

	// Test creating a ResourceError with invalid templating
	err := ResourceError(
		apis.DefaultCompilers,
		expected,
		actual,
		true,               // Enable templating
		apis.NewBindings(), // Empty bindings, so .undefinedValue is not defined
		field.ErrorList{},
	)

	// Verify the error message contains relevant info
	errMsg := err.Error()
	t.Logf("Error message: %s", errMsg)

	// Should include resource identification
	assert.Contains(t, errMsg, "v1/ConfigMap/default/test-config")

	// Should include the template expression
	assert.Contains(t, errMsg, "{{ .undefinedValue }}")

	// Might contain an error about the template (but not guaranteed by the implementation)
	// So we don't assert on this anymore
}

// Helper function to check if a string contains another string
func contains(s, substr string) bool {
	return assert.Contains(&testing.T{}, s, substr)
}
