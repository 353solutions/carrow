package carrow

import (
	"fmt"
	"runtime"
	"time"
	"unsafe"
)

/*
#cgo pkg-config: arrow plasma
#cgo LDFLAGS: -lcarrow -L./_lib

#include "./_lib/carrow.h"
#include <stdlib.h>
*/
import "C"

// DType is a data type
type DType C.int

// Supported data types
var (
	BoolType      = DType(C.BOOL_DTYPE)
	Float64Type   = DType(C.FLOAT64_DTYPE)
	Integer64Type = DType(C.INTEGER64_DTYPE)
	StringType    = DType(C.STRING_DTYPE)
	TimestampType = DType(C.TIMESTAMP_DTYPE)
)

func (dt DType) String() string {
	switch dt {
	case BoolType:
		return "bool"
	case Float64Type:
		return "float64"
	case Integer64Type:
		return "int64"
	case StringType:
		return "string"
	case TimestampType:
		return "timestamp"
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

// FieldList is a warpper around std::shared_ptr<arrow::Field>
type FieldList struct {
	ptr unsafe.Pointer
}

// NewFieldList returns a new Field List
func NewFieldList() (*FieldList, error) {

	ptr := C.fields_new()

	if ptr == nil {
		return nil, fmt.Errorf("can't create fields list")
	}

	fieldList := &FieldList{ptr}

	runtime.SetFinalizer(fieldList, func(f *FieldList) {
		C.field_free(f.ptr)
	})
	return fieldList, nil
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
func NewSchema(fields []*Field) (*Schema, error) {
	fieldsList, err := NewFieldList()
	if err != nil {
		return nil, fmt.Errorf("can't create schema,failed creating fields list")
	}
	cf := fieldsList.ptr

	for _, f := range fields {
		C.fields_append(cf, f.ptr)
	}
	ptr := C.schema_new(cf)
	if ptr == nil {
		return nil, fmt.Errorf("can't create schema")
	}
	schema := &Schema{ptr}
	runtime.SetFinalizer(schema, func(s *Schema) {
		C.schema_free(schema.ptr)
	})

	return schema, nil
}

type builder struct {
	ptr unsafe.Pointer
}

// BoolArrayBuilder used for building bool Arrays
type BoolArrayBuilder struct {
	builder
}

// NewBoolArrayBuilder returns a new BoolArrayBuilder
func NewBoolArrayBuilder() *BoolArrayBuilder {
	ptr := C.array_builder_new(C.int(BoolType))
	return &BoolArrayBuilder{builder{ptr}}
}

func (b *builder) Finish() (*Array, error) {
	out := C.array_builder_finish(b.ptr)
	if out.err != nil {
		err := fmt.Errorf(C.GoString(out.err))
		C.free(unsafe.Pointer(out.err))
		return nil, err
	}

	return &Array{out.obj}, nil
}

// Append appends a bool
func (b *BoolArrayBuilder) Append(val bool) error {
	var ival int
	if val {
		ival = 1
	}
	C.array_builder_append_bool(b.ptr, C.int(ival))
	return nil
}

// Float64ArrayBuilder used for building float Arrays
type Float64ArrayBuilder struct {
	builder
}

// NewFloat64ArrayBuilder returns a new Float64ArrayBuilder
func NewFloat64ArrayBuilder() *Float64ArrayBuilder {
	ptr := C.array_builder_new(C.int(Float64Type))
	return &Float64ArrayBuilder{builder{ptr}}
}

// Append appends an integer
func (b *Float64ArrayBuilder) Append(val float64) error {
	C.array_builder_append_float(b.ptr, C.double(val))
	return nil
}

// Int64ArrayBuilder used for building integer Arrays
type Int64ArrayBuilder struct {
	builder
}

// NewInt64ArrayBuilder returns a new Int64ArrayBuilder
func NewInt64ArrayBuilder() *Int64ArrayBuilder {
	ptr := C.array_builder_new(C.int(Integer64Type))
	return &Int64ArrayBuilder{builder{ptr}}
}

// Append appends an integer
func (b *Int64ArrayBuilder) Append(val int64) error {
	C.array_builder_append_int(b.ptr, C.longlong(val))
	return nil
}

// StringArrayBuilder used for building string Arrays
type StringArrayBuilder struct {
	builder
}

// NewStringArrayBuilder returns a new StringArrayBuilder
func NewStringArrayBuilder() *StringArrayBuilder {
	ptr := C.array_builder_new(C.int(StringType))
	return &StringArrayBuilder{builder{ptr}}
}

// Append appends a string
func (b *StringArrayBuilder) Append(val string) error {
	cStr := C.CString(val)
	length := C.ulong(len(val)) // len is in bytes
	C.array_builder_append_string(b.ptr, cStr, length)
	C.free(unsafe.Pointer(cStr))
	return nil
}

// TimestampArrayBuilder used for building bool Arrays
type TimestampArrayBuilder struct {
	builder
}

// NewTimestampArrayBuilder returns a new TimestampArrayBuilder
func NewTimestampArrayBuilder() *TimestampArrayBuilder {
	ptr := C.array_builder_new(C.int(TimestampType))
	return &TimestampArrayBuilder{builder{ptr}}
}

// Append appends a timestamp
func (b *TimestampArrayBuilder) Append(val time.Time) error {
	C.array_builder_append_timestamp(b.ptr, C.longlong(val.UnixNano()))
	return nil
}

// Array is arrow array
type Array struct {
	ptr unsafe.Pointer
}

// Column is an arrow colum
type Column struct {
	ptr unsafe.Pointer
}

// DType returns the Column data type
func (c *Column) DType() DType {
	return DType(C.column_dtype(c.ptr))
}

// NewColumn returns a new column
func NewColumn(field *Field, arr *Array) (*Column, error) {
	if field == nil || arr == nil {
		return nil, fmt.Errorf("nil pointer")
	}

	ptr := C.column_new(field.ptr, arr.ptr)
	c := &Column{ptr}
	if c.DType() != field.DType() {
		return nil, fmt.Errorf("column type doesn't match Field type")
	}

	return c, nil
}

// Field returns the column field
func (c *Column) Field() *Field {
	ptr := C.column_field(c.ptr)
	return &Field{ptr}
}

// Table is arrow table
type Table struct {
	ptr unsafe.Pointer
}

// NewTableFromColumns creates new Table from slice of columns
func NewTableFromColumns(columns []*Column) (*Table, error) {
	fields := make([]*Field, len(columns))
	cptr := C.columns_new()
	defer func() {
		// FIXME
		// C.columns_free(cptr)
	}()

	for i, col := range columns {
		fields[i] = col.Field()
		C.columns_append(cptr, col.ptr)
	}

	schema, err := NewSchema(fields)
	if err != nil {
		return nil, err
	}
	ptr := C.table_new(schema.ptr, cptr)
	table := &Table{ptr}

	if err := table.validate(); err != nil {
		// FIXME
		// C.table_free(ptr)
		return nil, err
	}

	return table, nil
}

// NewTableFromPtr creates a new table from underlying C pointer
// You probably shouldn't use this function
func NewTableFromPtr(ptr unsafe.Pointer) *Table {
	return &Table{ptr}
}

// NumRows returns the number of rows
func (t *Table) NumRows() int {
	return int(C.table_num_rows(t.ptr))
}

// NumCols returns the number of columns
func (t *Table) NumCols() int {
	return int(C.table_num_cols(t.ptr))
}

// Ptr returns the underlying C++ pointer
//  You probably shouldn't use this
func (t *Table) Ptr() unsafe.Pointer {
	return t.ptr
}

func (t *Table) validate() error {
	cp := C.table_validate(t.ptr)
	if cp == nil {
		return nil
	}

	err := fmt.Errorf(C.GoString(cp))
	C.free(unsafe.Pointer(cp))
	return err
}
