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
	Delete []ObjectReference `json:"delete,omitempty"`
}

type ObjectReference struct {
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// +optional
	v1.ObjectReference `json:",inline"`
}
