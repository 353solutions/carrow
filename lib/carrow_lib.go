package main

import (
	"C"

	"github.com/353solutions/carrow"
)

//export Mock
func Mock() {
	size := 100
	intBld := carrow.NewInt64ArrayBuilder()
	floatBld := carrow.NewFloat64ArrayBuilder()
	for i := 0; i < size; i++ {
		intBld.Append(int64(i))
		floatBld.Append(float64(i))
	}

}

func main() {}
