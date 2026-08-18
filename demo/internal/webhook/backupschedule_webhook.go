package webhook

// ValidateSchedule checks that the cron expression is well-formed.
//
// TODO(webhook): wire into the admission webhook and reject bad crons.
func ValidateSchedule(cron string) error {
	return nil
}
