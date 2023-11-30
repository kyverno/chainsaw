package v1alpha1

// TestSpec contains the test spec.
type TestSpec struct {
	// Description contains a description of the test.
	// +optional
	Description string `json:"description,omitempty"`

	// Timeouts for the test. Overrides the global timeouts set in the Configuration on a per operation basis.
	// +optional
	Timeouts *Timeouts `json:"timeouts,omitempty"`

	// Skip determines whether the test should skipped.
	// +optional
	Skip *bool `json:"skip,omitempty"`

	// Concurrent determines whether the test should run concurrently with other tests.
	// +optional
	Concurrent *bool `json:"concurrent,omitempty"`

	// SkipDelete determines whether the resources created by the test should be deleted after the test is executed.
	// +optional
	SkipDelete *bool `json:"skipDelete,omitempty"`

	// Namespace determines whether the test should run in a random ephemeral namespace or not.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Steps defining the test.
	Steps []TestSpecStep `json:"steps"`
}
