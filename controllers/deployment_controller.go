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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// DeploymentReconciler reconciles a Deployment object
type DeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// It's called each time a Deployment is created, updated, or deleted.
// When a Deployment is created or updated, it makes sure the Pod template features
// the desired sidecar container. When a Deployment is deleted, it ignores the deletion.
func (r *DeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the Deployment from the Kubernetes API.
	var deployment appsv1.Deployment
	if err := r.Get(ctx, req.NamespacedName, &deployment); err != nil {
		if apierrors.IsNotFound(err) {
			// Ignore not-found errors that occur when Deployments are deleted.
			return ctrl.Result{}, nil
		}

		log.Error(err, "unable to fetch Deployment")

		return ctrl.Result{}, err
	}

	// sidecar is a simple busybox-based container that sleeps for 36000.
	// The sidecar container is always named "<deploymentname>-sidecar".
	sidecar := corev1.Container{
		Name:    fmt.Sprintf("%s-sidecar", deployment.Name),
		Image:   "busybox",
		Command: []string{"sleep"},
		Args:    []string{"36000"},
	}

	// This is a crude way to ensure the controller doesn't attempt to add
	// redundant sidecar containers, which would result in an error a la:
	// Deployment.apps \"foo\" is invalid: spec.template.spec.containers[2].name: Duplicate value: \"foo-sidecar\"
	for _, c := range deployment.Spec.Template.Spec.Containers {
		if c.Name == sidecar.Name && c.Image == sidecar.Image {
			return ctrl.Result{}, nil
		}
	}

	// Otherwise, add the sidecar to the deployment's containers.
	deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, sidecar)

	if err := r.Update(ctx, &deployment); err != nil {
		// The Deployment has been updated or deleted since initially reading it.
		if apierrors.IsConflict(err) || apierrors.IsNotFound(err) {
			// Requeue the Deployment to try to reconciliate again.
			return ctrl.Result{Requeue: true}, nil
		}

		log.Error(err, "unable to update Deployment")

		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		Complete(r)
}
