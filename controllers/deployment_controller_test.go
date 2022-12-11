package controllers

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("sidecar injector controller", func() {

	const (
		deploymentName = "test-deployment"
		namespace      = "default"
		timeout        = time.Second * 10
		interval       = time.Millisecond * 250
	)

	Context("When a Deployment is created", func() {
		It("Should inject a sidecar container", func() {
			ctx := context.Background()
			deployment := &appsv1.Deployment{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "apps/v1",
					Kind:       "Deployment",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      deploymentName,
					Namespace: namespace,
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app": deploymentName,
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Name: deploymentName,
							Labels: map[string]string{
								"app": deploymentName,
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test-container",
									Image: "test-image",
								},
							},
							RestartPolicy: corev1.RestartPolicyAlways,
						},
					},
				},
			}

			Expect(k8sClient.Create(ctx, deployment)).Should(Succeed())

			var result appsv1.Deployment
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{
					Namespace: namespace,
					Name:      deploymentName,
				}, &result)
				if err != nil {
					return false
				}
				if len(result.Spec.Template.Spec.Containers) != 2 {
					return false
				}

				return true
			}, timeout, interval).Should(BeTrue())

			Expect(result.Spec.Template.Spec.Containers[1].Name).Should(Equal(fmt.Sprintf("%s-sidecar", deploymentName)))
			Expect(result.Spec.Template.Spec.Containers[1].Image).Should(Equal("busybox"))
		})
	})
})
