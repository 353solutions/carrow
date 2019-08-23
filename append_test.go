package carrow

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const arrSize = 10000

func buildInt(require *require.Assertions) *Array {
	bld := NewInteger64ArrayBuilder()
	require.NotNil(bld, "new")
	for i := int64(0); i < arrSize; i++ {
		err := bld.Append(i)
		require.NoErrorf(err, "append %d", i)
	}

	arr, err := bld.Finish()
	require.NoError(err, "finish")
	return arr
}

func BenchmarkAppendInt64(b *testing.B) {
	require := require.New(b)
	for i := 0; i < b.N; i++ {
		arr := buildInt(require)
		require.Equal(arrSize, arr.Length(), "length")
	}
}
