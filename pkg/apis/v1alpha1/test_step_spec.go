package v1alpha1

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
