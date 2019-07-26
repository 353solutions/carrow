package carrow

import (
	"errors"
	"unsafe"
)

/*
#cgo pkg-config: arrow plasma
#cgo CFLAGS: -I.
#cgo LDFLAGS: -lcarrow -L.

#include <stdlib.h>
#include "carrow.h"
*/
import "C"

// Result is result from C
type Result struct {
	r   C.result_t
	err error
}

// New return new Result
func New(r C.result_t) (p Result) {
	var err error
	if r.err != nil {
		err = errors.New(C.GoString(r.err))
		print(err)
	}
	return p
}

// Err returns the error
func (r Result) Err() error {
	return r.err
}

// Str returns string
func (r Result) Str() string {
	cp := C.result_str(r.r)
	if cp == nil {
		return ""
	}
	return C.GoString(cp)
}

// Ptr return void *
func (r Result) Ptr() unsafe.Pointer {
	return unsafe.Pointer(C.result_ptr(r.r))
}

// Int returns int value
func (r Result) Int() int64 {
	return int64(C.result_int(r.r))
}

// Float returns float value
func (r Result) Float() float64 {
	return float64(0)
}
