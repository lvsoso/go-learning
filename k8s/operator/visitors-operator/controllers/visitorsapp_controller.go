/*
Copyright 2022.

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

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/juju/errors"
	appv1alpha1 "github.com/lvsoso/visitor-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// VisitorsAppReconciler reconciles a VisitorsApp object
type VisitorsAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=app.example.com,resources=visitorsapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app.example.com,resources=visitorsapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=app.example.com,resources=visitorsapps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the VisitorsApp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *VisitorsAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.FromContext(ctx)
	reqLogger.Info("Reconciling VisitorsApp")

	// TODO(user): your logic here
	instance := &appv1alpha1.VisitorsApp{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	var result *ctrl.Result

	// == MySQL ==========
	result, err = r.ensureSecret(ctx, req, instance, r.mysqlAuthSecret(ctx, instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureDeployment(ctx, req, instance, r.mysqlDeployment(ctx, instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(ctx, req, instance, r.mysqlService(ctx, instance))
	if result != nil {
		return *result, err
	}

	mysqlRunning := r.isMysqlUp(ctx, instance)

	if !mysqlRunning {
		// If MySQL isn't running yet, requeue the reconcile
		// to run again after a delay
		delay := time.Second * time.Duration(5)

		reqLogger.Info(fmt.Sprintf("MySQL isn't running, waiting for %s", delay))
		return ctrl.Result{RequeueAfter: delay}, nil
	}

	// == Visitors Backend  ==========
	result, err = r.ensureDeployment(ctx, req, instance, r.backendDeployment(ctx, instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(ctx, req, instance, r.backendService(ctx, instance))
	if result != nil {
		return *result, err
	}

	err = r.updateBackendStatus(ctx, instance)
	if err != nil {
		// Requeue the request if the status could not be updated
		return ctrl.Result{}, err
	}

	result, err = r.handleBackendChanges(ctx, instance)
	if result != nil {
		return *result, err
	}

	// == Visitors Frontend ==========
	result, err = r.ensureDeployment(ctx, req, instance, r.frontendDeployment(ctx, instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(ctx, req, instance, r.frontendService(ctx, instance))
	if result != nil {
		return *result, err
	}

	err = r.updateFrontendStatus(ctx, instance)
	if err != nil {
		// Requeue the request
		return reconcile.Result{}, err
	}

	result, err = r.handleFrontendChanges(ctx, instance)
	if result != nil {
		return *result, err
	}

	//// with error
	return ctrl.Result{}, err
	//// without an error
	// return ctrl.Result{Requeue: true}, nil
	//// stop the reconcile
	// return ctrl.Result{}, nil
	/// /reconcile again after X time
	//  return ctrl.Result{RequeueAfter: nextRun.Sub(r.Now())}, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *VisitorsAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr). // provides a controller builder that allows various controller configurations.
							For(&appv1alpha1.VisitorsApp{}). // specifies the VisitorsApp type as the primary resource to watch.
							Owns(&appsv1.Deployment{}).      // specifies the Deployments type as the secondary resource to watch.
		// WithOptions(controller.Options{MaxConcurrentReconciles: 2}).
		Complete(r)
}
