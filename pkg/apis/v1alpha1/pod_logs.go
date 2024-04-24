package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodLogs defines how to collect pod logs.
type PodLogs struct {
	// Timeout for the operation. Overrides the global timeout set in the Configuration.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).
	// +optional
	Cluster string `json:"cluster,omitempty"`

	// ClusterConfig defines a reference to a cluster configuration (default cluster will be used if not specified and/or overridden).
	// +optional
	ClusterConfig *Cluster `json:"clusterConfig,omitempty"`

	// ObjectLabelsSelector determines the selection process of referenced objects.
	ObjectLabelsSelector `json:",inline"`

	// Container in pod to get logs from else --all-containers is used.
	// +optional
	Container string `json:"container,omitempty"`

	// Tail is the number of last lines to collect from pods. If omitted or zero,
	// then the default is 10 if you use a selector, or -1 (all) if you use a pod name.
	// This matches default behavior of `kubectl logs`.
	// +optional
	Tail *int `json:"tail,omitempty"`
}
