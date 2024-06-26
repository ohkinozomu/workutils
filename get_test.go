package workutils

import (
	"testing"

	testify "github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	workapiv1 "open-cluster-management.io/api/work/v1"
)

func TestGet(t *testing.T) {
	nsStr := `
apiVersion: v1
kind: Namespace
metadata:
  name: test-namespace
`
	raw, err := stringToRawExtension(nsStr)
	testify.Nil(t, err)
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
	client := NewWorkUtilsClient(scheme.Scheme)
	obj, err := client.Get(work, Resource{
		Group:   "",
		Version: "v1",
		Kind:    "Namespace",
		Name:    "test-namespace",
	})
	testify.Nil(t, err)
	namespace, ok := obj.(*corev1.Namespace)
	testify.Equal(t, true, ok)
	testify.Equal(t, "test-namespace", namespace.Name)
}
