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
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	v1beta1 "github.com/rossigee/provider-hostinger/apis/v1beta1"
)

func TestCreateAuthenticator_V1KeyAuth(t *testing.T) {
	// Create fake K8s client with secrets
	sch := fake.NewClientBuilder().Build().Scheme()

	secrets := []client.Object{
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hostinger-creds",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"api-key":     []byte("test-api-key-12345"),
				"customer-id": []byte("cust-67890"),
			},
		},
	}

	k8sClient := fake.NewClientBuilder().
		WithScheme(sch).
		WithObjects(secrets...).
		Build()

	config := &v1beta1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: "default",
		},
		Spec: v1beta1.ProviderConfigSpec{
			APIKeyAuth: &v1beta1.APIKeyAuthSpec{
				Endpoint: "https://api.hostinger.com/v1",
				APIKeySecretRef: xpv1.SecretKeySelector{
					SecretReference: xpv1.SecretReference{
						Name:      "hostinger-creds",
						Namespace: "default",
					},
					Key: "api-key",
				},
				CustomerIDSecretRef: xpv1.SecretKeySelector{
					SecretReference: xpv1.SecretReference{
						Name:      "hostinger-creds",
						Namespace: "default",
					},
					Key: "customer-id",
				},
			},
		},
	}

	auth, err := CreateAuthenticator(context.Background(), k8sClient, config)

	if err != nil {
		t.Errorf("CreateAuthenticator() error = %v, want nil", err)
	}

	if auth == nil {
		t.Fatal("CreateAuthenticator() returned nil authenticator")
	}

	if auth.Type() != "APIKeyAuth" {
		t.Errorf("Authenticator type = %v, want APIKeyAuth", auth.Type())
	}

	if auth.GetEndpoint() != "https://api.hostinger.com/v1" {
		t.Errorf("Endpoint = %v, want https://api.hostinger.com/v1", auth.GetEndpoint())
	}
}

func TestCreateAuthenticator_V2OAuthAuth(t *testing.T) {
	sch := fake.NewClientBuilder().Build().Scheme()

	secrets := []client.Object{
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hostinger-oauth",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"client-id":     []byte("oauth-client-123"),
				"client-secret": []byte("oauth-secret-456"),
			},
		},
	}

	k8sClient := fake.NewClientBuilder().
		WithScheme(sch).
		WithObjects(secrets...).
		Build()

	config := &v1beta1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: "default",
		},
		Spec: v1beta1.ProviderConfigSpec{
			OAuthAuth: &v1beta1.OAuthAuthSpec{
				ClientIDSecretRef: xpv1.SecretKeySelector{
					SecretReference: xpv1.SecretReference{
						Name:      "hostinger-oauth",
						Namespace: "default",
					},
					Key: "client-id",
				},
				ClientSecretSecretRef: xpv1.SecretKeySelector{
					SecretReference: xpv1.SecretReference{
						Name:      "hostinger-oauth",
						Namespace: "default",
					},
					Key: "client-secret",
				},
				Endpoint:      "https://api.hostinger.com/v2",
				TokenEndpoint: "https://auth.hostinger.com/oauth/token",
			},
		},
	}

	auth, err := CreateAuthenticator(context.Background(), k8sClient, config)

	if err != nil {
		t.Errorf("CreateAuthenticator() error = %v, want nil", err)
	}

	if auth == nil {
		t.Fatal("CreateAuthenticator() returned nil authenticator")
	}

	if auth.Type() != "OAuthAuth" {
		t.Errorf("Authenticator type = %v, want OAuthAuth", auth.Type())
	}
}

func TestCreateAuthenticator_NoAuthMethod(t *testing.T) {
	k8sClient := fake.NewClientBuilder().Build()

	config := &v1beta1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: "default",
		},
		Spec: v1beta1.ProviderConfigSpec{
			// Neither APIKeyAuth nor OAuthAuth specified
		},
	}

	auth, err := CreateAuthenticator(context.Background(), k8sClient, config)

	if err == nil {
		t.Error("CreateAuthenticator() expected error for no auth method, got nil")
	}

	if auth != nil {
		t.Errorf("CreateAuthenticator() expected nil authenticator, got %v", auth)
	}
}

