package v1alpha1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// FileRefOrResource represents a file reference or resource.
type FileRefOrResource struct {
	// FileRef provides a reference to the file containing the resources to be applied.
	// +optional
	FileRef `json:",inline"`

	// Resource provides a resource to be applied.
	// +kubebuilder:validation:XEmbeddedResource
	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	Resource *unstructured.Unstructured `json:"resource,omitempty"`
}
