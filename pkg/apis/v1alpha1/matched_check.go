package v1alpha1

// MatchedCheck represents a check be applied on the result of an operation
// with a match filter to determine if the verification should be considered.
type MatchedCheck struct {
	// Match defines the matching statement.
	Match *Check `json:"match,omitempty"`

	// Match defines the matching statement.
	Verify *Check `json:"verify,omitempty"`
}