func TestCreateV1KeyAuth_Success(t *testing.T) {
	sch := fake.NewClientBuilder().Build().Scheme()

	secrets := []client.Object{
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "v1-creds",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"key": []byte("api-key-xyz"),
				"id":  []byte("customer-123"),
			},
		},
	}

	k8sClient := fake.NewClientBuilder().
		WithScheme(sch).
		WithObjects(secrets...).
		Build()

	config := &v1beta1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: "default",
		},
		Spec: v1beta1.ProviderConfigSpec{
			APIKeyAuth: &v1beta1.APIKeyAuthSpec{
				APIKeySecretRef: xpv1.SecretKeySelector{
					SecretReference: xpv1.SecretReference{
						Name:      "v1-creds",
						Namespace: "default",
					},
					Key: "key",
				},
				CustomerIDSecretRef: xpv1.SecretKeySelector{
					SecretReference: xpv1.SecretReference{
						Name:      "v1-creds",
						Namespace: "default",
					},
					Key: "id",
				},
			},
		},
	}

	auth, err := createV1KeyAuth(context.Background(), k8sClient, config)

	if err != nil {
		t.Errorf("createV1KeyAuth() error = %v, want nil", err)
	}

	if auth == nil {
		t.Fatal("createV1KeyAuth() returned nil")
	}

	v1auth, ok := auth.(*V1KeyAuth)
	if !ok {
		t.Errorf("createV1KeyAuth() returned %T, want *V1KeyAuth", auth)
	}

	if v1auth.APIKey != "api-key-xyz" {
		t.Errorf("APIKey = %v, want api-key-xyz", v1auth.APIKey)
	}

	if v1auth.CustomerID != "customer-123" {
		t.Errorf("CustomerID = %v, want customer-123", v1auth.CustomerID)
	}
}

func TestCreateV1KeyAuth_DefaultEndpoint(t *testing.T) {
	sch := fake.NewClientBuilder().Build().Scheme()

	secrets := []client.Object{
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "v1-creds",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"key": []byte("api-key"),
				"id":  []byte("customer-id"),
			},
		},
	}

	k8sClient := fake.NewClientBuilder().
		WithScheme(sch).
		WithObjects(secrets...).
		Build()

	config := &v1beta1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: "default",
		},
		Spec: v1beta1.ProviderConfigSpec{
			APIKeyAuth: &v1beta1.APIKeyAuthSpec{
				APIKeySecretRef: xpv1.SecretKeySelector{
					SecretReference: xpv1.SecretReference{
						Name:      "v1-creds",
						Namespace: "default",
					},
					Key: "key",
				},
				CustomerIDSecretRef: xpv1.SecretKeySelector{
					SecretReference: xpv1.SecretReference{
						Name:      "v1-creds",
						Namespace: "default",
					},
					Key: "id",
				},
				// Endpoint not specified
			},
		},
	}

	auth, err := createV1KeyAuth(context.Background(), k8sClient, config)

	if err != nil {
		t.Errorf("createV1KeyAuth() error = %v, want nil", err)
	}

	if auth.GetEndpoint() != "https://api.hostinger.com/v1" {
		t.Errorf("GetEndpoint() = %v, want https://api.hostinger.com/v1", auth.GetEndpoint())
	}
}

