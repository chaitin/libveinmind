package cmd

import (
	"github.com/spf13/pflag"

	"github.com/chaitin/libveinmind/go/pkg/pflagext"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/tarball"
)

type tarballRoot struct {
	t *tarball.Tarball
}

func (r tarballRoot) ID() interface{} {
	return r.t
}

func (r tarballRoot) Mode() string {
	return "tarball"
}

func (r tarballRoot) Options() plugin.ExecOption {
	return plugin.WithExecOptions(plugin.WithPrependArgs(
		"--tarball-root", r.t.Root()))
}

var tarballFlags []tarball.NewOption

type tarballMode struct {
}

func (tarballMode) Name() string {
	return "tarball"
}

func (tarballMode) AddFlags(fset *pflag.FlagSet) {
	pflagext.StringVarF(fset, func(root string) error {
		tarballFlags = append(tarballFlags,
			tarball.WithRoot(root))
		return nil
	}, "tarball-root",
		"tarball manager system data root")
}

func (tarballMode) Invoke(c *Command, args []string, m ModeHandler) error {
	t, err := tarball.New(tarballFlags...)
	if err != nil {
		return err
	}
	defer func() { _ = t.Close() }()
	return m(c, args, t)
}

func init() {
	RegisterPartition(func(i *tarball.Image) (Root, string) {
		return tarballRoot{t: i.Runtime()}, i.ID()
	})
	RegisterMode(&tarballMode{})
}
