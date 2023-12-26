package v1alpha1

type Kubeconfig struct {
	// File is the path to the referenced file.
	File string `json:"file,omitempty"`
	// ClusterRef is the name of the cluster to use.
	ClusterRef string `json:"clusterRef,omitempty"`
}
