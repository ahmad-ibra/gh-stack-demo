package webhook

import (
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	backupv1 "github.com/ahmad-ibra/gh-stack-demo/demo/api/v1"
	"github.com/ahmad-ibra/gh-stack-demo/demo/internal/features"
)

// cronParser accepts standard five-field cron expressions ("m h dom mon dow").
var cronParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

// +kubebuilder:webhook:path=/validate-backup-example-com-v1-backupschedule,mutating=false,failurePolicy=fail,sideEffects=None,groups=backup.example.com,resources=backupschedules,verbs=create;update,versions=v1,name=vbackupschedule.kb.io,admissionReviewVersions=v1

// BackupScheduleValidator validates BackupSchedule resources on admission.
type BackupScheduleValidator struct{}

var _ webhook.CustomValidator = &BackupScheduleValidator{}

// SetupWebhookWithManager registers the validating webhook with the manager,
// but only when the cluster-backup feature is enabled in this build.
// adding a comment
// add another
// add a conflict
func SetupWebhookWithManager(mgr ctrl.Manager) error {
	if !features.Enabled(features.FeatureClusterBackup) {
		return nil
	}
	return ctrl.NewWebhookManagedBy(mgr).
		For(&backupv1.BackupSchedule{}).
		WithValidator(&BackupScheduleValidator{}).
		Complete()
}

func (v *BackupScheduleValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	return v.validate(obj)
}

func (v *BackupScheduleValidator) ValidateUpdate(_ context.Context, _, newObj runtime.Object) (admission.Warnings, error) {
	return v.validate(newObj)
}

func (v *BackupScheduleValidator) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func (v *BackupScheduleValidator) validate(obj runtime.Object) (admission.Warnings, error) {
	schedule, ok := obj.(*backupv1.BackupSchedule)
	if !ok {
		return nil, fmt.Errorf("expected a BackupSchedule but got %T", obj)
	}

	if _, err := cronParser.Parse(schedule.Spec.Schedule); err != nil {
		return nil, fmt.Errorf("spec.schedule %q is not a valid cron expression: %w", schedule.Spec.Schedule, err)
	}
	if schedule.Spec.Target.Bucket == "" {
		return nil, fmt.Errorf("spec.target.bucket must be set")
	}
	if schedule.Spec.Retention < 1 {
		return nil, fmt.Errorf("spec.retention must be at least 1, got %d", schedule.Spec.Retention)
	}

	return nil, nil
}
