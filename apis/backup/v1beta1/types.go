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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
)

// BackupScheduleType represents the backup schedule frequency
// +kubebuilder:validation:Enum=manual;daily;weekly;monthly
type BackupScheduleType string

const (
	// BackupScheduleManual means backups are created manually
	BackupScheduleManual BackupScheduleType = "manual"
	// BackupScheduleDaily means backups are created daily
	BackupScheduleDaily BackupScheduleType = "daily"
	// BackupScheduleWeekly means backups are created weekly
	BackupScheduleWeekly BackupScheduleType = "weekly"
	// BackupScheduleMonthly means backups are created monthly
	BackupScheduleMonthly BackupScheduleType = "monthly"
)

// BackupParameters are the configurable fields of a Hostinger VPS Backup.
type BackupParameters struct {
	// InstanceID is the ID of the VPS instance to backup.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	InstanceID string `json:"instanceId"`

	// Description is an optional description of the backup.
	// +kubebuilder:validation:Optional
	Description *string `json:"description,omitempty"`

	// Schedule is the backup schedule frequency.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=manual;daily;weekly;monthly
	Schedule *BackupScheduleType `json:"schedule,omitempty"`
}

// BackupObservation are the observable fields of a Hostinger VPS Backup.
type BackupObservation struct {
	// ID is the external backup resource ID.
	ID string `json:"id,omitempty"`

	// Status is the current status of the backup (pending, completed, failed, etc.).
	Status string `json:"status,omitempty"`

	// CreatedDate is when the backup was created.
	CreatedDate *metav1.Time `json:"createdDate,omitempty"`

	// Size is the backup size in MB.
	Size *int64 `json:"size,omitempty"`

	// ExpiryDate is when the backup will expire.
	ExpiryDate *metav1.Time `json:"expiryDate,omitempty"`

	// CurrentSchedule is the current backup schedule.
	CurrentSchedule *BackupScheduleType `json:"currentSchedule,omitempty"`
}

// BackupSpec defines the desired state of a Hostinger VPS Backup.
type BackupSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       BackupParameters `json:"forProvider"`
}

// BackupStatus defines the observed state of a Hostinger VPS Backup.
type BackupStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          BackupObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,hostinger}
// +kubebuilder:printcolumn:name="READY",type=string,JSONPath=.status.conditions[?(@.type=='Ready')].status
// +kubebuilder:printcolumn:name="SYNCED",type=string,JSONPath=.status.conditions[?(@.type=='Synced')].status
// +kubebuilder:printcolumn:name="AGE",type=date,JSONPath=.metadata.creationTimestamp
// +genclient

// Backup is the CRD type for Hostinger VPS backups.
type Backup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackupSpec   `json:"spec,omitempty"`
	Status BackupStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BackupList contains a list of Backup resources.
type BackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Backup `json:"items"`
}
