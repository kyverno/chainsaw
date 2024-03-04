package v1alpha1

// Condition represents parameters for waiting on a specific condition of a resource.
type Condition struct {
	// Name defines the specific condition to wait for, e.g., "Available", "Ready".
	Name string `json:"name"`

	// Value defines the specific condition status to wait for, e.g., "True", "False".
	// +optional
	Value *string `json:"value,omitempty"`
}
