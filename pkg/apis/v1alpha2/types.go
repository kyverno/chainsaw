package v1alpha2

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	_ "github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type (
	Any             = v1alpha1.Any
	Binding         = v1alpha1.Binding
	Check           = v1alpha1.Check
	Cluster         = v1alpha1.Cluster
	Clusters        = v1alpha1.Clusters
	Expectation     = v1alpha1.Expectation
	Format          = v1alpha1.Format
	ObjectName      = v1alpha1.ObjectName
	ObjectType      = v1alpha1.ObjectType
	Output          = v1alpha1.Output
	Timeouts        = v1alpha1.Timeouts
	DefaultTimeouts = v1alpha1.DefaultTimeouts
)

// ObjectReference represents one or more objects with a specific apiVersion and kind.
// For a single object name and namespace are used to identify the object.
// For multiple objects use labels.
// +k8s:conversion-gen=false
type ObjectReference struct {
	ObjectType `json:",inline"`
	ObjectName `json:",inline"`

	// Label selector to match objects to delete
	// +optional
	Labels *metav1.LabelSelector `json:"labelSelector,omitempty"`
}
