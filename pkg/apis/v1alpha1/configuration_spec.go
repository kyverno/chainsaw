package v1alpha1

type ReportFormatType string

const (
	JSONFormat ReportFormatType = "JSON"
	XMLFormat  ReportFormatType = "XML"
	NoReport   ReportFormatType = ""
)

// ConfigurationSpec contains the configuration used to run tests.
type ConfigurationSpec struct {
	// Global timeouts configuration. Applies to all tests/test steps if not overridden.
	// +optional
	Timeouts Timeouts `json:"timeouts"`

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
	// +kubebuilder:validation:Format:=int
	// +kubebuilder:validation:Minimum:=1
	// +optional
	Parallel *int `json:"parallel,omitempty"`

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
	// If not specified, every test will execute in a random ephemeral namespace
	// unless the namespace is overridden in a the test spec.
	// +optional
	Namespace string `json:"namespace,omitempty"`

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
	// +optional
	TestFile string `json:"testFile,omitempty"`
}
