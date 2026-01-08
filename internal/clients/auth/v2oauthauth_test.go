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
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewV2OAuthAuth(t *testing.T) {
	tests := []struct {
		name             string
		clientID         string
		clientSecret     string
		endpoint         string
		tokenEndpoint    string
		expectedEndpoint string
	}{
		{
			name:             "with custom token endpoint",
			clientID:         "client123",
			clientSecret:     "secret456",
			endpoint:         "https://api.hostinger.com/v2",
			tokenEndpoint:    "https://custom.auth.com/token",
			expectedEndpoint: "https://custom.auth.com/token",
		},
		{
			name:             "with empty token endpoint uses default",
			clientID:         "client123",
			clientSecret:     "secret456",
			endpoint:         "https://api.hostinger.com/v2",
			tokenEndpoint:    "",
			expectedEndpoint: "https://auth.hostinger.com/oauth/token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := NewV2OAuthAuth(tt.clientID, tt.clientSecret, tt.endpoint, tt.tokenEndpoint)

			if auth == nil {
				t.Fatal("NewV2OAuthAuth returned nil")
			}
			if auth.ClientID != tt.clientID {
				t.Errorf("ClientID = %v, want %v", auth.ClientID, tt.clientID)
			}
			if auth.ClientSecret != tt.clientSecret {
				t.Errorf("ClientSecret = %v, want %v", auth.ClientSecret, tt.clientSecret)
			}
			if auth.Endpoint != tt.endpoint {
				t.Errorf("Endpoint = %v, want %v", auth.Endpoint, tt.endpoint)
			}
			if auth.TokenEndpoint != tt.expectedEndpoint {
				t.Errorf("TokenEndpoint = %v, want %v", auth.TokenEndpoint, tt.expectedEndpoint)
			}
		})
	}
}

func TestV2OAuthAuthGetAuthHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %v", r.Method)
		}

		resp := OAuthTokenResponse{
			AccessToken: "test-token-12345",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Logf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	auth := NewV2OAuthAuth("client", "secret", "https://api.hostinger.com/v2", server.URL)
	header, err := auth.GetAuthHeader(context.Background())

	if err != nil {
		t.Errorf("GetAuthHeader() error = %v, want nil", err)
	}

	expectedHeader := "Bearer test-token-12345"
	if header != expectedHeader {
		t.Errorf("GetAuthHeader() = %v, want %v", header, expectedHeader)
	}
}

func TestV2OAuthAuthGetToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := OAuthTokenResponse{
			AccessToken: "test-token-xyz",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Logf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	auth := NewV2OAuthAuth("client", "secret", "https://api.hostinger.com/v2", server.URL)
	token, err := auth.GetToken(context.Background())

	if err != nil {
		t.Errorf("GetToken() error = %v, want nil", err)
	}

	if token != "test-token-xyz" {
		t.Errorf("GetToken() = %v, want test-token-xyz", token)
	}
}

