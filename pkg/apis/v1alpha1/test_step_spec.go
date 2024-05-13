package v1alpha1

// TestStepSpec defines the desired state and behavior for each test step.
type TestStepSpec struct {
	// Description contains a description of the test step.
	// +optional
	Description string `json:"description,omitempty"`

	// Timeouts for the test step. Overrides the global timeouts set in the Configuration and the timeouts eventually set in the Test.
	// +optional
	Timeouts *Timeouts `json:"timeouts,omitempty"`

	// Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).
	// +optional
	Cluster string `json:"cluster,omitempty"`

	// Clusters holds a registry to clusters to support multi-cluster tests.
	// +optional
	Clusters map[string]Cluster `json:"clusters,omitempty"`

	// SkipDelete determines whether the resources created by the step should be deleted after the test step is executed.
	// +optional
	SkipDelete *bool `json:"skipDelete,omitempty"`

	// Template determines whether resources should be considered for templating.
	// +optional
	Template *bool `json:"template,omitempty"`

	// Bindings defines additional binding key/values.
	// +optional
	Bindings []Binding `json:"bindings,omitempty"`

	// Try defines what the step will try to execute.
	// +kubebuilder:validation:MinItems:=1
	Try []Operation `json:"try"`

	// Catch defines what the step will execute when an error happens.
	// +optional
	Catch []Catch `json:"catch,omitempty"`

	// Finally defines what the step will execute after the step is terminated.
	// +optional
	Finally []Finally `json:"finally,omitempty"`
}
