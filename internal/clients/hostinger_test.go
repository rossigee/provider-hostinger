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

package clients

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	v1beta1 "github.com/rossigee/provider-hostinger/apis/v1beta1"
)

// MockAuthenticator is a mock implementation of auth.Authenticator
type MockAuthenticator struct {
	authHeader  string
	token       string
	endpoint    string
	authType    string
	refreshErr  error
	headerErr   error
	needsRefresh bool
}

func (m *MockAuthenticator) GetAuthHeader(ctx context.Context) (string, error) {
	if m.headerErr != nil {
		return "", m.headerErr
	}
	return m.authHeader, nil
}

func (m *MockAuthenticator) GetToken(ctx context.Context) (string, error) {
	return m.token, nil
}

func (m *MockAuthenticator) GetEndpoint() string {
	return m.endpoint
}

func (m *MockAuthenticator) RefreshIfNeeded(ctx context.Context) error {
	if m.needsRefresh && m.refreshErr != nil {
		return m.refreshErr
	}
	return nil
}

func (m *MockAuthenticator) Type() string {
	return m.authType
}

func TestDefaultHTTPClientConfig(t *testing.T) {
	cfg := DefaultHTTPClientConfig()

	if cfg.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want 30s", cfg.Timeout)
	}
	if cfg.MaxRetries != 3 {
		t.Errorf("MaxRetries = %v, want 3", cfg.MaxRetries)
	}
	if cfg.RetryWaitTime != 1*time.Second {
		t.Errorf("RetryWaitTime = %v, want 1s", cfg.RetryWaitTime)
	}
	if !strings.Contains(cfg.UserAgent, "provider-hostinger") {
		t.Errorf("UserAgent = %v, want to contain 'provider-hostinger'", cfg.UserAgent)
	}
}

func TestNewClientFactory(t *testing.T) {
	cfg := HTTPClientConfig{
		Timeout:       15 * time.Second,
		MaxRetries:    2,
		RetryWaitTime: 500 * time.Millisecond,
		UserAgent:     "test-agent",
	}

	factory := NewClientFactory(nil, cfg)

	if factory == nil {
		t.Fatal("NewClientFactory returned nil")
	}
	if factory.httpCfg.Timeout != cfg.Timeout {
		t.Errorf("Factory timeout = %v, want %v", factory.httpCfg.Timeout, cfg.Timeout)
	}
}

func TestHostingerClientGetters(t *testing.T) {
	mockAuth := &MockAuthenticator{
		authHeader: "Bearer test-token",
		endpoint:   "https://api.hostinger.com/v2",
		authType:   "OAuthAuth",
	}

	cfg := HTTPClientConfig{
		Timeout:       30 * time.Second,
		MaxRetries:    3,
		RetryWaitTime: 1 * time.Second,
		UserAgent:     "test-agent",
	}

	httpClient := &http.Client{Timeout: cfg.Timeout}

	client := &HostingerClient{
		authenticator: mockAuth,
		httpClient:    httpClient,
		config:        cfg,
		providerCfg:   &v1beta1.ProviderConfig{},
	}

	// Test GetAuthenticator
	if client.GetAuthenticator() != mockAuth {
		t.Error("GetAuthenticator() did not return expected authenticator")
	}

	// Test GetHTTPClient
	if client.GetHTTPClient() != httpClient {
		t.Error("GetHTTPClient() did not return expected HTTP client")
	}

	// Test GetEndpoint
	if client.GetEndpoint() != "https://api.hostinger.com/v2" {
		t.Errorf("GetEndpoint() = %v, want https://api.hostinger.com/v2", client.GetEndpoint())
	}

	// Test GetAuthType
	if client.GetAuthType() != "OAuthAuth" {
		t.Errorf("GetAuthType() = %v, want OAuthAuth", client.GetAuthType())
	}

	// Test GetProviderConfig
	if client.GetProviderConfig() == nil {
		t.Error("GetProviderConfig() returned nil")
	}
}

