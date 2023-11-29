package v1alpha1

// Apply represents a set of configurations or resources that
// should be applied during testing.
type Apply struct {
	// FileRefOrResource provides a reference to the resources to be applied.
	FileRefOrResource `json:",inline"`

	// DryRun determines whether the file should be applied in dry run mode.
	// +optional
	DryRun *bool `json:"dryRun,omitempty"`

	// Expect defines a list of matched checks to validate the operation outcome.
	// +optional
	Expect []Expectation `json:"expect,omitempty"`
}
