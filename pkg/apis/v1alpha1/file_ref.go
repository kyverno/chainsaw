package v1alpha1

// FileRef represents a file reference.
type FileRef struct {
	// File is the path to the referenced file. This can be a direct path to a file
	// or an expression that matches multiple files, such as "manifest/*.yaml" for all YAML
	// files within the "manifest" directory.
	File string `json:"file,omitempty"`
}
