package v1alpha1

// Create represents a set of resources that should be created.
// If a resource already exists in the cluster it will fail.
type Create struct {
	// FileRef provides a reference to the file containing the
	FileRef `json:",inline"`

	// ShouldFail determines whether applying the file is expected to fail.
	// +optional
	ShouldFail *bool `json:"shouldFail,omitempty"`
}
