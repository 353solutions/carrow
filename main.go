package main

import (
	"fmt"
)

// #cgo LDFLAGS: -larrow -lcarrow -L.
// #cgo CXXFLAGS: -I../arrow/cpp/src
// #include "carrow.h"
import "C"

func main() {
	field := C.field_new(C.CString("HI"), C.INTEGER)
	fmt.Println(field)
	name := C.GoString(C.field_name(field))
	fmt.Println(name)
	dtype := C.field_dtype(field)
	fmt.Println(dtype)
	C.field_free(field)
}
