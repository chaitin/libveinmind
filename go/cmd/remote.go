package cmd

import (
	"github.com/spf13/pflag"

	"github.com/chaitin/libveinmind/go/pkg/pflagext"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/remote"
)

type remoteRoot struct {
	runtime *remote.Runtime
}

func (r remoteRoot) ID() interface{} {
	return r.runtime
}

func (r remoteRoot) Mode() string {
	return "remote"
}

func (r remoteRoot) Options() plugin.ExecOption {
	return plugin.WithExecOptions(plugin.WithPrependArgs(
		"--remote-root", r.runtime.Root()))
}

var remoteRuntimeRoot string

type remoteMode struct {
}

func (remoteMode) Name() string {
	return "remote"
}

func (remoteMode) AddFlags(fset *pflag.FlagSet) {
	pflagext.StringVarF(fset, func(root string) error {
		remoteRuntimeRoot = root
		return nil
	}, "remote-root",
		"remote manager system data root")
}

func (remoteMode) Invoke(c *Command, args []string, m ModeHandler) error {
	t, err := remote.New(remoteRuntimeRoot)
	if err != nil {
		return err
	}
	defer func() { _ = t.Close() }()
	return m(c, args, t)
}

func init() {
	RegisterPartition(func(i *remote.Image) (Root, string) {
		return remoteRoot{runtime: i.Runtime()}, i.ID()
	})
	RegisterMode(&remoteMode{})
}
