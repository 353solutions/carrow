package carrow

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	testArrSize = 2079
)

func TestArrayBoolGet(t *testing.T) {
	require := require.New(t)
	b := NewBoolArrayBuilder()
	require.NotNil(b.ptr, "create")

	const mod = 7
	for i := 0; i < testArrSize; i++ {
		b.Append(i%mod == 0)
	}

	arr, err := b.Finish()
	require.NoError(err, "finish")

	for i := 0; i < testArrSize; i++ {
		v, err := arr.BoolAt(i)
		require.NoError(err, "bool at %d - error", i)
		require.Equalf(i%mod == 0, v, "bool at %d", i)
	}
}

func TestArrayIntGet(t *testing.T) {
	require := require.New(t)
	b := NewInt64ArrayBuilder()
	require.NotNil(b.ptr, "create")

	for i := int64(0); i < testArrSize; i++ {
		b.Append(i)
	}

	arr, err := b.Finish()
	require.NoError(err, "finish")

	for i := 0; i < testArrSize; i++ {
		v, err := arr.Int64At(i)
		require.NoError(err, "int at %d - error", i)
		require.Equalf(int64(i), v, "int at %d", i)
	}
}

func TestArrayFloatGet(t *testing.T) {
	require := require.New(t)
	b := NewFloat64ArrayBuilder()
	require.NotNil(b.ptr, "create")

	for i := 0; i < testArrSize; i++ {
		b.Append(float64(i))
	}

	arr, err := b.Finish()
	require.NoError(err, "finish")

	for i := 0; i < testArrSize; i++ {
		v, err := arr.Float64At(i)
		require.NoError(err, "float at %d - error", i)
		require.Equalf(float64(i), v, "float at %d", i)
	}
}

func TestArrayStringGet(t *testing.T) {
	require := require.New(t)
	b := NewStringArrayBuilder()
	require.NotNil(b.ptr, "create")

	ival := func(i int) string {
		return fmt.Sprintf("%d: value", i+1)
	}

	for i := 0; i < testArrSize; i++ {
		b.Append(ival(i))
	}

	arr, err := b.Finish()
	require.NoError(err, "finish")

	for i := 0; i < testArrSize; i++ {
		v, err := arr.StringAt(i)
		require.NoError(err, "string at %d - error", i)
		require.Equalf(ival(i), v, "string at %d", i)
	}
}

func TestArrayTimeGet(t *testing.T) {
	require := require.New(t)
	b := NewTimestampArrayBuilder()
	require.NotNil(b.ptr, "create")

	start := time.Now()

	tval := func(i int) time.Time {
		return start.Add(time.Duration(i) * time.Millisecond * 372)
	}

	for i := 0; i < testArrSize; i++ {
		b.Append(tval(i))
	}

	arr, err := b.Finish()
	require.NoError(err, "finish")

	for i := 0; i < testArrSize; i++ {
		v, err := arr.TimeAt(i)
		require.NoError(err, "time at %d - error", i)
		// .Equal is better than ==
		require.True(v.Equal(tval(i)), "time at %d", i)
	}
}
