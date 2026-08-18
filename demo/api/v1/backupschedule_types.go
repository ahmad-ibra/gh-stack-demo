package v1

// BackupScheduleSpec defines the desired state of a scheduled cluster backup.
type BackupScheduleSpec struct {
	// Schedule is a cron expression, e.g. "0 2 * * *".
	Schedule string `json:"schedule"`

	// TODO(api): retention policy, target storage location.
}

// BackupScheduleStatus captures the observed state.
type BackupScheduleStatus struct {
	LastRun string `json:"lastRun,omitempty"`
}
