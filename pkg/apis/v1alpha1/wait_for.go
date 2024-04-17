package v1alpha1

// For specifies the condition to wait for.
type For struct {
	// Deletion specifies parameters for waiting on a resource's deletion.
	// +optional
	Deletion *Deletion `json:"deletion,omitempty"`

	// Condition specifies the condition to wait for.
	// +optional
	Condition *Condition `json:"condition,omitempty"`

	// JsonPath specifies the json path condition to wait for.
	// +optional
	JsonPath *JsonPath `json:"jsonPath,omitempty"`
}
