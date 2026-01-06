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
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1beta1 "github.com/rossigee/provider-hostinger/apis/instance/v1beta1"
	"github.com/rossigee/provider-hostinger/internal/clients"
)

// Instance represents a Hostinger VPS instance
type Instance struct {
	ID              string
	Hostname        string
	Status          string
	IPAddress       string
	IPv6Address     string
	OSId            string
	CPUCount        int32
	RAM             int32
	DiskSize        int32
	Bandwidth       *int32
	CreationDate    *string
	ExpirationDate  *string
	RootPassword    *string
	IPv6Enabled     bool
	Inodes          *int32
}

// Client defines operations for managing Hostinger VPS instances
type Client interface {
	// Create creates a new VPS instance
	Create(ctx context.Context, params *v1beta1.InstanceParameters) (*Instance, error)

	// Get retrieves a VPS instance by ID
	Get(ctx context.Context, instanceID string) (*Instance, error)

	// Update modifies an existing VPS instance
	Update(ctx context.Context, instanceID string, params *v1beta1.InstanceParameters) error

	// Delete terminates a VPS instance
	Delete(ctx context.Context, instanceID string) error

	// List returns all VPS instances
	List(ctx context.Context) ([]*Instance, error)

	// GetObservation maps an Instance to the observation status
	GetObservation(instance *Instance) *v1beta1.InstanceObservation

	// LateInitialize updates unset fields from the remote instance
	LateInitialize(instance *Instance, params *v1beta1.InstanceParameters) bool

	// UpToDate checks if local spec matches remote instance
	UpToDate(instance *Instance, params *v1beta1.InstanceParameters) bool
}

// InstanceClient implements the Client interface
type InstanceClient struct {
	hostingerClient *clients.HostingerClient
}

// NewInstanceClient creates a new Instance client
func NewInstanceClient(hostingerClient *clients.HostingerClient) *InstanceClient {
	return &InstanceClient{
		hostingerClient: hostingerClient,
	}
}

// Create creates a new VPS instance
func (ic *InstanceClient) Create(ctx context.Context, params *v1beta1.InstanceParameters) (*Instance, error) {
	// Implementation stub - will call Hostinger API
	return nil, fmt.Errorf("not implemented yet")
}

// Get retrieves a VPS instance by ID
func (ic *InstanceClient) Get(ctx context.Context, instanceID string) (*Instance, error) {
	// Implementation stub - will call Hostinger API
	return nil, fmt.Errorf("not implemented yet")
}

// Update modifies an existing VPS instance
func (ic *InstanceClient) Update(ctx context.Context, instanceID string, params *v1beta1.InstanceParameters) error {
	// Implementation stub - will call Hostinger API
	return fmt.Errorf("not implemented yet")
}

// Delete terminates a VPS instance
func (ic *InstanceClient) Delete(ctx context.Context, instanceID string) error {
	// Implementation stub - will call Hostinger API
	return fmt.Errorf("not implemented yet")
}

// List returns all VPS instances
func (ic *InstanceClient) List(ctx context.Context) ([]*Instance, error) {
	// Implementation stub - will call Hostinger API
	return nil, fmt.Errorf("not implemented yet")
}

// GetObservation maps an Instance to the observation status
func (ic *InstanceClient) GetObservation(instance *Instance) *v1beta1.InstanceObservation {
	if instance == nil {
		return &v1beta1.InstanceObservation{}
	}

	obs := &v1beta1.InstanceObservation{
		ID:                 instance.ID,
		Status:             instance.Status,
		IPAddress:          instance.IPAddress,
		IPv6Address:        instance.IPv6Address,
		CurrentHostname:    instance.Hostname,
		CurrentCPUCount:    instance.CPUCount,
		CurrentRAM:         instance.RAM,
		CurrentDiskSize:    instance.DiskSize,
		CreationDate:       parseTime(instance.CreationDate),
		ExpirationDate:     parseTime(instance.ExpirationDate),
	}

	if instance.Bandwidth != nil {
		obs.CurrentBandwidth = instance.Bandwidth
	}

	if instance.Inodes != nil {
		obs.CurrentInodes = instance.Inodes
	}

	return obs
}

// parseTime parses an ISO 8601 time string to metav1.Time
func parseTime(timeStr *string) *metav1.Time {
	if timeStr == nil || *timeStr == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, *timeStr)
	if err != nil {
		return nil
	}
	mt := metav1.NewTime(t)
	return &mt
}

// LateInitialize updates unset fields from the remote instance
func (ic *InstanceClient) LateInitialize(instance *Instance, params *v1beta1.InstanceParameters) bool {
	if instance == nil {
		return false
	}

	changed := false

	// OSId - initialize if not set
	if params.OSId == "" && instance.OSId != "" {
		params.OSId = instance.OSId
		changed = true
	}

	// Bandwidth - initialize if not set
	if params.Bandwidth == nil && instance.Bandwidth != nil {
		params.Bandwidth = instance.Bandwidth
		changed = true
	}

	// IPv6Enabled - initialize if not set
	if params.IPv6Enabled == nil && instance.IPv6Enabled {
		enabled := true
		params.IPv6Enabled = &enabled
		changed = true
	}

	// Inodes - initialize if not set
	if params.Inodes == nil && instance.Inodes != nil {
		params.Inodes = instance.Inodes
		changed = true
	}

	return changed
}

// UpToDate checks if local spec matches remote instance
func (ic *InstanceClient) UpToDate(instance *Instance, params *v1beta1.InstanceParameters) bool {
	if instance == nil {
		return false
	}

	// Check hostname
	if params.Hostname != "" && params.Hostname != instance.Hostname {
		return false
	}

	// Check CPU count
	if params.CPUCount > 0 && params.CPUCount != instance.CPUCount {
		return false
	}

	// Check RAM
	if params.RAM > 0 && params.RAM != instance.RAM {
		return false
	}

	// Check disk size
	if params.DiskSize > 0 && params.DiskSize != instance.DiskSize {
		return false
	}

	// Check bandwidth
	if params.Bandwidth != nil && *params.Bandwidth != *instance.Bandwidth {
		return false
	}

	// Check IPv6
	if params.IPv6Enabled != nil && *params.IPv6Enabled != instance.IPv6Enabled {
		return false
	}

	return true
}
