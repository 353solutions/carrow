package carrow

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestField(t *testing.T) {
	require := require.New(t)
	name, dtype := "field-1", IntegerType
	field := NewField(name, dtype)
	require.Equal(field.Name(), name, "field name")
	require.Equal(field.DType(), dtype, "field dtype")
}
