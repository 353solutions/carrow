package carrow

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppendInt64(t *testing.T) {
	require := require.New(t)
	bld := NewInteger64ArrayBuilder()
	require.NotNil(bld, "create builder")

	const size = 20913
	for i := int64(0); i < size; i++ {
		err := bld.Append(i)
		require.NoErrorf(err, "append %d", i)
	}

	arr, err := bld.Finish()
	require.NoError(err, "finish")

	arrLen := arr.Length()
	require.Equal(arrLen, size, "length")

	for i := 0; i < size; i++ {
		val, err := arr.Int64At(i)
		require.NoErrorf(err, "Int64At %d", i)
		require.Equalf(int64(i), val, "value at %d", i)
	}

}

func BenchmarkAppendInt64(b *testing.B) {
	b.StopTimer()
	bld := NewInteger64ArrayBuilder()
	if bld == nil {
		b.Fatal("create builder")
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		bld.Append(int64(i))
	}
}
