package vmguestlib

/*
#cgo pkg-config: vmguestlib
#cgo CFLAGS: -I../native
#include <vmGuestLib.h>
#include <vmSessionId.h>
#include <vmGuestLibProxy.h>
*/
import "C"

// Session gives access to accessor functions through the usage of a
// particular Handle.
// Session is also used to detect changes in the "session" of a virtual
// machine. "Session" in this context refers to the particular running
// instance of this virtual machine on a given host. Moving a virtual
// machine to another host using VMotion will cause a session change
// as well as suspending and resuming a virtual machine or reverting
// to a snapshot.
//
// Events like VMotion, suspend/resume, snapshot revert are likely
// to render invalid any information previously retrieved through
// this API, so the intention of the session is to provide applications
// with a mechanism to detect those events and react accordingly,
// e.g. by refreshing and resetting any state that relies on validity
// of previously retrieved information.
//
type Session struct {
	// A handle object.
	Handle *Handle
	// The native VMSessionId object.
	NativeSession *C.VMSessionId
}

// Open fetch a new session from the given handle.
func Open(h *Handle) (*Session, error) {
	s := &Session{
		Handle:        h,
		NativeSession: new(C.VMSessionId),
	}
	_, err := s.RefreshInfo()
	return s, err
}

// Close tears down the underlying handle.
func (s *Session) Close() {
	s.Handle.Close()
}

// NewSession allocates a new handle and calls Open().
func NewSession() (*Session, error) {
	var h *Handle
	var s *Session
	var err error

	if h, err = NewHandle(); err != nil {
		return nil, err
	}
	if s, err = Open(h); err != nil {
		return nil, err
	}

	return s, nil
}

// Refresh updates the session if the virtual machine session
// changed on the host.
func (s *Session) Refresh() (bool, error) {
	newSession := new(C.VMSessionId)
	if e := C.VMGuestLib_GetSessionId(*s.Handle.NativeHandle, newSession); e != ErrorSuccess {
		return false, newError(e)
	}
	changed := (*newSession != *s.NativeSession)
	s.NativeSession = newSession
	return changed, nil
}

// RefreshInfo calls UpdateInfo() on the underlying handle
// then calls Refresh() on this session.
func (s *Session) RefreshInfo() (bool, error) {
	if err := s.Handle.UpdateInfo(); err != nil {
		return false, err
	}
	return s.Refresh()
}
