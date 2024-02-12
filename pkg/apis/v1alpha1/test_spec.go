package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestSpec contains the test spec.
type TestSpec struct {
	// Description contains a description of the test.
	// +optional
	Description string `json:"description,omitempty"`

	// Timeouts for the test. Overrides the global timeouts set in the Configuration on a per operation basis.
	// +optional
	Timeouts *Timeouts `json:"timeouts,omitempty"`

	// Skip determines whether the test should skipped.
	// +optional
	Skip *bool `json:"skip,omitempty"`

	// Concurrent determines whether the test should run concurrently with other tests.
	// +optional
	Concurrent *bool `json:"concurrent,omitempty"`

	// SkipDelete determines whether the resources created by the test should be deleted after the test is executed.
	// +optional
	SkipDelete *bool `json:"skipDelete,omitempty"`

	// Namespace determines whether the test should run in a random ephemeral namespace or not.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// NamespaceTemplate defines a template to create the test namespace.
	// +optional
	NamespaceTemplate *Any `json:"namespaceTemplate,omitempty"`

	// Steps defining the test.
	Steps []TestSpecStep `json:"steps"`

	// ForceTerminationGracePeriod forces the termination grace period on pods, statefulsets, daemonsets and deployments.
	// +optional
	ForceTerminationGracePeriod *metav1.Duration `json:"forceTerminationGracePeriod,omitempty"`

	// DelayBeforeCleanup adds a delay between the time a test ends and the time cleanup starts.
	// +optional
	DelayBeforeCleanup *metav1.Duration `json:"delayBeforeCleanup,omitempty"`
}
