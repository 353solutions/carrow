package carrow

import (
	"fmt"
	"runtime"
	"time"
	"unsafe"
)

/*
#cgo pkg-config: arrow plasma
#cgo LDFLAGS: -lcarrow -L.

#include <stdlib.h>
#include "carrow.h"
*/
import "C"

//go:generate go run gen.go
//go:generate go fmt carrow_generated.go

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
		C.schema_free(s.ptr)
	})

	return schema, nil
}

type builder struct {
	ptr unsafe.Pointer
}

// Finish returns array from builder
// You can't use the builder after calling Finish
func (b *builder) Finish() (*Array, error) {
	r := New(C.array_builder_finish(b.ptr))
	if err := r.Err(); err != nil {
		return nil, err
	}

	return &Array{r.Ptr()}, nil
}

// Append appends a bool
func (b *BoolArrayBuilder) Append(val bool) error {
	var ival int
	if val {
		ival = 1
	}
	r := New(C.array_builder_append_bool(b.ptr, C.int(ival)))
	return r.Err()
}

// Append appends an integer
func (b *Float64ArrayBuilder) Append(val float64) error {
	r := New(C.array_builder_append_float(b.ptr, C.double(val)))
	return r.Err()
}

// Append appends an integer
func (b *Integer64ArrayBuilder) Append(val int64) error {
	r := C.array_builder_append_int(b.ptr, C.long(val))
	if r.err != nil {
		return nil
	}
	return nil
}

// Append appends a string
func (b *StringArrayBuilder) Append(val string) error {
	cStr := C.CString(val)
	defer C.free(unsafe.Pointer(cStr))
	length := C.ulong(len(val)) // len is in bytes
	r := New(C.array_builder_append_string(b.ptr, cStr, length))
	return r.Err()
}

// Append appends a timestamp
func (b *TimestampArrayBuilder) Append(val time.Time) error {
	r := New(C.array_builder_append_timestamp(b.ptr, C.longlong(val.UnixNano())))
	return r.Err()
}

// Array is arrow array
type Array struct {
	ptr unsafe.Pointer
}

// Length returns the length of the array
func (a *Array) Length() int {
	i := C.array_length(a.ptr)
	return int(i)
}

// BoolAt returns bool at location
func (a *Array) BoolAt(i int) (bool, error) {
	val := C.array_bool_at(a.ptr, C.longlong(i))
	if val == -1 {
		return false, fmt.Errorf("can't get bool at %d", i)
	}

	if val == 0 {
		return false, nil
	}

	return true, nil
}

// Float64At returns float at location
func (a *Array) Float64At(i int) (float64, error) {
	val := C.array_float_at(a.ptr, C.longlong(i))
	return float64(val), nil
}

// Int64At returns integer at location
func (a *Array) Int64At(i int) (int64, error) {
	val := C.array_int_at(a.ptr, C.longlong(i))
	return int64(val), nil
}

// StringAt returns integer at location
func (a *Array) StringAt(i int) (string, error) {
	val := C.array_str_at(a.ptr, C.longlong(i))
	if val == nil {
		return "", fmt.Errorf("can't get string at %d", i)
	}

	s := C.GoString(val)
	C.free(unsafe.Pointer(val))
	return s, nil
}

// TimeAt returns time at location
func (a *Array) TimeAt(i int) (time.Time, error) {
	epochNano := int64(C.array_timestamp_at(a.ptr, C.longlong(i)))
	t := time.Unix(epochNano/1e9, epochNano%1e9)
	return t, nil
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

	/* FIXME
	   if err := table.validate(); err != nil {
	       C.table_free(ptr)
	       return nil, err
	   }
	*/

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
func (t *Table) Ptr() unsafe.Pointer {
	return t.ptr
}
