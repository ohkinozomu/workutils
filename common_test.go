package workutils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	rawExtension, err := objToRawExtension(obj)
	require.Nil(t, err)

	decoded, _, err := decode(rawExtension.Raw)
	require.Nil(t, err)

	assert.Equal(t, obj, decoded)
}
