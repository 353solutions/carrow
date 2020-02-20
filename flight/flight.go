package flight

/*
#cgo pkg-config: arrow plasma
#cgo LDFLAGS: -lcarrow
#cgo linux LDFLAGS: -L../bindings/linux-x86_64
#cgo CFLAGS: -I..

#include "carrow.h"
#include <stdlib.h>
*/
import "C"

// TODO: United with one in carrow (internal?)
func Start() error {
	return nil
}
