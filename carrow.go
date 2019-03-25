package carrow

/*
#cgo pkg-config: arrow
#cgo LDFLAGS: -lcarrow -L.

#include "carrow.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

// DType is a data type
type DType C.int

// Supported data types
var (
	IntegerType = DType(C.INTEGER_DTYPE)
	FloatType   = DType(C.FLOAT_DTYPE)
)

func (dt DType) String() string {
	switch dt {
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
func NewField(name string, dtype DType) (*Field, error) {
	cName := C.CString(name)
	defer func() { C.free(unsafe.Pointer(cName)) }()

	ptr := C.field_new(cName, C.int(dtype))
	if ptr == nil {
		return nil, fmt.Errorf("can't create field from %s:s", name, dtype)
	}

	field := &Field{ptr}
	runtime.SetFinalizer(field, func(f *Field) {
		C.field_free(f.ptr)
	})
	return field, nil
}

// Name returns the field name
func (f *Field) Name() string {
	return C.GoString(C.field_name(f.ptr))
}

// DType returns the field data type
func (f *Field) DType() DType {
	return DType(C.field_dtype(f.ptr))
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
	runtime.SetFinalizer(schema, func(s *Schema) {
		C.schema_free(schema.ptr)
	})

	return schema, nil
}

// Array of data
type Array struct {
	ptr unsafe.Pointer
}

// DType returns the array DType
func (a *Array) DType() DType {
	return 0 // FIXME
}

// Len is the length of the array
func (a *Array) Len() int {
	return 0 // FIXME
}
