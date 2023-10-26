package v1alpha1

// Assert represents a test condition that is expected to hold true
// during the testing process.
type Assert struct {
	// FileRef provides a reference to the file containing the assertion.
	FileRef `json:",inline"`

	// ContinueOnError determines whether a test should continue or not in case the operation was not successful.
	// Even if the test continues executing, it will still be reported as failed.
	// +optional
	ContinueOnError *bool `json:"continueOnError,omitempty"`
}
