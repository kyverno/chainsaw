package v1alpha1

// ConfigurationSpec contains the configuration used to run tests.
type ConfigurationSpec struct {
	Timeout int `json:"timeout"`
}
