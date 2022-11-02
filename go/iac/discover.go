package iac

import (
	"io/fs"
	"os"
	"path/filepath"
)

// discoverOption is the internal data to modify the way
// of discovering iacs.
type discoverOption struct {
}

// DiscoverOption specifies how to find and validate iacs.
type DiscoverOption func(*discoverOption)

// newDiscoverOption creates the discover option object.
func newDiscoverOption(opts ...DiscoverOption) *discoverOption {
	result := &discoverOption{}
	for _, opt := range opts {
		opt(result)
	}
	return result
}

func DiscoverIACs(root string, opts ...DiscoverOption) ([]IAC, error) {
	_ = newDiscoverOption(opts...)
	var results []IAC

	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		validators.Range(func(key, value interface{}) bool {
			t := key.(IACType)

			if v, ok := value.(Validator); ok {
				if v.Validate(path, info) {
					results = append(results, IAC{
						Path: path,
						Type: t,
					})
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
