package v1alpha1

// EventsCollector defines how to collects events.
type EventsCollector struct {
	// ObjectSelector determines the selection process of events to collect.
	ObjectSelector `json:",inline"`
}
