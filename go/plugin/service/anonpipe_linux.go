//go:build linux
// +build linux

package service

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/sync/errgroup"

	"github.com/chaitin/libveinmind/go/plugin"
)

// WithAnonymousPipeByProcFS attempt to create a pair of
// anonymous pipes in host process, and tell plugin processes
// to open pipes by visiting procfs mounted on specified path.
//
// This is the default option on linux, assuming the host and
// plugins are running inside the same process. However you
// cannot use this flag directly when the plugin is to be
// executed inside a container with CLONE_NEWPID set.
func WithAnonymousPipeByProcFS(procfs string) BindOption {
	return WithBindFunc(func(
		ctx context.Context, plug *plugin.Plugin, cmd *plugin.Command,
		reader io.ReadCloser, writer io.WriteCloser,
		next func(context.Context, ...plugin.ExecOption) error,
	) (rerr error) {
		group, groupCtx := errgroup.WithContext(ctx)
		defer func() {
			if err := group.Wait(); err != nil {
				rerr = err
			}
		}()
		inputReader, inputWriter, err := os.Pipe()
		if err != nil {
			return err
		}
		defer func() {
			_ = inputReader.Close()
			_ = inputWriter.Close()
			_ = reader.Close()
		}()
		outputReader, outputWriter, err := os.Pipe()
		if err != nil {
			return err
		}
		defer func() {
			_ = outputReader.Close()
			_ = outputWriter.Close()
			_ = writer.Close()
		}()
		group.Go(func() error {
			_, err := io.Copy(inputWriter, reader)
			return err
		})
		group.Go(func() error {
			_, err := io.Copy(writer, outputReader)
			return err
		})
		hostDir := filepath.Join(
			procfs, strconv.Itoa(os.Getpid()), "fd")
		readerFD := strconv.Itoa(int(inputReader.Fd()))
		writerFD := strconv.Itoa(int(outputWriter.Fd()))
		return next(groupCtx, WithFilePathPair(
			filepath.Join(hostDir, readerFD),
			filepath.Join(hostDir, writerFD)))
	})
}

// WithAnonymousPipe presumes the procfs is located on /proc, which
// is equivalent to service.WithAnonymousPipeByProcFS("/proc").
func WithAnonymousPipe() BindOption {
	return WithAnonymousPipeByProcFS("/proc")
}
