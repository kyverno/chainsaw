package client

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/json"
)

func PatchObject(actual, expected runtime.Object) ([]byte, error) {
	actualMeta, err := meta.Accessor(actual)
	if err != nil {
		return nil, err
	}
	expectedMeta, err := meta.Accessor(expected.DeepCopyObject())
	if err != nil {
		return nil, err
	}
	expectedMeta.SetResourceVersion(actualMeta.GetResourceVersion())
	return json.Marshal(expectedMeta)
}
