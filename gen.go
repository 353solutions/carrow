// The following directive is necessary to make the package coherent:

// +build ignore

// This program generates contributors.go. It can be invoked by running
// go generate
package main

import (
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

func main() {
	arrowTypes := []string{"Bool", "Float64", "Integer64", "String", "Timestamp"}
	f, err := os.Create("carrow_generated.go")
	die(err)
	defer f.Close()

	packageTemplate.Execute(f, struct {
		Timestamp  time.Time
		ArrowTypes []string
	}{
		Timestamp:  time.Now(),
		ArrowTypes: arrowTypes,
	})
}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CType(name string) string {
	switch name {
	case "Bool":
		return "C.int"
	case "Float64":
		return "C.double"
	case "Integer64":
		return "C.long"
	case "String":
		return "*C.char"
	case "Timestamp":
		return "C.longlong"
	}

	panic(name)
}

var funcMap = template.FuncMap{
	"ToUpper": strings.ToUpper,
	"CType":   CType,
}

var packageTemplate = template.Must(template.New("").Funcs(funcMap).Parse(`
// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// {{ .Timestamp }}
package carrow

/*
#cgo pkg-config: arrow plasma
#cgo LDFLAGS: -lcarrow -L.
#cgo CXXFLAGS: -I/src/arrow/cpp/src
// FIXME: plasma headers

#include "carrow.h"
#include <stdlib.h>
*/
import "C"

// DType is a data type
type DType C.int

const bufferSize = 1024 * 4


// Supported data types
var(
{{- range $val := .ArrowTypes}}
	{{$val}}Type = DType(C.{{$val | ToUpper }}_DTYPE)
{{- end}}
)

// Array Builders
{{- range $val := .ArrowTypes}}

	type {{$val}}ArrayBuilder struct {
		builder
		buffer [bufferSize]{{$val | CType}}
		bufferIdx int
	}

	// New{{$val}}ArrayBuilder returns a new {{$val}}ArrayBuilder
	func New{{$val}}ArrayBuilder() *{{$val}}ArrayBuilder {
		r := C.array_builder_new(C.int({{$val}}Type))
		if r.err != nil {
			return nil
		}
		return &{{$val}}ArrayBuilder{builder: builder{r.ptr}}
	}
{{- end}}

func (dt DType) String() string {
	switch dt {
{{- range $val := .ArrowTypes}}
	case {{$val}}Type:
		return "{{$val}}"
{{- end}}
	}

	return "<unknown>"
}

`))
