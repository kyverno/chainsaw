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

package v1alpha1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Apply) DeepCopyInto(out *Apply) {
	*out = *in
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Bindings != nil {
		in, out := &in.Bindings, &out.Bindings
		*out = make([]Binding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Outputs != nil {
		in, out := &in.Outputs, &out.Outputs
		*out = make([]Output, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Clusters != nil {
		in, out := &in.Clusters, &out.Clusters
		*out = make(map[string]Cluster, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.FileRefOrResource.DeepCopyInto(&out.FileRefOrResource)
	if in.Template != nil {
		in, out := &in.Template, &out.Template
		*out = new(bool)
		**out = **in
	}
	if in.DryRun != nil {
		in, out := &in.DryRun, &out.DryRun
		*out = new(bool)
		**out = **in
	}
	if in.Expect != nil {
		in, out := &in.Expect, &out.Expect
		*out = make([]Expectation, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Apply.
func (in *Apply) DeepCopy() *Apply {
	if in == nil {
		return nil
	}
	out := new(Apply)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Assert) DeepCopyInto(out *Assert) {
	*out = *in
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Bindings != nil {
		in, out := &in.Bindings, &out.Bindings
		*out = make([]Binding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Clusters != nil {
		in, out := &in.Clusters, &out.Clusters
		*out = make(map[string]Cluster, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.FileRefOrCheck.DeepCopyInto(&out.FileRefOrCheck)
	if in.Template != nil {
		in, out := &in.Template, &out.Template
		*out = new(bool)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Assert.
func (in *Assert) DeepCopy() *Assert {
	if in == nil {
		return nil
	}
	out := new(Assert)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Binding) DeepCopyInto(out *Binding) {
	*out = *in
	in.Value.DeepCopyInto(&out.Value)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Binding.
func (in *Binding) DeepCopy() *Binding {
	if in == nil {
		return nil
	}
	out := new(Binding)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Catch) DeepCopyInto(out *Catch) {
	*out = *in
	if in.PodLogs != nil {
		in, out := &in.PodLogs, &out.PodLogs
		*out = new(PodLogs)
		(*in).DeepCopyInto(*out)
	}
	if in.Events != nil {
		in, out := &in.Events, &out.Events
		*out = new(Events)
		(*in).DeepCopyInto(*out)
	}
	if in.Describe != nil {
		in, out := &in.Describe, &out.Describe
		*out = new(Describe)
		(*in).DeepCopyInto(*out)
	}
	if in.Wait != nil {
		in, out := &in.Wait, &out.Wait
		*out = new(Wait)
		(*in).DeepCopyInto(*out)
	}
	if in.Get != nil {
		in, out := &in.Get, &out.Get
		*out = new(Get)
		(*in).DeepCopyInto(*out)
	}
	if in.Delete != nil {
		in, out := &in.Delete, &out.Delete
		*out = new(Delete)
		(*in).DeepCopyInto(*out)
	}
	if in.Command != nil {
		in, out := &in.Command, &out.Command
		*out = new(Command)
		(*in).DeepCopyInto(*out)
	}
	if in.Script != nil {
		in, out := &in.Script, &out.Script
		*out = new(Script)
		(*in).DeepCopyInto(*out)
	}
	if in.Sleep != nil {
		in, out := &in.Sleep, &out.Sleep
		*out = new(Sleep)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Catch.
func (in *Catch) DeepCopy() *Catch {
	if in == nil {
		return nil
	}
	out := new(Catch)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Cluster) DeepCopyInto(out *Cluster) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Cluster.
func (in *Cluster) DeepCopy() *Cluster {
	if in == nil {
		return nil
	}
	out := new(Cluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Command) DeepCopyInto(out *Command) {
	*out = *in
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Bindings != nil {
		in, out := &in.Bindings, &out.Bindings
		*out = make([]Binding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Outputs != nil {
		in, out := &in.Outputs, &out.Outputs
		*out = make([]Output, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]Binding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Clusters != nil {
		in, out := &in.Clusters, &out.Clusters
		*out = make(map[string]Cluster, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Args != nil {
		in, out := &in.Args, &out.Args
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Check != nil {
		in, out := &in.Check, &out.Check
		*out = (*in).DeepCopy()
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Command.
func (in *Command) DeepCopy() *Command {
	if in == nil {
		return nil
	}
	out := new(Command)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Condition) DeepCopyInto(out *Condition) {
	*out = *in
	if in.Value != nil {
		in, out := &in.Value, &out.Value
		*out = new(string)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Condition.
func (in *Condition) DeepCopy() *Condition {
	if in == nil {
		return nil
	}
	out := new(Condition)
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
	in.Timeouts.DeepCopyInto(&out.Timeouts)
	if in.Template != nil {
		in, out := &in.Template, &out.Template
		*out = new(bool)
		**out = **in
	}
	if in.Parallel != nil {
		in, out := &in.Parallel, &out.Parallel
		*out = new(int)
		**out = **in
	}
	if in.DeletionPropagationPolicy != nil {
		in, out := &in.DeletionPropagationPolicy, &out.DeletionPropagationPolicy
		*out = new(v1.DeletionPropagation)
		**out = **in
	}
	if in.NamespaceTemplate != nil {
		in, out := &in.NamespaceTemplate, &out.NamespaceTemplate
		*out = (*in).DeepCopy()
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
	if in.DelayBeforeCleanup != nil {
		in, out := &in.DelayBeforeCleanup, &out.DelayBeforeCleanup
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Clusters != nil {
		in, out := &in.Clusters, &out.Clusters
		*out = make(map[string]Cluster, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Catch != nil {
		in, out := &in.Catch, &out.Catch
		*out = make([]Catch, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
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
func (in *Create) DeepCopyInto(out *Create) {
	*out = *in
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Bindings != nil {
		in, out := &in.Bindings, &out.Bindings
		*out = make([]Binding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Outputs != nil {
		in, out := &in.Outputs, &out.Outputs
		*out = make([]Output, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Clusters != nil {
		in, out := &in.Clusters, &out.Clusters
		*out = make(map[string]Cluster, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.FileRefOrResource.DeepCopyInto(&out.FileRefOrResource)
	if in.Template != nil {
		in, out := &in.Template, &out.Template
		*out = new(bool)
		**out = **in
	}
	if in.DryRun != nil {
		in, out := &in.DryRun, &out.DryRun
		*out = new(bool)
		**out = **in
	}
	if in.Expect != nil {
		in, out := &in.Expect, &out.Expect
		*out = make([]Expectation, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Create.
func (in *Create) DeepCopy() *Create {
	if in == nil {
		return nil
	}
	out := new(Create)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Delete) DeepCopyInto(out *Delete) {
	*out = *in
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Bindings != nil {
		in, out := &in.Bindings, &out.Bindings
		*out = make([]Binding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Clusters != nil {
		in, out := &in.Clusters, &out.Clusters
		*out = make(map[string]Cluster, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Template != nil {
		in, out := &in.Template, &out.Template
		*out = new(bool)
		**out = **in
	}
	in.ObjectReference.DeepCopyInto(&out.ObjectReference)
	if in.Expect != nil {
		in, out := &in.Expect, &out.Expect
		*out = make([]Expectation, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.DeletionPropagationPolicy != nil {
		in, out := &in.DeletionPropagationPolicy, &out.DeletionPropagationPolicy
		*out = new(v1.DeletionPropagation)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Delete.
func (in *Delete) DeepCopy() *Delete {
	if in == nil {
		return nil
	}
	out := new(Delete)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Deletion) DeepCopyInto(out *Deletion) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Deletion.
func (in *Deletion) DeepCopy() *Deletion {
	if in == nil {
		return nil
	}
	out := new(Deletion)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Describe) DeepCopyInto(out *Describe) {
	*out = *in
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Clusters != nil {
		in, out := &in.Clusters, &out.Clusters
		*out = make(map[string]Cluster, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	out.ResourceReference = in.ResourceReference
	out.ObjectLabelsSelector = in.ObjectLabelsSelector
	if in.ShowEvents != nil {
		in, out := &in.ShowEvents, &out.ShowEvents
		*out = new(bool)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Describe.
func (in *Describe) DeepCopy() *Describe {
	if in == nil {
		return nil
	}
	out := new(Describe)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Error) DeepCopyInto(out *Error) {
	*out = *in
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Bindings != nil {
		in, out := &in.Bindings, &out.Bindings
		*out = make([]Binding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Clusters != nil {
		in, out := &in.Clusters, &out.Clusters
		*out = make(map[string]Cluster, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.FileRefOrCheck.DeepCopyInto(&out.FileRefOrCheck)
	if in.Template != nil {
		in, out := &in.Template, &out.Template
		*out = new(bool)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Error.
func (in *Error) DeepCopy() *Error {
	if in == nil {
		return nil
	}
	out := new(Error)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Events) DeepCopyInto(out *Events) {
	*out = *in
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Clusters != nil {
		in, out := &in.Clusters, &out.Clusters
		*out = make(map[string]Cluster, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	out.ObjectLabelsSelector = in.ObjectLabelsSelector
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Events.
func (in *Events) DeepCopy() *Events {
	if in == nil {
		return nil
	}
	out := new(Events)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Expectation) DeepCopyInto(out *Expectation) {
	*out = *in
	if in.Match != nil {
		in, out := &in.Match, &out.Match
		*out = (*in).DeepCopy()
	}
	in.Check.DeepCopyInto(&out.Check)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Expectation.
func (in *Expectation) DeepCopy() *Expectation {
	if in == nil {
		return nil
	}
	out := new(Expectation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FileRef) DeepCopyInto(out *FileRef) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FileRef.
func (in *FileRef) DeepCopy() *FileRef {
	if in == nil {
		return nil
	}
	out := new(FileRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FileRefOrCheck) DeepCopyInto(out *FileRefOrCheck) {
	*out = *in
	out.FileRef = in.FileRef
	if in.Check != nil {
		in, out := &in.Check, &out.Check
		*out = (*in).DeepCopy()
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FileRefOrCheck.
func (in *FileRefOrCheck) DeepCopy() *FileRefOrCheck {
	if in == nil {
		return nil
	}
	out := new(FileRefOrCheck)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FileRefOrResource) DeepCopyInto(out *FileRefOrResource) {
	*out = *in
	out.FileRef = in.FileRef
	if in.Resource != nil {
		in, out := &in.Resource, &out.Resource
		*out = (*in).DeepCopy()
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FileRefOrResource.
func (in *FileRefOrResource) DeepCopy() *FileRefOrResource {
	if in == nil {
		return nil
	}
	out := new(FileRefOrResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Finally) DeepCopyInto(out *Finally) {
	*out = *in
	if in.PodLogs != nil {
		in, out := &in.PodLogs, &out.PodLogs
		*out = new(PodLogs)
		(*in).DeepCopyInto(*out)
	}
	if in.Events != nil {
		in, out := &in.Events, &out.Events
		*out = new(Events)
		(*in).DeepCopyInto(*out)
	}
	if in.Describe != nil {
		in, out := &in.Describe, &out.Describe
		*out = new(Describe)
		(*in).DeepCopyInto(*out)
	}
	if in.Wait != nil {
		in, out := &in.Wait, &out.Wait
		*out = new(Wait)
		(*in).DeepCopyInto(*out)
	}
	if in.Get != nil {
		in, out := &in.Get, &out.Get
		*out = new(Get)
		(*in).DeepCopyInto(*out)
	}
	if in.Delete != nil {
		in, out := &in.Delete, &out.Delete
		*out = new(Delete)
		(*in).DeepCopyInto(*out)
	}
	if in.Command != nil {
		in, out := &in.Command, &out.Command
		*out = new(Command)
		(*in).DeepCopyInto(*out)
	}
	if in.Script != nil {
		in, out := &in.Script, &out.Script
		*out = new(Script)
		(*in).DeepCopyInto(*out)
	}
	if in.Sleep != nil {
		in, out := &in.Sleep, &out.Sleep
		*out = new(Sleep)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Finally.
func (in *Finally) DeepCopy() *Finally {
	if in == nil {
		return nil
	}
	out := new(Finally)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *For) DeepCopyInto(out *For) {
	*out = *in
	if in.Deletion != nil {
		in, out := &in.Deletion, &out.Deletion
		*out = new(Deletion)
		**out = **in
	}
	if in.Condition != nil {
		in, out := &in.Condition, &out.Condition
		*out = new(Condition)
		(*in).DeepCopyInto(*out)
	}
	if in.JsonPath != nil {
		in, out := &in.JsonPath, &out.JsonPath
		*out = new(JsonPath)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new For.
func (in *For) DeepCopy() *For {
	if in == nil {
		return nil
	}
	out := new(For)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Get) DeepCopyInto(out *Get) {
	*out = *in
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Clusters != nil {
		in, out := &in.Clusters, &out.Clusters
		*out = make(map[string]Cluster, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	out.ResourceReference = in.ResourceReference
	out.ObjectLabelsSelector = in.ObjectLabelsSelector
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Get.
func (in *Get) DeepCopy() *Get {
	if in == nil {
		return nil
	}
	out := new(Get)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JsonPath) DeepCopyInto(out *JsonPath) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JsonPath.
func (in *JsonPath) DeepCopy() *JsonPath {
	if in == nil {
		return nil
	}
	out := new(JsonPath)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectLabelsSelector) DeepCopyInto(out *ObjectLabelsSelector) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectLabelsSelector.
func (in *ObjectLabelsSelector) DeepCopy() *ObjectLabelsSelector {
	if in == nil {
		return nil
	}
	out := new(ObjectLabelsSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectReference) DeepCopyInto(out *ObjectReference) {
	*out = *in
	out.ObjectType = in.ObjectType
	in.ObjectSelector.DeepCopyInto(&out.ObjectSelector)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectReference.
func (in *ObjectReference) DeepCopy() *ObjectReference {
	if in == nil {
		return nil
	}
	out := new(ObjectReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectSelector) DeepCopyInto(out *ObjectSelector) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectSelector.
func (in *ObjectSelector) DeepCopy() *ObjectSelector {
	if in == nil {
		return nil
	}
	out := new(ObjectSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectType) DeepCopyInto(out *ObjectType) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectType.
func (in *ObjectType) DeepCopy() *ObjectType {
	if in == nil {
		return nil
	}
	out := new(ObjectType)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Operation) DeepCopyInto(out *Operation) {
	*out = *in
	in.OperationBase.DeepCopyInto(&out.OperationBase)
	if in.Apply != nil {
		in, out := &in.Apply, &out.Apply
		*out = new(Apply)
		(*in).DeepCopyInto(*out)
	}
	if in.Assert != nil {
		in, out := &in.Assert, &out.Assert
		*out = new(Assert)
		(*in).DeepCopyInto(*out)
	}
	if in.Command != nil {
		in, out := &in.Command, &out.Command
		*out = new(Command)
		(*in).DeepCopyInto(*out)
	}
	if in.Create != nil {
		in, out := &in.Create, &out.Create
		*out = new(Create)
		(*in).DeepCopyInto(*out)
	}
	if in.Delete != nil {
		in, out := &in.Delete, &out.Delete
		*out = new(Delete)
		(*in).DeepCopyInto(*out)
	}
	if in.Error != nil {
		in, out := &in.Error, &out.Error
		*out = new(Error)
		(*in).DeepCopyInto(*out)
	}
	if in.Patch != nil {
		in, out := &in.Patch, &out.Patch
		*out = new(Patch)
		(*in).DeepCopyInto(*out)
	}
	if in.Script != nil {
		in, out := &in.Script, &out.Script
		*out = new(Script)
		(*in).DeepCopyInto(*out)
	}
	if in.Sleep != nil {
		in, out := &in.Sleep, &out.Sleep
		*out = new(Sleep)
		**out = **in
	}
	if in.Update != nil {
		in, out := &in.Update, &out.Update
		*out = new(Update)
		(*in).DeepCopyInto(*out)
	}
	if in.Wait != nil {
		in, out := &in.Wait, &out.Wait
		*out = new(Wait)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Operation.
func (in *Operation) DeepCopy() *Operation {
	if in == nil {
		return nil
	}
	out := new(Operation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OperationBase) DeepCopyInto(out *OperationBase) {
	*out = *in
	if in.ContinueOnError != nil {
		in, out := &in.ContinueOnError, &out.ContinueOnError
		*out = new(bool)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OperationBase.
func (in *OperationBase) DeepCopy() *OperationBase {
	if in == nil {
		return nil
	}
	out := new(OperationBase)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Output) DeepCopyInto(out *Output) {
	*out = *in
	in.Binding.DeepCopyInto(&out.Binding)
	if in.Match != nil {
		in, out := &in.Match, &out.Match
		*out = (*in).DeepCopy()
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Output.
func (in *Output) DeepCopy() *Output {
	if in == nil {
		return nil
	}
	out := new(Output)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Patch) DeepCopyInto(out *Patch) {
	*out = *in
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Bindings != nil {
		in, out := &in.Bindings, &out.Bindings
		*out = make([]Binding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Outputs != nil {
		in, out := &in.Outputs, &out.Outputs
		*out = make([]Output, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Clusters != nil {
		in, out := &in.Clusters, &out.Clusters
		*out = make(map[string]Cluster, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.FileRefOrResource.DeepCopyInto(&out.FileRefOrResource)
	if in.Template != nil {
		in, out := &in.Template, &out.Template
		*out = new(bool)
		**out = **in
	}
	if in.DryRun != nil {
		in, out := &in.DryRun, &out.DryRun
		*out = new(bool)
		**out = **in
	}
	if in.Expect != nil {
		in, out := &in.Expect, &out.Expect
		*out = make([]Expectation, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Patch.
func (in *Patch) DeepCopy() *Patch {
	if in == nil {
		return nil
	}
	out := new(Patch)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodLogs) DeepCopyInto(out *PodLogs) {
	*out = *in
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Clusters != nil {
		in, out := &in.Clusters, &out.Clusters
		*out = make(map[string]Cluster, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	out.ObjectLabelsSelector = in.ObjectLabelsSelector
	if in.Tail != nil {
		in, out := &in.Tail, &out.Tail
		*out = new(int)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodLogs.
func (in *PodLogs) DeepCopy() *PodLogs {
	if in == nil {
		return nil
	}
	out := new(PodLogs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceReference) DeepCopyInto(out *ResourceReference) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceReference.
func (in *ResourceReference) DeepCopy() *ResourceReference {
	if in == nil {
		return nil
	}
	out := new(ResourceReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Script) DeepCopyInto(out *Script) {
	*out = *in
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Bindings != nil {
		in, out := &in.Bindings, &out.Bindings
		*out = make([]Binding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Outputs != nil {
		in, out := &in.Outputs, &out.Outputs
		*out = make([]Output, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]Binding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Clusters != nil {
		in, out := &in.Clusters, &out.Clusters
		*out = make(map[string]Cluster, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Check != nil {
		in, out := &in.Check, &out.Check
		*out = (*in).DeepCopy()
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Script.
func (in *Script) DeepCopy() *Script {
	if in == nil {
		return nil
	}
	out := new(Script)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Sleep) DeepCopyInto(out *Sleep) {
	*out = *in
	out.Duration = in.Duration
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Sleep.
func (in *Sleep) DeepCopy() *Sleep {
	if in == nil {
		return nil
	}
	out := new(Sleep)
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
func (in *TestSpec) DeepCopyInto(out *TestSpec) {
	*out = *in
	if in.Timeouts != nil {
		in, out := &in.Timeouts, &out.Timeouts
		*out = new(Timeouts)
		(*in).DeepCopyInto(*out)
	}
	if in.Clusters != nil {
		in, out := &in.Clusters, &out.Clusters
		*out = make(map[string]Cluster, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Skip != nil {
		in, out := &in.Skip, &out.Skip
		*out = new(bool)
		**out = **in
	}
	if in.Concurrent != nil {
		in, out := &in.Concurrent, &out.Concurrent
		*out = new(bool)
		**out = **in
	}
	if in.SkipDelete != nil {
		in, out := &in.SkipDelete, &out.SkipDelete
		*out = new(bool)
		**out = **in
	}
	if in.Template != nil {
		in, out := &in.Template, &out.Template
		*out = new(bool)
		**out = **in
	}
	if in.NamespaceTemplate != nil {
		in, out := &in.NamespaceTemplate, &out.NamespaceTemplate
		*out = (*in).DeepCopy()
	}
	if in.Bindings != nil {
		in, out := &in.Bindings, &out.Bindings
		*out = make([]Binding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Steps != nil {
		in, out := &in.Steps, &out.Steps
		*out = make([]TestStep, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Catch != nil {
		in, out := &in.Catch, &out.Catch
		*out = make([]Catch, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ForceTerminationGracePeriod != nil {
		in, out := &in.ForceTerminationGracePeriod, &out.ForceTerminationGracePeriod
		*out = new(v1.Duration)
		**out = **in
	}
	if in.DelayBeforeCleanup != nil {
		in, out := &in.DelayBeforeCleanup, &out.DelayBeforeCleanup
		*out = new(v1.Duration)
		**out = **in
	}
	if in.DeletionPropagationPolicy != nil {
		in, out := &in.DeletionPropagationPolicy, &out.DeletionPropagationPolicy
		*out = new(v1.DeletionPropagation)
		**out = **in
	}
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

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TestStep) DeepCopyInto(out *TestStep) {
	*out = *in
	in.TestStepSpec.DeepCopyInto(&out.TestStepSpec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TestStep.
func (in *TestStep) DeepCopy() *TestStep {
	if in == nil {
		return nil
	}
	out := new(TestStep)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TestStepSpec) DeepCopyInto(out *TestStepSpec) {
	*out = *in
	if in.Timeouts != nil {
		in, out := &in.Timeouts, &out.Timeouts
		*out = new(Timeouts)
		(*in).DeepCopyInto(*out)
	}
	if in.DeletionPropagationPolicy != nil {
		in, out := &in.DeletionPropagationPolicy, &out.DeletionPropagationPolicy
		*out = new(v1.DeletionPropagation)
		**out = **in
	}
	if in.Clusters != nil {
		in, out := &in.Clusters, &out.Clusters
		*out = make(map[string]Cluster, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.SkipDelete != nil {
		in, out := &in.SkipDelete, &out.SkipDelete
		*out = new(bool)
		**out = **in
	}
	if in.Template != nil {
		in, out := &in.Template, &out.Template
		*out = new(bool)
		**out = **in
	}
	if in.Bindings != nil {
		in, out := &in.Bindings, &out.Bindings
		*out = make([]Binding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Try != nil {
		in, out := &in.Try, &out.Try
		*out = make([]Operation, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Catch != nil {
		in, out := &in.Catch, &out.Catch
		*out = make([]Catch, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Finally != nil {
		in, out := &in.Finally, &out.Finally
		*out = make([]Finally, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Cleanup != nil {
		in, out := &in.Cleanup, &out.Cleanup
		*out = make([]Finally, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TestStepSpec.
func (in *TestStepSpec) DeepCopy() *TestStepSpec {
	if in == nil {
		return nil
	}
	out := new(TestStepSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Timeouts) DeepCopyInto(out *Timeouts) {
	*out = *in
	if in.Apply != nil {
		in, out := &in.Apply, &out.Apply
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Assert != nil {
		in, out := &in.Assert, &out.Assert
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Cleanup != nil {
		in, out := &in.Cleanup, &out.Cleanup
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Delete != nil {
		in, out := &in.Delete, &out.Delete
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Error != nil {
		in, out := &in.Error, &out.Error
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Exec != nil {
		in, out := &in.Exec, &out.Exec
		*out = new(v1.Duration)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Timeouts.
func (in *Timeouts) DeepCopy() *Timeouts {
	if in == nil {
		return nil
	}
	out := new(Timeouts)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Update) DeepCopyInto(out *Update) {
	*out = *in
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Bindings != nil {
		in, out := &in.Bindings, &out.Bindings
		*out = make([]Binding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Outputs != nil {
		in, out := &in.Outputs, &out.Outputs
		*out = make([]Output, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Clusters != nil {
		in, out := &in.Clusters, &out.Clusters
		*out = make(map[string]Cluster, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.FileRefOrResource.DeepCopyInto(&out.FileRefOrResource)
	if in.Template != nil {
		in, out := &in.Template, &out.Template
		*out = new(bool)
		**out = **in
	}
	if in.DryRun != nil {
		in, out := &in.DryRun, &out.DryRun
		*out = new(bool)
		**out = **in
	}
	if in.Expect != nil {
		in, out := &in.Expect, &out.Expect
		*out = make([]Expectation, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Update.
func (in *Update) DeepCopy() *Update {
	if in == nil {
		return nil
	}
	out := new(Update)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Wait) DeepCopyInto(out *Wait) {
	*out = *in
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Clusters != nil {
		in, out := &in.Clusters, &out.Clusters
		*out = make(map[string]Cluster, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	out.ResourceReference = in.ResourceReference
	out.ObjectLabelsSelector = in.ObjectLabelsSelector
	in.For.DeepCopyInto(&out.For)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Wait.
func (in *Wait) DeepCopy() *Wait {
	if in == nil {
		return nil
	}
	out := new(Wait)
	in.DeepCopyInto(out)
	return out
}
