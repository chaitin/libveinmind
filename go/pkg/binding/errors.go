package binding

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

// #include "veinmind.h"
import "C"

// RuntimeError is the error that must be thrown out.
type RuntimeError ErrorType

const (
	ENOMEM   = RuntimeError(C.veinmind_ENOMEM)
	ENULLPTR = RuntimeError(C.veinmind_ENULLPTR)
	ENULLID  = RuntimeError(C.veinmind_ENULLID)
	ENOCAP   = RuntimeError(C.veinmind_ENOCAP)
	EPANIC   = RuntimeError(C.veinmind_EPANIC)
	ERANGE   = RuntimeError(C.veinmind_ERANGE)
	ERACE    = RuntimeError(C.veinmind_ERACE)
)

func (r RuntimeError) Error() string {
	switch r {
	case ENOMEM:
		return "out of memory"
	case ENULLPTR:
		return "null pointer receiver"
	case ENULLID:
		return "invalid ID specified"
	case ENOCAP:
		return "invalid OP specified"
	case EPANIC:
		return "panic in call"
	case ERACE:
		return "race detected"
	default:
		return "unknown error"
	}
}

func (h Handle) IsEOK() bool {
	return ErrorType(h) == C.veinmind_EOK
}

// WrapError is the error that stacks another error.
type WrapError struct {
	msg  string
	next error
}

func (w WrapError) Error() string {
	if w.next == nil {
		return w.msg
	}
	return fmt.Sprintf("%s: %v", w.msg, w.next)
}

func (w WrapError) Unwrap() error {
	return w.next
}

func (h Handle) WrapError() *WrapError {
	result := &WrapError{}
	result.msg = h.ErrorMessage()
	next := h.ErrorUnwrap()
	defer next.Free()
	result.next = handleErrorInternal(next)
	return result
}

func (h Handle) OSPathError() *os.PathError {
	result := &os.PathError{}
	result.Op = h.OSPathErrorGetOp()
	result.Path = h.OSPathErrorGetPath()
	next := h.ErrorUnwrap()
	defer next.Free()
	result.Err = handleErrorInternal(next)
	return result
}

// assertNoError is used under the condition of calling
// an API that will not generate exceptions.
//
// The only error that can be generated should be runtime
// error, and all of them must be reported by panicking.
func assertNoError(code ErrorType) {
	err := Handle(IDType(code))
	if err.IsEOK() {
		return
	}
	defer err.Free()
	panic(RuntimeError(code))
}

// Error formats the error by brutely return error string.
func (h Handle) Error() error {
	return errors.New(h.ErrorString())
}

// Errno formats the error by converting it into syscall.Errno.
func (h Handle) SyscallError() syscall.Errno {
	return syscall.Errno(uintptr(h.SyscallErrorGetErrno()))
}

// handleError is used under the condition of calling
// an API that might generate exceptions.
//
// The error other than runtime error must be reconstruted
// properly, in order to be handled properly.
func handleError(code ErrorType) error {
	err := Handle(IDType(code))
	if err.IsEOK() {
		return nil
	}
	defer err.Free()
	if code < runtimeErrorMaxValue {
		panic(RuntimeError(code))
	}
	return handleErrorInternal(err)
}

// handleWalkError is used for walk function
// error handler.
func handleWalkError(code ErrorType) error {
	err := Handle(IDType(code))
	if err.IsEOK() {
		return nil
	}
	return handleErrorInternal(err)
}
