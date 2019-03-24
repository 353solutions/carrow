package carrow

// #cgo pkg-config: arrow
// #cgo LDFLAGS: -lcarrow -L.
/* 
	#include "carrow.h"
   	#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// DType is a data type
type DType C.int

// Supported data types
	var IntegerType C.int = C.NC_INTEGER_DTYPE
	var FloatType   C.int = C.NC_FLOAT_DTYPE

func (t C.int) String() string {
	switch t {
	case IntegerType:
		return "int64"
	case FloatType:
		return "float64"
	}

	return "<unknown>"
}

// Field is a field description
type Field struct {
	ptr unsafe.Pointer
}

// NewField returns a new Field
func NewField(name string, dtype C.int) (*Field, error) {
	cName := C.CString(name)
	defer func() { C.free(unsafe.Pointer(cName)) }()

	ptr := C.field_new(cName, dtype)
	if ptr == nil {
		//return nil, fmt.Errorf("can't create field from %s:s", name, dtype)
		return nil, fmt.Errorf("can't create field")

	}

	field := Field{ptr}
	//runtime.SetFinalizer(field, func() {
	//	//C.field_free(ptr)
	//})
	return &field,nil
}

// Name returns the field name
func (f *Field) Name() string {
	return C.GoString(C.field_name(f.ptr))
}

// DType returns the field data type
func (f *Field) DType() C.int {
	return C.field_dtype(f.ptr)
}

// Schema is table schema
type Schema struct {
	ptr unsafe.Pointer
}

// NewSchema creates a new schema
func NewSchema(fields []Field) (*Schema, error) {
	ptr := C.schema_new()
	if ptr == nil {
		return nil, fmt.Errorf("can't create schema")
	}
	for _, field := range fields {
		C.schema_add_field(ptr, field.ptr)
	}

	schema := &Schema{ptr}
	//runtime.SetFinalizer(schema, func() {
	//	C.schema_free(schema.ptr)
	//})

	return schema, nil
}
