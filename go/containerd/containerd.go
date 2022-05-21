// Package containerd is the API implementation on containerd.
package containerd

import (
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/pkg/behaviour"
	"github.com/chaitin/libveinmind/go/pkg/binding"
)

// newArgs is the internal state that a containerd.NewOption can
// manipulate for creating a new containerd handle.
type newArgs struct {
	h binding.Handle
}

// NewOption is the option that can be used for initializing an
// containerd.Containerd object.
type NewOption func(*newArgs)

// WithConfigPath specifies the path of containerd's config file.
//
// Specifying this argument is semantically equivalent to specifying
// flag "--config" or "-c" to containerd, and its default search
// path is "/etc/containerd/config.toml".
//
// Both containerd and veinmind will render
// "/etc/containerd/config.toml" file as dispensible and fallback
// to use internal default config if unspecified. But once the
// argument is specified, it is no longer dispensible and error
// will be raised if the config is not found.
func WithConfigPath(path string) NewOption {
	return func(opts *newArgs) {
		opt := binding.ContainerdWithConfigPath(path)
		defer opt.Free()
		opts.h.Append(opt)
	}
}

// WithRootDir specifies the path of containerd's root directory.
//
// Specifying this argument is semantically equivalent to specifying
// flag "--root" to containerd, and its default value is
// "/var/lib/containerd".
func WithRootDir(path string) NewOption {
	return func(opts *newArgs) {
		opt := binding.ContainerdWithRootDir(path)
		defer opt.Free()
		opts.h.Append(opt)
	}
}

// WithUniqueDesc specifies the unique descriptor of dockerd.
//
// This argument must be result of
// containerd.(*Containerd).UniqueDesc() from another
// containerd.Containerd instance, potentially from another
// process. And the initialization might still fail if the API
// runtime context has not been set up properly.
func WithUniqueDesc(desc string) NewOption {
	return func(opts *newArgs) {
		opt := binding.ContainerdWithUniqueDesc(desc)
		defer opt.Free()
		opts.h.Append(opt)
	}
}

// Containerd is the connection established with a
// containerd runtime.
type Containerd struct {
	behaviour.Closer
	behaviour.Runtime
	behaviour.FileSystem
	runtime binding.Handle
}

// New a containerd runtime object.
func New(opts ...NewOption) (api.Runtime, error) {
	hopt := binding.ContainerdMakeNewOptionList()
	defer hopt.Free()
	args := &newArgs{h: hopt}
	for _, opt := range opts {
		opt(args)
	}
	h, err := binding.ContainerdNew(hopt)
	if err != nil {
		return nil, err
	}
	result := &Containerd{runtime: h}
	result.Closer = behaviour.NewCloser(&result.runtime)
	result.Runtime = behaviour.NewRuntime(&result.runtime)
	result.FileSystem = behaviour.NewFileSystem(&result.runtime)
	return result, nil
}

// UniqueDesc represents the containerd runtime's initialization
// arguments, which can be passed across process boundaries and
// initialize the same docker in another process.
func (d *Containerd) UniqueDesc() string {
	return d.runtime.ContainerdUniqueDesc()
}

// Image represents a containerd image, which is guaranteed to
// be the result of docker.Docker.OpenImageByID.
type Image struct {
	behaviour.Closer
	behaviour.Image
	behaviour.FileSystem
	runtime *Containerd
	image   binding.Handle
}

func (d *Containerd) OpenImageByID(id string) (api.Image, error) {
	h, err := d.runtime.RuntimeOpenImageByID(id)
	if err != nil {
		return nil, err
	}
	result := &Image{image: h}
	result.Closer = behaviour.NewCloser(&result.image)
	result.Image = behaviour.NewImage(&result.image)
	result.FileSystem = behaviour.NewFileSystem(&result.image)
	return result, nil
}

func (i *Image) Runtime() *Containerd {
	return i.runtime
}
