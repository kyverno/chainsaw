package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ReportFormatType string

const (
	JSONFormat ReportFormatType = "JSON"
	XMLFormat  ReportFormatType = "XML"
	NoReport   ReportFormatType = ""
)

// ConfigurationSpec contains the configuration used to run tests.
type ConfigurationSpec struct {
	// Timeout per test step.
	// +optional
	// +kubebuilder:default:="30s"
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// Directories containing test cases to run.
	// +optional
	TestDirs []string `json:"testDirs,omitempty"`

	// If set, do not delete the resources after running the tests (implies SkipClusterDelete).
	// +optional
	SkipDelete bool `json:"skipDelete,omitempty"`

	// StopOnFirstFailure determines whether the test should stop upon encountering the first failure.
	// +optional
	StopOnFirstFailure bool `json:"stopOnFirstFailure,omitempty"`

	// The maximum number of tests to run at once.
	// +kubebuilder:default:=8
	// +kubebuilder:validation:Format:=int
	Parallel int `json:"parallel,omitempty"`

	// ReportFormat determines test report format (JSON|XML|nil) nil == no report.
	// maps to report.Type, however we don't want generated.deepcopy to have reference to it.
	// +optional
	// +kubebuilder:validation:Enum=JSON;XML;
	ReportFormat ReportFormatType `json:"reportFormat,omitempty"`

	// ReportName defines the name of report to create. It defaults to "chainsaw-report".
	// +optional
	// +kubebuilder:default:="chainsaw-report"
	ReportName string `json:"reportName,omitempty"`

	// Namespace defines the namespace to use for tests.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Suppress is used to suppress logs.
	// +optional
	Suppress []string `json:"suppress,omitempty"`

	// FullName makes use of the full test case folder path instead of the folder name.
	// +optional
	FullName bool `json:"fullName,omitempty"`

	// SkipTestRegex is used to skip tests based on a regular expression.
	// +optional
	SkipTestRegex string `json:"skipTestRegex,omitempty"`
}
