package v1alpha1

// Apply represents a set of configurations or resources that
// should be applied during testing.
type Apply struct {
	// FileRef provides a reference to the file containing the
	FileRef `json:",inline"`
}
