package v1alpha1

// Expectation represents a check to be applied on the result of an operation
// with a match filter to determine if the verification should be considered.
type Expectation struct {
	// Match defines the matching statement.
	// +optional
	Match *Check `json:"match,omitempty"`

	// Match defines the matching statement.
	Check Check `json:"check"`
}
