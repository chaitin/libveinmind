package tarball

import (
	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/pkg/behaviour"
	"github.com/chaitin/libveinmind/go/pkg/binding"
)

// Image represents a tarball image.
type Image struct {
	behaviour.Closer
	behaviour.Image
	behaviour.FileSystem
	image   binding.Handle
	runtime *Tarball
}

type Layer struct {
	behaviour.Closer
	behaviour.FileSystem
	layer binding.Handle
	image *Image
}

func (i *Image) Runtime() *Tarball {
	return i.runtime
}

func (i *Image) NumLayers() int {
	return i.image.TarballImageNumLayers()
}

func (i *Image) OpenLayer(index int) (api.Layer, error) {
	h, err := i.image.TarballImageOpenLayer(index)
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
	return l.layer.TarballLayerId()
}

func (l *Layer) Image() *Image {
	return l.image
}

func (l *Layer) Opaques() ([]string, error) {
	return l.layer.TarballLayerOpaques()
}

func (l *Layer) Whiteouts() ([]string, error) {
	return l.layer.TarballLayerWhiteouts()
}
