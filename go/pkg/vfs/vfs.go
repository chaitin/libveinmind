package vfs

import (
	"os"
	"path/filepath"
)

func rootfs() string {
	if v := os.Getenv("LIBVEINMIND_HOST_ROOTFS"); v != "" {
		return v
	} else {
		return ""
	}
}

func Open(name string) (*os.File, error) {
	if path, err := filepath.Abs(name); err == nil {
		name = path
	}
	return os.Open(filepath.Join(rootfs(), name))
}

func Stat(name string) (os.FileInfo, error) {
	if path, err := filepath.Abs(name); err == nil {
		name = path
	}
	return os.Stat(filepath.Join(rootfs(), name))
}

func Lstat(name string) (os.FileInfo, error) {
	if path, err := filepath.Abs(name); err == nil {
		name = path
	}
	return os.Lstat(filepath.Join(rootfs(), name))
}

func Readlink(name string) (string, error) {
	if path, err := filepath.Abs(name); err == nil {
		name = path
	}
	return os.Readlink(filepath.Join(rootfs(), name))
}

func Readdir(name string) ([]os.DirEntry, error) {
	if path, err := filepath.Abs(name); err == nil {
		name = path
	}
	return os.ReadDir(filepath.Join(rootfs(), name))
}

func Walk(root string, f filepath.WalkFunc) error {
	return filepath.Walk(filepath.Join(rootfs(), root), f)
}
