package v1alpha1

// ObjectSelector represents a strategy to select objects.
// For a single object name and namespace are used to identify the object.
// For multiple objects use labels.
type ObjectSelector struct {
	// Namespace of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	// +optional
	Name string `json:"name,omitempty"`

	// Label selector to match objects to delete
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
}
