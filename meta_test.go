package carrow

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetadata(t *testing.T) {
	require := require.New(t)
	m := NewMetadata()

	size, err := m.Len()
	require.NoError(err, "empty len")
	require.Equal(0, size, "empty len")

	keyValOf := func(i int) (string, string) {
		return fmt.Sprintf("value-%d", i), fmt.Sprintf("key-%d", i)
	}

	n := 173
	for i := 0; i < n; i++ {
		key, value := keyValOf(i)
		m.Set(key, value)
	}
	size, err = m.Len()
	require.NoError(err, "len")
	require.Equal(n, size, "wrong len")

	for i := 0; i < n; i++ {
		k, v := keyValOf(i)
		key, err := m.Key(i)
		require.NoErrorf(err, "key %d", i)
		require.Equalf(k, key, "key %d", i)
		val, err := m.Value(i)
		require.NoErrorf(err, "value %d", i)
		require.Equalf(v, val, "value %d", i)
	}
}
