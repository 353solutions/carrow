package main

import "C"
import (
	"fmt"
	"unsafe"

	"github.com/353solutions/carrow"
)

//export CreateTable
func CreateTable() unsafe.Pointer {
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
		return nil
	}

	floatArr, err := floatBld.Finish()
	if err != nil {
		fmt.Printf("floatBld error: %s", err)
		return nil
	}

	intField, err := carrow.NewField("incCol", carrow.Integer64Type)
	if err != nil {
		fmt.Printf("intField error: %s", err)
		return nil
	}

	floatField, err := carrow.NewField("floatCol", carrow.Float64Type)
	if err != nil {
		fmt.Printf("floatField error: %s", err)
		return nil
	}

	intCol, err := carrow.NewColumn(intField, intArr)
	if err != nil {
		fmt.Printf("intCol error: %s", err)
		return nil
	}

	floatCol, err := carrow.NewColumn(floatField, floatArr)
	if err != nil {
		fmt.Printf("floatCol error: %s", err)
		return nil
	}

	cols := []*carrow.Column{intCol, floatCol}
	table, err := carrow.NewTableFromColumns(cols)
	if err != nil {
		fmt.Printf("table creation error: %s", err)
		return nil
	}

	fmt.Printf("num cols: %d\n", table.NumCols())
	fmt.Printf("num rows: %d\n", table.NumRows())

	return table.Pointer()
}

func main() {}
