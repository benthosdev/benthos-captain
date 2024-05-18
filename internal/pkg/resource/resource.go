package resource

import (
	captainv1 "github.com/benthosdev/benthos-captain/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const DefaultImage = "jeffail/benthos:4.22"

func NewDeployment(name string, namespace string, scope captainv1.PipelineSpec) *appsv1.Deployment {
	labels := map[string]string{
		"app.kubernetes.io/name":     "benthos",
		"app.kubernetes.io/instance": name,
	}

	image := DefaultImage
	if scope.Image != "" {
		image = scope.Image
	}

	podSpec := corev1.PodSpec{
		Containers: []corev1.Container{{
			Image:           image,
			ImagePullPolicy: corev1.PullAlways,
			Name:            "benthos",
			Ports: []corev1.ContainerPort{{
				ContainerPort: 4195,
				Name:          "http",
			}},
			Args: []string{
				"-c",
				"/config/benthos.yaml",
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "config",
					MountPath: "/config",
					ReadOnly:  true,
				},
			},
			Env: scope.Env,
		}},
		Volumes: []corev1.Volume{
			{
				Name: "config",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "benthos-" + name,
						},
					},
				},
			},
		},
	}

	if scope.ImagePullSecret != "" {
		podSpec.ImagePullSecrets = []corev1.LocalObjectReference{
			{Name: scope.ImagePullSecret},
		}
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &scope.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: podSpec,
			},
		},
	}
}
