// Package remote is the API implementation on remote format image.
package remote

import (
	"github.com/pkg/errors"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/pkg/behaviour"
	"github.com/chaitin/libveinmind/go/pkg/binding"
)

type Runtime struct {
	root string

	behaviour.Closer
	behaviour.Runtime
	behaviour.FileSystem
	runtime binding.Handle
}

func New(root string) (api.Runtime, error) {
	t := &Runtime{}

	h, err := binding.RemoteNew(t.root)
	if err != nil {
		return nil, err
	}
	t.runtime = h
	t.Closer = behaviour.NewCloser(&t.runtime)
	t.Runtime = behaviour.NewRuntime(&t.runtime)
	t.FileSystem = behaviour.NewFileSystem(&t.runtime)

	return t, nil
}

func (t *Runtime) OpenImageByID(id string) (api.Image, error) {
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

func (t *Runtime) OpenContainerByID(id string) (api.Container, error) {
	return nil, errors.New("remote: unsupported")
}

// Root return data root for remote system
func (t *Runtime) Root() string {
	return t.root
}

// Load image into remote manager system
func (t *Runtime) Load(imageRef string, opts ...LoadOption) ([]string, error) {
	options := &loadOptions{}
	for _, o := range opts {
		o(options)
	}
	return t.runtime.RemoteLoad(imageRef, options.username, options.password)
}

func (t *Runtime) Close() error {
	return t.runtime.Close()
}
