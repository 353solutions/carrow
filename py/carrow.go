package main

import (
	"log"
	"unsafe"

	"github.com/353solutions/carrow"
)

/*
#cgo pkg-config: arrow

#include <arrow/python/pyarrow.h>
*/
import "C"

//export build
func build() unsafe.Pointer {
	size := 100
	intBld := carrow.NewIntArrayBuilder()
	floatBld := carrow.NewFloatArrayBuilder()
	for i := 0; i < size; i++ {
		intBld.Append(i)
		floatBld.Append(float64(i))
	}

	intArr, err := intBld.Finish()
	if err != nil {
		log.Printf("intBld error: %s", err)
		return nil
	}

	floatArr, err := floatBld.Finish()
	if err != nil {
		log.Printf("floatBld error: %s", err)
		return nil
	}

	intField, err := carrow.NewField("incCol", carrow.IntegerType)
	if err != nil {
		log.Printf("intField error: %s", err)
		return nil
	}

	floatField, err := carrow.NewField("floatCol", carrow.FloatType)
	if err != nil {
		log.Printf("floatField error: %s", err)
		return nil
	}

	intCol, err := carrow.NewColumn(intField, intArr)
	if err != nil {
		log.Printf("intCol error: %s", err)
		return nil
	}

	floatCol, err := carrow.NewColumn(floatField, floatArr)
	if err != nil {
		log.Printf("floatCol error: %s", err)
		return nil
	}

	cols := []*carrow.Column{intCol, floatCol}
	table, err := carrow.NewTableFromColumns(cols)
	if err != nil {
		log.Printf("table creation error: %s", err)
		return nil
	}

	return table.Ptr()
}

func main() {}
