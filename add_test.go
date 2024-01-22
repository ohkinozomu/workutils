package workutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	workapiv1 "open-cluster-management.io/api/work/v1"
)

func TestAdd(t *testing.T) {
	nsStr := `
apiVersion: v1
kind: Namespace
metadata:
  name: test-namespace
`

	deploymentStr := `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
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
        ports:
        - containerPort: 80
`

	namespace, err := stringToRawExtension(nsStr)
	assert.Nil(t, err)
	deployment, err := stringToRawExtension(deploymentStr)
	assert.Nil(t, err)
	work := workapiv1.ManifestWork{
		Spec: workapiv1.ManifestWorkSpec{
			Workload: workapiv1.ManifestsTemplate{
				Manifests: []workapiv1.Manifest{
					{
						RawExtension: namespace,
					},
				},
			},
		},
	}
	updatedWork, err := Add(work, deployment.Object)
	assert.Nil(t, err)
	assert.Equal(t, len(updatedWork.Spec.Workload.Manifests), 2)
	obj, err := Get(updatedWork, Resource{
		Group:   "apps",
		Version: "v1",
		Kind:    "Deployment",
		Name:    "nginx-deployment",
	})
	assert.Nil(t, err)
	d, ok := obj.(*appsv1.Deployment)
	assert.Equal(t, true, ok)
	assert.Equal(t, "nginx-deployment", d.Name)
}
