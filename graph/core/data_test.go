package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestData(t *testing.T) {
	test_data := map[string]interface{}{
		"i": int(15),
		"i64": int64(1500),
		"f": float64(6.3),
		"s": "test",
		"b": false,
		"a": []string{"1", "2",},
	}
	d1 := newData()
	err := d1.SetDatas(test_data)
	require.Nil(t, err)
	for k, v := range test_data {
		vv, ok := d1.GetData(k)
		require.True(t, ok)
		require.Equal(t, v, vv)
	}
	d2 := d1.Clone()
	require.True(t, d1.IsEqual(d2))
	d2.SetData("i", float64(0))
	require.False(t, d1.IsEqual(d2))
}