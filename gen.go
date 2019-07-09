// The following directive is necessary to make the package coherent:

// +build ignore

// This program generates contributors.go. It can be invoked by running
// go generate
package main

import (
	"log"
	"os"
	"text/template"
	"strings"
	"time"
)

type data struct {
	Type string
}

func main() {
	arrowTypes := []string{"Bool", "Float64", "Integer64", "String", "Timestamp"}
	f, err := os.Create("lib_generated.go")
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

var funcMap = template.FuncMap{
	"ToUpper": strings.ToUpper,
}
var packageTemplate = template.Must(template.New("").Funcs(funcMap).Parse(`
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
	}

	// New{{$val}}ArrayBuilder returns a new {{$val}}ArrayBuilder
	func New{{$val}}ArrayBuilder() *{{$val}}ArrayBuilder {
		r := C.array_builder_new(C.int({{$val}}Type))
		if r.err != nil {
			return nil
		}
		return &{{$val}}ArrayBuilder{builder{r.ptr}}
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
