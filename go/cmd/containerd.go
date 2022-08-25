package cmd

import (
	"github.com/spf13/pflag"

	"github.com/chaitin/libveinmind/go/containerd"
	"github.com/chaitin/libveinmind/go/pkg/pflagext"
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
	return plugin.WithExecOptions(plugin.WithPrependArgs(
		"--containerd-unique-desc", r.c.UniqueDesc()))
}

var containerdFlags []containerd.NewOption

type containerdMode struct {
}

func (containerdMode) Name() string {
	return "containerd"
}

func (containerdMode) AddFlags(fset *pflag.FlagSet) {
	pflagext.StringVarF(fset, func(path string) error {
		containerdFlags = append(containerdFlags,
			containerd.WithConfigPath(path))
		return nil
	}, "containerd-config",
		`flag "--config" or "-c" specified to containerd command`)
	pflagext.StringVarF(fset, func(path string) error {
		containerdFlags = append(containerdFlags,
			containerd.WithRootDir(path))
		return nil
	}, "containerd-root",
		`flag "--root" specified to the containerd command`)
	pflagext.StringVarF(fset, func(desc string) error {
		containerdFlags = append(containerdFlags,
			containerd.WithUniqueDesc(desc))
		return nil
	}, "containerd-unique-desc",
		"unique descriptor of the containerd daemon")
}

func (containerdMode) Invoke(c *Command, args []string, m ModeHandler) error {
	r, err := containerd.New(containerdFlags...)
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
	RegisterPartition(func(c *containerd.Container) (Root, string) {
		return containerdRoot{c: c.Runtime()}, c.ID()
	})
	RegisterMode(&containerdMode{})
}
