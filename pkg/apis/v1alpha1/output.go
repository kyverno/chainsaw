package v1alpha1

// Output represents an output binding with a match to determine if the binding must be considered or not.
type Output struct {
	// Binding determines the binding to create when the match succeeds.
	Binding `json:",inline"`

	// Match defines the matching statement.
	// +optional
	Match *Match `json:"match,omitempty"`
}
