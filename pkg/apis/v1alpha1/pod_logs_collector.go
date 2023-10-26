package v1alpha1

// PodLogsCollector defines how to collects pod logs.
type PodLogsCollector struct {
	// ObjectSelector determines the selection process of pods to collect logs from.
	ObjectSelector `json:",inline"`

	// // Container in pod to get logs from else --all-containers is used.
	// Container string `json:"container,omitempty"`
	// // Tail is the number of last lines to collect from pods. If omitted or zero,
	// // then the default is 10 if you use a selector, or -1 (all) if you use a pod name.
	// // This matches default behavior of `kubectl logs`.
	// Tail int `json:"tail,omitempty"`
}
