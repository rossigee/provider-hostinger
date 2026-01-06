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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// SSHKeyParameters are the configurable fields of a Hostinger SSH Key.
type SSHKeyParameters struct {
	// Name is the name/label for the SSH key.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// PublicKeySecretRef is a reference to a secret containing the public key.
	// The secret key should be "public-key".
	// +kubebuilder:validation:Required
	PublicKeySecretRef xpv1.SecretKeySelector `json:"publicKeySecretRef"`

	// InstanceIDs are the instance IDs to attach this SSH key to.
	// +kubebuilder:validation:Optional
	InstanceIDs []string `json:"instanceIds,omitempty"`
}

// SSHKeyObservation are the observable fields of a Hostinger SSH Key.
type SSHKeyObservation struct {
	// ID is the external SSH key resource ID.
	ID string `json:"id,omitempty"`

	// Fingerprint is the SSH key fingerprint.
	Fingerprint string `json:"fingerprint,omitempty"`

	// CreatedDate is when the SSH key was created.
	CreatedDate *metav1.Time `json:"createdDate,omitempty"`

	// AttachedInstances is the list of instance IDs this key is attached to.
	AttachedInstances []string `json:"attachedInstances,omitempty"`

	// PublicKeyHash is a hash of the public key content.
	PublicKeyHash string `json:"publicKeyHash,omitempty"`
}

// SSHKeySpec defines the desired state of a Hostinger SSH Key.
type SSHKeySpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       SSHKeyParameters `json:"forProvider"`
}

// SSHKeyStatus defines the observed state of a Hostinger SSH Key.
type SSHKeyStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          SSHKeyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,hostinger}
// +kubebuilder:printcolumn:name="READY",type=string,JSONPath=.status.conditions[?(@.type=='Ready')].status
// +kubebuilder:printcolumn:name="SYNCED",type=string,JSONPath=.status.conditions[?(@.type=='Synced')].status
// +kubebuilder:printcolumn:name="AGE",type=date,JSONPath=.metadata.creationTimestamp

// SSHKey is the CRD type for Hostinger SSH keys.
type SSHKey struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SSHKeySpec   `json:"spec,omitempty"`
	Status SSHKeyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SSHKeyList contains a list of SSHKey resources.
type SSHKeyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SSHKey `json:"items"`
}