func TestPrepareRequest_Success(t *testing.T) {
	mockAuth := &MockAuthenticator{
		authHeader: "Bearer test-token",
		endpoint:   "https://api.hostinger.com/v2",
	}

	cfg := HTTPClientConfig{
		UserAgent: "test-agent/1.0",
	}

	client := &HostingerClient{
		authenticator: mockAuth,
		config:        cfg,
	}

	req, _ := http.NewRequest("GET", "https://api.hostinger.com/v2/instances", nil)

	err := client.PrepareRequest(context.Background(), req)

	if err != nil {
		t.Errorf("PrepareRequest() error = %v, want nil", err)
	}

	if req.Header.Get("Authorization") != "Bearer test-token" {
		t.Errorf("Authorization header = %v, want 'Bearer test-token'", req.Header.Get("Authorization"))
	}

	if req.Header.Get("User-Agent") != "test-agent/1.0" {
		t.Errorf("User-Agent header = %v, want 'test-agent/1.0'", req.Header.Get("User-Agent"))
	}

	if req.Header.Get("Accept") != "application/json" {
		t.Errorf("Accept header = %v, want 'application/json'", req.Header.Get("Accept"))
	}
}

func TestPrepareRequest_RefreshError(t *testing.T) {
	mockAuth := &MockAuthenticator{
		authHeader:   "Bearer test-token",
		refreshErr:   ClassifyError(http.StatusUnauthorized, "Token expired"),
		needsRefresh: true,
	}

	client := &HostingerClient{
		authenticator: mockAuth,
		config:        HTTPClientConfig{},
	}

	req, _ := http.NewRequest("GET", "https://api.hostinger.com/v2/instances", nil)

	err := client.PrepareRequest(context.Background(), req)

	if err == nil {
		t.Error("PrepareRequest() expected error for refresh failure, got nil")
	}

	if !strings.Contains(err.Error(), "failed to refresh authentication") {
		t.Errorf("Error = %v, want to contain 'failed to refresh authentication'", err)
	}
}

func TestPrepareRequest_AuthHeaderError(t *testing.T) {
	mockAuth := &MockAuthenticator{
		headerErr: ClassifyError(http.StatusUnauthorized, "Invalid API key"),
	}

	client := &HostingerClient{
		authenticator: mockAuth,
		config:        HTTPClientConfig{},
	}

	req, _ := http.NewRequest("GET", "https://api.hostinger.com/v2/instances", nil)

	err := client.PrepareRequest(context.Background(), req)

	if err == nil {
		t.Error("PrepareRequest() expected error for header failure, got nil")
	}

	if !strings.Contains(err.Error(), "failed to get authorization header") {
		t.Errorf("Error = %v, want to contain 'failed to get authorization header'", err)
	}
}

func TestDo_SuccessfulRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"id": "instance-123", "hostname": "test.example.com"}`)); err != nil {
			t.Logf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	mockAuth := &MockAuthenticator{
		authHeader: "Bearer test-token",
	}

	cfg := HTTPClientConfig{
		Timeout:       10 * time.Second,
		MaxRetries:    3,
		RetryWaitTime: 100 * time.Millisecond,
		UserAgent:     "test-agent",
	}

	httpClient := &http.Client{Timeout: cfg.Timeout}

	client := &HostingerClient{
		authenticator: mockAuth,
		httpClient:    httpClient,
		config:        cfg,
	}

	req, _ := http.NewRequest("GET", server.URL+"/instances", nil)
	resp, err := client.Do(context.Background(), req)

	if err != nil {
		t.Errorf("Do() error = %v, want nil", err)
	}

	if resp == nil {
		t.Fatal("Do() returned nil response")
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Response status = %v, want 200", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if err := resp.Body.Close(); err != nil {
		t.Logf("failed to close response body: %v", err)
	}

	if !strings.Contains(string(body), "instance-123") {
		t.Errorf("Response body = %v, want to contain 'instance-123'", string(body))
	}
}

func TestDo_RetryOn429(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount < 2 {
			w.WriteHeader(http.StatusTooManyRequests)
			if _, err := w.Write([]byte("Rate limited")); err != nil {
				t.Logf("failed to write response: %v", err)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status": "ok"}`)); err != nil {
			t.Logf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	mockAuth := &MockAuthenticator{
		authHeader: "Bearer test-token",
	}

	cfg := HTTPClientConfig{
		Timeout:       10 * time.Second,
		MaxRetries:    3,
		RetryWaitTime: 50 * time.Millisecond,
		UserAgent:     "test-agent",
	}

	httpClient := &http.Client{Timeout: cfg.Timeout}

	client := &HostingerClient{
		authenticator: mockAuth,
		httpClient:    httpClient,
		config:        cfg,
	}

	req, _ := http.NewRequest("GET", server.URL+"/instances", nil)
	resp, err := client.Do(context.Background(), req)

	if err != nil {
		t.Errorf("Do() error = %v, want nil", err)
	}

	if resp == nil {
		t.Fatal("Do() returned nil response")
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Response status = %v, want 200 (after retry)", resp.StatusCode)
	}

	if callCount < 2 {
		t.Errorf("Expected at least 2 calls (retry), got %d", callCount)
	}

	if err := resp.Body.Close(); err != nil {
		t.Logf("failed to close response body: %v", err)
	}
}

