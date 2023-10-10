package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigurationSpec contains the configuration used to run tests.
type ConfigurationSpec struct {
	// Timeout per test step.
	// +optional
	// +kubebuilder:default:="30s"
	Timeout *metav1.Duration `json:"timeout,omitempty"`
}
