//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func waitForDeploymentStatusCondition(t *testing.T, cl client.Client, deploymentName types.NamespacedName, conditions ...appsv1.DeploymentCondition) error {
	t.Helper()
	return wait.PollImmediate(1*time.Second, 5*time.Minute, func() (bool, error) {
		dep := &appsv1.Deployment{}
		if err := cl.Get(context.TODO(), deploymentName, dep); err != nil {
			t.Logf("failed to get deployment %s: %v (retrying)", deploymentName.Name, err)
			return false, nil
		}

		expected := deploymentConditionMap(conditions...)
		current := deploymentConditionMap(dep.Status.Conditions...)
		return conditionsMatchExpected(expected, current), nil
	})
}

func waitForingressStatusCondition(t *testing.T, cl client.Client, ingressName types.NamespacedName, hostnames []corev1.LoadBalancerIngress) error {
	t.Helper()
	return wait.PollImmediate(1*time.Second, 5*time.Minute, func() (bool, error) {
		ing := &networkingv1.Ingress{}
		if err := cl.Get(context.TODO(), ingressName, ing); err != nil {
			t.Logf("failed to get deployment %s: %v (retrying)", ingressName.Name, err)
			return false, nil
		}
		expected := ingressHostnamesMap(hostnames)
		current := ingressHostnamesMap(ing.Status.LoadBalancer.Ingress)
		return conditionsMatchExpected(expected, current), nil
	})
}

func ingressHostnamesMap(loadBalancerIngress []corev1.LoadBalancerIngress) map[string]int32 {
	names := map[string]int32{}
	for _, ing := range loadBalancerIngress {
		names[ing.Hostname] = 1
	}
	return names
}

func deploymentConditionMap(conditions ...appsv1.DeploymentCondition) map[string]string {
	conds := map[string]string{}
	for _, cond := range conditions {
		conds[string(cond.Type)] = string(cond.Status)
	}
	return conds
}

func conditionsMatchExpected(expected, actual map[string]string) bool {
	filtered := map[string]string{}
	for k := range actual {
		if _, comparable := expected[k]; comparable {
			filtered[k] = actual[k]
		}
	}
	return reflect.DeepEqual(expected, filtered)
}

// buildEchoPod returns a pod definition for an socat-based echo server.
func buildEchoPod(name, namespace string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app": "echo",
			},
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					// Note that HTTP/1.0 will strip the HSTS response header
					Args: []string{
						"TCP4-LISTEN:8080,reuseaddr,fork",
						`EXEC:'/bin/bash -c \"printf \\\"HTTP/1.0 200 OK\r\n\r\n\\\"; sed -e \\\"/^\r/q\\\"\"'`,
					},
					Command: []string{"/bin/socat"},
					Image:   "openshift/origin-node",
					Name:    "echo",
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: int32(8080),
							Protocol:      corev1.ProtocolTCP,
						},
					},
				},
			},
		},
	}
}

// buildEchoService returns a service definition for an HTTP service.
func buildEchoService(name, namespace string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       int32(80),
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(8080),
				},
			},
			Selector: map[string]string{
				"app": "echo",
			},
		},
	}
}

func buildDefaultEchoIngress(name, namespace, backendSvc string, backendSvcPort int32) networkingv1.Ingress {
	return networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: aws.String("alb"),
			DefaultBackend: &networkingv1.IngressBackend{
				Service: &networkingv1.IngressServiceBackend{
					Name: backendSvc,
					Port: networkingv1.ServiceBackendPort{
						Number: backendSvcPort,
					},
				},
			},
		},
	}
}
