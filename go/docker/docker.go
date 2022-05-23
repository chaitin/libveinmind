// Package docker is the API implementation on docker.
package docker

import (
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/pkg/behaviour"
	"github.com/chaitin/libveinmind/go/pkg/binding"
)

// newArgs is the internal state that a docker.NewOption can
// manipulate for creating a new docker handle.
type newArgs struct {
	h binding.Handle
}

// NewOption is the option that can be used for initializing an
// docker.Docker object.
type NewOption func(*newArgs)

// WithConfigPath specifies the path of dockerd's config file.
//
// Specifying this argument is semantically equivalent to specifying
// flag "--config-file" to dockerd, and its default search path is
// "/etc/docker/daemon.json".
//
// Both dockerd and veinmind will render "/etc/docker/daemon.json"
// file as dispensible and fallback to use internal default config
// if unspecified. But once the argument is specified, it is no
// longer dispensible and error will be raised if the config is
// not found.
func WithConfigPath(path string) NewOption {
	return func(opts *newArgs) {
		opt := binding.DockerWithConfigPath(path)
		defer opt.Free()
		opts.h.Append(opt)
	}
}

// WithDataRootDir specifies the path of dockerd's data directory.
//
// Specifying this argument is semantically equivalent to specifying
// flag "--data-root" to dockerd, and is default value is
// "/var/lib/docker".
func WithDataRootDir(path string) NewOption {
	return func(opts *newArgs) {
		opt := binding.DockerWithDataRootDir(path)
		defer opt.Free()
		opts.h.Append(opt)
	}
}

// WithUniqueDesc specifies the unique descriptor of dockerd.
//
// This argument must be result of docker.(*Docker).UniqueDesc()
// from another docker.Docker instance, potentially from another
// process. And the initialization might still fail if the API
// runtime context has not been set up properly.
func WithUniqueDesc(desc string) NewOption {
	return func(opts *newArgs) {
		opt := binding.DockerWithUniqueDesc(desc)
		defer opt.Free()
		opts.h.Append(opt)
	}
}

// Docker is the connection established with a docker runtime.
type Docker struct {
	behaviour.Closer
	behaviour.Runtime
	behaviour.FileSystem
	runtime binding.Handle
}

// New a docker runtime object.
func New(opts ...NewOption) (api.Runtime, error) {
	hopt := binding.DockerMakeNewOptionList()
	defer hopt.Free()
	args := &newArgs{h: hopt}
	for _, opt := range opts {
		opt(args)
	}
	h, err := binding.DockerNew(hopt)
	if err != nil {
		return nil, err
	}
	result := &Docker{runtime: h}
	result.Closer = behaviour.NewCloser(&result.runtime)
	result.Runtime = behaviour.NewRuntime(&result.runtime)
	result.FileSystem = behaviour.NewFileSystem(&result.runtime)
	return result, nil
}

// UniqueDesc represents the docker runtime's initialization
// arguments, which can be passed across process boundaries and
// initialize the same docker in another process.
func (d *Docker) UniqueDesc() string {
	return d.runtime.DockerUniqueDesc()
}

// Image represents a docker image, which is guaranteed to be
// the result of docker.Docker.OpenImageByID.
type Image struct {
	behaviour.Closer
	behaviour.Image
	behaviour.FileSystem
	runtime *Docker
	image   binding.Handle
}

func (d *Docker) OpenImageByID(id string) (api.Image, error) {
	h, err := d.runtime.RuntimeOpenImageByID(id)
	if err != nil {
		return nil, err
	}
	result := &Image{runtime: d, image: h}
	result.Closer = behaviour.NewCloser(&result.image)
	result.Image = behaviour.NewImage(&result.image)
	result.FileSystem = behaviour.NewFileSystem(&result.image)
	return result, nil
}

func (i *Image) Runtime() *Docker {
	return i.runtime
}

func (i *Image) NumLayers() int {
	return i.image.DockerImageNumLayers()
}

func (im *Image) GetLayerDiffID(i int) (string, error) {
	return im.image.DockerImageGetLayerDiffID(i)
}

// Layer represents a containerd layer, which is guaranteed to
// be the result of docker.Image.OpenLayer.
type Layer struct {
	behaviour.FileSystem
	image *Image
	layer binding.Handle
}

func (im *Image) OpenLayer(i int) (*Layer, error) {
	l, err := im.image.DockerImageOpenLayer(i)
	if err != nil {
		return nil, err
	}
	result := &Layer{image: im, layer: l}
	result.FileSystem = behaviour.NewFileSystem(&result.layer)
	return result, nil
}

func (l *Layer) Image() *Image {
	return l.image
}

func (l *Layer) ID() string {
	return l.layer.DockerLayerID()
}
