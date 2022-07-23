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

const frontendPort = 3000
const frontendServicePort = 30686
const frontendImage = "jdob/visitors-webui:1.0.0"

func frontendDeploymentName(v *appv1alpha1.VisitorsApp) string {
	return v.Name + "-frontend"
}

func frontendServiceName(v *appv1alpha1.VisitorsApp) string {
	return v.Name + "-frontend-service"
}

func (r *VisitorsAppReconciler) frontendDeployment(ctx context.Context, instance *appv1alpha1.VisitorsApp) *appsv1.Deployment {
	labels := labels(instance, "frontend")
	size := int32(1)

	// If the header was specified, add it as an env variable
	env := []corev1.EnvVar{}
	if instance.Spec.Title != "" {
		env = append(env, corev1.EnvVar{
			Name:  "REACT_APP_TITLE",
			Value: instance.Spec.Title,
		})
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      frontendDeploymentName(instance),
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
						Image: frontendImage,
						Name:  "visitors-webui",
						Ports: []corev1.ContainerPort{{
							ContainerPort: frontendPort,
							Name:          "visitors",
						}},
						Env: env,
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(instance, dep, r.Scheme)
	return dep
}

func (r *VisitorsAppReconciler) frontendService(ctx context.Context, instance *appv1alpha1.VisitorsApp) *corev1.Service {
	reqLogger := log.FromContext(ctx)
	labels := labels(instance, "frontend")

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      frontendServiceName(instance),
			Namespace: instance.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       frontendPort,
				TargetPort: intstr.FromInt(frontendPort),
				NodePort:   frontendServicePort,
			}},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	reqLogger.Info("Service Spec", "Service.Name", s.ObjectMeta.Name)

	controllerutil.SetControllerReference(instance, s, r.Scheme)
	return s
}

func (r *VisitorsAppReconciler) updateFrontendStatus(ctx context.Context, instance *appv1alpha1.VisitorsApp) error {
	instance.Status.FrontendImage = frontendImage
	err := r.Status().Update(ctx, instance)
	return err
}

func (r *VisitorsAppReconciler) handleFrontendChanges(ctx context.Context, instance *appv1alpha1.VisitorsApp) (*reconcile.Result, error) {
	reqLogger := log.FromContext(ctx)
	found := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      frontendDeploymentName(instance),
		Namespace: instance.Namespace,
	}, found)
	if err != nil {
		// The deployment may not have been created yet, so requeue
		return &reconcile.Result{RequeueAfter: 5 * time.Second}, err
	}

	title := instance.Spec.Title
	existing := (*found).Spec.Template.Spec.Containers[0].Env[0].Value

	if title != existing {
		(*found).Spec.Template.Spec.Containers[0].Env[0].Value = title
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
