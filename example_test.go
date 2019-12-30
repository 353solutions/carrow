// Tentative API
package carrow_test

import (
	"fmt"

	"github.com/353solutions/carrow"
)

func Example() {
	size := 100
	intBld := carrow.NewInteger64ArrayBuilder()
	floatBld := carrow.NewFloat64ArrayBuilder()
	for i := 0; i < size; i++ {
		if err := intBld.Append(int64(i)); err != nil {
			fmt.Printf("intBld.Append error: %s", err)
			return
		}
		if err := floatBld.Append(float64(i)); err != nil {
			fmt.Printf("floatBld.Append error: %s", err)
			return
		}
	}

	intArr, err := intBld.Finish()
	if err != nil {
		fmt.Printf("intBld.Finish error: %s", err)
		return
	}

	floatArr, err := floatBld.Finish()
	if err != nil {
		fmt.Printf("floatBld.Finish error: %s", err)
		return
	}

	intField, err := carrow.NewField("incCol", carrow.Integer64Type)
	if err != nil {
		fmt.Printf("intField error: %s", err)
		return
	}

	floatField, err := carrow.NewField("floatCol", carrow.Float64Type)
	if err != nil {
		fmt.Printf("floatField error: %s", err)
		return
	}

	schema, err := carrow.NewSchema([]*carrow.Field{intField, floatField})
	if err != nil {
		fmt.Printf("can't create schema: %s", err)
		return
	}
	arrs := []*carrow.Array{intArr, floatArr}

	table, err := carrow.NewTableFromArrays(schema, arrs)
	if err != nil {
		fmt.Printf("table creation error: %s", err)
		return
	}

	fmt.Printf("num cols: %d\n", table.NumCols())
	fmt.Printf("num rows: %d\n", table.NumRows())

	// Output:
	// num cols: 2
	// num rows: 100
}
