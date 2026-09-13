/*
Copyright 2025 Ross Golder.

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

package v1beta1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	xpv1 "github.com/crossplane/crossplane/apis/v2/core/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InstanceParameters are the configurable fields of a Hostinger VPS Instance.
type InstanceParameters struct {
	// Hostname is the hostname for the VPS instance.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Hostname string `json:"hostname"`

	// OSId is the operating system ID/template to use.
	// +kubebuilder:validation:Required
	OSId string `json:"osId"`

	// CPUCount is the number of CPU cores.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	CPUCount int32 `json:"cpuCount"`

	// RAM is the amount of RAM in MB.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=512
	RAM int32 `json:"ram"`

	// DiskSize is the disk size in GB.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=10
	DiskSize int32 `json:"diskSize"`

	// Bandwidth is the bandwidth in GB/month.
	// +kubebuilder:validation:Optional
	Bandwidth *int32 `json:"bandwidth,omitempty"`

	// IPv6Enabled specifies whether IPv6 should be enabled.
	// +kubebuilder:validation:Optional
	IPv6Enabled *bool `json:"ipv6Enabled,omitempty"`

	// Inodes is the number of inodes.
	// +kubebuilder:validation:Optional
	Inodes *int32 `json:"inodes,omitempty"`

	// RootPasswordSecretRef is a reference to a secret containing the root password.
	// +kubebuilder:validation:Optional
	RootPasswordSecretRef *xpv1.SecretKeySelector `json:"rootPasswordSecretRef,omitempty"`
}

// InstanceObservation are the observable fields of a Hostinger VPS Instance.
type InstanceObservation struct {
	// ID is the external resource ID.
	ID string `json:"id,omitempty"`

	// Status is the current status of the instance (active, pending, suspended, etc.).
	Status string `json:"status,omitempty"`

	// IPAddress is the primary IP address.
	IPAddress string `json:"ipAddress,omitempty"`

	// IPv6Address is the IPv6 address if enabled.
	IPv6Address string `json:"ipv6Address,omitempty"`

	// CreationDate is when the instance was created.
	CreationDate *metav1.Time `json:"creationDate,omitempty"`

	// ExpirationDate is when the instance will expire.
	ExpirationDate *metav1.Time `json:"expirationDate,omitempty"`

	// CurrentHostname is the current hostname set on the instance.
	CurrentHostname string `json:"currentHostname,omitempty"`

	// CurrentCPUCount is the current CPU count.
	CurrentCPUCount int32 `json:"currentCpuCount,omitempty"`

	// CurrentRAM is the current RAM in MB.
	CurrentRAM int32 `json:"currentRam,omitempty"`

	// CurrentDiskSize is the current disk size in GB.
	CurrentDiskSize int32 `json:"currentDiskSize,omitempty"`
}

// InstanceSpec defines the desired state of a Hostinger VPS Instance.
type InstanceSpec struct {
	xpv1.ManagedResourceSpec `json:",inline"`
	ForProvider                     InstanceParameters `json:"forProvider"`
}

// InstanceStatus defines the observed state of a Hostinger VPS Instance.
type InstanceStatus struct {
	xpv1.ManagedResourceStatus `json:",inline"`
	AtProvider                 InstanceObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,hostinger}
// +kubebuilder:printcolumn:name="READY",type=string,JSONPath=.status.conditions[?(@.type=='Ready')].status
// +kubebuilder:printcolumn:name="SYNCED",type=string,JSONPath=.status.conditions[?(@.type=='Synced')].status
// +kubebuilder:printcolumn:name="AGE",type=date,JSONPath=.metadata.creationTimestamp
// +genclient

// Instance is the CRD type for Hostinger VPS instances.
type Instance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InstanceSpec   `json:"spec,omitempty"`
	Status InstanceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// InstanceList contains a list of Instance resources.
type InstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Instance `json:"items"`
}


// GetCondition gets the condition from the resource status.
func (mg *Instance) GetCondition(ct xpv1.ConditionType) xpv1.Condition {
	return mg.Status.GetCondition(ct)
}

// SetConditions sets the conditions on the resource status.
func (mg *Instance) SetConditions(c ...xpv1.Condition) {
	mg.Status.SetConditions(c...)
}

// GetManagementPolicies gets the management policies for the resource.
func (mg *Instance) GetManagementPolicies() xpv1.ManagementPolicies {
	return mg.Spec.ManagementPolicies
}

// SetManagementPolicies sets the management policies for the resource.
func (mg *Instance) SetManagementPolicies(mp xpv1.ManagementPolicies) {
	mg.Spec.ManagementPolicies = mp
}

// DeepCopyObject returns a deep copy of this object as runtime.Object.
func (in *Instance) DeepCopyObject() runtime.Object {
	out := &Instance{}
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto fills DeepCopy receiver with DeepCopy of the provided receiver.
func (in *Instance) DeepCopyInto(out *Instance) {
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// GetProviderConfigReference returns the ProviderConfig reference.
func (mg *Instance) GetProviderConfigReference() *xpv1.Reference {
	if mg.Spec.ProviderConfigReference != nil && mg.Spec.ProviderConfigReference.Name != "" {
		return &xpv1.Reference{Name: mg.Spec.ProviderConfigReference.Name}
	}
	return nil
}

// GetItems returns the list items.
func (l *InstanceList) GetItems() []resource.Managed {
	items := make([]resource.Managed, len(l.Items))
	for i := range l.Items {
		items[i] = &l.Items[i]
	}
	return items
}

// DeepCopyObject returns a deep copy of this object as runtime.Object.
func (in *InstanceList) DeepCopyObject() runtime.Object {
	out := &InstanceList{}
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto fills DeepCopy receiver with DeepCopy of the provided receiver.
func (in *InstanceList) DeepCopyInto(out *InstanceList) {
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	out.Items = append([]Instance(nil), in.Items...)
}
