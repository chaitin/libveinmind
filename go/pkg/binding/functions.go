package binding

//#include "veinmind.h"
import "C"

func (h Handle) Close() error {
	return handleError(C.veinmind_Close(h.ID()))
}

func (h Handle) Append(item Handle) {
	assertNoError(C.veinmind_Append(h.ID(), item.ID()))
}

func (h Handle) ErrorString() string {
	var str Handle
	assertNoError(C.veinmind_ErrorString(str.Ptr(), h.ID()))
	defer str.Free()
	return str.String()
}

func (h Handle) IsWrapError() bool {
	return C.veinmind_IsWrapError(h.ID()) != C.int(0)
}

func (h Handle) ErrorMessage() string {
	var str Handle
	assertNoError(C.veinmind_ErrorMessage(str.Ptr(), h.ID()))
	defer str.Free()
	return str.String()
}

func (h Handle) ErrorUnwrap() Handle {
	var next Handle
	assertNoError(C.veinmind_ErrorUnwrap(next.Ptr(), h.ID()))
	return next
}

func (h Handle) IsSyscallError() bool {
	return C.veinmind_IsSyscallError(h.ID()) != C.int(0)
}

func (h Handle) SyscallErrorGetErrno() SizeType {
	var errno SizeType
	assertNoError(C.veinmind_SyscallErrorGetErrno(&errno, h.ID()))
	return errno
}

func (h Handle) IsOSPathError() bool {
	return C.veinmind_IsOSPathError(h.ID()) != C.int(0)
}

func (h Handle) OSPathErrorGetOp() string {
	var str Handle
	assertNoError(C.veinmind_OSPathErrorGetOp(str.Ptr(), h.ID()))
	defer str.Free()
	return str.String()
}

func (h Handle) IsEOFError() bool {
	return C.veinmind_IsEOFError(h.ID()) != C.int(0)
}

func (h Handle) OSPathErrorGetPath() string {
	var str Handle
	assertNoError(C.veinmind_OSPathErrorGetPath(str.Ptr(), h.ID()))
	defer str.Free()
	return str.String()
}

func (h Handle) Read(bytes Handle) (SizeType, error) {
	var n SizeType
	err := handleError(C.veinmind_Read(&n, h.ID(), bytes.ID()))
	return n, err
}

func (h Handle) ReadAt(bytes Handle, off int64) (SizeType, error) {
	var n SizeType
	err := handleError(C.veinmind_ReadAt(
		&n, h.ID(), bytes.ID(), C.int64_t(off)))
	return n, err
}

func (h Handle) Write(bytes Handle) (SizeType, error) {
	var n SizeType
	err := handleError(C.veinmind_Write(&n, h.ID(), bytes.ID()))
	return n, err
}

func (h Handle) WriteAt(bytes Handle, off int64) (SizeType, error) {
	var n SizeType
	err := handleError(C.veinmind_WriteAt(
		&n, h.ID(), bytes.ID(), C.int64_t(off)))
	return n, err
}

func (h Handle) FileStat() (Handle, error) {
	var result Handle
	if err := handleError(C.veinmind_FileStat(
		result.Ptr(), h.ID())); err != nil {
		return 0, err
	}
	return result, nil
}

func (h Handle) Open(path string) (Handle, error) {
	str := NewString(path)
	defer str.Free()
	var result Handle
	if err := handleError(C.veinmind_Open(
		result.Ptr(), h.ID(), str.ID())); err != nil {
		return 0, err
	}
	return result, nil
}

func (h Handle) Seek(offset int64, whence int) (int64, error) {
	var off C.int64_t
	err := handleError(C.veinmind_Seek(&off,
		h.ID(), C.int64_t(offset), C.int(whence)))
	return int64(off), err
}

func (h Handle) Stat(path string) (Handle, error) {
	str := NewString(path)
	defer str.Free()
	var result Handle
	if err := handleError(C.veinmind_Stat(
		result.Ptr(), h.ID(), str.ID())); err != nil {
		return 0, err
	}
	return result, nil
}

