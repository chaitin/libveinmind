package cmd

import (
	"github.com/spf13/pflag"

	"github.com/chaitin/libveinmind/go/docker"
	"github.com/chaitin/libveinmind/go/pkg/pflagext"
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
	return plugin.WithExecOptions(plugin.WithPrependArgs(
		"--docker-unique-desc", r.d.UniqueDesc()))
}

var dockerFlags []docker.NewOption

type dockerMode struct {
}

func (dockerMode) Name() string {
	return "docker"
}

func (dockerMode) AddFlags(fset *pflag.FlagSet) {
	pflagext.StringVarF(fset, func(path string) error {
		dockerFlags = append(dockerFlags,
			docker.WithConfigPath(path))
		return nil
	}, "docker-config-file",
		`flag "--config-file" specified to the dockerd command`)
	pflagext.StringVarF(fset, func(path string) error {
		dockerFlags = append(dockerFlags,
			docker.WithDataRootDir(path))
		return nil
	}, "docker-data-root",
		`flag "--data-root" specified to the dockerd command`)
	pflagext.StringVarF(fset, func(desc string) error {
		dockerFlags = append(dockerFlags,
			docker.WithUniqueDesc(desc))
		return nil
	}, "docker-unique-desc",
		"unique descriptor of the docker daemon")
}

func (dockerMode) Invoke(c *Command, args []string, m ModeHandler) error {
	d, err := docker.New(dockerFlags...)
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
