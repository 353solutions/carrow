// Tentative API
package carrow_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/353solutions/carrow"
)

func TestExample(t *testing.T) {
	require := require.New(t)
	size := 100
	intBld := carrow.NewIntArrayBuilder()
	floatBld := carrow.NewFloatArrayBuilder()
	for i := 0; i < size; i++ {
		intBld.Append(i)
		floatBld.Append(float64(i))
	}

	intArr, err := intBld.Finish()
	if err != nil {
		require.FailNow("intBld error: %s", err)
	}

	floatArr, err := floatBld.Finish()
	if err != nil {
		require.FailNow("floatBld error: %s", err)
	}

	intField, err := carrow.NewField("incCol", carrow.IntegerType)
	if err != nil {
		require.FailNow("intField error: %s", err)
	}

	floatField, err := carrow.NewField("floatCol", carrow.FloatType)
	if err != nil {
		require.FailNow("floatField error: %s", err)
	}

	intCol, err := carrow.NewColumn(intField, intArr)
	if err != nil {
		require.FailNow("intCol error: %s", err)
	}

	floatCol, err := carrow.NewColumn(floatField, floatArr)
	if err != nil {
		require.FailNow("floatCol error: %s", err)
	}

	cols := []*carrow.Column{intCol, floatCol}
	table, err := carrow.NewTableFromColumns(cols)
	if err != nil {
		require.FailNow("table creation error: %s", err)
	}

	require.Equal(2,table.NumCols(),"number of cols in table")
	require.Equal(100,table.NumRows(),"number of rows in table")
}
