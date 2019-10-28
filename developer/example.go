package main

import (
	"fmt"

	"github.com/353solutions/carrow"
)

func main() {
	name, dtype := "field-1", carrow.Integer64Type
	field, _ := carrow.NewField(name, dtype)
	fmt.Printf("%v\n", field)
}
