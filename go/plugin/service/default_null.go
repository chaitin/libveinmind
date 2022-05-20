//go:build !linux
// +build !linux

package service

import (
	"context"
	"io"

	"golang.org/x/xerrors"

	"github.com/chaitin/libveinmind/go/plugin"
)

func unspecifiedBind(
	ctx context.Context, plug *plugin.Plugin, cmd *plugin.Command,
	reader io.ReadCloser, writer io.WriteCloser,
	next func(context.Context, ...plugin.ExecOption) error,
) error {
	return xerrors.New("no default bind for this platform")
}

func newDefaultBindOption() *bindOption {
	return &bindOption{
		bind: unspecifiedBind,
	}
}
