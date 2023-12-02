package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Sleep represents a duration while nothing happens.
type Sleep struct {
	// Duration is the delay used for sleeping.
	Duration metav1.Duration `json:"duration"`
}
