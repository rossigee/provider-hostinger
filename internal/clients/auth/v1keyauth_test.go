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
	"strings"
	"testing"
)

func TestNewV1KeyAuth(t *testing.T) {
	apiKey := "test-api-key"
	customerID := "test-customer"
	endpoint := "https://api.hostinger.com/v1"

	auth := NewV1KeyAuth(apiKey, customerID, endpoint)

	if auth == nil {
		t.Fatal("NewV1KeyAuth returned nil")
	}
	if auth.APIKey != apiKey {
		t.Errorf("APIKey = %v, want %v", auth.APIKey, apiKey)
	}
	if auth.CustomerID != customerID {
		t.Errorf("CustomerID = %v, want %v", auth.CustomerID, customerID)
	}
	if auth.Endpoint != endpoint {
		t.Errorf("Endpoint = %v, want %v", auth.Endpoint, endpoint)
	}
}

func TestV1KeyAuthGetAuthHeader(t *testing.T) {
	tests := []struct {
		name       string
		apiKey     string
		customerID string
		wantErr    bool
		checkHeader func(string) bool
	}{
		{
			name:       "valid credentials",
			apiKey:     "key123",
			customerID: "customer456",
			wantErr:    false,
			checkHeader: func(header string) bool {
				// Check Basic auth format
				if !strings.HasPrefix(header, "Basic ") {
					return false
				}
				// Decode and verify
				encoded := strings.TrimPrefix(header, "Basic ")
				decoded, err := base64.StdEncoding.DecodeString(encoded)
				if err != nil {
					return false
				}
				expected := "customer456:key123"
				return string(decoded) == expected
			},
		},
		{
			name:       "empty credentials",
			apiKey:     "",
			customerID: "",
			wantErr:    false,
			checkHeader: func(header string) bool {
				// Even empty credentials should be encoded
				if !strings.HasPrefix(header, "Basic ") {
					return false
				}
				return true
			},
		},
		{
			name:       "special characters in credentials",
			apiKey:     "key@with:special#chars",
			customerID: "customer/with\\slashes",
			wantErr:    false,
			checkHeader: func(header string) bool {
				if !strings.HasPrefix(header, "Basic ") {
					return false
				}
				encoded := strings.TrimPrefix(header, "Basic ")
				decoded, _ := base64.StdEncoding.DecodeString(encoded)
				expected := "customer/with\\slashes:key@with:special#chars"
				return string(decoded) == expected
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := NewV1KeyAuth(tt.apiKey, tt.customerID, "https://api.hostinger.com/v1")
			header, err := auth.GetAuthHeader(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAuthHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.checkHeader(header) {
				t.Errorf("GetAuthHeader() = %v, header validation failed", header)
			}
		})
	}
}

func TestV1KeyAuthGetToken(t *testing.T) {
	auth := NewV1KeyAuth("key", "customer", "https://api.hostinger.com/v1")
	token, err := auth.GetToken(context.Background())

	if err != nil {
		t.Errorf("GetToken() error = %v, want nil", err)
	}
	if token != "" {
		t.Errorf("GetToken() = %v, want empty string (v1 doesn't use tokens)", token)
	}
}

func TestV1KeyAuthGetEndpoint(t *testing.T) {
	endpoint := "https://api.hostinger.com/v1"
	auth := NewV1KeyAuth("key", "customer", endpoint)

	if auth.GetEndpoint() != endpoint {
		t.Errorf("GetEndpoint() = %v, want %v", auth.GetEndpoint(), endpoint)
	}
}

func TestV1KeyAuthRefreshIfNeeded(t *testing.T) {
	auth := NewV1KeyAuth("key", "customer", "https://api.hostinger.com/v1")
	err := auth.RefreshIfNeeded(context.Background())

	if err != nil {
		t.Errorf("RefreshIfNeeded() error = %v, want nil", err)
	}
}

func TestV1KeyAuthType(t *testing.T) {
	auth := NewV1KeyAuth("key", "customer", "https://api.hostinger.com/v1")

	if auth.Type() != "APIKeyAuth" {
		t.Errorf("Type() = %v, want APIKeyAuth", auth.Type())
	}
}

func TestV1KeyAuthImplementsAuthenticator(t *testing.T) {
	// This is a compile-time check, but we verify it by creating the type assertion
	var _ Authenticator = (*V1KeyAuth)(nil)
}

func BenchmarkV1KeyAuthGetAuthHeader(b *testing.B) {
	auth := NewV1KeyAuth("key123", "customer456", "https://api.hostinger.com/v1")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = auth.GetAuthHeader(ctx)
	}
}
