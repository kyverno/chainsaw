package v1alpha1

// OnFailure defines actions to be executed on failure.
type OnFailure struct {
	// Collect define the collectors to run.
	// +optional
	Collect []Collector `json:"collect,omitempty"`
}
