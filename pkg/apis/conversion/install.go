package conversion

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha2"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"
)

func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1alpha2.Configuration)(nil), (*v1alpha1.Configuration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_Configuration_To_v1alpha1_Configuration(a.(*v1alpha2.Configuration), b.(*v1alpha1.Configuration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.Configuration)(nil), (*v1alpha2.Configuration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Configuration_To_v1alpha2_Configuration(a.(*v1alpha1.Configuration), b.(*v1alpha2.Configuration), scope)
	}); err != nil {
		return err
	}
	return nil
}
