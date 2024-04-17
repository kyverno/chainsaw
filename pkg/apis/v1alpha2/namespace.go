package v1alpha2

// Namespace contains info about the namespace used for testing.
type Namespace struct {
	// Name defines the namespace to use for tests.
	// If not specified, every test will execute in a random ephemeral namespace
	// unless the namespace is overridden in a the test spec.
	// +optional
	Name string `json:"name,omitempty"`

	// Template defines a template to create the test namespace.
	// +optional
	Template *Any `json:"template,omitempty"`
}
