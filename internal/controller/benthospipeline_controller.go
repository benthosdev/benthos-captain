package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	streamv1 "github.com/benthosdev/benthos-captain/api/v1alpha1"
	"github.com/benthosdev/benthos-captain/internal/pkg/resource"
)

// BenthosPipelineReconciler reconciles a BenthosPipeline object
type BenthosPipelineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

type PipelineScope struct {
	Log      *logr.Logger
	Ctx      context.Context
	Client   client.Client
	Pipeline *streamv1.BenthosPipeline
}

// +kubebuilder:rbac:groups=streaming.benthos.dev,resources=benthospipelines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=streaming.benthos.dev,resources=benthospipelines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=streaming.benthos.dev,resources=benthospipelines/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile is the main reconcile loop for the Benthos Pipeline
func (r *BenthosPipelineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.Log.WithName("pipeline")

	// fetch BenthosPipeline
	pipeline := &streamv1.BenthosPipeline{}
	err := r.Get(ctx, req.NamespacedName, pipeline)
	if err != nil {
		return reconcile.Result{}, err
	}

	scope := &PipelineScope{
		Log:      &log,
		Ctx:      ctx,
		Client:   r.Client,
		Pipeline: pipeline,
	}

	// handle pipeline deletion
	if !pipeline.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(scope)
	}

	// handle pipeline reconcile
	return r.reconcileNormal(scope)
}

// reconcileNormal handles normal reconciles
func (r *BenthosPipelineReconciler) reconcileNormal(scope *PipelineScope) (ctrl.Result, error) {
	// add finalizer to the BenthosPipeline
	controllerutil.AddFinalizer(scope.Pipeline, streamv1.PipelineFinalizer)

	// check if the Pipeline has already created a deployment
	_, err := r.createOrUpdatePipeline(scope)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// reconcileNormal handles a deletion during the reconcile
func (r *BenthosPipelineReconciler) reconcileDelete(scope *PipelineScope) (ctrl.Result, error) {
	// remove finalizer to allow the resource to be deleted
	controllerutil.RemoveFinalizer(scope.Pipeline, streamv1.PipelineFinalizer)

	return reconcile.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BenthosPipelineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&streamv1.BenthosPipeline{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}

// createOrUpdatePipeline orchestrates the updating and creation of the Benthos deployment
func (r *BenthosPipelineReconciler) createOrUpdatePipeline(scope *PipelineScope) (ctrl.Result, error) {
	pipeline := scope.Pipeline

	// create Benthos ConfigMap
	_, err := r.createOrPatchConfigMap(scope)
	if err != nil {
		return reconcile.Result{}, err
	}

	// create Benthos Deployment
	_, err = r.createOrPatchDeployment(scope)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Fetch deployment to get replicas
	deployment := &appsv1.Deployment{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: pipeline.Name, Namespace: pipeline.Namespace}, deployment)
	if err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "failed to get deployment %s", pipeline.Name)
	}

	// set pipeline phase if not set
	if pipeline.Status.Phase == "" {
		scope.status(false, "Provisioning")
	}

	// Update available replicas on the status
	scope.Log.Info("Found deployment", "status", deployment.Status)
	return r.setPipelineStatus(scope, deployment)
}

// setPipelineStatus sets the latest status of the deployment
func (r *BenthosPipelineReconciler) setPipelineStatus(scope *PipelineScope, deployment *appsv1.Deployment) (res ctrl.Result, resErr error) {
	// defer applying the status change
	defer func() {
		scope.Log.Info("Setting BenthosPipeline status.", "ready", scope.Pipeline.Status.Ready, "phase", scope.Pipeline.Status.Phase)

		err := r.Status().Update(context.Background(), scope.Pipeline)
		if err != nil {
			resErr = err
		}
	}()

	replicas := scope.Pipeline.Spec.Replicas
	status := deployment.Status

	// set available replicas on the pipeline
	scope.Pipeline.Status.AvailableReplicas = deployment.Status.AvailableReplicas

	available := getConditionStatus(status, appsv1.DeploymentAvailable) == corev1.ConditionTrue
	progressing := getConditionStatus(status, appsv1.DeploymentProgressing) == corev1.ConditionTrue

	// how long since deployment creation before reporting an error.
	failedDelay := time.Second * 30
	if !available && deployment.Status.UnavailableReplicas == replicas && deployment.CreationTimestamp.Sub(time.Now()) > failedDelay {
		scope.status(false, "Deployment failed")
		return reconcile.Result{}, errors.New("One or more pods failed to start. Check the logs of the pods.")
	}

	if progressing {
		if deployment.Status.UpdatedReplicas != replicas {
			scope.status(true, "Updating")
			return reconcile.Result{}, nil
		}
		if deployment.Status.ReadyReplicas > replicas {
			scope.status(true, "Scaling down")
			return reconcile.Result{}, nil
		}
		if deployment.Status.ReadyReplicas < replicas {
			scope.status(true, "Scaling up")
			return reconcile.Result{}, nil
		}
	}

	if available {
		scope.status(true, "Running")
	}
	return reconcile.Result{}, nil
}

