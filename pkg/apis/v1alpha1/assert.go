package v1alpha1

// Assert represents a test condition that is expected to hold true
// during the testing process.
type Assert struct {
	// FileRef provides a reference to the file containing the assertion.
	FileRef `json:",inline"`
}
