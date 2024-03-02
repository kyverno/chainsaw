package v1alpha1

// ResourceReference represents a resource (API), it can be represented with a resource or a kind.
// Optionally an apiVersion can be specified.
type ResourceReference struct {
	// API version of the referent.
	// +optional
	APIVersion string `json:"apiVersion,omitempty"`

	// Kind of the referent.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	Kind string `json:"kind,omitempty"`

	// Resource name of the referent.
	// +optional
	Resource string `json:"resource,omitempty"`
}
