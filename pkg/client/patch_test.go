package client

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestPatchObject(t *testing.T) {
	actualObj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata": map[string]interface{}{
				"name":            "test-pod",
				"resourceVersion": "12345",
			},
		},
	}

	expectedObj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata": map[string]interface{}{
				"name": "modified-pod",
			},
		},
	}

	patchBytes, err := PatchObject(actualObj, expectedObj)

	assert.Nil(t, err)

	var patchedMap map[string]interface{}
	err = json.Unmarshal(patchBytes, &patchedMap)
	assert.Nil(t, err)

	assert.Equal(t, "12345", patchedMap["metadata"].(map[string]interface{})["resourceVersion"])
	assert.Equal(t, "modified-pod", patchedMap["metadata"].(map[string]interface{})["name"])
}