func getConditionStatus(ds appsv1.DeploymentStatus, condition appsv1.DeploymentConditionType) corev1.ConditionStatus {
	for i := range ds.Conditions {
		c := ds.Conditions[i]
		if c.Type == condition {
			return c.Status
		}
	}
	return corev1.ConditionUnknown
}

// status is a wrapper for settings the pipeline status options
func (ps *PipelineScope) status(ready bool, phase string) {
	ps.Pipeline.Status.Ready = ready
	ps.Pipeline.Status.Phase = phase
}

// createOrPatchDeployment updates a benthos deployment or creates it if it doesn't exist
func (r *BenthosPipelineReconciler) createOrPatchDeployment(scope *PipelineScope) (ctrl.Result, error) {
	name := scope.Pipeline.Name
	namespace := scope.Pipeline.Namespace
	replicas := scope.Pipeline.Spec.Replicas

	scope.Log.Info("Creating deployment", "name", name)

	dp := &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace}}
	op, err := controllerutil.CreateOrPatch(scope.Ctx, scope.Client, dp, func() error {
		template := resource.NewDeployment(name, namespace, replicas)

		// Deployment selector is immutable we only change this value if we're creating a new resource.
		if dp.CreationTimestamp.IsZero() {
			dp.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: template.Spec.Selector.MatchLabels,
			}
			dp.Spec.Template.ObjectMeta = template.Spec.Template.ObjectMeta
		}

		// Update the template, ignore metadata
		dp.Spec.Template.Spec = template.Spec.Template.Spec
		dp.Spec.Replicas = template.Spec.Replicas

		err := controllerutil.SetControllerReference(scope.Pipeline, dp, r.Scheme)
		if err != nil {
			return errors.Wrapf(err, "Failed to set controller reference on deployment %s", name)
		}

		return nil
	})
	if err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "Failed to reconcile deployment %s", name)
	}
	scope.Log.Info("Succesfully reconciled deployment", "name", name, "operation", op)
	return reconcile.Result{}, nil
}

// createOrPatchConfigMap updates a config map or creates it if it doesn't exist
func (r *BenthosPipelineReconciler) createOrPatchConfigMap(scope *PipelineScope) (ctrl.Result, error) {
	name := scope.Pipeline.Name

	scope.Log.Info("Creating config map", "name", name)

	sc := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "benthos-" + name, Namespace: scope.Pipeline.Namespace}}
	op, err := controllerutil.CreateOrPatch(scope.Ctx, scope.Client, sc, func() error {
		sc.Data = map[string]string{
			"benthos.yaml": scope.Pipeline.Spec.Config,
		}
		err := controllerutil.SetControllerReference(scope.Pipeline, sc, r.Scheme)
		if err != nil {
			return errors.Wrapf(err, "Failed to set controller reference on configmap %s", name)
		}
		return nil
	})
	if err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "Failed to reconcile config map %s", name)
	}

	scope.Log.Info("Succesfully reconciled config map", "name", name, "operation", op)

	// rollout the deployment if the configmap changes
	if op == controllerutil.OperationResultUpdated {
		return r.rolloutDeployment(scope)
	}

	return reconcile.Result{}, nil
}

// rolloutDeployment rolls out a new Benthos deployment
func (r *BenthosPipelineReconciler) rolloutDeployment(scope *PipelineScope) (ctrl.Result, error) {
	name := scope.Pipeline.Name
	namespace := scope.Pipeline.Namespace

	scope.Log.Info("Rolling out deployment.", "name", name)

	dp := &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace}}
	body := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"%s":"%s"}}}}}`, "streaming.benthos.dev/restartedAt", time.Now().Format(time.RFC3339))
	err := r.Patch(scope.Ctx, dp, client.RawPatch(types.StrategicMergePatchType, []byte(body)))
	if err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "Failed to rollout deployment %s", name)
	}

	scope.Log.Info("Deployment rollout success", "name", name)
	return reconcile.Result{}, nil
}
