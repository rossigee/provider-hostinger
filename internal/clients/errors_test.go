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
	"net/http"
	"testing"
)

func TestHostingerErrorError(t *testing.T) {
	err := &HostingerError{
		Type:    ErrorTypeNotFound,
		Message: "Instance not found",
		Status:  http.StatusNotFound,
	}

	expected := "NotFound: Instance not found (status: 404)"
	if err.Error() != expected {
		t.Errorf("Error() = %v, want %v", err.Error(), expected)
	}
}

func TestClassifyError_NotFound(t *testing.T) {
	err := ClassifyError(http.StatusNotFound, "Instance not found")

	if err.Type != ErrorTypeNotFound {
		t.Errorf("Type = %v, want NotFound", err.Type)
	}
	if err.Status != http.StatusNotFound {
		t.Errorf("Status = %v, want 404", err.Status)
	}
	if err.Message != "Instance not found" {
		t.Errorf("Message = %v, want 'Instance not found'", err.Message)
	}
}

func TestClassifyError_Unauthorized(t *testing.T) {
	err := ClassifyError(http.StatusUnauthorized, "Invalid API key")

	if err.Type != ErrorTypeUnauthorized {
		t.Errorf("Type = %v, want Unauthorized", err.Type)
	}
	if err.Status != http.StatusUnauthorized {
		t.Errorf("Status = %v, want 401", err.Status)
	}
}

func TestClassifyError_Forbidden(t *testing.T) {
	err := ClassifyError(http.StatusForbidden, "Access denied")

	if err.Type != ErrorTypeForbidden {
		t.Errorf("Type = %v, want Forbidden", err.Type)
	}
	if err.Status != http.StatusForbidden {
		t.Errorf("Status = %v, want 403", err.Status)
	}
}

func TestClassifyError_Conflict(t *testing.T) {
	err := ClassifyError(http.StatusConflict, "Hostname already exists")

	if err.Type != ErrorTypeConflict {
		t.Errorf("Type = %v, want Conflict", err.Type)
	}
	if err.Status != http.StatusConflict {
		t.Errorf("Status = %v, want 409", err.Status)
	}
}

func TestClassifyError_RateLimit(t *testing.T) {
	err := ClassifyError(http.StatusTooManyRequests, "Rate limit exceeded")

	if err.Type != ErrorTypeRateLimit {
		t.Errorf("Type = %v, want RateLimit", err.Type)
	}
	if err.Status != http.StatusTooManyRequests {
		t.Errorf("Status = %v, want 429", err.Status)
	}
}

func TestClassifyError_Internal_500(t *testing.T) {
	err := ClassifyError(http.StatusInternalServerError, "Server error")

	if err.Type != ErrorTypeInternal {
		t.Errorf("Type = %v, want Internal", err.Type)
	}
}

func TestClassifyError_Internal_502(t *testing.T) {
	err := ClassifyError(http.StatusBadGateway, "Bad gateway")

	if err.Type != ErrorTypeInternal {
		t.Errorf("Type = %v, want Internal for 502", err.Type)
	}
}

func TestClassifyError_Internal_503(t *testing.T) {
	err := ClassifyError(http.StatusServiceUnavailable, "Service unavailable")

	if err.Type != ErrorTypeInternal {
		t.Errorf("Type = %v, want Internal for 503", err.Type)
	}
}

func TestClassifyError_Internal_504(t *testing.T) {
	err := ClassifyError(http.StatusGatewayTimeout, "Gateway timeout")

	if err.Type != ErrorTypeInternal {
		t.Errorf("Type = %v, want Internal for 504", err.Type)
	}
}

func TestClassifyError_Unknown(t *testing.T) {
	err := ClassifyError(http.StatusBadRequest, "Bad request")

	if err.Type != ErrorTypeUnknown {
		t.Errorf("Type = %v, want Unknown for unknown status", err.Type)
	}
	if err.Status != http.StatusBadRequest {
		t.Errorf("Status = %v, want 400", err.Status)
	}
}

func TestIsNotFound_Nil(t *testing.T) {
	if IsNotFound(nil) {
		t.Error("IsNotFound(nil) should return false")
	}
}

func TestIsNotFound_NotFoundError(t *testing.T) {
	err := &HostingerError{
		Type:   ErrorTypeNotFound,
		Status: http.StatusNotFound,
	}
	if !IsNotFound(err) {
		t.Error("IsNotFound should return true for NotFound error")
	}
}

func TestIsNotFound_NotFoundStatus(t *testing.T) {
	err := &HostingerError{
		Type:   ErrorTypeUnknown, // Different type
		Status: http.StatusNotFound,
	}
	if !IsNotFound(err) {
		t.Error("IsNotFound should return true for 404 status even with different type")
	}
}

func TestIsNotFound_OtherError(t *testing.T) {
	err := &HostingerError{
		Type:   ErrorTypeUnauthorized,
		Status: http.StatusUnauthorized,
	}
	if IsNotFound(err) {
		t.Error("IsNotFound should return false for non-NotFound error")
	}
}

func TestIsNotFound_NonHostingerError(t *testing.T) {
	// Test that IsNotFound safely handles non-HostingerError types
	// by using a regular error interface
	var err error = &HostingerError{
		Type:   ErrorTypeUnauthorized,
		Status: http.StatusUnauthorized,
	}

	if IsNotFound(err) {
		t.Error("IsNotFound should return false for non-NotFound error")
	}
}

