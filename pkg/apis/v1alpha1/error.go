package v1alpha1

// Error represents an anticipated error condition that may arise during testing.
// Instead of treating such an error as a test failure, it acknowledges it as expected.
type Error struct {
	// FileRef provides a reference to the file containing the expected error.
	FileRef `json:",inline"`
}
