//go:build nocompat
// +build nocompat

package binding

import (
	"io"
)

// #include "veinmind.h"
import "C"

// runtimeErrorMaxValue is specified to the enum constant
// when we will not consider forward compatibility.
const runtimeErrorMaxValue = C.veinmind_RuntimeErrorMax

// handleErrorInternal will convert the error recursively.
func handleErrorInternal(err Handle) error {
	if err.IsEOK() {
		return nil
	}
	switch err.Type() {
	case TypeWrapError:
		return err.WrapError()
	case TypeSyscallError:
		return err.SyscallError()
	case TypeOSPathError:
		return err.OSPathError()
	case TypeEOFError:
		return io.EOF
	default:
		return err.Error()
	}
}
