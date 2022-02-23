// Package docker is the API implementation on docker.
package docker

import (
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/pkg/behaviour"
	"github.com/chaitin/libveinmind/go/pkg/binding"
)

// Docker is the connection established with a docker runtime.
type Docker struct {
	behaviour.Closer
	behaviour.Runtime
	behaviour.FileSystem
	runtime binding.Handle
}

// New a docker runtime by parsing the configurations specified
// in "/etc/docker" directory, assuming it is defaultly installed.
func New() (api.Runtime, error) {
	h, err := binding.DockerNew()
	if err != nil {
		return nil, err
	}
	result := &Docker{runtime: h}
	result.Closer = behaviour.NewCloser(&result.runtime)
	result.Runtime = behaviour.NewRuntime(&result.runtime)
	result.FileSystem = behaviour.NewFileSystem(&result.runtime)
	return result, nil
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
