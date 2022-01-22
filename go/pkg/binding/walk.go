package binding

import (
	"path/filepath"
	"sync"
	"unsafe"
)

/*
#include "stdlib.h"
#include "veinmind.h"

extern int veinmind_WalkHandler(
	veinmind_id_t name, veinmind_id_t info,
	veinmind_err_t err, void* userdata);

static veinmind_err_t veinmind_InvokeWalkHandler(
	veinmind_id_t fs, veinmind_id_t root, void* userdata) {
	return veinmind_Walk(fs, root,
		veinmind_WalkHandler, userdata);
}
*/
import "C"

type WalkFunc func(name string, info Handle, err error) error

var walkFuncMap sync.Map

type walkState struct {
	err error
	f   WalkFunc
}

func (h Handle) Walk(path string, f WalkFunc) error {
	str := NewString(path)
	defer str.Free()
	for {
		state := &walkState{f: f}
		cookie := C.CBytes([]byte("cookie"))
		defer C.free(cookie)
		if _, load := walkFuncMap.LoadOrStore(cookie, state); load {
			continue
		}
		defer walkFuncMap.Delete(cookie)
		if err := handleError(C.veinmind_InvokeWalkHandler(
			h.ID(), str.ID(), unsafe.Pointer(cookie))); err != nil {
			return err
		}
		return state.err
	}
}

//export veinmind_WalkHandler
func veinmind_WalkHandler(
	nameID, infoID IDType, errID ErrorType,
	userdata unsafe.Pointer,
) C.int {
	object, _ := walkFuncMap.Load(userdata)
	state := object.(*walkState)
	name := Handle(nameID).String()
	info := Handle(infoID)
	err := handleError(errID)
	if walkErr := state.f(name, info, err); walkErr != nil {
		if walkErr == filepath.SkipDir {
			return C.int(1)
		}
		state.err = walkErr
		return C.int(2)
	}
	return C.int(0)
}
