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
	"encoding/base64"
	"fmt"

	"k8s.io/client-go/kubernetes"
)

// V1KeyAuth implements Authenticator for Hostinger API v1 (API key + customer ID)
type V1KeyAuth struct {
	APIKey     string
	CustomerID string
	Endpoint   string
}

// NewV1KeyAuth creates a new V1KeyAuth authenticator
func NewV1KeyAuth(apiKey, customerID, endpoint string) *V1KeyAuth {
	return &V1KeyAuth{
		APIKey:     apiKey,
		CustomerID: customerID,
		Endpoint:   endpoint,
	}
}

// GetAuthHeader returns the Authorization header value for v1 API key authentication
// API v1 uses Basic authentication with the format: base64(customer_id:api_key)
func (a *V1KeyAuth) GetAuthHeader(ctx context.Context) (string, error) {
	credentials := fmt.Sprintf("%s:%s", a.CustomerID, a.APIKey)
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
	return fmt.Sprintf("Basic %s", encoded), nil
}

// GetToken returns the bearer token (v1 doesn't use tokens, returns empty)
func (a *V1KeyAuth) GetToken(ctx context.Context) (string, error) {
	return "", nil // v1 API key uses Basic auth, not bearer tokens
}

// GetEndpoint returns the API endpoint
func (a *V1KeyAuth) GetEndpoint() string {
	return a.Endpoint
}

// RefreshIfNeeded performs any necessary refresh logic (v1 key auth doesn't need refresh)
func (a *V1KeyAuth) RefreshIfNeeded(ctx context.Context, k8sClient kubernetes.Interface) error {
	return nil // v1 API keys don't need refresh
}

// Type returns the authentication type
func (a *V1KeyAuth) Type() string {
	return "APIKeyAuth"
}
