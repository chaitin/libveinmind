package iac

import (
	"io/fs"
	"os"
	"path/filepath"
)

// discoverOption is the internal data to modify the way
// of discovering iacs.
type discoverOption struct {
	iacType IACType
}

// DiscoverOption specifies how to find and validate iacs.
type DiscoverOption func(*discoverOption)

func WithIaCType(iacType IACType) DiscoverOption {
	return func(o *discoverOption) {
		o.iacType = iacType
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
		// An IaC File Must Be An Regular File
		if !info.Mode().IsRegular() {
			return nil
		}
		validators.Range(func(key, value interface{}) bool {
			if v, ok := value.(Validator); ok {
				if discoverOpt.iacType == "" || discoverOpt.iacType == v.ID() {
					if v.Validate(path, info) {
						// add opt func
						results = append(results, IAC{
							Path: path,
							Type: v.ID(),
						})
						// if pass a validator, no need to check more
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
	info, err := os.Stat(path)
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
