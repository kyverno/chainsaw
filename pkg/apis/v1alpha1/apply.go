package v1alpha1

type Apply struct {
	// File containing the manifest to be applied.
	File string `json:"file"`
}
