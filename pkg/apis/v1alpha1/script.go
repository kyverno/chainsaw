package v1alpha1

// Script describes a script to run as a part of a test step.
type Script struct {
	// Content defines a shell script (run with "sh -c ...").
	// +optional
	Content string `json:"content,omitempty"`
}
