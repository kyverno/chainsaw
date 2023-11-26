package client

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func Namespace(name string) corev1.Namespace {
	return corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "Namespace",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

func PetNamespace() corev1.Namespace {
	return Namespace(Pet())
}

func ToUnstructured(obj interface{}) unstructured.Unstructured {
	data, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		panic(err)
	}
	return unstructured.Unstructured{Object: data}
}