func TestCreateV1KeyAuth_MissingAPIKeySecret(t *testing.T) {
	sch := fake.NewClientBuilder().Build().Scheme()
	k8sClient := fake.NewClientBuilder().
		WithScheme(sch).
		Build()

	config := &v1beta1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: "default",
		},
		Spec: v1beta1.ProviderConfigSpec{
			APIKeyAuth: &v1beta1.APIKeyAuthSpec{
				APIKeySecretRef: xpv1.SecretKeySelector{
					SecretReference: xpv1.SecretReference{
						Name:      "nonexistent",
						Namespace: "default",
					},
					Key: "key",
				},
				CustomerIDSecretRef: xpv1.SecretKeySelector{
					SecretReference: xpv1.SecretReference{
						Name:      "creds",
						Namespace: "default",
					},
					Key: "id",
				},
			},
		},
	}

	auth, err := createV1KeyAuth(context.Background(), k8sClient, config)

	if err == nil {
		t.Error("createV1KeyAuth() expected error for missing secret, got nil")
	}

	if auth != nil {
		t.Errorf("createV1KeyAuth() expected nil authenticator, got %v", auth)
	}
}

func TestCreateV1KeyAuth_MissingCustomerIDSecret(t *testing.T) {
	sch := fake.NewClientBuilder().Build().Scheme()

	secrets := []client.Object{
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "v1-creds",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"key": []byte("api-key"),
			},
		},
	}

	k8sClient := fake.NewClientBuilder().
		WithScheme(sch).
		WithObjects(secrets...).
		Build()

	config := &v1beta1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: "default",
		},
		Spec: v1beta1.ProviderConfigSpec{
			APIKeyAuth: &v1beta1.APIKeyAuthSpec{
				APIKeySecretRef: xpv1.SecretKeySelector{
					SecretReference: xpv1.SecretReference{
						Name:      "v1-creds",
						Namespace: "default",
					},
					Key: "key",
				},
				CustomerIDSecretRef: xpv1.SecretKeySelector{
					SecretReference: xpv1.SecretReference{
						Name:      "nonexistent",
						Namespace: "default",
					},
					Key: "id",
				},
			},
		},
	}

	auth, err := createV1KeyAuth(context.Background(), k8sClient, config)

	if err == nil {
		t.Error("createV1KeyAuth() expected error for missing customer ID secret, got nil")
	}

	if auth != nil {
		t.Errorf("createV1KeyAuth() expected nil authenticator, got %v", auth)
	}
}

func TestCreateV2OAuthAuth_Success(t *testing.T) {
	sch := fake.NewClientBuilder().Build().Scheme()

	secrets := []client.Object{
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "oauth-creds",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"client_id":     []byte("oauth-id-123"),
				"client_secret": []byte("oauth-secret-456"),
			},
		},
	}

	k8sClient := fake.NewClientBuilder().
		WithScheme(sch).
		WithObjects(secrets...).
		Build()

	config := &v1beta1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: "default",
		},
		Spec: v1beta1.ProviderConfigSpec{
			OAuthAuth: &v1beta1.OAuthAuthSpec{
				ClientIDSecretRef: xpv1.SecretKeySelector{
					SecretReference: xpv1.SecretReference{
						Name:      "oauth-creds",
						Namespace: "default",
					},
					Key: "client_id",
				},
				ClientSecretSecretRef: xpv1.SecretKeySelector{
					SecretReference: xpv1.SecretReference{
						Name:      "oauth-creds",
						Namespace: "default",
					},
					Key: "client_secret",
				},
				Endpoint:      "https://api.hostinger.com/v2",
				TokenEndpoint: "https://auth.hostinger.com/oauth/token",
			},
		},
	}

	auth, err := createV2OAuthAuth(context.Background(), k8sClient, config)

	if err != nil {
		t.Errorf("createV2OAuthAuth() error = %v, want nil", err)
	}

	if auth == nil {
		t.Fatal("createV2OAuthAuth() returned nil")
	}

	oauthAuth, ok := auth.(*V2OAuthAuth)
	if !ok {
		t.Errorf("createV2OAuthAuth() returned %T, want *V2OAuthAuth", auth)
	}

	if oauthAuth.ClientID != "oauth-id-123" {
		t.Errorf("ClientID = %v, want oauth-id-123", oauthAuth.ClientID)
	}

	if oauthAuth.ClientSecret != "oauth-secret-456" {
		t.Errorf("ClientSecret = %v, want oauth-secret-456", oauthAuth.ClientSecret)
	}
}

