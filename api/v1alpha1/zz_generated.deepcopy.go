//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2023.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ec2CostOptimizer) DeepCopyInto(out *Ec2CostOptimizer) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ec2CostOptimizer.
func (in *Ec2CostOptimizer) DeepCopy() *Ec2CostOptimizer {
	if in == nil {
		return nil
	}
	out := new(Ec2CostOptimizer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Ec2CostOptimizer) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ec2CostOptimizerList) DeepCopyInto(out *Ec2CostOptimizerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Ec2CostOptimizer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ec2CostOptimizerList.
func (in *Ec2CostOptimizerList) DeepCopy() *Ec2CostOptimizerList {
	if in == nil {
		return nil
	}
	out := new(Ec2CostOptimizerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Ec2CostOptimizerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ec2CostOptimizerSpec) DeepCopyInto(out *Ec2CostOptimizerSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ec2CostOptimizerSpec.
func (in *Ec2CostOptimizerSpec) DeepCopy() *Ec2CostOptimizerSpec {
	if in == nil {
		return nil
	}
	out := new(Ec2CostOptimizerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ec2CostOptimizerStatus) DeepCopyInto(out *Ec2CostOptimizerStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ec2CostOptimizerStatus.
func (in *Ec2CostOptimizerStatus) DeepCopy() *Ec2CostOptimizerStatus {
	if in == nil {
		return nil
	}
	out := new(Ec2CostOptimizerStatus)
	in.DeepCopyInto(out)
	return out
}