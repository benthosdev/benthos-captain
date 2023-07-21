package resource

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewDeployment(name string, namespace string, replicas int32) *appsv1.Deployment {
	labels := map[string]string{
		"app.kubernetes.io/name":     "benthos",
		"app.kubernetes.io/instance": name,
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           "jeffail/benthos:4.18",
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
				},
			},
		},
	}
}
