package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigurationSpec contains the configuration used to run tests.
type ConfigurationSpec struct {
	// Report configuration.
	// +optional
	Report *ReportConfigSpec `json:"report,omitempty"`

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

	// FailFast determines whether the test should stop upon encountering the first failure.
	// +optional
	FailFast bool `json:"failFast,omitempty"`

	// The maximum number of tests to run at once.
	// +kubebuilder:default:=8
	// +kubebuilder:validation:Format:=int
	Parallel int `json:"parallel,omitempty"`

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

	// RepeatCount indicates how many times the tests should be executed.
	// +kubebuilder:validation:Format:=int
	// +kubebuilder:validation:Minimum:=1
	// +optional
	RepeatCount *int `json:"repeatCount,omitempty"`
}
