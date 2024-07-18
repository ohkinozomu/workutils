package workutils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	workapiv1 "open-cluster-management.io/api/work/v1"
)

func TestDiff(t *testing.T) {
	client := &WorkUtilsClient{
		Scheme: scheme.Scheme,
	}

	manifest1 := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: config1
  namespace: default
data:
  key1: value1
`

	manifest2 := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: config2
  namespace: default
data:
  key2: value2
`

	manifest3 := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: config1
  namespace: default
data:
  key1: updatedValue
`

	manifest4 := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: config3
  namespace: default
data:
  key3: value3
`

	raw1, err := stringToRawExtension(manifest1)
	require.NoError(t, err)

	raw2, err := stringToRawExtension(manifest2)
	require.NoError(t, err)

	raw3, err := stringToRawExtension(manifest3)
	require.NoError(t, err)

	raw4, err := stringToRawExtension(manifest4)
	require.NoError(t, err)

	work1 := workapiv1.ManifestWork{
		Spec: workapiv1.ManifestWorkSpec{
			Workload: workapiv1.ManifestsTemplate{
				Manifests: []workapiv1.Manifest{{
					RawExtension: raw1,
				}, {
					RawExtension: raw4,
				}},
			},
		},
	}
	work2 := workapiv1.ManifestWork{
		Spec: workapiv1.ManifestWorkSpec{
			Workload: workapiv1.ManifestsTemplate{
				Manifests: []workapiv1.Manifest{{
					RawExtension: raw2,
				}, {
					RawExtension: raw3,
				}},
			},
		},
	}

	added, removed, updated, err := client.Diff(work1, work2)
	require.NoError(t, err)

	o2, _, err := decode([]byte(manifest2), scheme.Scheme)
	require.NoError(t, err)
	require.Equal(t, []runtime.Object{o2}, added)

	o4, _, err := decode([]byte(manifest4), scheme.Scheme)
	require.NoError(t, err)
	require.Equal(t, []runtime.Object{o4}, removed)

	o3, _, err := decode([]byte(manifest3), scheme.Scheme)
	require.NoError(t, err)
	require.Equal(t, []runtime.Object{o3}, updated)
}
