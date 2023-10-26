package v1alpha1

// ObjectReference represents one or more objects with a specific apiVersion and kind.
// For a single object name and namespace are used to identify the object.
// For multiple objects use labels.
type ObjectReference struct {
	// ObjectSelector determines the selection process of referenced objects.
	ObjectSelector `json:",inline"`

	// API version of the referent.
	APIVersion string `json:"apiVersion"`

	// Kind of the referent.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind string `json:"kind"`
}
