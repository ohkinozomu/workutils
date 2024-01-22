package workutils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
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
	nsStr := `
apiVersion: v1
kind: Namespace
metadata:
  name: test-namespace
`
	obj, _, err := decode([]byte(nsStr))
	require.Nil(t, err)

	_, err = objToRawExtension(obj)
	require.Nil(t, err)

	ns, ok := obj.(*v1.Namespace)
	require.True(t, ok)

	assert.Equal(t, "test-namespace", ns.Name)
}
