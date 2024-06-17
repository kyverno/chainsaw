package v1alpha2

import (
	_ "github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CleanupOptions contains the configuration used for cleaning up resources.
type CleanupOptions struct {
	// If set, do not delete the resources after running a test.
	// +optional
	SkipDelete bool `json:"skipDelete,omitempty"`

	// DelayBeforeCleanup adds a delay between the time a test ends and the time cleanup starts.
	// +optional
	DelayBeforeCleanup *metav1.Duration `json:"delayBeforeCleanup,omitempty"`
}

// DeletionOptions contains the configuration used for deleting resources.
type DeletionOptions struct {
	// Propagation decides if a deletion will propagate to the dependents of
	// the object, and how the garbage collector will handle the propagation.
	// +optional
	// +kubebuilder:validation:Enum:=Orphan;Background;Foreground
	// +kubebuilder:default:=Background
	Propagation metav1.DeletionPropagation `json:"propagation,omitempty"`
}

// DiscoveryOptions contains the discovery configuration used when discovering tests in folders.
type DiscoveryOptions struct {
	// ExcludeTestRegex is used to exclude tests based on a regular expression.
	// +optional
	ExcludeTestRegex string `json:"excludeTestRegex,omitempty"`

	// IncludeTestRegex is used to include tests based on a regular expression.
	// +optional
	IncludeTestRegex string `json:"includeTestRegex,omitempty"`

	// TestFile is the name of the file containing the test to run.
	// If no extension is provided, chainsaw will try with .yaml first and .yml if needed.
	// +optional
	// +kubebuilder:default:="chainsaw-test"
	TestFile string `json:"testFile,omitempty"`

	// FullName makes use of the full test case folder path instead of the folder name.
	// +optional
	FullName bool `json:"fullName,omitempty"`
}

// ErrorOptions contains the global error configuration.
type ErrorOptions struct {
	// Catch defines what the tests steps will execute when an error happens.
	// This will be combined with catch handlers defined at the test and step levels.
	// +optional
	Catch []Operation `json:"catch,omitempty"`
}

// ExecutionOptions determines how tests are run.
type ExecutionOptions struct {
	// FailFast determines whether the test should stop upon encountering the first failure.
	// +optional
	FailFast bool `json:"failFast,omitempty"`

	// The maximum number of tests to run at once.
	// +kubebuilder:validation:Format:=int
	// +kubebuilder:validation:Minimum:=1
	// +optional
	Parallel *int `json:"parallel,omitempty"`

	// RepeatCount indicates how many times the tests should be executed.
	// +kubebuilder:validation:Format:=int
	// +kubebuilder:validation:Minimum:=1
	// +optional
	RepeatCount *int `json:"repeatCount,omitempty"`

	// ForceTerminationGracePeriod forces the termination grace period on pods, statefulsets, daemonsets and deployments.
	// +optional
	ForceTerminationGracePeriod *metav1.Duration `json:"forceTerminationGracePeriod,omitempty"`
}

// NamespaceOptions contains the configuration used to allocate a namespace for each test.
type NamespaceOptions struct {
	// Name defines the namespace to use for tests.
	// If not specified, every test will execute in a random ephemeral namespace
	// unless the namespace is overridden in a the test spec.
	// +optional
	Name string `json:"name,omitempty"`

	// Template defines a template to create the test namespace.
	// +optional
	Template *Any `json:"template,omitempty"`
}

type ReportFormatType string

const (
	JSONFormat ReportFormatType = "JSON"
	XMLFormat  ReportFormatType = "XML"
)

// ReportOptions contains the configuration used for reporting.
type ReportOptions struct {
	// ReportFormat determines test report format (JSON|XML).
	// +optional
	// +kubebuilder:validation:Enum:=JSON;XML
	// +kubebuilder:default:="JSON"
	Format ReportFormatType `json:"format,omitempty"`

	// ReportPath defines the path.
	// +optional
	Path string `json:"path,omitempty"`

	// ReportName defines the name of report to create. It defaults to "chainsaw-report".
	// +optional
	// +kubebuilder:default:="chainsaw-report"
	Name string `json:"name,omitempty"`
}

// TemplatingOptions contains the templating configuration.
type TemplatingOptions struct {
	// Enabled determines whether resources should be considered for templating.
	// +optional
	// +kubebuilder:default:=true
	Enabled bool `json:"enabled"`
}
