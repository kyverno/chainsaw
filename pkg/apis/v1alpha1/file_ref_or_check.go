package v1alpha1

// FileRefOrCheck represents a file reference or resource.
type FileRefOrCheck struct {
	// FileRef provides a reference to the file containing the resources to be applied.
	// +optional
	FileRef `json:",inline"`

	// Check provides a check used in assertions.
	// +optional
	Check *Check `json:"resource,omitempty"`
}
