package v1alpha1

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
