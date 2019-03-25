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
		return nil, fmt.Errorf("can't create field from %s: %s", name, dtype)
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

// FloatArrayBuilder used for building float Arrays
type FloatArrayBuilder struct {
	ptr unsafe.Pointer
}

// NewFloatArrayBuilder returns a new FloatArrayBuilder
func NewFloatArrayBuilder() *FloatArrayBuilder {
	ptr := C.array_builder_new(C.int(FloatType))
	return &FloatArrayBuilder{ptr}
}

// Append appends an integer
func (b *FloatArrayBuilder) Append(val float64) error {
	C.array_builder_append_float(b.ptr, C.double(val))
	return nil
}

// Finish creates the array
func (b *FloatArrayBuilder) Finish() (*Array, error) {
	return builderFinish(b.ptr)
}

// IntArrayBuilder used for building integer Arrays
type IntArrayBuilder struct {
	ptr unsafe.Pointer
}

// NewIntArrayBuilder returns a new IntArrayBuilder
func NewIntArrayBuilder() *IntArrayBuilder {
	ptr := C.array_builder_new(C.int(IntegerType))
	return &IntArrayBuilder{ptr}
}

// Append appends an integer
func (b *IntArrayBuilder) Append(val int) error {
	C.array_builder_append_int(b.ptr, C.longlong(val))
	return nil
}

// Finish creates the array
func (b *IntArrayBuilder) Finish() (*Array, error) {
	return builderFinish(b.ptr)
}

func builderFinish(ptr unsafe.Pointer) (*Array, error) {
	out := C.array_builder_finish(ptr)
	if out.err != nil {
		err := fmt.Errorf(C.GoString(out.err))
		C.free(unsafe.Pointer(out.err))
		return nil, err
	}

	return &Array{out.arr}, nil
}

// Array is arrow array
type Array struct {
	ptr unsafe.Pointer
}
