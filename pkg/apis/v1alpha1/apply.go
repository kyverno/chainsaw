package v1alpha1

// Apply represents a manifest to be applied.
type Apply struct {
	// Manifest to be applied.
	FileRef `json:",inline"`
}
