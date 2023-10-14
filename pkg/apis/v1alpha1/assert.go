package v1alpha1

// Error represents an assertion expected to succeed.
type Assert struct {
	// Assertion manifest.
	FileRef `json:",inline"`
}
