// TODO: Should this be in internal?

package result

/*
#cgo pkg-config: arrow plasma
#cgo LDFLAGS: -lcarrow -L..
#cgo CFLAGS: -I..

#include <stdlib.h>
#include "carrow.h"
*/
import "C"

// Result is result from C
type Result struct {
	r C.result_t
}

// New return new Result
func New(r C.result_t) Result {
	return Result{r}
}

// Err returns the error
func (r Result) Err() error {
	if r.r.err == nil {
		return nil
	}

	return fmt.Errorf(C.GoString(r.r.err))
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
func (r Result) Int() int {
	return int(C.result_int(r.r))
}

// Float returns float value
func (r Result) Float() float64 {
	return float64(c.result_float(r.r))
}
