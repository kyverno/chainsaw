//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by conversion-gen. DO NOT EDIT.

package v1alpha2

import (
	unsafe "unsafe"

	v1alpha1 "github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	policyv1alpha1 "github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*Configuration)(nil), (*v1alpha1.Configuration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_Configuration_To_v1alpha1_Configuration(a.(*Configuration), b.(*v1alpha1.Configuration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.Configuration)(nil), (*Configuration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Configuration_To_v1alpha2_Configuration(a.(*v1alpha1.Configuration), b.(*Configuration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*Test)(nil), (*v1alpha1.Test)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_Test_To_v1alpha1_Test(a.(*Test), b.(*v1alpha1.Test), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.Test)(nil), (*Test)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Test_To_v1alpha2_Test(a.(*v1alpha1.Test), b.(*Test), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*TestSpec)(nil), (*v1alpha1.TestSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_TestSpec_To_v1alpha1_TestSpec(a.(*TestSpec), b.(*v1alpha1.TestSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.TestSpec)(nil), (*TestSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_TestSpec_To_v1alpha2_TestSpec(a.(*v1alpha1.TestSpec), b.(*TestSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1alpha1.ConfigurationSpec)(nil), (*ConfigurationSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ConfigurationSpec_To_v1alpha2_ConfigurationSpec(a.(*v1alpha1.ConfigurationSpec), b.(*ConfigurationSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*ConfigurationSpec)(nil), (*v1alpha1.ConfigurationSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_ConfigurationSpec_To_v1alpha1_ConfigurationSpec(a.(*ConfigurationSpec), b.(*v1alpha1.ConfigurationSpec), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha2_Configuration_To_v1alpha1_Configuration(in *Configuration, out *v1alpha1.Configuration, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha2_ConfigurationSpec_To_v1alpha1_ConfigurationSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha2_Configuration_To_v1alpha1_Configuration is an autogenerated conversion function.
func Convert_v1alpha2_Configuration_To_v1alpha1_Configuration(in *Configuration, out *v1alpha1.Configuration, s conversion.Scope) error {
	return autoConvert_v1alpha2_Configuration_To_v1alpha1_Configuration(in, out, s)
}

func autoConvert_v1alpha1_Configuration_To_v1alpha2_Configuration(in *v1alpha1.Configuration, out *Configuration, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_ConfigurationSpec_To_v1alpha2_ConfigurationSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_Configuration_To_v1alpha2_Configuration is an autogenerated conversion function.
func Convert_v1alpha1_Configuration_To_v1alpha2_Configuration(in *v1alpha1.Configuration, out *Configuration, s conversion.Scope) error {
	return autoConvert_v1alpha1_Configuration_To_v1alpha2_Configuration(in, out, s)
}

func autoConvert_v1alpha2_Test_To_v1alpha1_Test(in *Test, out *v1alpha1.Test, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha2_TestSpec_To_v1alpha1_TestSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha2_Test_To_v1alpha1_Test is an autogenerated conversion function.
func Convert_v1alpha2_Test_To_v1alpha1_Test(in *Test, out *v1alpha1.Test, s conversion.Scope) error {
	return autoConvert_v1alpha2_Test_To_v1alpha1_Test(in, out, s)
}

func autoConvert_v1alpha1_Test_To_v1alpha2_Test(in *v1alpha1.Test, out *Test, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_TestSpec_To_v1alpha2_TestSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_Test_To_v1alpha2_Test is an autogenerated conversion function.
func Convert_v1alpha1_Test_To_v1alpha2_Test(in *v1alpha1.Test, out *Test, s conversion.Scope) error {
	return autoConvert_v1alpha1_Test_To_v1alpha2_Test(in, out, s)
}

func autoConvert_v1alpha2_TestSpec_To_v1alpha1_TestSpec(in *TestSpec, out *v1alpha1.TestSpec, s conversion.Scope) error {
	out.Description = in.Description
	out.Timeouts = (*v1alpha1.Timeouts)(unsafe.Pointer(in.Timeouts))
	out.Cluster = in.Cluster
	out.Clusters = *(*v1alpha1.Clusters)(unsafe.Pointer(&in.Clusters))
	out.Skip = (*bool)(unsafe.Pointer(in.Skip))
	out.Concurrent = (*bool)(unsafe.Pointer(in.Concurrent))
	out.SkipDelete = (*bool)(unsafe.Pointer(in.SkipDelete))
	out.Template = (*bool)(unsafe.Pointer(in.Template))
	out.Namespace = in.Namespace
	out.NamespaceTemplate = (*policyv1alpha1.Any)(unsafe.Pointer(in.NamespaceTemplate))
	out.Bindings = *(*[]v1alpha1.Binding)(unsafe.Pointer(&in.Bindings))
	out.Steps = *(*[]v1alpha1.TestStep)(unsafe.Pointer(&in.Steps))
	out.Catch = *(*[]v1alpha1.CatchFinally)(unsafe.Pointer(&in.Catch))
	out.ForceTerminationGracePeriod = (*v1.Duration)(unsafe.Pointer(in.ForceTerminationGracePeriod))
	out.DelayBeforeCleanup = (*v1.Duration)(unsafe.Pointer(in.DelayBeforeCleanup))
	out.DeletionPropagationPolicy = (*v1.DeletionPropagation)(unsafe.Pointer(in.DeletionPropagationPolicy))
	return nil
}

// Convert_v1alpha2_TestSpec_To_v1alpha1_TestSpec is an autogenerated conversion function.
func Convert_v1alpha2_TestSpec_To_v1alpha1_TestSpec(in *TestSpec, out *v1alpha1.TestSpec, s conversion.Scope) error {
	return autoConvert_v1alpha2_TestSpec_To_v1alpha1_TestSpec(in, out, s)
}

func autoConvert_v1alpha1_TestSpec_To_v1alpha2_TestSpec(in *v1alpha1.TestSpec, out *TestSpec, s conversion.Scope) error {
	out.Description = in.Description
	out.Timeouts = (*v1alpha1.Timeouts)(unsafe.Pointer(in.Timeouts))
	out.Cluster = in.Cluster
	out.Clusters = *(*v1alpha1.Clusters)(unsafe.Pointer(&in.Clusters))
	out.Skip = (*bool)(unsafe.Pointer(in.Skip))
	out.Concurrent = (*bool)(unsafe.Pointer(in.Concurrent))
	out.SkipDelete = (*bool)(unsafe.Pointer(in.SkipDelete))
	out.Template = (*bool)(unsafe.Pointer(in.Template))
	out.Namespace = in.Namespace
	out.NamespaceTemplate = (*policyv1alpha1.Any)(unsafe.Pointer(in.NamespaceTemplate))
	out.Bindings = *(*[]v1alpha1.Binding)(unsafe.Pointer(&in.Bindings))
	out.Steps = *(*[]v1alpha1.TestStep)(unsafe.Pointer(&in.Steps))
	out.Catch = *(*[]v1alpha1.CatchFinally)(unsafe.Pointer(&in.Catch))
	out.ForceTerminationGracePeriod = (*v1.Duration)(unsafe.Pointer(in.ForceTerminationGracePeriod))
	out.DelayBeforeCleanup = (*v1.Duration)(unsafe.Pointer(in.DelayBeforeCleanup))
	out.DeletionPropagationPolicy = (*v1.DeletionPropagation)(unsafe.Pointer(in.DeletionPropagationPolicy))
	return nil
}

// Convert_v1alpha1_TestSpec_To_v1alpha2_TestSpec is an autogenerated conversion function.
func Convert_v1alpha1_TestSpec_To_v1alpha2_TestSpec(in *v1alpha1.TestSpec, out *TestSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_TestSpec_To_v1alpha2_TestSpec(in, out, s)
}
