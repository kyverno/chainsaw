package client

import (
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	apijson "k8s.io/apimachinery/pkg/util/json"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func ObjectKey(obj metav1.Object) ctrlclient.ObjectKey {
	return ctrlclient.ObjectKey{
		Name:      obj.GetName(),
		Namespace: obj.GetNamespace(),
	}
}

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
	return apijson.Marshal(expectedMeta)
}
