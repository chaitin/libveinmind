package behaviour

import (
	"encoding/json"

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
