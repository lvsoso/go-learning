package controllers

import (
	"context"
	"time"

	appv1alpha1 "github.com/lvsoso/visitor-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const backendPort = 8000
const backendServicePort = 30685
const backendImage = "jdob/visitors-service:1.0.0"

func backendDeploymentName(v *appv1alpha1.VisitorsApp) string {
	return v.Name + "-backend"
}

func backendServiceName(v *appv1alpha1.VisitorsApp) string {
	return v.Name + "-backend-service"
}

func (r *VisitorsAppReconciler) backendDeployment(ctx context.Context, instance *appv1alpha1.VisitorsApp) *appsv1.Deployment {
	labels := labels(instance, "backend")
	size := instance.Spec.Size

	userSecret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: mysqlAuthName()},
			Key:                  "username",
		},
	}

	passwordSecret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: mysqlAuthName()},
			Key:                  "password",
		},
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      backendDeploymentName(instance),
			Namespace: instance.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           backendImage,
						ImagePullPolicy: corev1.PullAlways,
						Name:            "visitors-service",
						Ports: []corev1.ContainerPort{{
							ContainerPort: backendPort,
							Name:          "visitors",
						}},
						Env: []corev1.EnvVar{
							{
								Name:  "MYSQL_DATABASE",
								Value: "visitors",
							},
							{
								Name:  "MYSQL_SERVICE_HOST",
								Value: mysqlServiceName(),
							},
							{
								Name:      "MYSQL_USERNAME",
								ValueFrom: userSecret,
							},
							{
								Name:      "MYSQL_PASSWORD",
								ValueFrom: passwordSecret,
							},
						},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(instance, dep, r.Scheme)
	return dep
}

func (r *VisitorsAppReconciler) backendService(ctx context.Context, instance *appv1alpha1.VisitorsApp) *corev1.Service {
	labels := labels(instance, "backend")

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      backendServiceName(instance),
			Namespace: instance.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       backendPort,
				TargetPort: intstr.FromInt(backendPort),
				NodePort:   30685,
			}},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	controllerutil.SetControllerReference(instance, s, r.Scheme)
	return s
}

func (r *VisitorsAppReconciler) updateBackendStatus(ctx context.Context, instance *appv1alpha1.VisitorsApp) error {
	instance.Status.BackendImage = backendImage
	err := r.Status().Update(ctx, instance)
	return err
}

func (r *VisitorsAppReconciler) handleBackendChanges(ctx context.Context, instance *appv1alpha1.VisitorsApp) (*reconcile.Result, error) {
	reqLogger := log.FromContext(ctx)

	found := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      backendDeploymentName(instance),
		Namespace: instance.Namespace,
	}, found)
	if err != nil {
		// The deployment may not have been created yet, so requeue
		return &reconcile.Result{RequeueAfter: 5 * time.Second}, err
	}

	size := instance.Spec.Size

	if size != *found.Spec.Replicas {
		found.Spec.Replicas = &size
		err = r.Update(ctx, found)
		if err != nil {
			reqLogger.Error(err, "Failed to update Deployment.", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return &reconcile.Result{}, err
		}
		// Spec updated - return and requeue
		return &reconcile.Result{Requeue: true}, nil
	}

	return nil, nil
}
