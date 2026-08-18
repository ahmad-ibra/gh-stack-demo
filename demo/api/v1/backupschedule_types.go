package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BackupScheduleSpec defines the desired state of a scheduled cluster backup.
type BackupScheduleSpec struct {
	// Schedule is a standard five-field cron expression (for example "0 2 * * *")
	// that determines when backups run.
	// +kubebuilder:validation:Required
	Schedule string `json:"schedule"`

	// Retention is the number of successful backups to keep. Older backups are
	// pruned once this limit is exceeded.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=7
	Retention int32 `json:"retention,omitempty"`

	// Target describes where backup artifacts are written.
	Target BackupTarget `json:"target"`

	// Suspend pauses the schedule without deleting it.
	// +optional
	Suspend bool `json:"suspend,omitempty"`
}

// BackupTarget describes an object-storage destination for backups.
type BackupTarget struct {
	// Bucket is the object-storage bucket name.
	// +kubebuilder:validation:Required
	Bucket string `json:"bucket"`

	// Prefix is an optional key prefix within the bucket.
	// +optional
	Prefix string `json:"prefix,omitempty"`
}

// BackupScheduleStatus captures the observed state of a BackupSchedule.
type BackupScheduleStatus struct {
	// LastScheduleTime is the last time a backup job was successfully created.
	// +optional
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty"`

	// Active lists backup jobs that are currently running.
	// +optional
	Active []corev1.ObjectReference `json:"active,omitempty"`

	// Conditions represents the latest observations of the schedule's state.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Schedule",type=string,JSONPath=`.spec.schedule`
// +kubebuilder:printcolumn:name="Suspended",type=boolean,JSONPath=`.spec.suspend`

// BackupSchedule is the Schema for the backupschedules API.
type BackupSchedule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackupScheduleSpec   `json:"spec,omitempty"`
	Status BackupScheduleStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BackupScheduleList contains a list of BackupSchedule.
type BackupScheduleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BackupSchedule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BackupSchedule{}, &BackupScheduleList{})
}
