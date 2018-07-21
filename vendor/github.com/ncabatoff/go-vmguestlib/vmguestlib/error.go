package vmguestlib

/*
#cgo pkg-config: vmguestlib
#cgo CFLAGS: -I../native
#include <vmGuestLib.h>
*/
import "C"

// Error codes returned by the native VMGuestLib API.
const (
	ErrorSuccess            = iota // No error.
	ErrorOther                     // Other error.
	ErrorNotRunningInVm            // Not running in a VM.
	ErrorNotEnabled                // GuestLib not enabled on the host.
	ErrorNotAvailable              // This stat is not available on the host.
	ErrorNoInfo                    // UpdateInfo() has never been called.
	ErrorMemory                    // Not enough memory.
	ErrorBufferTooSmall            // Buffer too small.
	ErrorInvalidHandle             // Handle is invalid.
	ErrorInvalidArg                // One or more arguments were invalid.
	ErrorUnsupportedVersion        // The host doesn't support this request.
)

// Error messages associated to the native VMGuestLib API error codes.
const (
	ErrorMsgSuccess            = "The function has completed successfully."
	ErrorMsgOther              = "An error has occurred. No additional information about the type of the error is available."
	ErrorMsgNotRunningInVm     = "The program making this call is not running on a VMware virtualmachine."
	ErrorMsgNotEnabled         = "The vSphere Guest API is not enabled on this host so these functions cannot be used."
	ErrorMsgNotAvailable       = "The information requested is not available on this host."
	ErrorMsgNoInfo             = "The handle data structure does not contain any information."
	ErrorMsgMemory             = "There is not enough memory available to complete the call."
	ErrorMsgBufferTooSmall     = "The buffer is too small to accommodate the function call."
	ErrorMsgInvalidHandle      = "The handle that you used is invalid."
	ErrorMsgInvalidArg         = "One or more of the arguments passed to the function were invalid."
	ErrorMsgUnsupportedVersion = "The host does not support the requested statistic."
)

var errors = map[uint]string{
	ErrorSuccess:            ErrorMsgSuccess,
	ErrorOther:              ErrorMsgOther,
	ErrorNotRunningInVm:     ErrorMsgNotRunningInVm,
	ErrorNotEnabled:         ErrorMsgNotEnabled,
	ErrorNotAvailable:       ErrorMsgNotAvailable,
	ErrorNoInfo:             ErrorMsgNoInfo,
	ErrorMemory:             ErrorMsgMemory,
	ErrorBufferTooSmall:     ErrorMsgBufferTooSmall,
	ErrorInvalidHandle:      ErrorMsgInvalidHandle,
	ErrorInvalidArg:         ErrorMsgInvalidArg,
	ErrorUnsupportedVersion: ErrorMsgUnsupportedVersion,
}

// Error is a native error code with an explanatory message.
type Error struct {
	// The native VMGuestLibError object.
	NativeError *C.VMGuestLibError
	// The error message corresponding to the native error.
	Message string
}

// NewError creates a VMGuestLib error from a native error.
func newError(e C.VMGuestLibError) *Error {
	return &Error{
		NativeError: &e,
		Message:     errors[uint(e)],
	}
}

// Error returns the message associated with the error.
func (e *Error) Error() string {
	return e.Message
}
