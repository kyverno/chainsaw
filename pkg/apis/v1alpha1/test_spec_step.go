package v1alpha1

// TestStep contains the test step definition used in a test spec.
type TestStep struct {
	// Name of the step.
	// +optional
	Name string `json:"name,omitempty"`

	// TestStepSpec of the step.
	TestStepSpec `json:",inline"`
}
