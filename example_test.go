// +build ignore
// Tentative API
package carrow_test

import (
	"fmt"

	"github.com/353solutions/carrow"
)

func Example() {
	size := 100
	intField := carrow.NewField("incCol", carrow.IntegerType)
	intBld := carrow.NewArrayBuilder(intField.DType())
	floatField := carrow.NewField("floatCol", carrow.FloatType)
	floatBld := carrow.NewArrayBuilder(floatField.DType())
	for i := 0; i < size; i++ {
		intBld.AppendInt(i)
		floatBld.AppendFloat(float64(i))
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
