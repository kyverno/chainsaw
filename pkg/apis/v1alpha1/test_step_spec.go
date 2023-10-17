package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
)

type TestStepSpec struct {
	// +optional
	Assert []Assert `json:"assert,omitempty"`
	// +optional
	Apply []Apply `json:"apply,omitempty"`
	// +optional
	Error []Error `json:"error,omitempty"`

	// List of object that should be deleted before the test step is executed
	// +optional
	Delete []v1.ObjectReference `json:"delete,omitempty"`
}
