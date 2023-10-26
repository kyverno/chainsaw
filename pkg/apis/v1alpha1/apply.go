package v1alpha1

// Apply represents a set of configurations or resources that
// should be applied during testing.
type Apply struct {
	// FileRef provides a reference to the file containing the
	FileRef `json:",inline"`

	// ContinueOnError determines whether a test should continue or not in case the operation was not successful.
	// Even if the test continues executing, it will still be reported as failed.
	// +optional
	ContinueOnError *bool `json:"continueOnError,omitempty"`
}
