package workutils

import (
	"context"
	"testing"
	"time"

	sh "github.com/ohkinozomu/go-sh"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/pointer"
	workclient "open-cluster-management.io/api/client/work/clientset/versioned"
	workapiv1 "open-cluster-management.io/api/work/v1"
)

func TestE2E(t *testing.T) {
	getResult, err := sh.RunR("kind get clusters")
	require.Nil(t, err)
	if getResult == "No kind clusters found." {
		t.Log("Creating kind cluster...")
		err := sh.Run("kind create cluster")
		require.Nil(t, err)
	}

	if err := sh.Run("which clusteradm"); err != nil {
		t.Log("Installing clusteradm...")
		err := sh.Run("curl -L https://raw.githubusercontent.com/open-cluster-management-io/clusteradm/main/install.sh | bash")
		require.Nil(t, err)
	}

	t.Log("Initializing cluster with clusteradm...")
	err = sh.Run("clusteradm init")
	require.Nil(t, err)

	// hack
	time.Sleep(40 * time.Second)

	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	require.Nil(t, err)

	workClientset, err := workclient.NewForConfig(config)
	require.Nil(t, err)

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

	rawExtension, err := objToRawExtension(deployment)
	require.Nil(t, err)
	manifestWork := &workapiv1.ManifestWork{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx",
			Namespace: "default",
		},
		Spec: workapiv1.ManifestWorkSpec{
			Workload: workapiv1.ManifestsTemplate{
				Manifests: []workapiv1.Manifest{
					{
						RawExtension: rawExtension,
					},
				},
			},
		},
	}

	_, err = workClientset.WorkV1().ManifestWorks("default").Create(context.Background(), manifestWork, metav1.CreateOptions{})
	require.Nil(t, err)
}
