package v1alpha1

import corev1 "k8s.io/api/core/v1"

// ObjectReference is a reference to an object that should be deleted
type ObjectReference struct {
	// Label selector to match objects to delete
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// Object reference to delete
	// +optional
	corev1.ObjectReference `json:",inline"`
}
