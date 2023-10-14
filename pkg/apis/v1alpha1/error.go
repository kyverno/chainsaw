package v1alpha1

// Error represents an error condition.
type Error struct {
	// Error manifest.
	FileRef `json:",inline"`
}
