package service

import (
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/chaitin/libveinmind/go/plugin"
)

func filePathFromURI(p string) string {
	var components []string
	var dir, file string
	dir = p
	for {
		dir, file = path.Split(dir)
		if file == "" {
			break
		}
		components = append(components, file)
	}
	if dir != "" {
		components = append(components, dir)
	}
	var result []string
	for i := len(components); i > 0; i-- {
		result = append(result, components[i-1])
	}
	return filepath.Join(result...)
}

func filePathToURI(p string) string {
	var components []string
	var dir, file string
	dir = p
	for {
		dir, file = filepath.Split(dir)
		if file == "" {
			break
		}
		components = append(components, file)
	}
	if dir != "" {
		components = append(components, dir)
	}
	var result []string
	for i := len(components); i > 0; i-- {
		result = append(result, components[i-1])
	}
	return path.Join(result...)
}

// WithFilePath is used when the input and output stream can
// be open by directly opening file with filepath.
func WithFilePath(p string) plugin.ExecOption {
	u := url.URL{
		Scheme: "file",
		Path:   filePathToURI(p),
	}
	return plugin.WithPrependArgs("--host", u.String())
}

// WithFilePathPair is just like WithFilePath but is used when
// the input and output stream should use different files.
func WithFilePathPair(input, output string) plugin.ExecOption {
	return plugin.WithExecOptions(
		WithFilePath(input), WithFilePath(output))
}

func openFile(url *url.URL, flag int) (io.ReadWriteCloser, error) {
	f, err := os.OpenFile(filePathFromURI(url.Path), flag, 0)
	return f, err
}

func init() {
	RegisterFileOpener("file", openFile)
}
