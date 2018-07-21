package vmguestlib

/*
#cgo pkg-config: vmguestlib
#cgo CFLAGS: -I../native
#include <vmGuestLib.h>
#include <vmGuestLibProxy.h>
*/
import "C"

// Handle provides a context for accessing all GuestLib
// state. Use OpenHandle() to get a handle for use with other
// GuestLib functions, and use CloseHandle() to release a handle
// previously acquired with OpenHandle().
//
// All of the statistics and session state are maintained per GuestLib
// handle, so operating on one GuestLib handle will not affect the
// state of another handle.
//
type Handle struct {
	// The native VMGuestLibHandle object.
	NativeHandle *C.VMGuestLibHandle
}

// NewHandle opens a new handle.
func NewHandle() (*Handle, error) {
	h := &Handle{
		NativeHandle: new(C.VMGuestLibHandle),
	}
	if e := C.VMGuestLib_OpenHandle(h.NativeHandle); e != ErrorSuccess {
		return h, newError(e)
	}
	return h, nil
}

// Close releases a previously opened handle.
func (h *Handle) Close() {
	C.VMGuestLib_CloseHandle(*h.NativeHandle)
}

// UpdateInfo updates the session state for the current handle.
// No locking is done internally, so multithreading environments
// should use separate handles. This is a fairly heavyweight
// function on the native API side so environments concerned
// about performance should minimize the number of calls to
// this function.
func (h *Handle) UpdateInfo() error {
	if e := C.VMGuestLib_UpdateInfo(*h.NativeHandle); e != ErrorSuccess {
		return newError(e)
	}
	return nil
}

// getUint32Value retrieves one uint32 native accessor value
// by calling the native accessor routine wrapped into a
// native proxy function pointer type and converting the C
// uint32 value to golang's relevant type.
func (h *Handle) getUint32Value(p C.p_uint32_f) (uint32, error) {
	nativeUint32 := new(C.uint32)
	if e := C.proxy_uint32_f(p, *h.NativeHandle, nativeUint32); e != ErrorSuccess {
		return 0, newError(e)
	}
	return uint32(*nativeUint32), nil
}

// getUint64Value retrieves one uint64 native accessor value
// by calling the native accessor routine wrapped into a
// native proxy function pointer type and converting the C
// uint64 value to golang's relevant type.
func (h *Handle) getUint64Value(p C.p_uint64_f) (uint64, error) {
	nativeUint64 := new(C.uint64)
	if e := C.proxy_uint64_f(p, *h.NativeHandle, nativeUint64); e != ErrorSuccess {
		return 0, newError(e)
	}
	return uint64(*nativeUint64), nil
}
