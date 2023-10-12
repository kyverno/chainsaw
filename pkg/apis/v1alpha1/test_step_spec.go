package v1alpha1

type TestStepSpec struct {
	// +optional
	Assert []Assert `json:"assert"`
	// +optional
	Apply []Apply `json:"apply"`
	// +optional
	Error []Error `json:"error"`
}
