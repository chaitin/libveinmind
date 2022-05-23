package behaviour

import (
	"os"
	"path/filepath"
	"time"
	"unsafe"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/pkg/binding"
)

// fileInfo is the mocking object implementing os.FileInfo.
type fileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	sys     *fileStat
}

func (i *fileInfo) Name() string {
	return i.name
}

func (i *fileInfo) Size() int64 {
	return i.size
}

func (i *fileInfo) Mode() os.FileMode {
	return i.mode
}

func (i *fileInfo) ModTime() time.Time {
	return i.modTime
}

func (i *fileInfo) IsDir() bool {
	return i.mode.IsDir()
}

func (i *fileInfo) Sys() interface{} {
	return i.sys
}

// NewFileInfo attempt to fetch and convert the handle
// into an os.FileInfo interface.
func NewFileInfo(h binding.Handle) os.FileInfo {
	if h.IsNil() {
		return nil
	}
	result := &fileInfo{
		name:    h.FileInfoName(),
		size:    int64(h.FileInfoSize()),
		mode:    os.FileMode(h.FileInfoMode()),
		modTime: time.Unix(0, h.FileInfoModTime()),
	}
	sys := h.FileInfoSys()
	defer sys.Free()
	var stat fileStat
	buf := (*(*[unsafe.Sizeof(fileStat{})]byte)(
		unsafe.Pointer(&stat)))[:]
	copy(buf, sys.Bytes())
	result.sys = &stat
	return result
}

// file specifies the behaviour exhibited by those objects
// implementing api.File interface.
type file struct {
	Closer
	Reader
	ReaderAt
	Writer
	WriterAt
	Seeker
	file binding.Handle
}

func (f *file) Stat() (os.FileInfo, error) {
	info, err := f.file.FileStat()
	if err != nil {
		return nil, err
	}
	defer info.Free()
	return NewFileInfo(info), nil
}

// NewFile creates a file handle specified for operations.
//
// The resource represented by the handle will be transferred
// to the underlying file object.
func NewFile(h binding.Handle) api.File {
	f := &file{file: h}
	f.Closer = NewCloser(&f.file)
	f.Reader = NewReader(&f.file)
	f.ReaderAt = NewReaderAt(&f.file)
	f.Writer = NewWriter(&f.file)
	f.WriterAt = NewWriterAt(&f.file)
	f.Seeker = NewSeeker(&f.file)
	return f
}

// FileSystem specifies the behaviour exhibited by those
// objects implementing api.FileSystem interface.
type FileSystem struct {
	h *binding.Handle
}

func (fs *FileSystem) Open(path string) (api.File, error) {
	f, err := fs.h.Open(path)
	if err != nil {
		return nil, err
	}
	return NewFile(f), nil
}

func (fs *FileSystem) Stat(path string) (os.FileInfo, error) {
	info, err := fs.h.Stat(path)
	if err != nil {
		return nil, err
	}
	defer info.Free()
	return NewFileInfo(info), nil
}

func (fs *FileSystem) Lstat(path string) (os.FileInfo, error) {
	info, err := fs.h.Lstat(path)
	if err != nil {
		return nil, err
	}
	defer info.Free()
	return NewFileInfo(info), nil
}

func (fs *FileSystem) Readlink(path string) (string, error) {
	return fs.h.Readlink(path)
}

func (fs *FileSystem) EvalSymlink(path string) (string, error) {
	return fs.h.EvalSymlink(path)
}

func (fs *FileSystem) Readdir(path string) ([]os.FileInfo, error) {
	fileInfos, err := fs.h.Readdir(path)
	if err != nil {
		return nil, err
	}
	defer fileInfos.Free()
	var result []os.FileInfo
	length := fileInfos.Length()
	for i := 0; i < length; i++ {
		func() {
			fileInfo := fileInfos.Index(i)
			defer fileInfo.Free()
			result = append(result, NewFileInfo(fileInfo))
		}()
	}
	return result, nil
}

func (fs *FileSystem) Walk(root string, f filepath.WalkFunc) error {
	return fs.h.Walk(root, binding.WalkFunc(func(
		name string, info binding.Handle, err error,
	) error {
		return f(name, NewFileInfo(info), err)
	}))
}

func NewFileSystem(h *binding.Handle) FileSystem {
	return FileSystem{h: h}
}
