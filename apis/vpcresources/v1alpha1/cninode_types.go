// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FeatureName is a type of feature name supported by AWS VPC CNI. It can be Security Group for Pods, custom networking, or others
type FeatureName string

const (
	SecurityGroupsForPods FeatureName = "SecurityGroupsForPods"
	CustomNetworking      FeatureName = "CustomNetworking"
)

// Feature is a type of feature being supported by VPC resource controller and other AWS Services
type Feature struct {
	Name  FeatureName `json:"name,omitempty"`
	Value string      `json:"value,omitempty"`
}

// Important: Run "make" to regenerate code after modifying this file
// CNINodeSpec defines the desired state of CNINode
type CNINodeSpec struct {
	Features []Feature `json:"features,omitempty"`
}

type WarmBranchENI struct {
	// BranchENId is the network interface id of the branch interface
	ID string `json:"id,omitempty"`
	// MacAdd is the MAC address of the network interface
	MACAddr string `json:"macAddr,omitempty"`
	// IPv4 and/or IPv6 address assigned to the branch Network interface
	IPV4Addr string `json:"ipv4Addr,omitempty"`
	IPV6Addr string `json:"ipv6Addr,omitempty"`
	// VlanId is the VlanId of the branch network interface
	VlanID int `json:"vlanId,omitempty"`
	// SubnetCIDR is the CIDR block of the subnet
	SubnetCIDR   string `json:"subnetCIDR,omitempty"`
	SubnetV6CIDR string `json:"subnetV6CIDR,omitempty"`
}

// CNINodeStatus defines the managed VPC resources.
type CNINodeStatus struct {
	//TODO: add VPC resources which will be managed by this CRD and its finalizer
	BranchENIs []WarmBranchENI `json:"branchenis,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Features",type=string,JSONPath=`.spec.features`,description="The features delegated to VPC resource controller"
// +kubebuilder:resource:shortName=cnd,scope=Cluster

// +kubebuilder:object:root=true
type CNINode struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CNINodeSpec   `json:"spec,omitempty"`
	Status            CNINodeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// CNINodeList contains a list of CNINodeList
type CNINodeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CNINode `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CNINode{}, &CNINodeList{})
}