func TestCreateV2OAuthAuth_DefaultEndpoints(t *testing.T) {
	sch := fake.NewClientBuilder().Build().Scheme()

	secrets := []client.Object{
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "oauth-creds",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"id":     []byte("client-id"),
				"secret": []byte("client-secret"),
			},
		},
	}

	k8sClient := fake.NewClientBuilder().
		WithScheme(sch).
		WithObjects(secrets...).
		Build()

	config := &v1beta1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: "default",
		},
		Spec: v1beta1.ProviderConfigSpec{
			OAuthAuth: &v1beta1.OAuthAuthSpec{
				ClientIDSecretRef: xpv1.SecretKeySelector{
					SecretReference: xpv1.SecretReference{
						Name:      "oauth-creds",
						Namespace: "default",
					},
					Key: "id",
				},
				ClientSecretSecretRef: xpv1.SecretKeySelector{
					SecretReference: xpv1.SecretReference{
						Name:      "oauth-creds",
						Namespace: "default",
					},
					Key: "secret",
				},
				// Endpoints not specified - should use defaults
			},
		},
	}

	auth, err := createV2OAuthAuth(context.Background(), k8sClient, config)

	if err != nil {
		t.Errorf("createV2OAuthAuth() error = %v, want nil", err)
	}

	if auth.GetEndpoint() != "https://api.hostinger.com/v2" {
		t.Errorf("Endpoint = %v, want https://api.hostinger.com/v2", auth.GetEndpoint())
	}

	oauthAuth := auth.(*V2OAuthAuth)
	if oauthAuth.TokenEndpoint != "https://auth.hostinger.com/oauth/token" {
		t.Errorf("TokenEndpoint = %v, want https://auth.hostinger.com/oauth/token", oauthAuth.TokenEndpoint)
	}
}

func TestCreateV2OAuthAuth_MissingClientIDSecret(t *testing.T) {
	sch := fake.NewClientBuilder().Build().Scheme()
	k8sClient := fake.NewClientBuilder().
		WithScheme(sch).
		Build()

	config := &v1beta1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: "default",
		},
		Spec: v1beta1.ProviderConfigSpec{
			OAuthAuth: &v1beta1.OAuthAuthSpec{
				ClientIDSecretRef: xpv1.SecretKeySelector{
					SecretReference: xpv1.SecretReference{
						Name:      "nonexistent",
						Namespace: "default",
					},
					Key: "id",
				},
				ClientSecretSecretRef: xpv1.SecretKeySelector{
					SecretReference: xpv1.SecretReference{
						Name:      "creds",
						Namespace: "default",
					},
					Key: "secret",
				},
			},
		},
	}

	auth, err := createV2OAuthAuth(context.Background(), k8sClient, config)

	if err == nil {
		t.Error("createV2OAuthAuth() expected error for missing secret, got nil")
	}

	if auth != nil {
		t.Errorf("createV2OAuthAuth() expected nil authenticator, got %v", auth)
	}
}

func TestCreateV2OAuthAuth_MissingClientSecretSecret(t *testing.T) {
	sch := fake.NewClientBuilder().Build().Scheme()

	secrets := []client.Object{
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "oauth-creds",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"id": []byte("client-id"),
			},
		},
	}

	k8sClient := fake.NewClientBuilder().
		WithScheme(sch).
		WithObjects(secrets...).
		Build()

	config := &v1beta1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: "default",
		},
		Spec: v1beta1.ProviderConfigSpec{
			OAuthAuth: &v1beta1.OAuthAuthSpec{
				ClientIDSecretRef: xpv1.SecretKeySelector{
					SecretReference: xpv1.SecretReference{
						Name:      "oauth-creds",
						Namespace: "default",
					},
					Key: "id",
				},
				ClientSecretSecretRef: xpv1.SecretKeySelector{
					SecretReference: xpv1.SecretReference{
						Name:      "nonexistent",
						Namespace: "default",
					},
					Key: "secret",
				},
			},
		},
	}

	auth, err := createV2OAuthAuth(context.Background(), k8sClient, config)

	if err == nil {
		t.Error("createV2OAuthAuth() expected error for missing client secret, got nil")
	}

	if auth != nil {
		t.Errorf("createV2OAuthAuth() expected nil authenticator, got %v", auth)
	}
}

