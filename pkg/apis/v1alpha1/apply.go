package v1alpha1

// Apply represents a set of configurations or resources that
// should be applied during testing.
type Apply struct {
	// FileRefOrResource provides a reference to the file containing the resources to be applied.
	FileRefOrResource `json:",inline"`

	// ShouldFail determines whether applying the file is expected to fail.
	// +optional
	ShouldFail *bool `json:"shouldFail,omitempty"`
}
