package v1alpha1

// FileRef represents a file reference.
type FileRef struct {
	// File is the path to the referenced file.
	File string `json:"file"`
}
