package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	test_namespace = "test"
)

var test_data = map[string]interface{}{
	"i": int(15),
	"i64": int64(1500),
	"f": float64(6.3),
	"s": "test",
	"b": false,
}

func TestNodes(t *testing.T) {
	n1_name := "n1"
	n1 := NewNode(test_namespace, n1_name)
	err := n1.Data.SetDatas(test_data)
	require.Nil(t, err)
	require.Equal(t, test_namespace, n1.ID.Namespace)
	require.Equal(t, test_namespace, n1.Name.Namespace)
	require.Equal(t, n1_name, n1.Name.Identifier)
	n2 := n1.Clone()
	require.True(t, n1.ID.IsEqual(n2.ID))
	require.Equal(t, test_namespace, n2.ID.Namespace)
	require.Equal(t, test_namespace, n2.Name.Namespace)
	require.Equal(t, n1_name, n2.Name.Identifier)
}

func TestEncDec(t *testing.T) {
	n1_name := "n1"
	n1 := NewNode(test_namespace, n1_name)
	n1.Data.SetDatas(test_data)
	dd, err := Encode(n1)
	require.Nil(t, err)
	n2, err := DecodeNode(dd)
	require.Nil(t, err)
	require.True(t, n1.ID.IsEqual(n2.ID))
	require.Equal(t, test_namespace, n2.ID.Namespace)
	require.Equal(t, test_namespace, n2.Name.Namespace)
	require.Equal(t, n1_name, n2.Name.Identifier)
}