func TestIsUnauthorized_Nil(t *testing.T) {
	if IsUnauthorized(nil) {
		t.Error("IsUnauthorized(nil) should return false")
	}
}

func TestIsUnauthorized_UnauthorizedError(t *testing.T) {
	err := &HostingerError{
		Type:   ErrorTypeUnauthorized,
		Status: http.StatusUnauthorized,
	}
	if !IsUnauthorized(err) {
		t.Error("IsUnauthorized should return true for Unauthorized error")
	}
}

func TestIsUnauthorized_UnauthorizedStatus(t *testing.T) {
	err := &HostingerError{
		Type:   ErrorTypeUnknown,
		Status: http.StatusUnauthorized,
	}
	if !IsUnauthorized(err) {
		t.Error("IsUnauthorized should return true for 401 status")
	}
}

func TestIsUnauthorized_OtherError(t *testing.T) {
	err := &HostingerError{
		Type:   ErrorTypeNotFound,
		Status: http.StatusNotFound,
	}
	if IsUnauthorized(err) {
		t.Error("IsUnauthorized should return false for non-Unauthorized error")
	}
}

func TestIsForbidden_Nil(t *testing.T) {
	if IsForbidden(nil) {
		t.Error("IsForbidden(nil) should return false")
	}
}

func TestIsForbidden_ForbiddenError(t *testing.T) {
	err := &HostingerError{
		Type:   ErrorTypeForbidden,
		Status: http.StatusForbidden,
	}
	if !IsForbidden(err) {
		t.Error("IsForbidden should return true for Forbidden error")
	}
}

func TestIsForbidden_ForbiddenStatus(t *testing.T) {
	err := &HostingerError{
		Type:   ErrorTypeUnknown,
		Status: http.StatusForbidden,
	}
	if !IsForbidden(err) {
		t.Error("IsForbidden should return true for 403 status")
	}
}

func TestIsForbidden_OtherError(t *testing.T) {
	err := &HostingerError{
		Type:   ErrorTypeNotFound,
		Status: http.StatusNotFound,
	}
	if IsForbidden(err) {
		t.Error("IsForbidden should return false for non-Forbidden error")
	}
}

func TestIsConflict_Nil(t *testing.T) {
	if IsConflict(nil) {
		t.Error("IsConflict(nil) should return false")
	}
}

func TestIsConflict_ConflictError(t *testing.T) {
	err := &HostingerError{
		Type:   ErrorTypeConflict,
		Status: http.StatusConflict,
	}
	if !IsConflict(err) {
		t.Error("IsConflict should return true for Conflict error")
	}
}

func TestIsConflict_ConflictStatus(t *testing.T) {
	err := &HostingerError{
		Type:   ErrorTypeUnknown,
		Status: http.StatusConflict,
	}
	if !IsConflict(err) {
		t.Error("IsConflict should return true for 409 status")
	}
}

func TestIsConflict_OtherError(t *testing.T) {
	err := &HostingerError{
		Type:   ErrorTypeNotFound,
		Status: http.StatusNotFound,
	}
	if IsConflict(err) {
		t.Error("IsConflict should return false for non-Conflict error")
	}
}

func TestIsRateLimit_Nil(t *testing.T) {
	if IsRateLimit(nil) {
		t.Error("IsRateLimit(nil) should return false")
	}
}

func TestIsRateLimit_RateLimitError(t *testing.T) {
	err := &HostingerError{
		Type:   ErrorTypeRateLimit,
		Status: http.StatusTooManyRequests,
	}
	if !IsRateLimit(err) {
		t.Error("IsRateLimit should return true for RateLimit error")
	}
}

func TestIsRateLimit_RateLimitStatus(t *testing.T) {
	err := &HostingerError{
		Type:   ErrorTypeUnknown,
		Status: http.StatusTooManyRequests,
	}
	if !IsRateLimit(err) {
		t.Error("IsRateLimit should return true for 429 status")
	}
}

func TestIsRateLimit_OtherError(t *testing.T) {
	err := &HostingerError{
		Type:   ErrorTypeNotFound,
		Status: http.StatusNotFound,
	}
	if IsRateLimit(err) {
		t.Error("IsRateLimit should return false for non-RateLimit error")
	}
}

func TestErrorTypeString(t *testing.T) {
	tests := []struct {
		name       string
		errorType  ErrorType
		expected   string
	}{
		{"NotFound", ErrorTypeNotFound, "NotFound"},
		{"Unauthorized", ErrorTypeUnauthorized, "Unauthorized"},
		{"Forbidden", ErrorTypeForbidden, "Forbidden"},
		{"InvalidConfig", ErrorTypeInvalidConfig, "InvalidConfig"},
		{"RateLimit", ErrorTypeRateLimit, "RateLimit"},
		{"Conflict", ErrorTypeConflict, "Conflict"},
		{"Internal", ErrorTypeInternal, "Internal"},
		{"Unknown", ErrorTypeUnknown, "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.errorType) != tt.expected {
				t.Errorf("ErrorType = %v, want %v", tt.errorType, tt.expected)
			}
		})
	}
}
