package v1alpha1

// Finally defines actions to be executed at the end of a test.
type Finally struct {
	// Collect define the collectors to run.
	// +optional
	Collect *Collect `json:"collect,omitempty"`

	// Exec define the commands and/or scripts to run.
	// +optional
	Exec *Exec `json:"exec,omitempty"`
}
