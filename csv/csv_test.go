package csv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	require := require.New(t)
	file, err := os.Open("cart.csv")
	require.NoError(err, "open cart.csv")

	table, err := Read(file)
	require.NoError(err, "read csv")

	require.Equal(4, table.NumCols(), "columns")
	require.Equal(4, table.NumRows(), "rows")
}
