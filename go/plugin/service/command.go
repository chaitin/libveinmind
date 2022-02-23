package service

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"sync"

	"github.com/spf13/pflag"
	"golang.org/x/xerrors"
)

// hostFiles stores the URL of host communication file.
//
// Either one or two of URL might be specified. They will be
// open as input and output stream and then the plugin will
// attempt to communicate with host on it.
var hostFiles []string

// AddHostFlags attempts to add flags for opening host file.
func AddHostFlags(flag *pflag.FlagSet) {
	flag.StringArrayVar(&hostFiles, "host", nil,
		"the URL of host communication file")
}

// FileOpener specifies the opener of host communication file.
type FileOpener func(url *url.URL, flag int) (io.ReadWriteCloser, error)

var openerRegistry sync.Map

// RegisterFileOpener registers the opener of file with scheme.
func RegisterFileOpener(scheme string, f FileOpener) {
	if v, ok := openerRegistry.LoadOrStore(scheme, f); ok {
		fmt.Printf("conflict opener %q: %v, %v", scheme, v, f)
	}
}

// openHostFile attempt to open host file with openers.
func openHostFile(hostFile string, flag int) (io.ReadWriteCloser, error) {
	u, err := url.Parse(hostFile)
	if err != nil {
		return nil, err
	}
	f, ok := openerRegistry.Load(u.Scheme)
	if !ok {
		return nil, xerrors.Errorf("unknown scheme %q", u.Scheme)
	}
	return (f.(FileOpener))(u, flag)
}

// openHostFiles attempt to open the files as input and output.
func openHostFiles() (io.ReadCloser, io.WriteCloser, error) {
	if len(hostFiles) == 0 {
		return nil, nil, xerrors.New("process not hosted")
	}
	if len(hostFiles) == 1 {
		f, err := openHostFile(hostFiles[0], os.O_RDWR)
		return ioutil.NopCloser(f), f, err
	}
	reader, err := openHostFile(hostFiles[0], os.O_RDONLY)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if reader != nil {
			_ = reader.Close()
		}
	}()
	writer, err := openHostFile(hostFiles[1], os.O_WRONLY)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if writer != nil {
			_ = writer.Close()
		}
	}()
	r, w := reader, writer
	reader, writer = nil, nil
	return r, w, nil
}

// Hosted returns whether current process is hosted.
//
// With this argument, each plugin process can choose to write
// to the standard output directly if it is not hosted, or to
// communicate with the host process for output write back.
func Hosted() bool {
	return len(hostFiles) > 0
}
