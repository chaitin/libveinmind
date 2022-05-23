package api

import (
	"io"
	"os"
	"path/filepath"
)

// File abstracts an open file from container. Some behaviours
// might be masked due to potential incompatibility.
type File interface {
	io.ReadWriteCloser
	io.ReaderAt
	io.WriterAt
	io.Seeker

	Stat() (os.FileInfo, error)
}

// FileSystem abstracts the property of an object to visit its
// internal file system structure, which is usually prepared
// by the container runtime.
type FileSystem interface {
	Open(path string) (File, error)
	Stat(path string) (os.FileInfo, error)
	Lstat(path string) (os.FileInfo, error)
	Readlink(path string) (string, error)
	EvalSymlink(path string) (string, error)
	Readdir(path string) ([]os.FileInfo, error)
	Walk(root string, walkFn filepath.WalkFunc) error
}
