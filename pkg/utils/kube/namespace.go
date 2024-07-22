package kube

import (
	"fmt"

	petname "github.com/dustinkirkland/golang-petname"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	return Namespace(fmt.Sprintf("chainsaw-%s", petname.Generate(2, "-")))
}
