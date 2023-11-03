package v1alpha1

// Catch defines actions to be executed on failure.
type Catch struct {
	// Collect define the collectors to run.
	// +optional
	Collect *Collect `json:"collect,omitempty"`

	// Exec define the commands and/or scripts to run.
	// +optional
	Exec *Exec `json:"exec,omitempty"`
}
