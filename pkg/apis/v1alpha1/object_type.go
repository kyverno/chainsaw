package v1alpha1

// ObjectType represents a specific apiVersion and kind.
type ObjectType struct {
	// API version of the referent.
	APIVersion string `json:"apiVersion"`

	// Kind of the referent.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind string `json:"kind"`
}
