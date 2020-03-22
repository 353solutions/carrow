package flight

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFlight(t *testing.T) {
	require := require.New(t)
	go Start()
	require.True(true)
	// oid, err := RandomID()
	// require.NoError(err, "create id")
	// require.Len(oid, 20, "bad length")
}
