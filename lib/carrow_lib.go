package main

import (
	"C"
	"fmt"

	"github.com/353solutions/carrow"
)

//export Build
func Build() {
	size := 100
	intBld := carrow.NewInt64ArrayBuilder()
	floatBld := carrow.NewFloat64ArrayBuilder()
	for i := 0; i < size; i++ {
		intBld.Append(int64(i))
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
	fmt.Printf("num rows: %d\n", table.NumRows())
}

func main() {}
