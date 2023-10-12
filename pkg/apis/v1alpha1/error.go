package v1alpha1

type Error struct {
	// File containing the assertion manifest.
	// It is expected that the assertion fails.
	File string `json:"file"`
}
