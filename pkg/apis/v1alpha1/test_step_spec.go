package v1alpha1

// TestStepSpec defines the desired state and behavior for each test step.
type TestStepSpec struct {
	// Timeouts for the test step. Overrides the global timeouts set in the Configuration and the timeouts eventually set in the Test.
	// +optional
	Timeouts Timeouts `json:"timeouts"`

	// SkipDelete determines whether the resources created by the step should be deleted after the test step is executed.
	// +optional
	SkipDelete *bool `json:"skipDelete,omitempty"`

	// Operations defines the order in which the test step should be executed.
	// +optional
	Operations Operations `json:"ordering,omitempty"`

	// OnFailure defines actions to be executed in case of step failure.
	// +optional
	OnFailure []OnFailure `json:"onFailure,omitempty"`
}
