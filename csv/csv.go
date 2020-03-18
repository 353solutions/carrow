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
func Read(rdr io.Reader, po *ParseOptions) (*carrow.Table, error) {
	is := &inStream{rdr: rdr}
	id := reg.Alloc(is)
	defer reg.Release(id)
	res := C.csv_read(C.longlong(id), po.c)
	if res.err != nil {
		// TODO: Free res.err?
		return nil, fmt.Errorf(C.GoString(res.err))
	}

	ptr := unsafe.Pointer(res.table)
	return carrow.NewTableFromPtr(ptr), nil
}

// ParseOptions used by ParseOption
// used for not exposing C internals in the API
type ParseOptions struct {
	c C.parse_options_t
}

// NewParseOptions return parse options
func NewParseOptions(opts ...ParseOption) *ParseOptions {
	p := &ParseOptions{c: C.default_parse_options()}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// ParseOption is a parsing option
type ParseOption func(*ParseOptions)

func WithDelimiter(char byte) ParseOption {
	return func(p *ParseOptions) {
		p.c.delimiter = C.char(char)
	}
}

func WithoutQuoting(p *ParseOptions) {
	p.c.quoting = 0
}

func WithoutQuoteChar(char byte) ParseOption {
	return func(p *ParseOptions) {
		p.c.quote_char = C.char(char)
	}
}

func WithoutDoubleQuote(p *ParseOptions) {
	p.c.double_quote = 0
}

func WithEscaping(p *ParseOptions) {
	p.c.escaping = 1
}

func WithoutEscapeChar(char byte) ParseOption {
	return func(p *ParseOptions) {
		p.c.escape_char = C.char(char)
	}
}

func WithNewlinesInValues(p *ParseOptions) {
	p.c.newlines_in_values = 1
}

func WithNoIgnoreEmptyLines(p *ParseOptions) {
	p.c.ignore_empty_lines = 0
}
