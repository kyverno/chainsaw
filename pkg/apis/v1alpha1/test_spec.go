package v1alpha1

// TestSpec contains the test spec.
type TestSpec struct {
	// Steps defining the test.
	Steps []TestStepSpec `json:"steps"`
}
