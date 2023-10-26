package v1alpha1

// EventsCollector defines how to collects events.
type EventsCollector struct {
	// Namespace of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	// +optional
	Name string `json:"name,omitempty"`

	// Selector defines labels selector.
	// +optional
	Selector string `json:"selector,omitempty"`
}
