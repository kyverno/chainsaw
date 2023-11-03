package v1alpha1

// TestStepSpec defines the desired state and behavior for each test step.
type TestStepSpec struct {
	// Timeouts for the test step. Overrides the global timeouts set in the Configuration and the timeouts eventually set in the Test.
	// +optional
	Timeouts Timeouts `json:"timeouts"`

	// SkipDelete determines whether the resources created by the step should be deleted after the test step is executed.
	// +optional
	SkipDelete *bool `json:"skipDelete,omitempty"`

	// Try defines what the step will try to execute.
	Try []Operation `json:"try"`

	// Catch defines what the step will execute when an error happens.
	// +optional
	Catch []Catch `json:"catch,omitempty"`

	// Finally defines what the step will execute after the step is terminated.
	// +optional
	Finally []Finally `json:"finally,omitempty"`
}
