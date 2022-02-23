package cmd

import (
	"github.com/spf13/pflag"

	"github.com/chaitin/libveinmind/go/containerd"
	"github.com/chaitin/libveinmind/go/plugin"
)

type containerdRoot struct {
	c *containerd.Containerd
}

func (r containerdRoot) ID() interface{} {
	return r.c
}

func (r containerdRoot) Mode() string {
	return "containerd"
}

func (r containerdRoot) Options() plugin.ExecOption {
	return plugin.WithExecOptions()
}

type containerdMode struct {
}

func (containerdMode) Name() string {
	return "containerd"
}

func (containerdMode) AddFlags(pflag *pflag.FlagSet) {
}

func (containerdMode) Invoke(c *Command, args []string, m ModeHandler) error {
	r, err := containerd.New()
	if err != nil {
		return err
	}
	defer func() { _ = r.Close() }()
	return m(c, args, r)
}

func init() {
	RegisterPartition(func(c *containerd.Containerd) Root {
		return containerdRoot{c: c}
	})
	RegisterPartition(func(i *containerd.Image) (Root, string) {
		return containerdRoot{c: i.Runtime()}, i.ID()
	})
	RegisterMode(&containerdMode{})
}
