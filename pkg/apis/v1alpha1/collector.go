package v1alpha1

// Collector defines a set of collectors.
type Collector struct {
	// PodLogs determines the pod logs collector to execute.
	// +optional
	PodLogs *PodLogsCollector `json:"podLogs,omitempty"`

	// Events determines the events collector to execute.
	// +optional
	Events *EventsCollector `json:"events,omitempty"`
}
