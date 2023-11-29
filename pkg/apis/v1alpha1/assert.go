package v1alpha1

// Assert represents a test condition that is expected to hold true
// during the testing process.
type Assert struct {
	// FileRefOrResource provides a reference to the assertion.
	FileRefOrResource `json:",inline"`
}
