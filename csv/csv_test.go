package csv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	testCases := []struct {
		fileName string
		opts     []ParseOption
	}{
		{"cart.csv", nil},
		{"cart.tsv", []ParseOption{WithDelimiter('\t')}},
	}
	for _, tc := range testCases {
		t.Run(tc.fileName, func(t *testing.T) {
			require := require.New(t)
			file, err := os.Open(tc.fileName)
			require.NoErrorf(err, "open")

			po := NewParseOptions(tc.opts...)
			table, err := Read(file, po)
			require.NoError(err, "read csv")

			require.Equal(4, table.NumCols(), "columns")
			require.Equal(4, table.NumRows(), "rows")
		})
	}
}
