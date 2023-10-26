package v1alpha1

// OnFailure defines actions to be executed on failure.
type OnFailure struct {
	// Collect define the collectors to run.
	// +optional
	Collect *Collect `json:"collect,omitempty"`

	// Exec define the commands and/or scripts to run.
	// +optional
	Exec *Exec `json:"exec,omitempty"`
}
