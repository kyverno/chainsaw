package v1alpha1

// Format determines the output format (json or yaml).
// +kubebuilder:validation:Enum=json;yaml;name;go-template;go-template-file;template;templatefile;jsonpath;jsonpath-as-json;jsonpath-file
type Format string
