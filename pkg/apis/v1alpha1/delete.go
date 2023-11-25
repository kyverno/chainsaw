package v1alpha1

// Delete is a reference to an object that should be deleted
type Delete struct {
	// ObjectReference determines objects to be deleted.
	ObjectReference `json:"ref"`

	// Check is an assertion tree to validate the operation outcome.
	// +optional
	Check *Check `json:"check,omitempty"`
}
