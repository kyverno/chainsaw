package v1alpha1

// Create represents a set of resources that should be created.
// If a resource already exists in the cluster it will fail.
type Create struct {
	// FileRefOrResource provides a reference to the file containing the resources to be created.
	FileRefOrResource `json:",inline"`

	// DryRun determines whether the file should be applied in dry run mode.
	// +optional
	DryRun *bool `json:"dryRun,omitempty"`

	// Expect defines a list of matched checks to validate the operation outcome.
	// +optional
	Expect []MatchedCheck `json:"expect,omitempty"`
}
