package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// Configuration is the resource that contains the configuration used to run tests.
type Configuration struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Configuration spec.
	Spec ConfigurationSpec `json:"spec"`
}

// ConfigurationSpec contains the configuration used to run tests.
type ConfigurationSpec struct {
	// Global timeouts configuration. Applies to all tests/test steps if not overridden.
	// +optional
	// +kubebuilder:default:={}
	Timeouts DefaultTimeouts `json:"timeouts"`

	// If set, do not delete the resources after running the tests (implies SkipClusterDelete).
	// +optional
	SkipDelete bool `json:"skipDelete,omitempty"`

	// Template determines whether resources should be considered for templating.
	// +optional
	// +kubebuilder:default:=true
	Template bool `json:"template"`

	// FailFast determines whether the test should stop upon encountering the first failure.
	// +optional
	FailFast bool `json:"failFast,omitempty"`

	// The maximum number of tests to run at once.
	// +kubebuilder:validation:Format:=int
	// +kubebuilder:validation:Minimum:=1
	// +optional
	Parallel *int `json:"parallel,omitempty"`

	// DeletionPropagationPolicy decides if a deletion will propagate to the dependents of
	// the object, and how the garbage collector will handle the propagation.
	// +optional
	// +kubebuilder:validation:Enum:=Orphan;Background;Foreground
	// +kubebuilder:default:=Background
	DeletionPropagationPolicy metav1.DeletionPropagation `json:"deletionPropagationPolicy,omitempty"`

	// ReportFormat determines test report format (JSON|XML|JUNIT-TEST|JUNIT-STEP|JUNIT-OPERATION|nil) nil == no report.
	// maps to report.Type, however we don't want generated.deepcopy to have reference to it.
	// +optional
	// +kubebuilder:validation:Enum:=JSON;XML;JUNIT-TEST;JUNIT-STEP;JUNIT-OPERATION;
	ReportFormat ReportFormatType `json:"reportFormat,omitempty"`

	// ReportPath defines the path.
	// +optional
	ReportPath string `json:"reportPath,omitempty"`

	// ReportName defines the name of report to create. It defaults to "chainsaw-report".
	// +optional
	// +kubebuilder:default:="chainsaw-report"
	ReportName string `json:"reportName,omitempty"`

	// Namespace defines the namespace to use for tests.
	// If not specified, every test will execute in a random ephemeral namespace
	// unless the namespace is overridden in a the test spec.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// NamespaceTemplate defines a template to create the test namespace.
	// +optional
	NamespaceTemplate *Projection `json:"namespaceTemplate,omitempty"`

	// FullName makes use of the full test case folder path instead of the folder name.
	// +optional
	FullName bool `json:"fullName,omitempty"`

	// ExcludeTestRegex is used to exclude tests based on a regular expression.
	// +optional
	ExcludeTestRegex string `json:"excludeTestRegex,omitempty"`

	// IncludeTestRegex is used to include tests based on a regular expression.
	// +optional
	IncludeTestRegex string `json:"includeTestRegex,omitempty"`

	// RepeatCount indicates how many times the tests should be executed.
	// +kubebuilder:validation:Format:=int
	// +kubebuilder:validation:Minimum:=1
	// +optional
	RepeatCount *int `json:"repeatCount,omitempty"`

	// TestFile is the name of the file containing the test to run.
	// If no extension is provided, chainsaw will try with .yaml first and .yml if needed.
	// +kubebuilder:default:="chainsaw-test"
	// +optional
	TestFile string `json:"testFile,omitempty"`

	// ForceTerminationGracePeriod forces the termination grace period on pods, statefulsets, daemonsets and deployments.
	// +optional
	ForceTerminationGracePeriod *metav1.Duration `json:"forceTerminationGracePeriod,omitempty"`

	// DelayBeforeCleanup adds a delay between the time a test ends and the time cleanup starts.
	// +optional
	DelayBeforeCleanup *metav1.Duration `json:"delayBeforeCleanup,omitempty"`

	// Clusters holds a registry to clusters to support multi-cluster tests.
	// +optional
	Clusters Clusters `json:"clusters,omitempty"`

	// Catch defines what the tests steps will execute when an error happens.
	// This will be combined with catch handlers defined at the test and step levels.
	// +optional
	Catch []CatchFinally `json:"catch,omitempty"`
}

type ReportFormatType string

const (
	JSONFormat           ReportFormatType = "JSON"
	XMLFormat            ReportFormatType = "XML"
	JUnitTestFormat      ReportFormatType = "JUNIT-TEST"
	JUnitStepFormat      ReportFormatType = "JUNIT-STEP"
	JUnitOperationFormat ReportFormatType = "JUNIT-OPERATION"
	NoReport             ReportFormatType = ""
)
