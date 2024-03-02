package v1alpha1

// ObjectReference represents one or more objects with a specific apiVersion and kind.
// For a single object name and namespace are used to identify the object.
// For multiple objects use labels.
type ObjectReference struct {
	// ObjectType determines the type of referenced objects.
	ObjectType `json:",inline"`

	// ObjectSelector determines the selection process of referenced objects.
	ObjectSelector `json:",inline"`
}
