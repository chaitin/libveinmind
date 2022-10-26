package iac

import (
	"io/fs"
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
		// An IaC File Must Be An Regular File
		if !info.Mode().IsRegular() {
			return nil
		}
		for _, validator := range validators {
			if validator.Validate(path, info) {
				results = append(results, IAC{
					Path: path,
					Type: validator.ID(),
				})
				// if pass a validator, no need to check more
				break
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}