func (h Handle) Lstat(path string) (Handle, error) {
	str := NewString(path)
	defer str.Free()
	var result Handle
	if err := handleError(C.veinmind_Lstat(
		result.Ptr(), h.ID(), str.ID())); err != nil {
		return 0, err
	}
	return result, nil
}

func (h Handle) Readlink(path string) (string, error) {
	str := NewString(path)
	defer str.Free()
	var result Handle
	if err := handleError(C.veinmind_Readlink(
		result.Ptr(), h.ID(), str.ID())); err != nil {
		return "", err
	}
	defer result.Free()
	return result.String(), nil
}

func (h Handle) EvalSymlink(path string) (string, error) {
	str := NewString(path)
	defer str.Free()
	var result Handle
	if err := handleError(C.veinmind_EvalSymlink(
		result.Ptr(), h.ID(), str.ID())); err != nil {
		return "", err
	}
	defer result.Free()
	return result.String(), nil
}

func (h Handle) Readdir(path string) (Handle, error) {
	str := NewString(path)
	defer str.Free()
	var result Handle
	if err := handleError(C.veinmind_Readdir(
		result.Ptr(), h.ID(), str.ID())); err != nil {
		return 0, err
	}
	return result, nil
}

func (h Handle) FileInfoName() string {
	var str Handle
	assertNoError(C.veinmind_FileInfoName(str.Ptr(), h.ID()))
	defer str.Free()
	return str.String()
}

func (h Handle) FileInfoSize() SizeType {
	var size SizeType
	assertNoError(C.veinmind_FileInfoSize(&size, h.ID()))
	return size
}

func (h Handle) FileInfoMode() uint32 {
	var mode C.uint32_t
	assertNoError(C.veinmind_FileInfoMode(&mode, h.ID()))
	return uint32(mode)
}

func (h Handle) FileInfoModTime() int64 {
	var modTime C.int64_t
	assertNoError(C.veinmind_FileInfoModTime(&modTime, h.ID()))
	return int64(modTime)
}

func (h Handle) FileInfoSys() Handle {
	var result Handle
	assertNoError(C.veinmind_FileInfoSys(result.Ptr(), h.ID()))
	return result
}

func (h Handle) RuntimeListImageIDs() ([]string, error) {
	result := new(Handle)
	if err := handleError(C.veinmind_RuntimeListImageIDs(
		result.Ptr(), h.ID())); err != nil {
		return nil, err
	}
	defer result.Free()
	return result.StringArray(), nil
}

func (h Handle) RuntimeFindImageIDs(pattern string) ([]string, error) {
	str := NewString(pattern)
	defer str.Free()
	var result Handle
	if err := handleError(C.veinmind_RuntimeFindImageIDs(
		result.Ptr(), h.ID(), str.ID())); err != nil {
		return nil, err
	}
	defer result.Free()
	return result.StringArray(), nil
}

func (h Handle) RuntimeOpenImageByID(id string) (Handle, error) {
	str := NewString(id)
	defer str.Free()
	var result Handle
	if err := handleError(C.veinmind_RuntimeOpenImageByID(
		result.Ptr(), h.ID(), str.ID())); err != nil {
		return 0, err
	}
	return result, nil
}

func (h Handle) ImageID() string {
	var str Handle
	assertNoError(C.veinmind_ImageID(str.Ptr(), h.ID()))
	defer str.Free()
	return str.String()
}

func (h Handle) ImageRepos() ([]string, error) {
	var result Handle
	if err := handleError(C.veinmind_ImageRepos(
		result.Ptr(), h.ID())); err != nil {
		return nil, err
	}
	defer result.Free()
	return result.StringArray(), nil
}

func (h Handle) ImageRepoRefs() ([]string, error) {
	var result Handle
	if err := handleError(C.veinmind_ImageRepoRefs(
		result.Ptr(), h.ID())); err != nil {
		return nil, err
	}
	defer result.Free()
	return result.StringArray(), nil
}

