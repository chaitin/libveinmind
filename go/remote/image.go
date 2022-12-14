package remote

import (
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/pkg/behaviour"
	"github.com/chaitin/libveinmind/go/pkg/binding"
)

// Image represents a remote image.
type Image struct {
	behaviour.Closer
	behaviour.Image
	behaviour.FileSystem
	image   binding.Handle
	runtime *Runtime
}

type Layer struct {
	behaviour.Closer
	behaviour.FileSystem
	layer binding.Handle
	image *Image
}

func (i *Image) Runtime() *Runtime {
	return i.runtime
}

func (i *Image) NumLayers() int {
	return i.image.RemoteImageNumLayers()
}

func (i *Image) OpenLayer(index int) (api.Layer, error) {
	h, err := i.image.RemoteImageOpenLayer(index)
	if err != nil {
		return nil, err
	}

	return &Layer{
		Closer:     behaviour.NewCloser(&h),
		FileSystem: behaviour.NewFileSystem(&h),
		layer:      h,
		image:      i,
	}, nil
}

func (l *Layer) ID() string {
	return l.layer.RemoteLayerId()
}

func (l *Layer) Image() *Image {
	return l.image
}
