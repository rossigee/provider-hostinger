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

	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
)

// FirewallProtocol represents the network protocol
// +kubebuilder:validation:Enum=tcp;udp;icmp
type FirewallProtocol string

const (
	// FirewallProtocolTCP is TCP protocol
	FirewallProtocolTCP FirewallProtocol = "tcp"
	// FirewallProtocolUDP is UDP protocol
	FirewallProtocolUDP FirewallProtocol = "udp"
	// FirewallProtocolICMP is ICMP protocol
	FirewallProtocolICMP FirewallProtocol = "icmp"
)

// FirewallDirection represents the traffic direction
// +kubebuilder:validation:Enum=inbound;outbound
type FirewallDirection string

const (
	// FirewallDirectionInbound is inbound traffic
	FirewallDirectionInbound FirewallDirection = "inbound"
	// FirewallDirectionOutbound is outbound traffic
	FirewallDirectionOutbound FirewallDirection = "outbound"
)

// FirewallAction represents the action for matching traffic
// +kubebuilder:validation:Enum=allow;deny
type FirewallAction string

const (
	// FirewallActionAllow allows matching traffic
	FirewallActionAllow FirewallAction = "allow"
	// FirewallActionDeny denies matching traffic
	FirewallActionDeny FirewallAction = "deny"
)

// FirewallRuleSpec represents a single firewall rule
type FirewallRuleSpec struct {
	// Port is the port number or port range (e.g., "80" or "8000-9000").
	// +kubebuilder:validation:Required
	Port string `json:"port"`

	// Protocol is the network protocol.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=tcp;udp;icmp
	Protocol FirewallProtocol `json:"protocol"`

	// Direction is the traffic direction.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=inbound;outbound
	Direction FirewallDirection `json:"direction"`

	// Action is the action for matching traffic.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=allow;deny
	Action *FirewallAction `json:"action,omitempty"`

	// Source is the source IP/CIDR for inbound rules.
	// +kubebuilder:validation:Optional
	Source *string `json:"source,omitempty"`

	// Destination is the destination IP/CIDR for outbound rules.
	// +kubebuilder:validation:Optional
	Destination *string `json:"destination,omitempty"`
}

// FirewallRuleParameters are the configurable fields of a Hostinger Firewall Rule.
type FirewallRuleParameters struct {
	// InstanceID is the ID of the VPS instance to configure firewall for.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	InstanceID string `json:"instanceId"`

	// Rules is the list of firewall rules.
	// +kubebuilder:validation:Optional
	Rules []FirewallRuleSpec `json:"rules,omitempty"`

	// DefaultAction is the default action for traffic not matching any rules.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=allow;deny
	DefaultAction *FirewallAction `json:"defaultAction,omitempty"`
}

// FirewallRuleObservation are the observable fields of a Hostinger Firewall Rule.
type FirewallRuleObservation struct {
	// ID is the external firewall configuration ID.
	ID string `json:"id,omitempty"`

	// Status is the current status of the firewall (active, pending, etc.).
	Status string `json:"status,omitempty"`

	// AppliedDate is when the firewall rules were last applied.
	AppliedDate *metav1.Time `json:"appliedDate,omitempty"`

	// RuleCount is the number of active rules.
	RuleCount *int32 `json:"ruleCount,omitempty"`

	// CurrentDefaultAction is the current default action.
	CurrentDefaultAction *FirewallAction `json:"currentDefaultAction,omitempty"`
}

// FirewallRuleSpec defines the desired state of a Hostinger Firewall Rule.
type FirewallSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       FirewallRuleParameters `json:"forProvider"`
}

// FirewallRuleStatus defines the observed state of a Hostinger Firewall Rule.
type FirewallStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          FirewallRuleObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,hostinger}
// +kubebuilder:printcolumn:name="READY",type=string,JSONPath=.status.conditions[?(@.type=='Ready')].status
// +kubebuilder:printcolumn:name="SYNCED",type=string,JSONPath=.status.conditions[?(@.type=='Synced')].status
// +kubebuilder:printcolumn:name="AGE",type=date,JSONPath=.metadata.creationTimestamp
// +genclient

// FirewallRule is the CRD type for Hostinger firewall rules.
type FirewallRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FirewallSpec   `json:"spec,omitempty"`
	Status FirewallStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FirewallRuleList contains a list of FirewallRule resources.
type FirewallRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FirewallRule `json:"items"`
}
