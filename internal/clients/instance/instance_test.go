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

package instance

import (
	"context"
	"testing"
	"time"

	v1beta1 "github.com/rossigee/provider-hostinger/apis/instance/v1beta1"
	"github.com/rossigee/provider-hostinger/internal/clients"
)

func TestNewInstanceClient(t *testing.T) {
	mockHostingerClient := &clients.HostingerClient{}
	client := NewInstanceClient(mockHostingerClient)

	if client == nil {
		t.Fatal("NewInstanceClient returned nil")
	}
	if client.hostingerClient != mockHostingerClient {
		t.Error("HostingerClient not set correctly")
	}
}

func TestGetObservation_NilInstance(t *testing.T) {
	client := NewInstanceClient(nil)
	obs := client.GetObservation(nil)

	if obs == nil {
		t.Fatal("GetObservation returned nil for nil instance")
	}

	// Should return empty observation
	if obs.ID != "" || obs.Status != "" {
		t.Error("Expected empty observation for nil instance")
	}
}

func TestGetObservation_ValidInstance(t *testing.T) {
	creationDate := "2024-01-08T10:00:00Z"
	expirationDate := "2025-01-08T10:00:00Z"
	bandwidth := int32(1000)

	instance := &Instance{
		ID:              "inst-123",
		Hostname:        "vps.example.com",
		Status:          "active",
		IPAddress:       "192.168.1.100",
		IPv6Address:     "2001:db8::1",
		OSId:            "ubuntu20",
		CPUCount:        4,
		RAM:             8,
		DiskSize:        100,
		Bandwidth:       &bandwidth,
		CreationDate:    &creationDate,
		ExpirationDate:  &expirationDate,
		IPv6Enabled:     true,
	}

	client := NewInstanceClient(nil)
	obs := client.GetObservation(instance)

	if obs.ID != "inst-123" {
		t.Errorf("ID = %v, want inst-123", obs.ID)
	}
	if obs.Status != "active" {
		t.Errorf("Status = %v, want active", obs.Status)
	}
	if obs.IPAddress != "192.168.1.100" {
		t.Errorf("IPAddress = %v, want 192.168.1.100", obs.IPAddress)
	}
	if obs.IPv6Address != "2001:db8::1" {
		t.Errorf("IPv6Address = %v, want 2001:db8::1", obs.IPv6Address)
	}
	if obs.CurrentHostname != "vps.example.com" {
		t.Errorf("CurrentHostname = %v, want vps.example.com", obs.CurrentHostname)
	}
	if obs.CurrentCPUCount != 4 {
		t.Errorf("CurrentCPUCount = %v, want 4", obs.CurrentCPUCount)
	}
	if obs.CurrentRAM != 8 {
		t.Errorf("CurrentRAM = %v, want 8", obs.CurrentRAM)
	}
	if obs.CurrentDiskSize != 100 {
		t.Errorf("CurrentDiskSize = %v, want 100", obs.CurrentDiskSize)
	}

	// Check parsed times are not nil
	if obs.CreationDate == nil {
		t.Error("CreationDate should not be nil")
	}
	if obs.ExpirationDate == nil {
		t.Error("ExpirationDate should not be nil")
	}
}

func TestGetObservation_NilDates(t *testing.T) {
	instance := &Instance{
		ID:     "inst-123",
		Status: "active",
	}

	client := NewInstanceClient(nil)
	obs := client.GetObservation(instance)

	if obs.CreationDate != nil {
		t.Error("CreationDate should be nil for nil date string")
	}
	if obs.ExpirationDate != nil {
		t.Error("ExpirationDate should be nil for nil date string")
	}
}

func TestParseTime_Valid(t *testing.T) {
	dateStr := "2024-01-08T10:30:45Z"
	result := parseTime(&dateStr)

	if result == nil {
		t.Fatal("parseTime returned nil for valid date")
	}

	expected := time.Date(2024, 1, 8, 10, 30, 45, 0, time.UTC)
	if !result.Time.Equal(expected) {
		t.Errorf("parseTime = %v, want %v", result.Time, expected)
	}
}

func TestParseTime_Nil(t *testing.T) {
	result := parseTime(nil)
	if result != nil {
		t.Errorf("parseTime(nil) = %v, want nil", result)
	}
}

func TestParseTime_Empty(t *testing.T) {
	empty := ""
	result := parseTime(&empty)
	if result != nil {
		t.Errorf("parseTime(empty) = %v, want nil", result)
	}
}

func TestParseTime_Invalid(t *testing.T) {
	invalid := "not-a-date"
	result := parseTime(&invalid)
	if result != nil {
		t.Errorf("parseTime(invalid) = %v, want nil", result)
	}
}

