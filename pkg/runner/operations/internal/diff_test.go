package internal

import (
	"fmt"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func mockCtrlClientObject(data map[string]interface{}) ctrlclient.Object {
	return &unstructured.Unstructured{Object: data}
}

func TestDiff(t *testing.T) {
	tests := []struct {
		name         string
		expectedData map[string]interface{}
		actualData   map[string]interface{}
		expectedDiff string
	}{
		{
			name: "Test simple diff",
			expectedData: map[string]interface{}{
				"kind":       "Pod",
				"apiVersion": "v1",
				"metadata": map[string]interface{}{
					"name": "test-pod",
				},
				"spec": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{
							"name":  "test-container",
							"image": "test-image:v1",
						},
					},
				},
			},
			actualData: map[string]interface{}{
				"kind":       "Pod",
				"apiVersion": "v1",
				"metadata": map[string]interface{}{
					"name": "test-pod",
				},
				"spec": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{
							"name":  "test-container",
							"image": "test-image:v2",
						},
					},
				},
			},
			expectedDiff: `--- Expected
+++ Actual
@@ -4,6 +4,6 @@
   name: test-pod
 spec:
   containers:
-  - image: test-image:v1
+  - image: test-image:v2
     name: test-container`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := mockCtrlClientObject(tt.expectedData)
			actual := mockCtrlClientObject(tt.actualData)

			gotDiff, err := diff(expected, actual)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assert.EqualValues(t, tt.expectedDiff, gotDiff)
		})
	}
}

func Test_diffHelper_Error(t *testing.T) {
	mockData := &unstructured.Unstructured{}

	testCases := []struct {
		name              string
		expected          interface{}
		actual            interface{}
		expectedMarshaler yamlMarshaler
		actualMarshaler   yamlMarshaler
		getDiffString     diffLibInterface
	}{
		{
			name:     "error in expected marshaler",
			expected: mockData,
			actual:   mockData,
			expectedMarshaler: func(obj interface{}) ([]byte, error) {
				return nil, fmt.Errorf("expected marshal error")
			},
		},
		{
			name:     "error in actual marshaler",
			expected: mockData,
			actual:   mockData,
			actualMarshaler: func(obj interface{}) ([]byte, error) {
				return nil, fmt.Errorf("actual marshal error")
			},
		},
		{
			name:          "error in getting diff string",
			expected:      mockData,
			actual:        mockData,
			getDiffString: func(diff difflib.UnifiedDiff) (string, error) { return "", fmt.Errorf("diff string error") },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := diffHelper(tc.expected, tc.actual, tc.expectedMarshaler, tc.actualMarshaler, tc.getDiffString)
			assert.Error(t, err)
		})
	}
}
