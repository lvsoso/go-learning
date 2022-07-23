package controllers

import (
	"context"

	appv1alpha1 "github.com/lvsoso/visitor-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *VisitorsAppReconciler) ensureDeployment(ctx context.Context, req ctrl.Request,
	instance *appv1alpha1.VisitorsApp,
	dep *appsv1.Deployment,
) (*ctrl.Result, error) {
	reqLogger := log.FromContext(ctx)
	// See if deployment already exists and create if it doesn't
	found := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      dep.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the deployment
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)

		if err != nil {
			// Deployment failed
			reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return &ctrl.Result{}, err
		} else {
			// Deployment was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the deployment not existing
		reqLogger.Error(err, "Failed to get Deployment")
		return &ctrl.Result{}, err
	}

	return nil, nil
}

func (r *VisitorsAppReconciler) ensureService(ctx context.Context, req ctrl.Request,
	instance *appv1alpha1.VisitorsApp,
	s *corev1.Service,
) (*ctrl.Result, error) {
	reqLogger := log.FromContext(ctx)
	found := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      s.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the service
		reqLogger.Info("Creating a new Service", "Service.Namespace", s.Namespace, "Service.Name", s.Name)
		err = r.Create(ctx, s)

		if err != nil {
			// Creation failed
			reqLogger.Error(err, "Failed to create new Service", "Service.Namespace", s.Namespace, "Service.Name", s.Name)
			return &ctrl.Result{}, err
		} else {
			// Creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the service not existing
		reqLogger.Error(err, "Failed to get Service")
		return &ctrl.Result{}, err
	}

	return nil, nil
}

func (r *VisitorsAppReconciler) ensureSecret(ctx context.Context, req ctrl.Request,
	instance *appv1alpha1.VisitorsApp,
	s *corev1.Secret,
) (*ctrl.Result, error) {
	reqLogger := log.FromContext(ctx)
	found := &corev1.Secret{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      s.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {
		// Create the secret
		reqLogger.Info("Creating a new secret", "Secret.Namespace", s.Namespace, "Secret.Name", s.Name)
		err = r.Create(ctx, s)

		if err != nil {
			// Creation failed
			reqLogger.Error(err, "Failed to create new Secret", "Secret.Namespace", s.Namespace, "Secret.Name", s.Name)
			return &ctrl.Result{}, err
		} else {
			// Creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the secret not existing
		reqLogger.Error(err, "Failed to get Secret")
		return &ctrl.Result{}, err
	}

	return nil, nil
}

func labels(v *appv1alpha1.VisitorsApp, tier string) map[string]string {
	return map[string]string{
		"app":             "visitors",
		"visitorssite_cr": v.Name,
		"tier":            tier,
	}
}
