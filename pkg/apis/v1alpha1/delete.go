package v1alpha1

// Delete is a reference to an object that should be deleted
type Delete struct {
	// ObjectReference determines objects to be deleted.
	ObjectReference `json:",inline"`

	// ContinueOnError determines whether a test should continue or not in case the operation was not successful.
	// Even if the test continues executing, it will still be reported as failed.
	// +optional
	ContinueOnError *bool `json:"continueOnError,omitempty"`
}
