package controllers

import (
	"context"
	"strconv"
	"time"

	costoptimizerv1alpha1 "github.com/KubeInBox/aws-utility-controller/api/v1alpha1"
	"github.com/KubeInBox/aws-utility-controller/pkg/utils"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Ec2CostOptimizerReconciler reconciles a Ec2CostOptimizer object
type Ec2CostOptimizerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Logger *logrus.Logger
}

// SetupWithManager sets up the controller with the Manager.
func (r *Ec2CostOptimizerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
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
	logger := log.FromContext(ctx)

	// your logic here
	logger.WithValues("object", req.NamespacedName)
	r.Logger.WithField("object", req.NamespacedName)
	logger.Info("Reconciling Ec2CostOptimizer ...")

	ec2CostOptimizer := &costoptimizerv1alpha1.Ec2CostOptimizer{}
	err := r.Get(context.TODO(), req.NamespacedName, ec2CostOptimizer)
	if err != nil {
		if errors.IsNotFound(err) {
			// object not found, could have been deleted after
			// reconcile request, hence don't requeue
			logger.V(1).Info("object %s not found")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "unable to fetch object ")
		// error reading the object, requeue the request
		return ctrl.Result{}, err
	}

	switch ec2CostOptimizer.Spec.WindowType {
	case costoptimizerv1alpha1.OnDemand:
		r.Logger.Infof("Handling ondemand ec2 with operation type %v ", ec2CostOptimizer.Spec.Operation)
		if err = r.handleOnDemandEc2Oprn(ec2CostOptimizer); err != nil {
			logger.V(1).Error(err, "error processing onDemand ec2 operation")
			return ctrl.Result{}, err
		}
	case costoptimizerv1alpha1.Scheduled:
		r.Logger.Infof("Handling scheduled ec2 with operation type %v ", ec2CostOptimizer.Spec.Operation)
		if err = r.handleScheduledEc2Oprn(ec2CostOptimizer); err != nil {
			logger.V(1).Error(err, "error processing onDemand ec2 operation")
			return ctrl.Result{}, err
		}
	default:
		logger.V(1).Info("invalid window type specified")
	}

	// TODO: update status

	return ctrl.Result{}, nil
}

// handleOnDemandEc2Oprn will start/stop ec2 instances right away, upon error it will keep retrying
// the operation until it gets succeeded.
func (r *Ec2CostOptimizerReconciler) handleOnDemandEc2Oprn(ec2CostOptimizer *costoptimizerv1alpha1.Ec2CostOptimizer) error {
	switch ec2CostOptimizer.Spec.Operation {
	case costoptimizerv1alpha1.Start:
		if err := utils.StartEc2Instance(ec2CostOptimizer.Spec.InstanceIDs); err != nil {
			return err
		}
	case costoptimizerv1alpha1.Stop:
		if err := utils.StopEc2Instance(ec2CostOptimizer.Spec.InstanceIDs); err != nil {
			return err
		}
	default:
		logrus.Info("invalid ec2 operation type specified for resource %s/%s", ec2CostOptimizer.GetNamespace(), ec2CostOptimizer.GetName())
	}
	return nil
}

func (r *Ec2CostOptimizerReconciler) handleScheduledEc2Oprn(ec2CostOptimizer *costoptimizerv1alpha1.Ec2CostOptimizer) error {
	if !isInTimeWindow(ec2CostOptimizer.Spec.StartTimeWindow, ec2CostOptimizer.Spec.EndTimeWindow) {
		r.Logger.Debugf("ignoring as it is not in scheduled time window")
		return nil
	}
	r.Logger.Infof("current time is within the time window, starting operations")

	// start/stop if it is in given time window
	if err := r.handleScheduledEc2Oprn(ec2CostOptimizer); err != nil {
		return err
	}
	return nil
}

func isInTimeWindow(startTimeWindow, endTimeWindow string) bool {
	if startTimeWindow == "" || endTimeWindow == "" {
		return false
	}

	// load current IST time
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		logrus.Errorf("failed to load timezone location %v", err)
		return false
	}
	now := time.Now().In(loc)

	// return true if it is a weekends
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return true
	}

	// parse timestamp from current IST time
	timeFormat := "15:04:00"
	currentTime, _ := time.Parse(timeFormat, strconv.Itoa(now.Hour())+":"+strconv.Itoa(now.Minute()))
	startTime, _ := time.Parse(timeFormat, startTimeWindow)
	endTime, _ := time.Parse(timeFormat, endTimeWindow)

	if currentTime.After(startTime) && currentTime.Before(endTime) {
		return true
	}

	return false
}
