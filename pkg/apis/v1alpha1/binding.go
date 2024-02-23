package v1alpha1

// Binding represents a key/value set as a binding in an executing test.
type Binding struct {
	// Name the name of the binding.
	Name string `json:"name"`

	// Value value of the binding.
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:pruning:PreserveUnknownFields
	Value Any `json:"value"`
}
