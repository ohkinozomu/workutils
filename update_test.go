package workutils

import (
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	workapiv1 "open-cluster-management.io/api/work/v1"
)

func TestUpdate(t *testing.T) {
	deploymentStr := `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  namespace: nginx
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        env:
        - name: FOO
          value: "BAR"
        ports:
        - containerPort: 80
`
	raw, err := stringToRawExtension(deploymentStr)
	require.Nil(t, err)
	work := workapiv1.ManifestWork{
		Spec: workapiv1.ManifestWorkSpec{
			Workload: workapiv1.ManifestsTemplate{
				Manifests: []workapiv1.Manifest{
					{
						RawExtension: raw,
					},
				},
			},
		},
	}

	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-deployment",
			Namespace: "nginx",
			Labels: map[string]string{
				"app": "nginx",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32(3),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "nginx",
							Image: "nginx:1.14.2",
							Env: []v1.EnvVar{
								{
									Name:  "FOO",
									Value: "UPDATED",
								},
							},
							Ports: []v1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	updatedWork, err := Update(work, deployment)
	require.Nil(t, err)
	require.Equal(t, len(updatedWork.Spec.Workload.Manifests), 1)

	obj, err := Get(updatedWork, Resource{
		Group:     "apps",
		Version:   "v1",
		Kind:      "Deployment",
		Name:      "nginx-deployment",
		Namespace: "nginx",
	})
	require.Nil(t, err)
	updatedDeployment := obj.(*appsv1.Deployment)
	require.Equal(t, "UPDATED", updatedDeployment.Spec.Template.Spec.Containers[0].Env[0].Value)
}
