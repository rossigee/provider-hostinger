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
)

// Authenticator handles authentication for Hostinger API requests
type Authenticator interface {
	// GetAuthHeader returns the Authorization header value for API requests
	GetAuthHeader(ctx context.Context) (string, error)

	// GetToken returns the bearer token (if applicable)
	GetToken(ctx context.Context) (string, error)

	// GetEndpoint returns the API endpoint
	GetEndpoint() string

	// RefreshIfNeeded checks if credentials need refreshing and updates them
	RefreshIfNeeded(ctx context.Context) error

	// Type returns the authentication type name
	Type() string
}
