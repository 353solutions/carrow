package carrow

import (
	"fmt"
	"runtime"
	"time"
	"unsafe"
)

/*
#cgo pkg-config: arrow plasma
#cgo LDFLAGS: -lcarrow -larrow_flight
#cgo linux LDFLAGS: -L./bindings/linux-x86_64
#cgo CXXFLAGS: -I/src/arrow/cpp/src
// FIXME: plasma headers

#include "carrow.h"
#include <stdlib.h>
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
func NewSchema(fields []*Field) (*Schema, error) {
	arr := make([]unsafe.Pointer, 0, len(fields))
	for _, fld := range fields {
		arr = append(arr, fld.ptr)
	}
	cf := (unsafe.Pointer)(&arr[0])
	count := len(fields)
	ptr := C.schema_new(cf, C.size_t(count))
	if ptr == nil {
		return nil, fmt.Errorf("can't create schema")
	}
	schema := &Schema{ptr}
	runtime.SetFinalizer(schema, func(s *Schema) {
		C.schema_free(s.ptr)
	})

	return schema, nil
}

// Metadata returns the schema metadata
func (s *Schema) Metadata() (*Metadata, error) {
	r := C.schema_meta(s.ptr)
	if err := errFromResult(r); err != nil {
		return nil, err
	}

	return &Metadata{r.ptr}, nil
}

// SetMetadata sets the metadata
func (s *Schema) SetMetadata(m *Metadata) error {
	r := C.schema_set_meta(s.ptr, m.ptr)
	if err := errFromResult(r); err != nil {
		return err
	}
	s.ptr = r.ptr
	return nil
}

type flusher interface {
	flush() error
}

type builder struct {
	ptr unsafe.Pointer
	fl  flusher
}

func errFromResult(r C.result_t) error {
	if r.err == nil {
		return nil
	}
	err := fmt.Errorf(C.GoString(r.err))
	C.free(unsafe.Pointer(r.err))
	return err
}

// Finish returns array from builder
// You can't use the builder after calling Finish
func (b *builder) Finish() (*Array, error) {
	if err := b.fl.flush(); err != nil {
		return nil, err
	}

	r := C.array_builder_finish(b.ptr)
	if err := errFromResult(r); err != nil {
		return nil, err
	}

	return &Array{r.ptr}, nil
}

// TODO: Templetaize Append & flush
// Append appends a bool
func (b *BoolArrayBuilder) Append(val bool) error {
	var ival C.uint8_t = 0
	if val {
		ival = 1
	}
	b.buffer[b.bufferIdx] = ival
	b.bufferIdx++
	if b.bufferIdx < bufferSize {
		return nil
	}
	return b.flush()
}

func (b *BoolArrayBuilder) flush() error {
	cSize := C.long(b.bufferIdx)
	b.bufferIdx = 0
	r := C.array_builder_append_bools(b.ptr, (*C.uint8_t)(&b.buffer[0]), cSize)
	return errFromResult(r)
}

// Append appends an float
func (b *Float64ArrayBuilder) Append(val float64) error {
	b.buffer[b.bufferIdx] = C.double(val)
	b.bufferIdx++
	if b.bufferIdx < bufferSize {
		return nil
	}

	return b.flush()
}

func (b *Float64ArrayBuilder) flush() error {
	cSize := C.long(b.bufferIdx)
	b.bufferIdx = 0
	r := C.array_builder_append_floats(b.ptr, (*C.double)(&b.buffer[0]), cSize)
	return errFromResult(r)
}

// Append appends an integer
func (b *Integer64ArrayBuilder) Append(val int64) error {
	b.buffer[b.bufferIdx] = C.long(val)
	b.bufferIdx++
	if b.bufferIdx < bufferSize {
		return nil
	}

	return b.flush()
}

func (b *Integer64ArrayBuilder) flush() error {
	cSize := C.long(b.bufferIdx)
	b.bufferIdx = 0
	r := C.array_builder_append_ints(b.ptr, (*C.long)(&b.buffer[0]), cSize)
	return errFromResult(r)
}

// Finish creates the array from the builder
func (b *Integer64ArrayBuilder) Finish() (*Array, error) {
	if b.bufferIdx > 0 {
		if err := b.flush(); err != nil {
			return nil, err
		}
	}

	return b.builder.Finish()
}

// Append appends a string
func (b *StringArrayBuilder) Append(val string) error {
	b.buffer[b.bufferIdx] = C.CString(val)
	b.bufferIdx++
	if b.bufferIdx < bufferSize {
		return nil
	}

	return b.flush()
}

func (b *StringArrayBuilder) flush() error {
	cSize := C.long(b.bufferIdx)
	b.bufferIdx = 0
	r := C.array_builder_append_strings(b.ptr, (**C.char)(&b.buffer[0]), cSize)
	for _, cp := range b.buffer {
		C.free(unsafe.Pointer(cp))
	}
	return errFromResult(r)
}

// Append appends a timestamp
func (b *TimestampArrayBuilder) Append(val time.Time) error {
	b.buffer[b.bufferIdx] = C.long(val.UnixNano())
	b.bufferIdx++
	if b.bufferIdx < bufferSize {
		return nil
	}

	return b.flush()
}

func (b *TimestampArrayBuilder) flush() error {
	cSize := C.long(b.bufferIdx)
	b.bufferIdx = 0
	r := C.array_builder_append_timestamps(b.ptr, (*C.long)(&b.buffer[0]), cSize)
	return errFromResult(r)
}

// Array is arrow array
type Array struct {
	ptr unsafe.Pointer
}

// DType returns the array data type
func (a *Array) DType() DType {
	return DType(C.array_dtype(a.ptr))
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

// Table is arrow table
type Table struct {
	ptr unsafe.Pointer
}

// NewTableFromArrays creates new Table from slice of arrays
func NewTableFromArrays(schema *Schema, arrays []*Array) (*Table, error) {
	arrs := make([]unsafe.Pointer, 0, len(arrays))
	for _, arr := range arrays {
		arrs = append(arrs, arr.ptr)
	}
	aptr := (unsafe.Pointer)(&arrs[0])
	ncols := len(arrays)
	ptr := C.table_new(schema.ptr, aptr, C.size_t(ncols))
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

// Schema returns the table Schema
func (t *Table) Schema() *Schema {
	ptr := C.table_schema(t.ptr)
	if ptr == nil {
		return nil
	}

	return &Schema{ptr}
}

// Column returns the nth column (Array)
func (t *Table) Column(i int) (*Array, error) {
	ptr := C.table_column(t.ptr, C.int(i))
	if ptr == nil {
		return nil, fmt.Errorf("can't find column %d", i)
	}

	return &Array{ptr}, nil
}

// ColumnByName returns column by name
func (t *Table) ColumnByName(name string) (*Array, error) {
	for i := 0; i < t.NumCols(); i++ {
		fld, err := t.Field(i)
		if err != nil {
			return nil, err
		}
		if fld.Name() == name {
			return t.Column(i)
		}
	}

	return nil, fmt.Errorf("column %q not found", name)
}

// ColumnNames names returns names of columns
func (t *Table) ColumnNames() ([]string, error) {
	ncols := t.NumCols()
	names := make([]string, 0, ncols)
	for i := 0; i < ncols; i++ {
		fld, err := t.Field(i)
		if err != nil {
			return nil, err
		}
		names = append(names, fld.Name())
	}
	return names, nil
}

// Slice returns a 0 copy slize of t
// If length is -1 will return until end of table
func (t *Table) Slice(offset int, length int) *Table {
	if length == -1 {
		length = t.NumRows()
	}

	ptr := C.table_slice(t.ptr, C.int64_t(offset), C.int64_t(length))
	return &Table{ptr}
}

// Field returns the nth field
func (t *Table) Field(i int) (*Field, error) {
	ptr := C.table_field(t.ptr, C.int(i))
	if ptr == nil {
		return nil, fmt.Errorf("can't find field %d", i)
	}

	return &Field{ptr}, nil
}

// Ptr returns the underlying C++ pointer
func (t *Table) Ptr() unsafe.Pointer {
	return t.ptr
}

// Metadata in schema
type Metadata struct {
	ptr unsafe.Pointer
}

// NewMetadata creates new Metadata
func NewMetadata() *Metadata {
	return &Metadata{C.meta_new()}
}

// Set sets a key/value
func (m *Metadata) Set(key, value string) error {
	cKey, cVal := C.CString(key), C.CString(value)
	defer C.free(unsafe.Pointer(cKey))
	defer C.free(unsafe.Pointer(cVal))

	r := C.meta_set(m.ptr, cKey, cVal)
	return errFromResult(r)
}

// Len returns number of elements
func (m *Metadata) Len() (int, error) {
	r := C.meta_size(m.ptr)
	if err := errFromResult(r); err != nil {
		return 0, err
	}

	return int(r.i), nil
}

// Key returns key at index i
func (m *Metadata) Key(i int) (string, error) {
	r := C.meta_key(m.ptr, C.long(i))
	if err := errFromResult(r); err != nil {
		return "", err
	}

	key := C.GoString((*C.char)(r.ptr))
	C.free(r.ptr)
	return key, nil
}

// Value returns value at index i
func (m *Metadata) Value(i int) (string, error) {
	r := C.meta_value(m.ptr, C.long(i))
	if err := errFromResult(r); err != nil {
		return "", err
	}

	key := C.GoString((*C.char)(r.ptr))
	C.free(r.ptr)
	return key, nil
}