func TestGetSecretValue_Success(t *testing.T) {
	sch := fake.NewClientBuilder().Build().Scheme()

	secrets := []client.Object{
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-secret",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"username": []byte("admin"),
				"password": []byte("secret123"),
			},
		},
	}

	k8sClient := fake.NewClientBuilder().
		WithScheme(sch).
		WithObjects(secrets...).
		Build()

	secretRef := &xpv1.SecretKeySelector{
		SecretReference: xpv1.SecretReference{
			Name:      "test-secret",
			Namespace: "default",
		},
		Key: "password",
	}

	value, err := getSecretValue(context.Background(), k8sClient, "default", secretRef)

	if err != nil {
		t.Errorf("getSecretValue() error = %v, want nil", err)
	}

	if value != "secret123" {
		t.Errorf("getSecretValue() = %v, want secret123", value)
	}
}

func TestGetSecretValue_MissingSecret(t *testing.T) {
	sch := fake.NewClientBuilder().Build().Scheme()
	k8sClient := fake.NewClientBuilder().
		WithScheme(sch).
		Build()

	secretRef := &xpv1.SecretKeySelector{
		SecretReference: xpv1.SecretReference{
			Name:      "nonexistent",
			Namespace: "default",
		},
		Key: "key",
	}

	value, err := getSecretValue(context.Background(), k8sClient, "default", secretRef)

	if err == nil {
		t.Error("getSecretValue() expected error for missing secret, got nil")
	}

	if value != "" {
		t.Errorf("getSecretValue() = %v, want empty string on error", value)
	}
}

func TestGetSecretValue_MissingKey(t *testing.T) {
	sch := fake.NewClientBuilder().Build().Scheme()

	secrets := []client.Object{
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-secret",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"username": []byte("admin"),
			},
		},
	}

	k8sClient := fake.NewClientBuilder().
		WithScheme(sch).
		WithObjects(secrets...).
		Build()

	secretRef := &xpv1.SecretKeySelector{
		SecretReference: xpv1.SecretReference{
			Name:      "test-secret",
			Namespace: "default",
		},
		Key: "nonexistent-key",
	}

	value, err := getSecretValue(context.Background(), k8sClient, "default", secretRef)

	if err == nil {
		t.Error("getSecretValue() expected error for missing key, got nil")
	}

	if value != "" {
		t.Errorf("getSecretValue() = %v, want empty string on error", value)
	}
}

func TestGetSecretValue_NilSecretRef(t *testing.T) {
	k8sClient := fake.NewClientBuilder().Build()

	value, err := getSecretValue(context.Background(), k8sClient, "default", nil)

	if err == nil {
		t.Error("getSecretValue() expected error for nil secret ref, got nil")
	}

	if value != "" {
		t.Errorf("getSecretValue() = %v, want empty string on error", value)
	}
}

func TestGetSecretValue_DifferentNamespace(t *testing.T) {
	sch := fake.NewClientBuilder().Build().Scheme()

	secrets := []client.Object{
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-secret",
				Namespace: "production",
			},
			Data: map[string][]byte{
				"key": []byte("prod-value"),
			},
		},
	}

	k8sClient := fake.NewClientBuilder().
		WithScheme(sch).
		WithObjects(secrets...).
		Build()

	secretRef := &xpv1.SecretKeySelector{
		SecretReference: xpv1.SecretReference{
			Name:      "test-secret",
			Namespace: "production",
		},
		Key: "key",
	}

	value, err := getSecretValue(context.Background(), k8sClient, "production", secretRef)

	if err != nil {
		t.Errorf("getSecretValue() error = %v, want nil", err)
	}

	if value != "prod-value" {
		t.Errorf("getSecretValue() = %v, want prod-value", value)
	}
}