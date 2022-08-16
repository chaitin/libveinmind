//go:build nocompat
// +build nocompat

package binding

//#include "veinmind.h"
import "C"

type Type = C.veinmind_type_t

const (
	TypeNull = Type(C.veinmind_TypeNull)

	TypeBytes       = Type(C.veinmind_TypeBytes)
	TypeString      = Type(C.veinmind_TypeString)
	TypeStringArray = Type(C.veinmind_TypeStringArray)

	TypeError        = Type(C.veinmind_TypeError)
	TypeWrapError    = Type(C.veinmind_TypeWrapError)
	TypeRuntimeError = Type(C.veinmind_TypeRuntimeError)
	TypeSyscallError = Type(C.veinmind_TypeSyscallError)

	TypeOSFile          = Type(C.veinmind_TypeOSFile)
	TypeOSFileInfo      = Type(C.veinmind_TypeOSFileInfo)
	TypeOSFileInfoArray = Type(C.veinmind_TypeOSFileInfoArray)
	TypeOSPathError     = Type(C.veinmind_TypeOSPathError)
	TypeEOFError        = Type(C.veinmind_TypeEOFError)

	TypeRuntime   = Type(C.veinmind_TypeRuntime)
	TypeContainer = Type(C.veinmind_TypeContainer)
	TypeImage     = Type(C.veinmind_TypeImage)
	TypeFile      = Type(C.veinmind_TypeFile)

	TypeDockerRuntime   = Type(C.veinmind_TypeDockerRuntime)
	TypeDockerContainer = Type(C.veinmind_TypeDockerContainer)
	TypeDockerImage     = Type(C.veinmind_TypeDockerImage)
	TypeDockerLayer     = Type(C.veinmind_TypeDockerLayer)

	TypeContainerdRuntime   = Type(C.veinmind_TypeContainerdRuntime)
	TypeContainerdContainer = Type(C.veinmind_TypeContainerdContainer)
	TypeContainerdImage     = Type(C.veinmind_TypeContainerdImage)
)

func (h Handle) Type() Type {
	return C.veinmind_Type(h.ID())
}
