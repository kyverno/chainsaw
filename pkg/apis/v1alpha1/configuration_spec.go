package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ReportFormatType defines supported report formats.
type ReportFormatType string

const (
	// JSONFormat indicates a report in json format.
	JSONFormat ReportFormatType = "JSON"
	// XMLFormat indicates a report in xml format.
	XMLFormat ReportFormatType = "XML"
)

// ReportConfigSpec contains the configuration related to reports.
type ReportConfigSpec struct {
	// Format determines test report format (JSON|XML).
	// +optional
	// +kubebuilder:validation:Enum=JSON;XML
	Format ReportFormatType `json:"reportFormat,omitempty"`

	// Name defines the name of report to create.
	// +optional
	Name string `json:"reportName,omitempty"`
}

// ConfigurationSpec contains the configuration used to run tests.
type ConfigurationSpec struct {
	// Directories containing test cases to run.
	// +optional
	TestDirs []string `json:"testDirs,omitempty"`

	// FailFast determines whether the test should stop upon encountering the first failure.
	// +optional
	FailFast bool `json:"failFast,omitempty"`

	// The maximum number of tests to run at once.
	// +kubebuilder:validation:Format:=int
	// +kubebuilder:validation:Minimum:=1
	// +optional
	Parallel *int `json:"parallel,omitempty"`

	// Repeat indicates how many times the tests should be executed.
	// +kubebuilder:validation:Format:=int
	// +kubebuilder:validation:Minimum:=1
	// +optional
	Repeat *int `json:"repeat,omitempty"`

	// Report configuration.
	// +optional
	Report *ReportConfigSpec `json:"report,omitempty"`

	// Timeout per test step.
	// +optional
	// +kubebuilder:default:="30s"
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// If set, do not delete the resources after running the tests (implies SkipClusterDelete).
	// +optional
	SkipDelete bool `json:"skipDelete,omitempty"`

	// Namespace defines the namespace to use for tests.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Suppress is used to suppress logs.
	// +optional
	Suppress []string `json:"suppress,omitempty"`

	// FullName makes use of the full test case folder path instead of the folder name.
	// +optional
	FullName bool `json:"fullName,omitempty"`

	// ExcludeTestRegex is used to exclude tests based on a regular expression.
	// +optional
	ExcludeTestRegex string `json:"excludeTestRegex,omitempty"`

	// IncludeTestRegex is used to include tests based on a regular expression.
	// +optional
	IncludeTestRegex string `json:"includeTestRegex,omitempty"`
}
