package kube

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestNamespace(t *testing.T) {
	expectedName := "testNamespace"
	ns := Namespace(expectedName)
	assert.Equal(t, corev1.SchemeGroupVersion.String(), ns.APIVersion)
	assert.Equal(t, "Namespace", ns.Kind)
	assert.Equal(t, expectedName, ns.Name)
}
