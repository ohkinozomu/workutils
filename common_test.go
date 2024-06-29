package workutils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

func TestStringToRawExtension(t *testing.T) {
	nsStr := `
apiVersion: v1
kind: Namespace
metadata:
  name: test-namespace
`

	raw, err := stringToRawExtension(nsStr)
	require.Nil(t, err)
	rawStr := string(raw.Raw)
	assert.Equal(t, strings.TrimSpace(nsStr), strings.TrimSpace(rawStr))
}

func TestObjToRawExtension(t *testing.T) {
	s1 := scheme.Scheme
	s2 := scheme.Scheme
	apiextensionsv1.AddToScheme(s2)

	testCases := []struct {
		name   string
		str    string
		scheme *runtime.Scheme
	}{
		{
			name: "Namespace",
			str: `
apiVersion: v1
kind: Namespace
metadata:
  name: test-namespace
`,
			scheme: s1,
		},
		{
			name: "CRD",
			str: `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: examplecrd.sample.com
spec:
  group: sample.com
  versions:
  - name: v1
    served: true
    storage: true
  scope: Namespaced
  names:
  plural: examplecrds
  singular: examplecrd
  kind: ExampleCRD
  shortNames:
    - ec
`,
			scheme: s2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			obj, _, err := decode([]byte(tc.str), tc.scheme)
			require.Nil(t, err)

			rawExtension, err := objToRawExtension(obj, tc.scheme)
			require.Nil(t, err)

			decoded, _, err := decode(rawExtension.Raw, tc.scheme)
			require.Nil(t, err)

			assert.Equal(t, obj, decoded)
		})
	}
}