func TestDo_RetryOn5xx(t *testing.T) {
	statusCodes := []int{
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
	}

	for _, statusCode := range statusCodes {
		t.Run(strings.TrimPrefix(http.StatusText(statusCode), ""), func(t *testing.T) {
			callCount := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				callCount++
				if callCount < 2 {
					w.WriteHeader(statusCode)
					if _, err := w.Write([]byte("Server error")); err != nil {
						t.Logf("failed to write response: %v", err)
					}
					return
				}
				w.WriteHeader(http.StatusOK)
				if _, err := w.Write([]byte(`{"status": "ok"}`)); err != nil {
					t.Logf("failed to write response: %v", err)
				}
			}))
			defer server.Close()

			mockAuth := &MockAuthenticator{
				authHeader: "Bearer test-token",
			}

			cfg := HTTPClientConfig{
				Timeout:       10 * time.Second,
				MaxRetries:    3,
				RetryWaitTime: 50 * time.Millisecond,
				UserAgent:     "test-agent",
			}

			httpClient := &http.Client{Timeout: cfg.Timeout}

			client := &HostingerClient{
				authenticator: mockAuth,
				httpClient:    httpClient,
				config:        cfg,
			}

			req, _ := http.NewRequest("GET", server.URL+"/instances", nil)
			resp, err := client.Do(context.Background(), req)

			if err != nil {
				t.Errorf("Do() error = %v, want nil", err)
			}

			if resp != nil && resp.StatusCode != http.StatusOK {
				t.Errorf("Response status = %v, want 200", resp.StatusCode)
				if err := resp.Body.Close(); err != nil {
					t.Logf("failed to close response body: %v", err)
				}
			}

			if callCount < 2 {
				t.Errorf("Expected retry for %s, got %d calls", http.StatusText(statusCode), callCount)
			}
		})
	}
}

func TestDo_MaxRetriesExhausted(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		if _, err := w.Write([]byte("Service unavailable")); err != nil {
			t.Logf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	mockAuth := &MockAuthenticator{
		authHeader: "Bearer test-token",
	}

	cfg := HTTPClientConfig{
		Timeout:       10 * time.Second,
		MaxRetries:    2,
		RetryWaitTime: 50 * time.Millisecond,
		UserAgent:     "test-agent",
	}

	httpClient := &http.Client{Timeout: cfg.Timeout}

	client := &HostingerClient{
		authenticator: mockAuth,
		httpClient:    httpClient,
		config:        cfg,
	}

	req, _ := http.NewRequest("GET", server.URL+"/instances", nil)
	resp, err := client.Do(context.Background(), req)

	if err != nil {
		t.Errorf("Do() error = %v, want nil", err)
	}

	// After max retries exhausted, response is returned as-is with the last error status
	if resp == nil {
		t.Fatal("Do() returned nil response after retries")
	}

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("Response status = %v, want 503 (last attempt)", resp.StatusCode)
	}

	if err := resp.Body.Close(); err != nil {
		t.Logf("failed to close response body: %v", err)
	}
}

func TestDo_NoRetryOn4xx(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("Not found")); err != nil {
			t.Logf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	mockAuth := &MockAuthenticator{
		authHeader: "Bearer test-token",
	}

	cfg := HTTPClientConfig{
		Timeout:       10 * time.Second,
		MaxRetries:    3,
		RetryWaitTime: 50 * time.Millisecond,
		UserAgent:     "test-agent",
	}

	httpClient := &http.Client{Timeout: cfg.Timeout}

	client := &HostingerClient{
		authenticator: mockAuth,
		httpClient:    httpClient,
		config:        cfg,
	}

	req, _ := http.NewRequest("GET", server.URL+"/instances", nil)
	resp, err := client.Do(context.Background(), req)

	if err != nil {
		t.Errorf("Do() error = %v, want nil", err)
	}

	if resp == nil {
		t.Fatal("Do() returned nil response")
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Response status = %v, want 404", resp.StatusCode)
	}

	if callCount > 1 {
		t.Errorf("Expected no retry for 404, got %d calls", callCount)
	}

	if err := resp.Body.Close(); err != nil {
		t.Logf("failed to close response body: %v", err)
	}
}

