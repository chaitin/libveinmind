package cmd

import (
	"github.com/spf13/pflag"

	"github.com/chaitin/libveinmind/go/docker"
	"github.com/chaitin/libveinmind/go/plugin"
)

type dockerRoot struct {
	d *docker.Docker
}

func (r dockerRoot) ID() interface{} {
	return r.d
}

func (r dockerRoot) Mode() string {
	return "docker"
}

func (r dockerRoot) Options() plugin.ExecOption {
	return plugin.WithExecOptions()
}

type dockerMode struct {
}

func (dockerMode) Name() string {
	return "docker"
}

func (dockerMode) AddFlags(pflag *pflag.FlagSet) {
}

func (dockerMode) Invoke(c *Command, args []string, m ModeHandler) error {
	d, err := docker.New()
	if err != nil {
		return err
	}
	defer func() { _ = d.Close() }()
	return m(c, args, d)
}

func init() {
	RegisterPartition(func(d *docker.Docker) Root {
		return dockerRoot{d: d}
	})
	RegisterPartition(func(i *docker.Image) (Root, string) {
		return dockerRoot{d: i.Runtime()}, i.ID()
	})
	RegisterMode(&dockerMode{})
}
