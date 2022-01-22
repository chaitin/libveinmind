// Package containerd is the API implementation on containerd.
package containerd

import (
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/pkg/behaviour"
	"github.com/chaitin/libveinmind/go/pkg/binding"
)

// Containerd is the connection established with a
// containerd runtime.
type Containerd struct {
	behaviour.Closer
	behaviour.Runtime
	behaviour.FileSystem
	runtime binding.Handle
}

// New a containerd runtime by parsing the configurations
// specified in "/var/lib/containerd" directory, assuming it is
// defaultly installed.
func New() (api.Runtime, error) {
	h, err := binding.DockerNew()
	if err != nil {
		return nil, err
	}
	result := &Containerd{runtime: h}
	result.Closer = behaviour.NewCloser(&result.runtime)
	result.Runtime = behaviour.NewRuntime(&result.runtime)
	result.FileSystem = behaviour.NewFileSystem(&result.runtime)
	return result, nil
}

// Image represents a containerd image, which is guaranteed to
// be the result of docker.Docker.OpenImageByID.
type Image struct {
	behaviour.Closer
	behaviour.Image
	behaviour.FileSystem
	image binding.Handle
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
