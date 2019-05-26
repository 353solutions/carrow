package plasma

import (
    "unsafe"
    "fmt"
    "runtime"
    "math/rand"
    "time"
    "encoding/hex"
    
    "github.com/353solutions/carrow"
)

/*
#cgo pkg-config: arrow plasma
#cgo LDFLAGS: -lcarrow -L..
#cgo CFLAGS: -I..
// FIXME: plasma headers

#include "carrow.h"
#include <stdlib.h>
*/
import "C"

const (
    // IDLength is length of ObjectID in bytes
    IDLength = 20
)

var (
    idRnd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// Client is a client to Arrow's plasma store
type Client struct {
    ptr unsafe.Pointer
}

// ObjectID is store ID for an object
type ObjectID [IDLength]byte

// Connect connects to plasma store
func Connect(path string) (*Client, error) {
    cStr := C.CString(path)
    ptr := C.plasma_connect(cStr)
    C.free(unsafe.Pointer(cStr))
    if ptr == nil {
        return nil, fmt.Errorf("can't connect to %s", path)
    }

    client := &Client{ptr}
    runtime.SetFinalizer(client, func(c *Client) {
        c.Disconnect()
    })

    return client, nil
}

// WriteTable write a table to plasma store
// If id is empty, a new random id will be generated
func (c *Client) WriteTable(t *carrow.Table, id ObjectID) error {
    cID := C.CString(string(id[:]))
    n := C.plasma_write(c.ptr, t.Ptr(), cID)
    C.free(unsafe.Pointer(cID))

    if n == -1 {
        return fmt.Errorf("can't write table") // TODO
    }

    return nil
}

// ReadTable reads a table from plasma store
func (c *Client) ReadTable(id ObjectID, timeout time.Duration) (*carrow.Table, error) {
    cID := C.CString(string(id[:]))
    msec := C.int64_t(timeout / time.Millisecond)
    ptr := C.plasma_read(c.ptr, cID, msec)
    C.free(unsafe.Pointer(cID))

    if ptr == nil {
        return nil, fmt.Errorf("can't read %s", id)
    }

    return carrow.NewTableFromPtr(ptr), nil
}

// Release releases (deletes) object from plasma store
func (c *Client) Release(id ObjectID) error {
    cID := C.CString(string(id[:]))
    out := C.plasma_release(c.ptr, cID)
    C.free(unsafe.Pointer(cID))

    if out != 0 {
        return fmt.Errorf("can't release object %s", id)

    }

    return nil
}

// Disconnect disconnects from plasma store
func (c *Client) Disconnect() {
    if c.ptr == nil {
        return
    }
    C.plasma_disconnect(c.ptr)
    c.ptr = nil
}	

func (oid ObjectID) String() string {
    return hex.EncodeToString(oid[:])
}

// RandomID return a new random plasma ID
func RandomID() (ObjectID, error) {
    var oid ObjectID
    _, err := idRnd.Read(oid[:])
    if err != nil {
        return oid, err
    }

    oid[8] = (oid[8] | 0x80) & 0xBF
    oid[6] = (oid[6] | 0x40) & 0x4F
    return oid, nil
}

// IDFromString converts a string to ObjectID
func IDFromString(s string) (ObjectID, error) {
    data := s[:]
    var oid ObjectID
    if len(data) != IDLength {
        return oid, fmt.Errorf("wrong length, expected %d, got %d", IDLength, len(data))
    }
    copy(oid[:], data)
    return oid, nil
}