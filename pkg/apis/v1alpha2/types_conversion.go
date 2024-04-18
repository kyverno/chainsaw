package v1alpha2

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	conversion "k8s.io/apimachinery/pkg/conversion"
)

func Convert_v1alpha2_ConfigurationSpec_To_v1alpha1_ConfigurationSpec(in *ConfigurationSpec, out *v1alpha1.ConfigurationSpec, s conversion.Scope) error {
	return autoConvert_v1alpha2_ConfigurationSpec_To_v1alpha1_ConfigurationSpec(in, out, s)
}

func Convert_v1alpha1_ConfigurationSpec_To_v1alpha2_ConfigurationSpec(in *v1alpha1.ConfigurationSpec, out *ConfigurationSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_ConfigurationSpec_To_v1alpha2_ConfigurationSpec(in, out, s)
}
