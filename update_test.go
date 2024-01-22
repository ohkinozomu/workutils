package workutils

import (
	"testing"

	testify "github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	workapiv1 "open-cluster-management.io/api/work/v1"
)

func TestUpdate(t *testing.T) {
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
	updatedWork, err := Update(work, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "updated",
		},
	})
	testify.Nil(t, err)

	_, err = Get(updatedWork, Resource{
		Group:   "",
		Version: "v1",
		Kind:    "Namespace",
		Name:    "updated",
	})
	testify.Nil(t, err)
}
