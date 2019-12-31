package csv

import (
	"fmt"
	"io"
	"unsafe"

	"github.com/353solutions/carrow"
)

/*
#cgo pkg-config: arrow plasma

#include "csv.h"
*/
import "C"

var (
	reg = &Registry{reg: make(map[int]*inStream)}
)

type inStream struct {
	rdr    io.Reader
	pos    int
	buf    []byte
	closed bool
}

type Registry struct {
	reg    map[int]*inStream
	nextID int
}

func (r *Registry) Alloc(is *inStream) int {
	id := r.nextID
	r.nextID++
	r.reg[id] = is
	return id
}

func (r *Registry) Get(id int) *inStream {
	return r.reg[id]
}

func (r *Registry) Release(id int) {
	delete(r.reg, id)
}

//export istream_read
func istream_read(id int, size int) C.csv_res_t {
	res := C.csv_res_t{nil, 0, nil}

	is := reg.Get(id)
	if is == nil {
		err := fmt.Sprintf("%d: unknown id", id)
		res.err = C.CString(err)
		return res
	}

	if size > len(is.buf) {
		is.buf = make([]byte, size)
	}

	n, err := is.rdr.Read(is.buf)
	if err != nil {
		if err == io.EOF {
			is.closed = true
		} else {
			res.err = C.CString(err.Error())
			return res
		}
	}
	res.size = C.ulonglong(n)
	res.data = C.CBytes(is.buf[:n])

	return res
}

//export istream_tell
func istream_tell(id int) C.csv_res_t {
	res := C.csv_res_t{nil, 0, nil}

	is := reg.Get(id)
	if is == nil {
		err := fmt.Sprintf("%d: unknown id", id)
		res.err = C.CString(err)
		return res
	}

	res.size = C.ulonglong(is.pos)

	return res
}

//export istream_closed
func istream_closed(id int) C.csv_res_t {
	res := C.csv_res_t{nil, 0, nil}

	is := reg.Get(id)
	if is == nil {
		err := fmt.Sprintf("%d: unknown id", id)
		res.err = C.CString(err)
		return res
	}

	if is.closed {
		res.size = 1
	}
	return res
}

// Reads a CSV data from rdr, returns a *carrow.Table
func Read(rdr io.Reader) (*carrow.Table, error) {
	is := &inStream{rdr: rdr}
	id := reg.Alloc(is)
	defer reg.Release(id)
	res := C.csv_read(C.longlong(id))
	if res.err != nil {
		// TODO: Free res.err?
		return nil, fmt.Errorf(C.GoString(res.err))
	}

	ptr := unsafe.Pointer(res.table)
	return carrow.NewTableFromPtr(ptr), nil
}