func TestLateInitialize_NilInstance(t *testing.T) {
	params := &v1beta1.InstanceParameters{}
	client := NewInstanceClient(nil)

	changed := client.LateInitialize(nil, params)

	if changed {
		t.Error("LateInitialize(nil instance) should return false")
	}
}

func TestLateInitialize_NoChanges(t *testing.T) {
	instance := &Instance{
		OSId:        "ubuntu20",
		IPv6Enabled: false,
	}
	params := &v1beta1.InstanceParameters{
		OSId: "ubuntu20",
	}
	client := NewInstanceClient(nil)

	changed := client.LateInitialize(instance, params)

	if changed {
		t.Error("LateInitialize should return false when nothing to initialize")
	}
}

func TestLateInitialize_OSId(t *testing.T) {
	instance := &Instance{
		OSId: "ubuntu20",
	}
	params := &v1beta1.InstanceParameters{
		OSId: "", // Not set
	}
	client := NewInstanceClient(nil)

	changed := client.LateInitialize(instance, params)

	if !changed {
		t.Error("LateInitialize should return true when OSId is initialized")
	}
	if params.OSId != "ubuntu20" {
		t.Errorf("OSId = %v, want ubuntu20", params.OSId)
	}
}

func TestLateInitialize_Bandwidth(t *testing.T) {
	bandwidth := int32(1000)
	instance := &Instance{
		Bandwidth: &bandwidth,
	}
	params := &v1beta1.InstanceParameters{
		Bandwidth: nil, // Not set
	}
	client := NewInstanceClient(nil)

	changed := client.LateInitialize(instance, params)

	if !changed {
		t.Error("LateInitialize should return true when Bandwidth is initialized")
	}
	if params.Bandwidth == nil || *params.Bandwidth != 1000 {
		t.Errorf("Bandwidth = %v, want 1000", params.Bandwidth)
	}
}

func TestLateInitialize_IPv6Enabled(t *testing.T) {
	instance := &Instance{
		IPv6Enabled: true,
	}
	params := &v1beta1.InstanceParameters{
		IPv6Enabled: nil, // Not set
	}
	client := NewInstanceClient(nil)

	changed := client.LateInitialize(instance, params)

	if !changed {
		t.Error("LateInitialize should return true when IPv6Enabled is initialized")
	}
	if params.IPv6Enabled == nil || !*params.IPv6Enabled {
		t.Errorf("IPv6Enabled = %v, want true", params.IPv6Enabled)
	}
}

func TestLateInitialize_Inodes(t *testing.T) {
	inodes := int32(512000)
	instance := &Instance{
		Inodes: &inodes,
	}
	params := &v1beta1.InstanceParameters{
		Inodes: nil, // Not set
	}
	client := NewInstanceClient(nil)

	changed := client.LateInitialize(instance, params)

	if !changed {
		t.Error("LateInitialize should return true when Inodes is initialized")
	}
	if params.Inodes == nil || *params.Inodes != 512000 {
		t.Errorf("Inodes = %v, want 512000", params.Inodes)
	}
}

func TestLateInitialize_Multiple(t *testing.T) {
	bandwidth := int32(1000)
	instance := &Instance{
		OSId:        "ubuntu20",
		Bandwidth:   &bandwidth,
		IPv6Enabled: true,
	}
	params := &v1beta1.InstanceParameters{}
	client := NewInstanceClient(nil)

	changed := client.LateInitialize(instance, params)

	if !changed {
		t.Error("LateInitialize should return true when multiple fields are initialized")
	}
	if params.OSId != "ubuntu20" {
		t.Errorf("OSId not initialized correctly")
	}
	if params.Bandwidth == nil {
		t.Errorf("Bandwidth not initialized correctly")
	}
	if params.IPv6Enabled == nil || !*params.IPv6Enabled {
		t.Errorf("IPv6Enabled not initialized correctly")
	}
}

func TestUpToDate_NilInstance(t *testing.T) {
	params := &v1beta1.InstanceParameters{}
	client := NewInstanceClient(nil)

	upToDate := client.UpToDate(nil, params)

	if upToDate {
		t.Error("UpToDate(nil instance) should return false")
	}
}

func TestUpToDate_AllMatching(t *testing.T) {
	bandwidth := int32(1000)
	ipv6 := true

	instance := &Instance{
		Hostname:    "vps.example.com",
		CPUCount:    4,
		RAM:         8,
		DiskSize:    100,
		Bandwidth:   &bandwidth,
		IPv6Enabled: ipv6,
	}
	params := &v1beta1.InstanceParameters{
		Hostname:    "vps.example.com",
		CPUCount:    4,
		RAM:         8,
		DiskSize:    100,
		Bandwidth:   &bandwidth,
		IPv6Enabled: &ipv6,
	}
	client := NewInstanceClient(nil)

	upToDate := client.UpToDate(instance, params)

	if !upToDate {
		t.Error("UpToDate should return true when all fields match")
	}
}

