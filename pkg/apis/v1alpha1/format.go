package v1alpha1

// Format determines the output format (json or yaml).
// +kubebuilder:validation:Pattern=`^(?:json|yaml|\(.+\))$`
type Format string
