package v1alpha1

// Cluster defines cluster config and context.
type Cluster struct {
	// Kubeconfig is the path to the referenced file.
	Kubeconfig string `json:"kubeconfig"`

	// Context is the name of the context to use.
	// +optional
	Context string `json:"context,omitempty"`
}
