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

	// Directories containing test cases to run.
	// +optional
	TestDirs []string `json:"testDirs,omitempty"`

	// If set, do not delete the resources after running the tests (implies SkipClusterDelete).
	// +optional
	// +kubebuilder:default:=false
	SkipDelete bool `json:"skipDelete,omitempty"`

	// StopOnFirstFailure determines whether the test should stop upon encountering the first failure.
	// +optional
	// +kubebuilder:default:=false
	StopOnFirstFailure bool `json:"stopOnFirstFailure,omitempty"`

	// The maximum number of tests to run at once.
	// +kubebuilder:default:=8
	// +kubebuilder:validation:Format:=int64
	Parallel int `json:"parallel,omitempty"`

	// ReportFormat determines test report format (JSON|XML|nil) nil == no report.
	// maps to report.Type, however we don't want generated.deepcopy to have reference to it.
	// +optional
	// +kubebuilder:default:=""
	ReportFormat string `json:"reportFormat,omitempty"`

	// ReportName defines the name of report to create. It defaults to "kuttl-report".
	// +optional
	// +kubebuilder:default:="kuttl-report"
	ReportName string `json:"reportName,omitempty"`

	// Namespace defines the namespace to use for tests.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Suppress is used to suppress logs.
	// +optional
	Suppress []string `json:"suppress,omitempty"`

	// FullName makes use of the full test case folder path instead of the folder name.
	// +optional
	// +kubebuilder:default:=false
	FullName bool `json:"fullName,omitempty"`

	// SkipTestRegex is used to skip tests based on a regular expression.
	// +optional
	// +kubebuilder:default:=""
	SkipTestRegex string `json:"skipTestRegex,omitempty"`
}
