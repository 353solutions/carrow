// The following directive is necessary to make the package coherent:

// +build ignore

// This program generates contributors.go. It can be invoked by running
// go generate
package main

import (
	"log"
	"os"
	"text/template"
	"time"
)

type data struct {
	Type string
}

func main() {
	arrowTypes := []string{"Bool", "FLOAT64", "INTEGER64", "STRING", "TIMESTAMP"}

	f, err := os.Create("gen/lib_generated.go")
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

var packageTemplate = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at {{ .Timestamp.Format "2006-01-02T15:04:05" }}
package carrow

/*
#cgo pkg-config: arrow plasma
#cgo LDFLAGS: -lcarrow -L.
#cgo CXXFLAGS: -I/src/arrow/cpp/src
// FIXME: plasma headers

#include "../carrow.h"
#include <stdlib.h>
*/
import "C"

// DType is a data type
type DType C.int


// Supported data types
var(
{{- range $val := .ArrowTypes}}
	{{$val}}Type = DType(C.{{$val}}_DTYPE)
{{- end}}
)

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
