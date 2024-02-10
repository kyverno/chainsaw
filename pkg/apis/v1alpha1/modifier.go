package v1alpha1

// Modifier represents an object mutation.
type Modifier struct {
	// Match defines the matching statement.
	// +optional
	Match *Match `json:"match,omitempty"`

	// Annotate defines a mutation of object annotations.
	// +optional
	Annotate *Any `json:"annotate,omitempty"`

	// Label defines a mutation of object labels.
	// +optional
	Label *Any `json:"label,omitempty"`

	// Merge defines an arbitrary merge mutation.
	// +optional
	Merge *Any `json:"merge,omitempty"`
}
