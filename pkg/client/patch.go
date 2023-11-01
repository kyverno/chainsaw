package client

import (
	"errors"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

func PatchObject(actual, expected runtime.Object) (runtime.Object, error) {
	if actual == nil || expected == nil {
		return nil, errors.New("actual and expected objects must not be nil")
	}
	actualMeta, err := meta.Accessor(actual)
	if err != nil {
		return nil, err
	}
	copy := expected.DeepCopyObject()
	expectedMeta, err := meta.Accessor(copy)
	if err != nil {
		return nil, err
	}
	expectedMeta.SetResourceVersion(actualMeta.GetResourceVersion())
	return copy, nil
}
