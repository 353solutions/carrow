package carrow

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	intColName   = "intCol"
	floatColName = "floatCol"
)

func TestField(t *testing.T) {
	require := require.New(t)
	name, dtype := "field-1", Integer64Type
	field, _ := NewField(name, dtype)
	require.Equal(field.Name(), name, "field name")
	require.Equal(field.DType(), dtype, "field dtype")
}

func TestSchema(t *testing.T) {
	require := require.New(t)
	name, dtype := "field-1", Integer64Type
	field, _ := NewField(name, dtype)
	schema, _ := NewSchema([]*Field{field})
	require.Equal(field.Name(), name, "field name")
	require.Equal(field.DType(), dtype, "field dtype")
	require.NotNil(schema)
}

func TestBoolBuilder(t *testing.T) {
	require := require.New(t)
	b := NewBoolArrayBuilder()
	require.NotNil(b.ptr, "create")
	b.Append(true)
}

func TestFloatBuilder(t *testing.T) {
	require := require.New(t)
	b := NewFloat64ArrayBuilder()
	require.NotNil(b.ptr, "create")
	b.Append(7.2)
}

func TestIntBuilder(t *testing.T) {
	require := require.New(t)
	b := NewInteger64ArrayBuilder()
	require.NotNil(b.ptr, "create")
	b.Append(7)
}

func TestStringBuilder(t *testing.T) {
	require := require.New(t)
	b := NewStringArrayBuilder()
	require.NotNil(b.ptr, "create")
	b.Append("hello")
}

func TestTimestampBuilder(t *testing.T) {
	require := require.New(t)
	b := NewTimestampArrayBuilder()
	require.NotNil(b.ptr, "create")
	b.Append(time.Now())
}

func TestTable(t *testing.T) {
	require := require.New(t)
	nrows := 117

	table := buildTable(require, nrows)
	require.Equal(nrows, table.NumRows(), "rows")
	require.Equal(2, table.NumCols(), "columns")

	names, err := table.ColumnNames()
	require.NoError(err, "ColumnNames")
	require.Equal([]string{intColName, floatColName}, names, "column names")

	arr, err := table.Column(0)
	require.NoError(err, "Column(0)")
	require.Equal(Integer64Type, arr.DType(), "int dtype")

	arr, err = table.Column(1)
	require.NoError(err, "Column(1)")
	require.Equal(Float64Type, arr.DType(), "float dtype")

	arr, err = table.ColumnByName(floatColName)
	require.NoError(err, "col by name")
	require.Equal(Float64Type, arr.DType(), "float dtype")

	offset, length := 10, 37
	s := table.Slice(offset, length)
	names, err = s.ColumnNames()
	require.NoError(err, "ColumnNames")
	require.Equal([]string{intColName, floatColName}, names, "slice column names")
	require.Equal(length, s.NumRows(), "slice rows")
}

func buildTable(require *require.Assertions, nrows int) *Table {
	intBld := NewInteger64ArrayBuilder()
	floatBld := NewFloat64ArrayBuilder()
	for i := 0; i < nrows; i++ {
		err := intBld.Append(int64(i))
		require.NoErrorf(err, "append int %d", i)
		err = floatBld.Append(float64(i))
		require.NoErrorf(err, "append floatj %d", i)
	}

	intArr, err := intBld.Finish()
	require.NoError(err, "build int")
	floatArr, err := floatBld.Finish()
	require.NoError(err, "build float")

	intField, err := NewField(intColName, Integer64Type)
	require.NoError(err, "int field")
	floatField, err := NewField(floatColName, Float64Type)
	require.NoError(err, "float field")
	schema, err := NewSchema([]*Field{intField, floatField})
	require.NoError(err, "schema")
	arrs := []*Array{intArr, floatArr}
	table, err := NewTableFromArrays(schema, arrs)
	require.NoError(err, "build table")
	return table
}
