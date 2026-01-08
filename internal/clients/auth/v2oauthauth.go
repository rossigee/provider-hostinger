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
	"net/url"
	"strings"
	"sync"
	"time"
)

// V2OAuthAuth implements Authenticator for Hostinger API v2 (OAuth)
type V2OAuthAuth struct {
	ClientID      string
	ClientSecret  string
	Endpoint      string
	TokenEndpoint string

	// Token caching
	mu              sync.RWMutex
	cachedToken     string
	cachedExpiresAt time.Time
}

// OAuthTokenResponse represents the response from the OAuth token endpoint
type OAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// NewV2OAuthAuth creates a new V2OAuthAuth authenticator
func NewV2OAuthAuth(clientID, clientSecret, endpoint, tokenEndpoint string) *V2OAuthAuth {
	if tokenEndpoint == "" {
		// Default OAuth token endpoint for Hostinger v2 API
		tokenEndpoint = "https://auth.hostinger.com/oauth/token"
	}
	return &V2OAuthAuth{
		ClientID:      clientID,
		ClientSecret:  clientSecret,
		Endpoint:      endpoint,
		TokenEndpoint: tokenEndpoint,
	}
}

// GetAuthHeader returns the Authorization header value for v2 OAuth
func (a *V2OAuthAuth) GetAuthHeader(ctx context.Context) (string, error) {
	token, err := a.getToken(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Bearer %s", token), nil
}

// GetToken returns the bearer token for v2 OAuth
func (a *V2OAuthAuth) GetToken(ctx context.Context) (string, error) {
	return a.getToken(ctx)
}

// getToken gets or refreshes the OAuth token
func (a *V2OAuthAuth) getToken(ctx context.Context) (string, error) {
	a.mu.RLock()
	if a.cachedToken != "" && time.Now().Before(a.cachedExpiresAt) {
		defer a.mu.RUnlock()
		return a.cachedToken, nil
	}
	a.mu.RUnlock()

	// Token is expired or missing, refresh it
	token, expiresAt, err := a.refreshToken(ctx)
	if err != nil {
		return "", err
	}

	a.mu.Lock()
	a.cachedToken = token
	a.cachedExpiresAt = expiresAt
	a.mu.Unlock()

	return token, nil
}

// refreshToken performs the OAuth token refresh request
func (a *V2OAuthAuth) refreshToken(ctx context.Context) (string, time.Time, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", a.ClientID)
	data.Set("client_secret", a.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", a.TokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to request token: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", time.Time{}, fmt.Errorf("token request failed with status %d", resp.StatusCode)
	}

	var tokenResp OAuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", time.Time{}, fmt.Errorf("failed to decode token response: %w", err)
	}

	// Set expiry with 5-minute buffer to prevent using expired tokens
	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn-300) * time.Second)

	return tokenResp.AccessToken, expiresAt, nil
}

// GetEndpoint returns the API endpoint
func (a *V2OAuthAuth) GetEndpoint() string {
	return a.Endpoint
}

// RefreshIfNeeded checks if the token needs refreshing and updates it
func (a *V2OAuthAuth) RefreshIfNeeded(ctx context.Context) error {
	a.mu.RLock()
	needsRefresh := a.cachedToken == "" || time.Now().After(a.cachedExpiresAt.Add(-5*time.Minute))
	a.mu.RUnlock()

	if !needsRefresh {
		return nil
	}

	_, _, err := a.refreshToken(ctx)
	return err
}

// Type returns the authentication type
func (a *V2OAuthAuth) Type() string {
	return "OAuthAuth"
}
