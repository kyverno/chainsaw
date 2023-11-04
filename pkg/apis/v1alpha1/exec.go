package v1alpha1

// Exec describes a command or script.
type Exec struct {
	// Command defines a command to run.
	// +optional
	Command *Command `json:"command,omitempty"`

	// Script defines a script to run.
	// +optional
	Script *Script `json:"script,omitempty"`
}
