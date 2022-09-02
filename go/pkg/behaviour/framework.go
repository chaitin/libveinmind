package behaviour

import (
	"encoding/json"

	"github.com/opencontainers/runtime-spec/specs-go"

	imageV1 "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/chaitin/libveinmind/go/pkg/binding"
)

type Runtime struct {
	h *binding.Handle
}

func (r *Runtime) ListImageIDs() ([]string, error) {
	return r.h.RuntimeListImageIDs()
}

func (r *Runtime) FindImageIDs(pattern string) ([]string, error) {
	return r.h.RuntimeFindImageIDs(pattern)
}

func (r *Runtime) ListContainerIDs() ([]string, error) {
	return r.h.RuntimeListContainerIDs()
}

func (r *Runtime) FindContainerIDs(pattern string) ([]string, error) {
	return r.h.RuntimeFindContainerIDs(pattern)
}

func NewRuntime(h *binding.Handle) Runtime {
	return Runtime{h: h}
}

type Image struct {
	h *binding.Handle
}

func (i *Image) ID() string {
	return i.h.ImageID()
}

func (i *Image) Repos() ([]string, error) {
	return i.h.ImageRepos()
}

func (i *Image) RepoRefs() ([]string, error) {
	return i.h.ImageRepoRefs()
}

func (i *Image) OCISpecV1() (*imageV1.Image, error) {
	bytes, err := i.h.ImageOCISpecV1MarshalJSON()
	if err != nil {
		return nil, err
	}
	result := &imageV1.Image{}
	if err := json.Unmarshal(bytes, result); err != nil {
		return nil, err
	}
	return result, nil
}

func NewImage(h *binding.Handle) Image {
	return Image{h: h}
}

type Container struct {
	h *binding.Handle
}

func (c *Container) ID() string {
	return c.h.ContainerID()
}

func (c *Container) Name() string {
	return c.h.ContainerName()
}

func (c *Container) ImageID() string {
	return c.h.ContainerImageID()
}

func (c *Container) OCISpec() (*specs.Spec, error) {
	bytes, err := c.h.ContainerOCISpecMarshalJSON()
	if err != nil {
		return nil, err
	}
	result := &specs.Spec{}
	if err := json.Unmarshal(bytes, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Container) OCIState() (*specs.State, error) {
	bytes, err := c.h.ContainerOCIStateMarshalJSON()
	if err != nil {
		return nil, err
	}
	result := &specs.State{}
	if err := json.Unmarshal(bytes, result); err != nil {
		return nil, err
	}
	return result, nil
}

func NewContainer(h *binding.Handle) Container {
	return Container{h: h}
}
