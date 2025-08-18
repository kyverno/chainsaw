package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:storageversion

// Configuration is the resource that contains the configuration used to run tests.
type Configuration struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata"`

	// Configuration spec.
	Spec ConfigurationSpec `json:"spec"`
}

// ConfigurationSpec contains the configuration used to run tests.
// +k8s:conversion-gen=false
type ConfigurationSpec struct {
	// Cleanup contains cleanup configuration.
	// +optional
	// +kubebuilder:default:={}
	Cleanup CleanupOptions `json:"cleanup"`

	// Clusters holds a registry to clusters to support multi-cluster tests.
	// +optional
	Clusters Clusters `json:"clusters"`

	// Deletion contains the global deletion configuration.
	// +optional
	// +kubebuilder:default:={}
	Deletion DeletionOptions `json:"deletion"`

	// Discovery contains tests discovery configuration.
	// +optional
	// +kubebuilder:default:={}
	Discovery DiscoveryOptions `json:"discovery"`

	// Error contains the global error configuration.
	// +optional
	// +kubebuilder:default:={}
	Error ErrorOptions `json:"error"`

	// Execution contains tests execution configuration.
	// +optional
	// +kubebuilder:default:={}
	Execution ExecutionOptions `json:"execution"`

	// Namespace contains properties for the namespace to use for tests.
	// +optional
	// +kubebuilder:default:={}
	Namespace NamespaceOptions `json:"namespace"`

	// Report contains properties for the report.
	// +optional
	Report *ReportOptions `json:"report,omitempty"`

	// Templating contains the templating config.
	// +optional
	// +kubebuilder:default:={}
	Templating TemplatingOptions `json:"templating"`

	// Global timeouts configuration. Applies to all tests/test steps if not overridden.
	// +optional
	// +kubebuilder:default:={}
	Timeouts DefaultTimeouts `json:"timeouts"`
}
