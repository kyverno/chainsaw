package v1alpha1

type Modifier struct {
	// Match defines the matching statement.
	// +optional
	Match *Match `json:"match,omitempty"`

	Annotate *Any `json:"annotate,omitempty"`
	Label    *Any `json:"label,omitempty"`
	Merge    *Any `json:"merge,omitempty"`
}

// type Merge struct {
// 	With Any `json:"with"`
// }
