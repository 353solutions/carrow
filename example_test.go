// Tentative API
package carrow_test

import (
	"fmt"

	"github.com/353solutions/carrow"
)

func Example() {
	size := 100
	intBld := carrow.NewIntArrayBuilder()
	floatBld := carrow.NewFloatArrayBuilder()
	for i := 0; i < size; i++ {
		intBld.Append(i)
		floatBld.Append(float64(i))
	}

	intArr, err := intBld.Finish()
	if err != nil {
		fmt.Printf("intBld error: %s", err)
		return
	}

	floatArr, err := floatBld.Finish()
	if err != nil {
		fmt.Printf("floatBld error: %s", err)
		return
	}

	intField, err := carrow.NewField("incCol", carrow.IntegerType)
	if err != nil {
		fmt.Printf("intField error: %s", err)
		return
	}

	floatField, err := carrow.NewField("floatCol", carrow.FloatType)
	if err != nil {
		fmt.Printf("floatField error: %s", err)
		return
	}

	intCol, err := carrow.NewColumn(intField, intArr)
	if err != nil {
		fmt.Printf("intCol error: %s", err)
		return
	}

	floatCol, err := carrow.NewColumn(floatField, floatArr)
	if err != nil {
		fmt.Printf("floatCol error: %s", err)
		return
	}

	cols := []*carrow.Column{intCol, floatCol}
	table, err := carrow.NewTableFromColumns(cols)
	if err != nil {
		fmt.Printf("table creation error: %s", err)
		return
	}

	fmt.Printf("num cols: %d\n", table.NumCols())
	// Output: num cols: 2

	fmt.Printf("num rows: %d\n", table.NumRows())
	// Output: num rows: 100
}
