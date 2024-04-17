package v1alpha1

// JsonPath represents parameters for waiting on a json path of a resource.
type JsonPath struct {
	// Path defines the json path to wait for, e.g. '{.status.phase}'.
	Path string `json:"path"`

	// Value defines the expected value to wait for, e.g., "Running".
	Value string `json:"value"`
}
