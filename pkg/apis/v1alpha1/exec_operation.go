package v1alpha1

// Exec describes a command and/or script operation.
type ExecOperation struct {
	// Exec defines the command and/or script.
	// +optional
	Exec `json:",inline"`

	// SkipLogOutput removes the output from the command. Useful for sensitive logs or to reduce noise.
	// +optional
	SkipLogOutput bool `json:"skipLogOutput,omitempty"`
}
