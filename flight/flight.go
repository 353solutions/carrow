package flight

/*
#cgo pkg-config: arrow plasma
#cgo LDFLAGS: -lcarrow -larrow_flight -lstdc++
#cgo linux LDFLAGS: -L../bindings/linux-x86_64
#cgo CFLAGS: -I..

#include "carrow.h"
#include <stdlib.h>
*/
import "C"

func Start() error {
	C.flight_server_start()
	return nil
}
