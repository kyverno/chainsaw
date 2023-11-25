package v1alpha1

// Apply represents a set of configurations or resources that
// should be applied during testing.
type Apply struct {
	// FileRefOrResource provides a reference to the file containing the resources to be applied.
	FileRefOrResource `json:",inline"`

	// DryRun determines whether the file should be applied in dry run mode.
	// +optional
	DryRun *bool `json:"dryRun,omitempty"`

	// Check is an assertion tree to validate outcome.
	// +optional
	Check *Check `json:"check,omitempty"`
}
