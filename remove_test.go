package workutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/scheme"
	workapiv1 "open-cluster-management.io/api/work/v1"
)

func TestRemove(t *testing.T) {
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
					{
						RawExtension: deployment,
					},
				},
			},
		},
	}
	client := NewWorkUtilsClient(scheme.Scheme)
	updatedWork, err := client.Remove(work, Resource{
		Group:   "",
		Version: "v1",
		Kind:    "Namespace",
		Name:    "test-namespace",
	})
	assert.Nil(t, err)

	assert.Equal(t, 1, len(updatedWork.Spec.Workload.Manifests))
}
