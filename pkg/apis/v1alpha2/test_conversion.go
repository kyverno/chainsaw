package v1alpha2

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	conversion "k8s.io/apimachinery/pkg/conversion"
)

func Convert_v1alpha2_TestSpec_To_v1alpha1_TestSpec(in *TestSpec, out *v1alpha1.TestSpec, s conversion.Scope) error {
	return nil
}

func Convert_v1alpha1_TestSpec_To_v1alpha2_TestSpec(in *v1alpha1.TestSpec, out *TestSpec, s conversion.Scope) error {
	return nil
}
