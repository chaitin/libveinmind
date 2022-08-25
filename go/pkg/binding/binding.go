// Package binding is the actual package binding part that
// requires the libveinmind library through pkg-config, and
// and attempt to reconstruct the API interface from it.
package binding

import (
	"unsafe"
)

// #cgo pkg-config: libveinmind
// #include "veinmind.h"
// #include <stdlib.h>
import "C"

type (
	IDType    = C.veinmind_id_t
	ErrorType = C.veinmind_err_t
	SizeType  = C.size_t
)

// Handle to the resource managed by the API binding.
type Handle IDType

// Ptr returns the underlying pointer type.
func (h *Handle) Ptr() *IDType {
	return (*IDType)(h)
}

// ID returns the underlying id value type.
func (h Handle) ID() IDType {
	return (IDType)(h)
}

// Free the object registry occupied by the.
func (h *Handle) Free() {
	C.veinmind_Free(h.ID())
	*h = Handle(IDType(0))
}

// IsNil judges whether the binding is null handle.
func (h Handle) IsNil() bool {
	return IDType(h) == 0
}

// Length returns the length of the specified Handle.
func (h Handle) Length() int {
	var length SizeType
	assertNoError(C.veinmind_Length(&length, h.ID()))
	return int(length)
}

// Index create a new object at specified index.
func (h Handle) Index(i int) Handle {
	var obj Handle
	assertNoError(C.veinmind_Index(obj.Ptr(), h.ID(), SizeType(i)))
	return obj
}

// Crop creates a new slice at specified range.
func (h Handle) Crop(begin, end int) Handle {
	var obj Handle
	assertNoError(C.veinmind_Crop(obj.Ptr(), h.ID(),
		SizeType(begin), SizeType(end)))
	return obj
}

// RawPointer returns the raw pointer of the object.
func (h Handle) RawPointer() unsafe.Pointer {
	var ptr unsafe.Pointer
	assertNoError(C.veinmind_RawPointer(&ptr, h.ID()))
	return ptr
}

// Bytes returns the underying byte slice.
func (h Handle) Bytes() []byte {
	length := h.Length()
	data := h.RawPointer()
	return C.GoBytes(data, C.int(length))
}

// BytesToString creates a string from bytes.
func (h Handle) BytesToString() Handle {
	var obj Handle
	assertNoError(C.veinmind_BytesToString(obj.Ptr(), h.ID()))
	return obj
}

// StringToBytes creates a bytes from strin.
func (h Handle) StringToBytes() Handle {
	var obj Handle
	assertNoError(C.veinmind_StringToBytes(obj.Ptr(), h.ID()))
	return obj
}

// String returns the underlying string object.
func (h Handle) String() string {
	bytes := h.StringToBytes()
	defer bytes.Free()
	return string(bytes.Bytes())
}

// StringArray returns the parsed array of string object.
func (h Handle) StringArray() []string {
	length := h.Length()
	var item Handle
	defer func() { item.Free() }()
	var result []string
	for i := 0; i < length; i++ {
		item = h.Index(i)
		result = append(result, item.String())
		item.Free()
	}
	return result
}

// Int32 return the parsed of int32 object.
func (h Handle) Int32() int32 {
	var result C.int32_t
	assertNoError(C.veinmind_Int32(&result, h.ID()))
	return int32(result)
}

// Int32Array return the parsed array of int32 object.
func (h Handle) Int32Array() []int32 {
	length := h.Length()
	var item Handle
	defer func() { item.Free() }()
	var result []int32
	for i := 0; i < length; i++ {
		item = h.Index(i)
		result = append(result, item.Int32())
		item.Free()
	}
	return result
}

// NewBytes pushes a bytes buffer and creates its handle.
func NewBytes(b []byte) Handle {
	bytes := C.CBytes(b)
	defer C.free(bytes)
	var result Handle
	assertNoError(C.veinmind_Bytes(result.Ptr(),
		bytes, SizeType(len(b))))
	return result
}

// NewString pushes a string and creates its handle.
func NewString(str string) Handle {
	bytes := NewBytes([]byte(str))
	defer bytes.Free()
	return bytes.BytesToString()
}
