// Package tarball is the API implementation on tarball format image.
package tarball

import (
	"github.com/pkg/errors"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/pkg/behaviour"
	"github.com/chaitin/libveinmind/go/pkg/binding"
)

// NewOption is the option that can be used for initializing an
// tarball.Tarball object.
type NewOption func(tarball *Tarball)

func WithRoot(root string) NewOption {
	return func(tarball *Tarball) {
		tarball.root = root
	}
}

type Tarball struct {
	root string

	behaviour.Closer
	behaviour.Runtime
	behaviour.FileSystem
	runtime binding.Handle
}

func New(options ...NewOption) (api.Runtime, error) {
	t := &Tarball{}
	for _, opt := range options {
		opt(t)
	}

	h, err := binding.TarballNew(t.root)
	if err != nil {
		return nil, err
	}
	t.runtime = h
	t.Closer = behaviour.NewCloser(&t.runtime)
	t.Runtime = behaviour.NewRuntime(&t.runtime)
	t.FileSystem = behaviour.NewFileSystem(&t.runtime)

	return t, nil
}

func (t *Tarball) OpenImageByID(id string) (api.Image, error) {
	h, err := t.runtime.RuntimeOpenImageByID(id)
	if err != nil {
		return nil, err
	}
	result := &Image{runtime: t, image: h}
	result.Closer = behaviour.NewCloser(&result.image)
	result.Image = behaviour.NewImage(&result.image)
	result.FileSystem = behaviour.NewFileSystem(&result.image)
	return result, nil
}

func (t *Tarball) OpenContainerByID(id string) (api.Container, error) {
	return nil, errors.New("tarball: unsupported")
}

// Root return data root for tarball system
func (t *Tarball) Root() string {
	return t.root
}

// Load image into tarball manager system
func (t *Tarball) Load(tarPath string) ([]string, error) {
	return t.runtime.TarballLoad(tarPath)
}

func (t *Tarball) RemoveImageByID(id string) error {
	return t.runtime.TarballRemoveImageByID(id)
}

func (t *Tarball) Close() error {
	return t.runtime.Close()
}
