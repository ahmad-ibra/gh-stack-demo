package controller

import (
	"context"
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	backupv1 "github.com/ahmad-ibra/gh-stack-demo/demo/api/v1"
)

const backupImage = "ghcr.io/ahmad-ibra/cluster-backup:latest"

// BackupScheduleReconciler reconciles a BackupSchedule into a CronJob.
type BackupScheduleReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=backup.example.com,resources=backupschedules,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=backup.example.com,resources=backupschedules/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=batch,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete

// Reconcile ensures a CronJob exists that matches the desired BackupSchedule.
func (r *BackupScheduleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var schedule backupv1.BackupSchedule
	if err := r.Get(ctx, req.NamespacedName, &schedule); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	desired := r.cronJobFor(&schedule)
	if err := ctrl.SetControllerReference(&schedule, desired, r.Scheme); err != nil {
		return ctrl.Result{}, fmt.Errorf("setting owner reference: %w", err)
	}

	var existing batchv1.CronJob
	err := r.Get(ctx, client.ObjectKeyFromObject(desired), &existing)
	switch {
	case apierrors.IsNotFound(err):
		logger.Info("creating backup CronJob", "cronjob", desired.Name)
		if err := r.Create(ctx, desired); err != nil {
			return ctrl.Result{}, fmt.Errorf("creating CronJob: %w", err)
		}
	case err != nil:
		return ctrl.Result{}, fmt.Errorf("fetching CronJob: %w", err)
	default:
		existing.Spec = desired.Spec
		if err := r.Update(ctx, &existing); err != nil {
			return ctrl.Result{}, fmt.Errorf("updating CronJob: %w", err)
		}
	}

	return ctrl.Result{}, nil
}

// cronJobFor renders the CronJob that backs a BackupSchedule.
func (r *BackupScheduleReconciler) cronJobFor(schedule *backupv1.BackupSchedule) *batchv1.CronJob {
	history := int32(schedule.Spec.Retention)
	return &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      schedule.Name + "-backup",
			Namespace: schedule.Namespace,
		},
		Spec: batchv1.CronJobSpec{
			Schedule:                   schedule.Spec.Schedule,
			Suspend:                    &schedule.Spec.Suspend,
			SuccessfulJobsHistoryLimit: &history,
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							RestartPolicy: corev1.RestartPolicyOnFailure,
							Containers: []corev1.Container{{
								Name:  "backup",
								Image: backupImage,
								Args: []string{
									"--bucket", schedule.Spec.Target.Bucket,
									"--prefix", schedule.Spec.Target.Prefix,
								},
							}},
						},
					},
				},
			},
		},
	}
}

// SetupWithManager wires the reconciler into the manager.
func (r *BackupScheduleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&backupv1.BackupSchedule{}).
		Owns(&batchv1.CronJob{}).
		Complete(r)
}
