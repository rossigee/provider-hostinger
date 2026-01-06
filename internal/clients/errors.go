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
	"fmt"
	"net/http"
)

// Error types for Hostinger API errors
type ErrorType string

const (
	ErrorTypeNotFound      ErrorType = "NotFound"
	ErrorTypeUnauthorized  ErrorType = "Unauthorized"
	ErrorTypeForbidden     ErrorType = "Forbidden"
	ErrorTypeInvalidConfig ErrorType = "InvalidConfig"
	ErrorTypeRateLimit     ErrorType = "RateLimit"
	ErrorTypeConflict      ErrorType = "Conflict"
	ErrorTypeInternal      ErrorType = "Internal"
	ErrorTypeUnknown       ErrorType = "Unknown"
)

// HostingerError wraps Hostinger API errors with context
type HostingerError struct {
	Type    ErrorType
	Message string
	Status  int
	Err     error
}

func (e *HostingerError) Error() string {
	return fmt.Sprintf("%s: %s (status: %d)", e.Type, e.Message, e.Status)
}

// IsNotFound checks if an error is a 404 Not Found error
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	if he, ok := err.(*HostingerError); ok {
		return he.Type == ErrorTypeNotFound || he.Status == http.StatusNotFound
	}
	return false
}

// IsUnauthorized checks if an error is a 401 Unauthorized error
func IsUnauthorized(err error) bool {
	if err == nil {
		return false
	}
	if he, ok := err.(*HostingerError); ok {
		return he.Type == ErrorTypeUnauthorized || he.Status == http.StatusUnauthorized
	}
	return false
}

// IsForbidden checks if an error is a 403 Forbidden error
func IsForbidden(err error) bool {
	if err == nil {
		return false
	}
	if he, ok := err.(*HostingerError); ok {
		return he.Type == ErrorTypeForbidden || he.Status == http.StatusForbidden
	}
	return false
}

// IsConflict checks if an error is a 409 Conflict error
func IsConflict(err error) bool {
	if err == nil {
		return false
	}
	if he, ok := err.(*HostingerError); ok {
		return he.Type == ErrorTypeConflict || he.Status == http.StatusConflict
	}
	return false
}

// IsRateLimit checks if an error is a rate limit error
func IsRateLimit(err error) bool {
	if err == nil {
		return false
	}
	if he, ok := err.(*HostingerError); ok {
		return he.Type == ErrorTypeRateLimit || he.Status == http.StatusTooManyRequests
	}
	return false
}

// ClassifyError converts HTTP status codes to HostingerError types
func ClassifyError(status int, message string) *HostingerError {
	var errType ErrorType

	switch status {
	case http.StatusNotFound:
		errType = ErrorTypeNotFound
	case http.StatusUnauthorized:
		errType = ErrorTypeUnauthorized
	case http.StatusForbidden:
		errType = ErrorTypeForbidden
	case http.StatusConflict:
		errType = ErrorTypeConflict
	case http.StatusTooManyRequests:
		errType = ErrorTypeRateLimit
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		errType = ErrorTypeInternal
	default:
		errType = ErrorTypeUnknown
	}

	return &HostingerError{
		Type:    errType,
		Message: message,
		Status:  status,
	}
}
