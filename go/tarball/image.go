package tarball

import (
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

func (i *Image) Runtime() *Tarball {
	return i.runtime
}