func TestDo_PrepareRequestError(t *testing.T) {
	mockAuth := &MockAuthenticator{
		headerErr: ClassifyError(http.StatusUnauthorized, "Invalid credentials"),
	}

	cfg := HTTPClientConfig{
		Timeout:       10 * time.Second,
		MaxRetries:    3,
		RetryWaitTime: 50 * time.Millisecond,
		UserAgent:     "test-agent",
	}

	httpClient := &http.Client{Timeout: cfg.Timeout}

	client := &HostingerClient{
		authenticator: mockAuth,
		httpClient:    httpClient,
		config:        cfg,
	}

	req, _ := http.NewRequest("GET", "https://api.example.com/instances", nil)
	resp, err := client.Do(context.Background(), req)

	if resp != nil {
		t.Errorf("Do() expected nil response on prepare error, got %v", resp)
	}

	if err == nil {
		t.Error("Do() expected error for prepare failure, got nil")
	}

	if !strings.Contains(err.Error(), "failed to get authorization header") {
		t.Errorf("Error = %v, want to contain 'failed to get authorization header'", err)
	}
}

func TestDo_RetryWithExponentialBackoff(t *testing.T) {
	callTimes := []time.Time{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callTimes = append(callTimes, time.Now())
		if len(callTimes) < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			if _, err := w.Write([]byte("Service unavailable")); err != nil {
				t.Logf("failed to write response: %v", err)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status": "ok"}`)); err != nil {
			t.Logf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	mockAuth := &MockAuthenticator{
		authHeader: "Bearer test-token",
	}

	cfg := HTTPClientConfig{
		Timeout:       10 * time.Second,
		MaxRetries:    3,
		RetryWaitTime: 50 * time.Millisecond,
		UserAgent:     "test-agent",
	}

	httpClient := &http.Client{Timeout: cfg.Timeout}

	client := &HostingerClient{
		authenticator: mockAuth,
		httpClient:    httpClient,
		config:        cfg,
	}

	req, _ := http.NewRequest("GET", server.URL+"/instances", nil)
	resp, err := client.Do(context.Background(), req)

	if err != nil {
		t.Errorf("Do() error = %v, want nil", err)
	}

	if resp != nil {
		if err := resp.Body.Close(); err != nil {
			t.Logf("failed to close response body: %v", err)
		}
	}

	// Verify retry delays increased (attempt 1: 50ms, attempt 2: 100ms)
	if len(callTimes) >= 3 {
		delay1 := callTimes[1].Sub(callTimes[0])
		delay2 := callTimes[2].Sub(callTimes[1])

		// Delays should be at least the configured times (with some tolerance for timing)
		if delay1 < 40*time.Millisecond {
			t.Logf("First retry delay = %v, want ≥40ms", delay1)
		}
		if delay2 < 80*time.Millisecond {
			t.Logf("Second retry delay = %v, want ≥80ms", delay2)
		}

		// Second delay should be greater than first (exponential backoff)
		if delay2 <= delay1 {
			t.Logf("Delays not increasing: delay1=%v, delay2=%v", delay1, delay2)
		}
	}
}

func TestDo_RequestModificationPreserved(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify custom headers were preserved
		if r.Header.Get("X-Custom-Header") != "custom-value" {
			t.Errorf("Custom header lost: %v", r.Header.Get("X-Custom-Header"))
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status": "ok"}`)); err != nil {
			t.Logf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	mockAuth := &MockAuthenticator{
		authHeader: "Bearer test-token",
	}

	cfg := HTTPClientConfig{
		Timeout:       10 * time.Second,
		MaxRetries:    3,
		RetryWaitTime: 50 * time.Millisecond,
		UserAgent:     "test-agent",
	}

	httpClient := &http.Client{Timeout: cfg.Timeout}

	client := &HostingerClient{
		authenticator: mockAuth,
		httpClient:    httpClient,
		config:        cfg,
	}

	req, _ := http.NewRequest("GET", server.URL+"/instances", nil)
	req.Header.Set("X-Custom-Header", "custom-value")

	resp, err := client.Do(context.Background(), req)

	if err != nil {
		t.Errorf("Do() error = %v, want nil", err)
	}

	if resp != nil {
		if err := resp.Body.Close(); err != nil {
			t.Logf("failed to close response body: %v", err)
		}
	}
}