func TestUpToDate_HostnameMismatch(t *testing.T) {
	instance := &Instance{
		Hostname: "vps1.example.com",
	}
	params := &v1beta1.InstanceParameters{
		Hostname: "vps2.example.com",
	}
	client := NewInstanceClient(nil)

	upToDate := client.UpToDate(instance, params)

	if upToDate {
		t.Error("UpToDate should return false when hostname doesn't match")
	}
}

func TestUpToDate_CPUMismatch(t *testing.T) {
	instance := &Instance{
		CPUCount: 4,
	}
	params := &v1beta1.InstanceParameters{
		CPUCount: 8,
	}
	client := NewInstanceClient(nil)

	upToDate := client.UpToDate(instance, params)

	if upToDate {
		t.Error("UpToDate should return false when CPU count doesn't match")
	}
}

func TestUpToDate_RAMMismatch(t *testing.T) {
	instance := &Instance{
		RAM: 8,
	}
	params := &v1beta1.InstanceParameters{
		RAM: 16,
	}
	client := NewInstanceClient(nil)

	upToDate := client.UpToDate(instance, params)

	if upToDate {
		t.Error("UpToDate should return false when RAM doesn't match")
	}
}

func TestUpToDate_DiskMismatch(t *testing.T) {
	instance := &Instance{
		DiskSize: 100,
	}
	params := &v1beta1.InstanceParameters{
		DiskSize: 200,
	}
	client := NewInstanceClient(nil)

	upToDate := client.UpToDate(instance, params)

	if upToDate {
		t.Error("UpToDate should return false when disk size doesn't match")
	}
}

func TestUpToDate_BandwidthMismatch(t *testing.T) {
	bandwidth1 := int32(1000)
	bandwidth2 := int32(2000)

	instance := &Instance{
		Bandwidth: &bandwidth1,
	}
	params := &v1beta1.InstanceParameters{
		Bandwidth: &bandwidth2,
	}
	client := NewInstanceClient(nil)

	upToDate := client.UpToDate(instance, params)

	if upToDate {
		t.Error("UpToDate should return false when bandwidth doesn't match")
	}
}

func TestUpToDate_IPv6Mismatch(t *testing.T) {
	ipv6False := false

	instance := &Instance{
		IPv6Enabled: true,
	}
	params := &v1beta1.InstanceParameters{
		IPv6Enabled: &ipv6False,
	}
	client := NewInstanceClient(nil)

	upToDate := client.UpToDate(instance, params)

	if upToDate {
		t.Error("UpToDate should return false when IPv6Enabled doesn't match")
	}
}

func TestUpToDate_OptionalFieldsNotSet(t *testing.T) {
	instance := &Instance{
		Hostname: "vps.example.com",
		CPUCount: 4,
	}
	params := &v1beta1.InstanceParameters{
		Hostname:    "vps.example.com",
		CPUCount:    4,
		Bandwidth:   nil,      // Not set in params
		IPv6Enabled: nil,      // Not set in params
	}
	client := NewInstanceClient(nil)

	upToDate := client.UpToDate(instance, params)

	if !upToDate {
		t.Error("UpToDate should return true when optional fields are not set in params")
	}
}

func TestInstanceClientImplementsInterface(t *testing.T) {
	// This is a compile-time check
	var _ Client = (*InstanceClient)(nil)
}

func TestCreateNotImplemented(t *testing.T) {
	client := NewInstanceClient(nil)
	_, err := client.Create(context.Background(), &v1beta1.InstanceParameters{})

	if err == nil {
		t.Error("Create should return error (not implemented)")
	}
}

func TestGetNotImplemented(t *testing.T) {
	client := NewInstanceClient(nil)
	_, err := client.Get(context.Background(), "inst-123")

	if err == nil {
		t.Error("Get should return error (not implemented)")
	}
}

func TestUpdateNotImplemented(t *testing.T) {
	client := NewInstanceClient(nil)
	err := client.Update(context.Background(), "inst-123", &v1beta1.InstanceParameters{})

	if err == nil {
		t.Error("Update should return error (not implemented)")
	}
}

func TestDeleteNotImplemented(t *testing.T) {
	client := NewInstanceClient(nil)
	err := client.Delete(context.Background(), "inst-123")

	if err == nil {
		t.Error("Delete should return error (not implemented)")
	}
}

func TestListNotImplemented(t *testing.T) {
	client := NewInstanceClient(nil)
	_, err := client.List(context.Background())

	if err == nil {
		t.Error("List should return error (not implemented)")
	}
}
