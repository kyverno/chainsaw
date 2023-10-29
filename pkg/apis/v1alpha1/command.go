package v1alpha1

// Command describes a command to run as a part of a test step.
type Command struct {
	// Entrypoint is the command entry point to run.
	Entrypoint string `json:"entrypoint"`

	// Args is the command arguments.
	// +optional
	Args []string `json:"args,omitempty"`
}
