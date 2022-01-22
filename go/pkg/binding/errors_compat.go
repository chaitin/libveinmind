//go:build !nocompat
// +build !nocompat

package binding

import (
	"io"
)

// #include "veinmind.h"
import "C"

// runtimeErrorMaxValue is max value of runtime error.
var runtimeErrorMaxValue ErrorType

func init() {
	// Retrieve the value dynamically to ensure backward
	// compatibility when new errors are added.
	runtimeErrorMaxValue = C.veinmind_RuntimeErrorMaxValue()
}

// handleErrorInternal will convert the error recursively.
func handleErrorInternal(err Handle) error {
	switch {
	case err.IsEOK():
		return nil
	case err.IsWrapError():
		return err.WrapError()
	case err.IsSyscallError():
		return err.SyscallError()
	case err.IsOSPathError():
		return err.OSPathError()
	case err.IsEOFError():
		return io.EOF
	default:
		return err.Error()
	}
}
