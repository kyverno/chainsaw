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
	out.FileRef = in.FileRef
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		*out = new(v1.Duration)
		**out = **in
	}
	if in.ShouldFail != nil {
		in, out := &in.ShouldFail, &out.ShouldFail
		*out = new(bool)
		**out = **in
	}
	if in.ContinueOnError != nil {
		in, out := &in.ContinueOnError, &out.ContinueOnError
		*out = new(bool)
		**out = **in
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
	out.FileRef = in.FileRef
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		*out = new(v1.Duration)
		**out = **in
	}
	if in.ContinueOnError != nil {
		in, out := &in.ContinueOnError, &out.ContinueOnError
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
func (in *Collect) DeepCopyInto(out *Collect) {
	*out = *in
	if in.PodLogs != nil {
		in, out := &in.PodLogs, &out.PodLogs
		*out = new(PodLogs)
		(*in).DeepCopyInto(*out)
	}
	if in.Events != nil {
		in, out := &in.Events, &out.Events
		*out = new(Events)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Collect.
func (in *Collect) DeepCopy() *Collect {
	if in == nil {
		return nil
	}
	out := new(Collect)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Command) DeepCopyInto(out *Command) {
	*out = *in
	if in.Args != nil {
		in, out := &in.Args, &out.Args
		*out = make([]string, len(*in))
		copy(*out, *in)
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
	if in.TestDirs != nil {
		in, out := &in.TestDirs, &out.TestDirs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
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
func (in *Delete) DeepCopyInto(out *Delete) {
	*out = *in
	in.ObjectReference.DeepCopyInto(&out.ObjectReference)
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		*out = new(v1.Duration)
		**out = **in
	}
	if in.ContinueOnError != nil {
		in, out := &in.ContinueOnError, &out.ContinueOnError
		*out = new(bool)
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
func (in *Error) DeepCopyInto(out *Error) {
	*out = *in
	out.FileRef = in.FileRef
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		*out = new(v1.Duration)
		**out = **in
	}
	if in.ContinueOnError != nil {
		in, out := &in.ContinueOnError, &out.ContinueOnError
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
func (in *Exec) DeepCopyInto(out *Exec) {
	*out = *in
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Command != nil {
		in, out := &in.Command, &out.Command
		*out = new(Command)
		(*in).DeepCopyInto(*out)
	}
	if in.Script != nil {
		in, out := &in.Script, &out.Script
		*out = new(Script)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Exec.
func (in *Exec) DeepCopy() *Exec {
	if in == nil {
		return nil
	}
	out := new(Exec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExecOperation) DeepCopyInto(out *ExecOperation) {
	*out = *in
	in.Exec.DeepCopyInto(&out.Exec)
	if in.ContinueOnError != nil {
		in, out := &in.ContinueOnError, &out.ContinueOnError
		*out = new(bool)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExecOperation.
func (in *ExecOperation) DeepCopy() *ExecOperation {
	if in == nil {
		return nil
	}
	out := new(ExecOperation)
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
func (in *ObjectReference) DeepCopyInto(out *ObjectReference) {
	*out = *in
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
func (in *OnFailure) DeepCopyInto(out *OnFailure) {
	*out = *in
	if in.Collect != nil {
		in, out := &in.Collect, &out.Collect
		*out = new(Collect)
		(*in).DeepCopyInto(*out)
	}
	if in.Exec != nil {
		in, out := &in.Exec, &out.Exec
		*out = new(Exec)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OnFailure.
func (in *OnFailure) DeepCopy() *OnFailure {
	if in == nil {
		return nil
	}
	out := new(OnFailure)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Operations) DeepCopyInto(out *Operations) {
	*out = *in
	if in.Assert != nil {
		in, out := &in.Assert, &out.Assert
		*out = make([]Assert, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Apply != nil {
		in, out := &in.Apply, &out.Apply
		*out = make([]Apply, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Error != nil {
		in, out := &in.Error, &out.Error
		*out = make([]Error, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Delete != nil {
		in, out := &in.Delete, &out.Delete
		*out = make([]Delete, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Exec != nil {
		in, out := &in.Exec, &out.Exec
		*out = make([]ExecOperation, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Operations.
func (in *Operations) DeepCopy() *Operations {
	if in == nil {
		return nil
	}
	out := new(Operations)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodLogs) DeepCopyInto(out *PodLogs) {
	*out = *in
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
func (in *Script) DeepCopyInto(out *Script) {
	*out = *in
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
	in.Timeouts.DeepCopyInto(&out.Timeouts)
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
	if in.Steps != nil {
		in, out := &in.Steps, &out.Steps
		*out = make([]TestSpecStep, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
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
func (in *TestSpecStep) DeepCopyInto(out *TestSpecStep) {
	*out = *in
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TestSpecStep.
func (in *TestSpecStep) DeepCopy() *TestSpecStep {
	if in == nil {
		return nil
	}
	out := new(TestSpecStep)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TestStep) DeepCopyInto(out *TestStep) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
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

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TestStep) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TestStepSpec) DeepCopyInto(out *TestStepSpec) {
	*out = *in
	in.Timeouts.DeepCopyInto(&out.Timeouts)
	if in.SkipDelete != nil {
		in, out := &in.SkipDelete, &out.SkipDelete
		*out = new(bool)
		**out = **in
	}
	in.Operations.DeepCopyInto(&out.Operations)
	if in.OnFailure != nil {
		in, out := &in.OnFailure, &out.OnFailure
		*out = make([]OnFailure, len(*in))
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
	if in.Error != nil {
		in, out := &in.Error, &out.Error
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Delete != nil {
		in, out := &in.Delete, &out.Delete
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Cleanup != nil {
		in, out := &in.Cleanup, &out.Cleanup
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
