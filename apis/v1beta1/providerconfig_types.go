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

// APIKeyAuthSpec contains API v1 authentication credentials.
type APIKeyAuthSpec struct {
	// Endpoint is the Hostinger API v1 endpoint URL.
	// +kubebuilder:validation:Required
	Endpoint string `json:"endpoint"`

	// APIKeySecretRef is a reference to a secret containing the API key.
	// The secret key should be "api-key".
	// +kubebuilder:validation:Required
	APIKeySecretRef xpv1.SecretKeySelector `json:"apiKeySecretRef"`

	// CustomerIDSecretRef is a reference to a secret containing the customer ID.
	// The secret key should be "customer-id".
	// +kubebuilder:validation:Required
	CustomerIDSecretRef xpv1.SecretKeySelector `json:"customerIdSecretRef"`
}

// OAuthAuthSpec contains API v2 OAuth authentication credentials.
type OAuthAuthSpec struct {
	// Endpoint is the Hostinger API v2 endpoint URL.
	// +kubebuilder:validation:Required
	Endpoint string `json:"endpoint"`

	// ClientIDSecretRef is a reference to a secret containing the OAuth client ID.
	// The secret key should be "client-id".
	// +kubebuilder:validation:Required
	ClientIDSecretRef xpv1.SecretKeySelector `json:"clientIdSecretRef"`

	// ClientSecretSecretRef is a reference to a secret containing the OAuth client secret.
	// The secret key should be "client-secret".
	// +kubebuilder:validation:Required
	ClientSecretSecretRef xpv1.SecretKeySelector `json:"clientSecretSecretRef"`

	// TokenEndpoint is the OAuth token endpoint URL.
	// +kubebuilder:validation:Required
	TokenEndpoint string `json:"tokenEndpoint"`
}

// ProviderConfigSpec defines the desired state of a ProviderConfig.
type ProviderConfigSpec struct {
	// Credentials specifies the authentication method to use.
	// Only one of the following may be specified.
	// If none of the following are specified, the identity will be assumed
	// from the environment.

	// APIKeyAuth contains API v1 (API key) authentication credentials.
	// +kubebuilder:validation:Optional
	APIKeyAuth *APIKeyAuthSpec `json:"apiKeyAuth,omitempty"`

	// OAuthAuth contains API v2 (OAuth) authentication credentials.
	// +kubebuilder:validation:Optional
	OAuthAuth *OAuthAuthSpec `json:"oauthAuth,omitempty"`
}

// ProviderConfigStatus defines the observed state of a ProviderConfig.
type ProviderConfigStatus struct {
	xpv1.ProviderConfigStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="AGE",type=date,JSONPath=.metadata.creationTimestamp
// +kubebuilder:printcolumn:name="READY",type=string,JSONPath=.status.conditions[?(@.type=='Ready')].status
// +genclient
// +genclient:nonNamespaced

// ProviderConfig is the CRD type for Hostinger API provider configurations.
type ProviderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProviderConfigSpec   `json:"spec,omitempty"`
	Status ProviderConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProviderConfigList contains a list of ProviderConfig.
type ProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProviderConfig `json:"items"`
}