func TestV2OAuthAuthTokenCaching(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		resp := OAuthTokenResponse{
			AccessToken: fmt.Sprintf("token-%d", callCount),
			TokenType:   "Bearer",
			ExpiresIn:   3600,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Logf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	auth := NewV2OAuthAuth("client", "secret", "https://api.hostinger.com/v2", server.URL)

	// First call should hit the server
	token1, err := auth.GetToken(context.Background())
	if err != nil {
		t.Fatalf("First GetToken() error = %v, want nil", err)
	}
	if token1 != "token-1" {
		t.Errorf("First GetToken() = %v, want token-1", token1)
	}
	if callCount != 1 {
		t.Errorf("Expected 1 server call, got %d", callCount)
	}

	// Second call should use cached token
	token2, err := auth.GetToken(context.Background())
	if err != nil {
		t.Fatalf("Second GetToken() error = %v, want nil", err)
	}
	if token2 != "token-1" {
		t.Errorf("Second GetToken() = %v, want token-1 (cached)", token2)
	}
	if callCount != 1 {
		t.Errorf("Expected 1 server call (cached), got %d", callCount)
	}
}

func TestV2OAuthAuthTokenRefresh(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		resp := OAuthTokenResponse{
			AccessToken: fmt.Sprintf("token-%d", callCount),
			TokenType:   "Bearer",
			ExpiresIn:   1, // 1 second expiry
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Logf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	auth := NewV2OAuthAuth("client", "secret", "https://api.hostinger.com/v2", server.URL)

	// First token
	token1, _ := auth.GetToken(context.Background())
	if token1 != "token-1" {
		t.Errorf("First token = %v, want token-1", token1)
	}

	// Wait for token to expire
	time.Sleep(1100 * time.Millisecond)

	// Second token should be refreshed
	token2, _ := auth.GetToken(context.Background())
	if token2 != "token-2" {
		t.Errorf("Second token = %v, want token-2 (refreshed)", token2)
	}

	if callCount != 2 {
		t.Errorf("Expected 2 server calls (refresh), got %d", callCount)
	}
}

func TestV2OAuthAuthGetEndpoint(t *testing.T) {
	endpoint := "https://api.hostinger.com/v2"
	auth := NewV2OAuthAuth("client", "secret", endpoint, "")

	if auth.GetEndpoint() != endpoint {
		t.Errorf("GetEndpoint() = %v, want %v", auth.GetEndpoint(), endpoint)
	}
}

func TestV2OAuthAuthType(t *testing.T) {
	auth := NewV2OAuthAuth("client", "secret", "https://api.hostinger.com/v2", "")

	if auth.Type() != "OAuthAuth" {
		t.Errorf("Type() = %v, want OAuthAuth", auth.Type())
	}
}

func TestV2OAuthAuthRefreshIfNeeded(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := OAuthTokenResponse{
			AccessToken: "test-token",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Logf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	auth := NewV2OAuthAuth("client", "secret", "https://api.hostinger.com/v2", server.URL)
	err := auth.RefreshIfNeeded(context.Background())

	if err != nil {
		t.Errorf("RefreshIfNeeded() error = %v, want nil", err)
	}
}

func TestV2OAuthAuthServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		if _, err := w.Write([]byte("unauthorized")); err != nil {
			t.Logf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	auth := NewV2OAuthAuth("client", "secret", "https://api.hostinger.com/v2", server.URL)
	_, err := auth.GetToken(context.Background())

	if err == nil {
		t.Error("GetToken() expected error for 401 response, got nil")
	}

	if !strings.Contains(err.Error(), "401") {
		t.Errorf("GetToken() error = %v, want error containing 401", err)
	}
}

func TestV2OAuthAuthInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte("invalid json")); err != nil {
			t.Logf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	auth := NewV2OAuthAuth("client", "secret", "https://api.hostinger.com/v2", server.URL)
	_, err := auth.GetToken(context.Background())

	if err == nil {
		t.Error("GetToken() expected error for invalid JSON, got nil")
	}
}

func TestV2OAuthAuthGetAuthHeader_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("server error")); err != nil {
			t.Logf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	auth := NewV2OAuthAuth("client", "secret", "https://api.hostinger.com/v2", server.URL)
	header, err := auth.GetAuthHeader(context.Background())

	if err == nil {
		t.Error("GetAuthHeader() expected error for 500 response, got nil")
	}

	if header != "" {
		t.Errorf("GetAuthHeader() = %v, want empty string on error", header)
	}
}

func TestV2OAuthAuthRequestCreation(t *testing.T) {
	requestReceived := false
	var receivedAuth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestReceived = true

		// Check method
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %v", r.Method)
		}

		// Check content type
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Errorf("Expected application/x-www-form-urlencoded, got %v", r.Header.Get("Content-Type"))
		}

		// Check auth header
		receivedAuth = r.Header.Get("Accept")

		resp := OAuthTokenResponse{
			AccessToken: "test-token",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Logf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	auth := NewV2OAuthAuth("client", "secret", "https://api.hostinger.com/v2", server.URL)
	_, _ = auth.GetToken(context.Background())

	if !requestReceived {
		t.Error("OAuth token request was not made")
	}

	if receivedAuth != "application/json" {
		t.Errorf("Accept header = %v, want application/json", receivedAuth)
	}
}

func TestV2OAuthAuthContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow server
		<-time.After(100 * time.Millisecond)
		resp := OAuthTokenResponse{
			AccessToken: "test-token",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Logf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	auth := NewV2OAuthAuth("client", "secret", "https://api.hostinger.com/v2", server.URL)

	// Create a context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := auth.GetToken(ctx)

	if err == nil {
		t.Error("GetToken() expected error on context deadline, got nil")
	}
}

func TestV2OAuthAuthImplementsAuthenticator(t *testing.T) {
	// This is a compile-time check
	var _ Authenticator = (*V2OAuthAuth)(nil)
}
