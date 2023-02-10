package iac

import (
	"io/fs"
	"path/filepath"

	"github.com/chaitin/libveinmind/go/pkg/vfs"
)

// discoverOption is the internal data to modify the way
// of discovering iacs.
type discoverOption struct {
	iacType IACType
	limit   int64
}

// DiscoverOption specifies how to find and validate iacs.
type DiscoverOption func(*discoverOption)

func WithIACType(iacType IACType) DiscoverOption {
	return func(o *discoverOption) {
		o.iacType = iacType
	}
}

func WithIACLimitSize(limit int64) DiscoverOption {
	return func(o *discoverOption) {
		o.limit = limit
	}
}

// newDiscoverOption creates the discover option object.
func newDiscoverOption(opts ...DiscoverOption) *discoverOption {
	result := &discoverOption{}
	for _, opt := range opts {
		opt(result)
	}
	return result
}

func DiscoverIACs(root string, opts ...DiscoverOption) ([]IAC, error) {
	discoverOpt := newDiscoverOption(opts...)
	var results []IAC

	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if !info.Mode().IsRegular() {
			return nil
		}
		if discoverOpt.limit > 0 && info.Size() > discoverOpt.limit {
			return nil
		}
		validators.Range(func(key, value interface{}) bool {
			if v, ok := value.(Validator); ok {
				if discoverOpt.iacType == "" || discoverOpt.iacType == v.ID() {
					if v.Validate(path, info) {
						results = append(results, IAC{
							Path: path,
							Type: v.ID(),
						})
						// if pass oen of validator, no need to check more.
						return false
					}
				}
			}
			return true
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

func DiscoverType(path string, opts ...DiscoverOption) (IACType, error) {
	_ = newDiscoverOption(opts...)
	info, err := vfs.Stat(path)
	if err != nil {
		return Unknown, err
	}

	result := Unknown
	validators.Range(func(key, value interface{}) bool {
		t := key.(IACType)

		if v, ok := value.(Validator); ok {
			if v.Validate(path, info) {
				result = t
				return false
			}
		}

		return true
	})
	return result, nil
}
