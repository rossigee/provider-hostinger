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

package auth

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	v1beta1 "github.com/rossigee/provider-hostinger/apis/v1beta1"
)

// CreateAuthenticator creates an Authenticator from ProviderConfig credentials
func CreateAuthenticator(ctx context.Context, k8sClient client.Client, config *v1beta1.ProviderConfig) (Authenticator, error) {
	// Determine which auth method is configured
	if config.Spec.APIKeyAuth != nil {
		return createV1KeyAuth(ctx, k8sClient, config)
	} else if config.Spec.OAuthAuth != nil {
		return createV2OAuthAuth(ctx, k8sClient, config)
	}

	return nil, fmt.Errorf("no authentication method configured in ProviderConfig")
}

// createV1KeyAuth creates a V1KeyAuth authenticator from ProviderConfig
func createV1KeyAuth(ctx context.Context, k8sClient client.Client, config *v1beta1.ProviderConfig) (Authenticator, error) {
	authSpec := config.Spec.APIKeyAuth

	// Get API key from secret
	apiKey, err := getSecretValue(ctx, k8sClient, config.Namespace, &authSpec.APIKeySecretRef)
	if err != nil {
		return nil, fmt.Errorf("failed to get API key from secret: %w", err)
	}

	// Get customer ID from secret
	customerID, err := getSecretValue(ctx, k8sClient, config.Namespace, &authSpec.CustomerIDSecretRef)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer ID from secret: %w", err)
	}

	// Get endpoint (default to public API if not specified)
	endpoint := authSpec.Endpoint
	if endpoint == "" {
		endpoint = "https://api.hostinger.com/v1"
	}

	return NewV1KeyAuth(apiKey, customerID, endpoint), nil
}

// createV2OAuthAuth creates a V2OAuthAuth authenticator from ProviderConfig
func createV2OAuthAuth(ctx context.Context, k8sClient client.Client, config *v1beta1.ProviderConfig) (Authenticator, error) {
	authSpec := config.Spec.OAuthAuth

	// Get client ID from secret
	clientID, err := getSecretValue(ctx, k8sClient, config.Namespace, &authSpec.ClientIDSecretRef)
	if err != nil {
		return nil, fmt.Errorf("failed to get client ID from secret: %w", err)
	}

	// Get client secret from secret
	clientSecret, err := getSecretValue(ctx, k8sClient, config.Namespace, &authSpec.ClientSecretSecretRef)
	if err != nil {
		return nil, fmt.Errorf("failed to get client secret from secret: %w", err)
	}

	// Get endpoint (default to public API if not specified)
	endpoint := authSpec.Endpoint
	if endpoint == "" {
		endpoint = "https://api.hostinger.com/v2"
	}

	// Token endpoint (default to Hostinger's OAuth endpoint)
	tokenEndpoint := authSpec.TokenEndpoint
	if tokenEndpoint == "" {
		tokenEndpoint = "https://auth.hostinger.com/oauth/token"
	}

	return NewV2OAuthAuth(clientID, clientSecret, endpoint, tokenEndpoint), nil
}

// getSecretValue retrieves a value from a Kubernetes secret
func getSecretValue(ctx context.Context, k8sClient client.Client, namespace string, secretRef *xpv1.SecretKeySelector) (string, error) {
	if secretRef == nil {
		return "", fmt.Errorf("secret reference is nil")
	}

	// Get the secret from Kubernetes
	secret := &corev1.Secret{}
	err := k8sClient.Get(ctx, types.NamespacedName{
		Namespace: namespace,
		Name:      secretRef.Name,
	}, secret)
	if err != nil {
		return "", fmt.Errorf("failed to get secret %s/%s: %w", namespace, secretRef.Name, err)
	}

	// Get the specific key from the secret
	value, ok := secret.Data[secretRef.Key]
	if !ok {
		return "", fmt.Errorf("key %q not found in secret %s/%s", secretRef.Key, namespace, secretRef.Name)
	}

	return string(value), nil
}
