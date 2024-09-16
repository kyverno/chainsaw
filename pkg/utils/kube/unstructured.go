package kube

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func ToUnstructured(obj any) unstructured.Unstructured {
	data, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		panic(err)
	}
	return unstructured.Unstructured{Object: data}
}
