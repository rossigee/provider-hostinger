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
	"fmt"
	"net/http"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client"

	v1beta1 "github.com/rossigee/provider-hostinger/apis/v1beta1"
	"github.com/rossigee/provider-hostinger/internal/clients/auth"
)

// HTTPClientConfig contains configuration for the HTTP client
type HTTPClientConfig struct {
	Timeout        time.Duration
	MaxRetries     int
	RetryWaitTime  time.Duration
	UserAgent      string
}

// DefaultHTTPClientConfig returns the default HTTP client configuration
func DefaultHTTPClientConfig() HTTPClientConfig {
	return HTTPClientConfig{
		Timeout:       30 * time.Second,
		MaxRetries:    3,
		RetryWaitTime: 1 * time.Second,
		UserAgent:     "provider-hostinger/v0.1.0",
	}
}

// HostingerClient represents the Hostinger API client
type HostingerClient struct {
	authenticator auth.Authenticator
	httpClient    *http.Client
	config        HTTPClientConfig
	k8sClient     client.Client
	providerCfg   *v1beta1.ProviderConfig
}

// ClientFactory creates Hostinger API clients
type ClientFactory struct {
	k8sClient client.Client
	httpCfg   HTTPClientConfig
}

// NewClientFactory creates a new Hostinger client factory
func NewClientFactory(k8sClient client.Client, cfg HTTPClientConfig) *ClientFactory {
	return &ClientFactory{
		k8sClient: k8sClient,
		httpCfg:   cfg,
	}
}

// CreateHostingerClient creates a new Hostinger API client from ProviderConfig
func (cf *ClientFactory) CreateHostingerClient(ctx context.Context, config *v1beta1.ProviderConfig) (*HostingerClient, error) {
	// Create authenticator based on ProviderConfig
	authenticator, err := auth.CreateAuthenticator(ctx, cf.k8sClient, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create authenticator: %w", err)
	}

	// Create HTTP client
	httpClient := &http.Client{
		Timeout: cf.httpCfg.Timeout,
	}

	return &HostingerClient{
		authenticator: authenticator,
		httpClient:    httpClient,
		config:        cf.httpCfg,
		k8sClient:     cf.k8sClient,
		providerCfg:   config,
	}, nil
}

// GetAuthenticator returns the configured authenticator
func (hc *HostingerClient) GetAuthenticator() auth.Authenticator {
	return hc.authenticator
}

// GetHTTPClient returns the configured HTTP client
func (hc *HostingerClient) GetHTTPClient() *http.Client {
	return hc.httpClient
}

// GetEndpoint returns the API endpoint
func (hc *HostingerClient) GetEndpoint() string {
	return hc.authenticator.GetEndpoint()
}

// GetAuthType returns the authentication type
func (hc *HostingerClient) GetAuthType() string {
	return hc.authenticator.Type()
}

// PrepareRequest prepares an HTTP request with authentication headers
func (hc *HostingerClient) PrepareRequest(ctx context.Context, req *http.Request) error {
	// Refresh authentication if needed
	if err := hc.authenticator.RefreshIfNeeded(ctx); err != nil {
		return fmt.Errorf("failed to refresh authentication: %w", err)
	}

	// Get authorization header
	authHeader, err := hc.authenticator.GetAuthHeader(ctx)
	if err != nil {
		return fmt.Errorf("failed to get authorization header: %w", err)
	}

	// Set authorization and user agent headers
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("User-Agent", hc.config.UserAgent)
	req.Header.Set("Accept", "application/json")

	return nil
}

// Do performs an HTTP request with error handling and retry logic
func (hc *HostingerClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	// Prepare request with authentication
	if err := hc.PrepareRequest(ctx, req); err != nil {
		return nil, err
	}

	// Perform request with retry logic
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= hc.config.MaxRetries; attempt++ {
		resp, err = hc.httpClient.Do(req)
		if err != nil {
			if attempt < hc.config.MaxRetries {
				time.Sleep(hc.config.RetryWaitTime * time.Duration(attempt+1))
				continue
			}
			return nil, fmt.Errorf("request failed after %d retries: %w", hc.config.MaxRetries, err)
		}

		// Check if response indicates a retryable error
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= http.StatusInternalServerError {
			resp.Body.Close()
			if attempt < hc.config.MaxRetries {
				time.Sleep(hc.config.RetryWaitTime * time.Duration(attempt+1))
				continue
			}
		}

		// Success or non-retryable error
		break
	}

	return resp, nil
}

// GetProviderConfig returns the ProviderConfig used to create this client
func (hc *HostingerClient) GetProviderConfig() *v1beta1.ProviderConfig {
	return hc.providerCfg
}
