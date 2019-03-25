// +build ignore
// Tentative API
package carrow_test

import (
	"fmt"

	"github.com/353solutions/carrow"
)

func Example() {
	size := 100
	intBld := carrow.NewIntBuilder()
	floatBld := carrow.NewFloatBuilder()
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

	intField := carrow.NewField("incCol", carrow.IntegerType)
	floatField := carrow.NewField("floatCol", carrow.FloatType)
	intCol, err := carrow.NewColum(intField, intArr)
	if err != nil {
		fmt.Printf("intCol error: %s", err)
		return
	}

	floatCol, err := carrow.NewColum(floatField, floatArr)
	if err != nil {
		fmt.Printf("floatCol error: %s", err)
		return
	}

	tbl := carrow.NewTable(nil) // Can pass []*carrow.Column
	tbl.Append(intCol)
	tbl.Append(floatCol)

	if err := tbl.Validate(); err != nil {
		fmt.Printf("table validate error: %s", err)
		return
	}
}
