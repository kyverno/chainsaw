package v1alpha1

// Command describes a command to run as a part of a test step.
type Command struct {
	// EntryPoint is the command entry point to run.
	// +optional
	EntryPoint string `json:"entryPoint,omitempty"`

	// Args is the command arguments.
	// +optional
	Args []string `json:"args,omitempty"`
}
