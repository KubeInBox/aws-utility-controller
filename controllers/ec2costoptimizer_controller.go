package controllers

import (
	"context"
	"fmt"
	"time"

	costoptimizerv1alpha1 "github.com/KubeInBox/aws-utility-controller/api/v1alpha1"
	"github.com/KubeInBox/aws-utility-controller/pkg/utils"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const (
	inProgress      = "InProgress"
	failed          = "Failed"
	complete        = "Completed"
	inTimeWindow    = "InTimeWindow"
	outOfTimeWindow = "OutOfTimeWindow"
)

// Ec2CostOptimizerReconciler reconciles a Ec2CostOptimizer object
type Ec2CostOptimizerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	logger logr.Logger
}

// SetupWithManager sets up the controller with the Manager.
func (r *Ec2CostOptimizerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	pred := predicate.GenerationChangedPredicate{}
	return ctrl.NewControllerManagedBy(mgr).
		WithEventFilter(pred).
		For(&costoptimizerv1alpha1.Ec2CostOptimizer{}).
		Complete(r)
}

//+kubebuilder:rbac:groups=kubeinbox.io.kubeinbox.io,resources=ec2costoptimizers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kubeinbox.io.kubeinbox.io,resources=ec2costoptimizers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kubeinbox.io.kubeinbox.io,resources=ec2costoptimizers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.1/pkg/reconcile
func (r *Ec2CostOptimizerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.logger = log.FromContext(ctx)
	// your logic here
	r.logger.WithValues("object", req.NamespacedName)
	r.logger.Info("Reconciling Ec2CostOptimizer ...")

	ec2CostOptimizer := &costoptimizerv1alpha1.Ec2CostOptimizer{}
	err := r.Get(context.TODO(), req.NamespacedName, ec2CostOptimizer)
	if err != nil {
		if errors.IsNotFound(err) {
			// object not found, could have been deleted after
			// reconcile request, hence don't requeue
			r.logger.V(1).Info("object not found")
			return ctrl.Result{}, nil
		}
		r.logger.Error(err, "unable to fetch object ")
		// error reading the object, requeue the request
		return ctrl.Result{}, err
	}

	switch ec2CostOptimizer.Spec.WindowType {
	case costoptimizerv1alpha1.OnDemand:
		if ec2CostOptimizer.Status.State == fmt.Sprintf("%s/%s", costoptimizerv1alpha1.OnDemand, complete) {
			r.logger.V(1).Info("ignoring already processed onDemand object")
			return ctrl.Result{}, nil
		}
		r.logger.V(1).Info("Handling onDemand ec2 with operation", "type", ec2CostOptimizer.Spec.Operation)
		r.UpdateStatus(ctx, ec2CostOptimizer, inProgress)
		if err = r.handleOnDemandEc2Oprn(ec2CostOptimizer); err != nil {
			r.logger.Error(err, "error processing onDemand ec2 operation")
			r.UpdateStatus(ctx, ec2CostOptimizer, failed)
			return ctrl.Result{}, err
		}
		r.UpdateStatus(ctx, ec2CostOptimizer, complete)
	case costoptimizerv1alpha1.Scheduled:
		r.logger.V(1).Info("Handling scheduled ec2 with operation", "type", ec2CostOptimizer.Spec.Operation)
		err = r.handleScheduledEc2Oprn(ctx, ec2CostOptimizer)
		if err != nil {
			r.logger.Error(err, "error processing scheduled ec2 operation")
		}
		return ctrl.Result{RequeueAfter: wait.Jitter(1*time.Minute, 0.5)}, err
	default:
		r.logger.V(1).Info("invalid window type specified")
	}

	// TODO: update status

	return ctrl.Result{}, nil
}

// handleOnDemandEc2Oprn will start/stop ec2 instances right away, upon error it will keep retrying
// the operation until it gets succeeded.
func (r *Ec2CostOptimizerReconciler) handleOnDemandEc2Oprn(ec2CostOptimizer *costoptimizerv1alpha1.Ec2CostOptimizer) error {
	switch ec2CostOptimizer.Spec.Operation {
	case costoptimizerv1alpha1.Start:
		if err := utils.StartEc2Instance(r.logger, ec2CostOptimizer.Spec.InstanceIDs); err != nil {
			return err
		}
	case costoptimizerv1alpha1.Stop:
		if err := utils.StopEc2Instance(r.logger, ec2CostOptimizer.Spec.InstanceIDs); err != nil {
			return err
		}
	default:
		r.logger.Info("specified invalid ec2 operation type")
	}
	return nil
}

func (r *Ec2CostOptimizerReconciler) handleScheduledEc2Oprn(ctx context.Context, ec2CostOptimizer *costoptimizerv1alpha1.Ec2CostOptimizer) error {
	if !isInTimeWindow(r.logger, ec2CostOptimizer.Spec.StartTimeWindow, ec2CostOptimizer.Spec.EndTimeWindow) {
		r.UpdateStatus(ctx, ec2CostOptimizer, outOfTimeWindow)
		r.logger.Info("ignoring as it is not in scheduled time window")
		// perform counter operation, if it was stopped in time window then start or vice-versa.
		return nil
	}
	r.UpdateStatus(ctx, ec2CostOptimizer, inTimeWindow)
	r.logger.Info("current time is within the time window, starting operations")

	// start/stop if it is in given time window
	if err := r.handleOnDemandEc2Oprn(ec2CostOptimizer); err != nil {
		return err
	}
	return nil
}

func (r *Ec2CostOptimizerReconciler) UpdateStatus(ctx context.Context, obj *costoptimizerv1alpha1.Ec2CostOptimizer, msg string) {
	// create patches for the object and its possible status
	statusPatch := client.MergeFrom(obj.DeepCopy())

	obj.Status.State = fmt.Sprintf("%s/%s", obj.Spec.WindowType, msg)
	data, err := statusPatch.Data(obj)
	if err != nil {
		return
	}

	// no change in the patch we get length of 2.
	if len(data) <= 2 {
		return
	}

	// patch status of a given object.
	err = r.Status().Patch(ctx, obj, statusPatch)
	if err != nil {
		r.logger.Error(err, "failed to update status")
		return
	}

	r.logger.Info(fmt.Sprintf("updated status with state %s", obj.Status.State))
}

func isInTimeWindow(logger logr.Logger, startTimeWindow, endTimeWindow string) bool {
	if startTimeWindow == "" || endTimeWindow == "" {
		return false
	}

	// load current IST time
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		logger.Error(err, "failed to load timezone location")
		return false
	}
	now := time.Now().In(loc)

	// parse timestamp from current IST time
	timeFormat := "15:04:05"
	currTimeFormat := fmt.Sprintf("%02d:%02d:%02d", now.Hour(), now.Minute(), now.Second())
	currentTime, err := time.Parse(timeFormat, currTimeFormat)
	if err != nil {
		logger.Error(err, "failed to parse current time")
		return false
	}
	startTime, err := time.Parse(timeFormat, startTimeWindow)
	if err != nil {
		logger.Error(err, "invalid start time")
		return false
	}
	endTime, err := time.Parse(timeFormat, endTimeWindow)
	if err != nil {
		logger.Error(err, "invalid end time")
		return false
	}

	logger.V(1).Info("", "curr time", currTimeFormat, "start time", startTimeWindow, "end time", endTimeWindow)
	if currentTime.After(startTime) && currentTime.Before(endTime) {
		return true
	}

	return false
}
