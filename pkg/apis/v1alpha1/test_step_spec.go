package v1alpha1

type TestStepSpec struct {
	// +optional
	Assert []Assert `json:"assert,omitempty"`
	// +optional
	Apply []Apply `json:"apply,omitempty"`
	// +optional
	Error []Error `json:"error,omitempty"`

	// Delete provides a list of objects that should be deleted before this test step is executed.
	// This helps in ensuring that the environment is set up correctly before the test step runs.
	// +optional
	Delete []ObjectReference `json:"delete,omitempty"`
}
