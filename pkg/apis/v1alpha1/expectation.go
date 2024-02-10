package v1alpha1

// Expectation represents a check to be applied on the result of an operation
// with a match filter to determine if the verification should be considered.
type Expectation struct {
	// Match defines the matching statement.
	// +optional
	Match *Match `json:"match,omitempty"`

	// Check defines the verification statement.
	Check Check `json:"check"`
}