func (h Handle) ImageOCISpecV1MarshalJSON() ([]byte, error) {
	var result Handle
	if err := handleError(C.veinmind_ImageOCISpecV1MarshalJSON(
		result.Ptr(), h.ID())); err != nil {
		return nil, err
	}
	defer result.Free()
	return result.Bytes(), nil
}

func DockerMakeNewOptionList() Handle {
	var result Handle
	assertNoError(C.veinmind_DockerMakeNewOptionList(result.Ptr()))
	return result
}

func DockerNew(opts Handle) (Handle, error) {
	var result Handle
	if err := handleError(C.veinmind_DockerNew(
		result.Ptr(), opts.ID())); err != nil {
		return 0, err
	}
	return result, nil
}

func DockerWithConfigPath(path string) Handle {
	str := NewString(path)
	defer str.Free()
	var result Handle
	assertNoError(C.veinmind_DockerWithConfigPath(
		result.Ptr(), str.ID()))
	return result
}

func DockerWithDataRootDir(path string) Handle {
	str := NewString(path)
	defer str.Free()
	var result Handle
	assertNoError(C.veinmind_DockerWithDataRootDir(
		result.Ptr(), str.ID()))
	return result
}

func DockerWithUniqueDesc(desc string) Handle {
	str := NewString(desc)
	defer str.Free()
	var result Handle
	assertNoError(C.veinmind_DockerWithUniqueDesc(
		result.Ptr(), str.ID()))
	return result
}

func (h Handle) DockerUniqueDesc() string {
	var str Handle
	assertNoError(C.veinmind_DockerUniqueDesc(
		str.Ptr(), h.ID()))
	defer str.Free()
	return str.String()
}

func (h Handle) DockerImageOpenLayer(i int) (Handle, error) {
	var result Handle
	if err := handleError(C.veinmind_DockerImageOpenLayer(
		result.Ptr(), h.ID(), C.size_t(i))); err != nil {
		return 0, err
	}
	return result, nil
}

func (h Handle) DockerImageGetLayerDiffID(i int) (string, error) {
	var result Handle
	if err := handleError(C.veinmind_DockerImageGetLayerDiffID(
		result.Ptr(), h.ID(), C.size_t(i))); err != nil {
		return "", nil
	}
	defer result.Free()
	return result.String(), nil
}

func (h Handle) DockerImageNumLayers() int {
	var numLayers C.size_t
	assertNoError(C.veinmind_DockerImageNumLayers(&numLayers, h.ID()))
	return int(numLayers)
}

func (h Handle) DockerLayerID() string {
	var result Handle
	assertNoError(C.veinmind_DockerLayerID(result.Ptr(), h.ID()))
	defer result.Free()
	return result.String()
}

func ContainerdMakeNewOptionList() Handle {
	var result Handle
	assertNoError(C.veinmind_ContainerdMakeNewOptionList(result.Ptr()))
	return result
}

func ContainerdNew(opts Handle) (Handle, error) {
	var result Handle
	if err := handleError(C.veinmind_ContainerdNew(
		result.Ptr(), opts.ID())); err != nil {
		return 0, err
	}
	return result, nil
}

func ContainerdWithConfigPath(path string) Handle {
	str := NewString(path)
	defer str.Free()
	var result Handle
	assertNoError(C.veinmind_ContainerdWithConfigPath(
		result.Ptr(), str.ID()))
	return result
}

func ContainerdWithRootDir(path string) Handle {
	str := NewString(path)
	defer str.Free()
	var result Handle
	assertNoError(C.veinmind_ContainerdWithRootDir(
		result.Ptr(), str.ID()))
	return result
}

func ContainerdWithUniqueDesc(desc string) Handle {
	str := NewString(desc)
	defer str.Free()
	var result Handle
	assertNoError(C.veinmind_ContainerdWithUniqueDesc(
		result.Ptr(), str.ID()))
	return result
}

func (h Handle) ContainerdUniqueDesc() string {
	var str Handle
	assertNoError(C.veinmind_ContainerdUniqueDesc(
		str.Ptr(), h.ID()))
	defer str.Free()
	return str.String()
}
