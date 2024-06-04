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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha2

import (
	v1alpha1 "github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CleanupOptions) DeepCopyInto(out *CleanupOptions) {
	*out = *in
	if in.DelayBeforeCleanup != nil {
		in, out := &in.DelayBeforeCleanup, &out.DelayBeforeCleanup
		*out = new(v1.Duration)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CleanupOptions.
func (in *CleanupOptions) DeepCopy() *CleanupOptions {
	if in == nil {
		return nil
	}
	out := new(CleanupOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Configuration) DeepCopyInto(out *Configuration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Configuration.
func (in *Configuration) DeepCopy() *Configuration {
	if in == nil {
		return nil
	}
	out := new(Configuration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Configuration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigurationSpec) DeepCopyInto(out *ConfigurationSpec) {
	*out = *in
	if in.Cleanup != nil {
		in, out := &in.Cleanup, &out.Cleanup
		*out = new(CleanupOptions)
		(*in).DeepCopyInto(*out)
	}
	if in.Clusters != nil {
		in, out := &in.Clusters, &out.Clusters
		*out = make(v1alpha1.Clusters, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	out.Deletion = in.Deletion
	out.Discovery = in.Discovery
	if in.Error != nil {
		in, out := &in.Error, &out.Error
		*out = new(ErrorOptions)
		(*in).DeepCopyInto(*out)
	}
	if in.Execution != nil {
		in, out := &in.Execution, &out.Execution
		*out = new(ExecutionOptions)
		(*in).DeepCopyInto(*out)
	}
	if in.Namespace != nil {
		in, out := &in.Namespace, &out.Namespace
		*out = new(NamespaceOptions)
		(*in).DeepCopyInto(*out)
	}
	if in.Report != nil {
		in, out := &in.Report, &out.Report
		*out = new(ReportOptions)
		**out = **in
	}
	out.Templating = in.Templating
	in.Timeouts.DeepCopyInto(&out.Timeouts)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigurationSpec.
func (in *ConfigurationSpec) DeepCopy() *ConfigurationSpec {
	if in == nil {
		return nil
	}
	out := new(ConfigurationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeletionOptions) DeepCopyInto(out *DeletionOptions) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeletionOptions.
func (in *DeletionOptions) DeepCopy() *DeletionOptions {
	if in == nil {
		return nil
	}
	out := new(DeletionOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DiscoveryOptions) DeepCopyInto(out *DiscoveryOptions) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DiscoveryOptions.
func (in *DiscoveryOptions) DeepCopy() *DiscoveryOptions {
	if in == nil {
		return nil
	}
	out := new(DiscoveryOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ErrorOptions) DeepCopyInto(out *ErrorOptions) {
	*out = *in
	if in.Catch != nil {
		in, out := &in.Catch, &out.Catch
		*out = make([]v1alpha1.CatchFinally, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ErrorOptions.
func (in *ErrorOptions) DeepCopy() *ErrorOptions {
	if in == nil {
		return nil
	}
	out := new(ErrorOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExecutionOptions) DeepCopyInto(out *ExecutionOptions) {
	*out = *in
	if in.Parallel != nil {
		in, out := &in.Parallel, &out.Parallel
		*out = new(int)
		**out = **in
	}
	if in.RepeatCount != nil {
		in, out := &in.RepeatCount, &out.RepeatCount
		*out = new(int)
		**out = **in
	}
	if in.ForceTerminationGracePeriod != nil {
		in, out := &in.ForceTerminationGracePeriod, &out.ForceTerminationGracePeriod
		*out = new(v1.Duration)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExecutionOptions.
func (in *ExecutionOptions) DeepCopy() *ExecutionOptions {
	if in == nil {
		return nil
	}
	out := new(ExecutionOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespaceOptions) DeepCopyInto(out *NamespaceOptions) {
	*out = *in
	if in.Template != nil {
		in, out := &in.Template, &out.Template
		*out = (*in).DeepCopy()
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespaceOptions.
func (in *NamespaceOptions) DeepCopy() *NamespaceOptions {
	if in == nil {
		return nil
	}
	out := new(NamespaceOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReportOptions) DeepCopyInto(out *ReportOptions) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReportOptions.
func (in *ReportOptions) DeepCopy() *ReportOptions {
	if in == nil {
		return nil
	}
	out := new(ReportOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TemplatingOptions) DeepCopyInto(out *TemplatingOptions) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TemplatingOptions.
func (in *TemplatingOptions) DeepCopy() *TemplatingOptions {
	if in == nil {
		return nil
	}
	out := new(TemplatingOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Test) DeepCopyInto(out *Test) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Test.
func (in *Test) DeepCopy() *Test {
	if in == nil {
		return nil
	}
	out := new(Test)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Test) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TestExecutionOptions) DeepCopyInto(out *TestExecutionOptions) {
	*out = *in
	if in.TerminationGracePeriod != nil {
		in, out := &in.TerminationGracePeriod, &out.TerminationGracePeriod
		*out = new(v1.Duration)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TestExecutionOptions.
func (in *TestExecutionOptions) DeepCopy() *TestExecutionOptions {
	if in == nil {
		return nil
	}
	out := new(TestExecutionOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TestSpec) DeepCopyInto(out *TestSpec) {
	*out = *in
	in.Cleanup.DeepCopyInto(&out.Cleanup)
	if in.Clusters != nil {
		in, out := &in.Clusters, &out.Clusters
		*out = make(v1alpha1.Clusters, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.Execution.DeepCopyInto(&out.Execution)
	if in.Bindings != nil {
		in, out := &in.Bindings, &out.Bindings
		*out = make([]v1alpha1.Binding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	out.Deletion = in.Deletion
	in.Error.DeepCopyInto(&out.Error)
	in.Namespace.DeepCopyInto(&out.Namespace)
	if in.Steps != nil {
		in, out := &in.Steps, &out.Steps
		*out = make([]v1alpha1.TestStep, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	out.Templating = in.Templating
	in.Timeouts.DeepCopyInto(&out.Timeouts)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TestSpec.
func (in *TestSpec) DeepCopy() *TestSpec {
	if in == nil {
		return nil
	}
	out := new(TestSpec)
	in.DeepCopyInto(out)
	return out
}
