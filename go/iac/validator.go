package iac

import (
	"io/fs"
)

var validators []Validator

type Validator interface {
	ID() IACType
	Validate(path string, info fs.FileInfo) bool
}

func RegisterValidator(validator Validator) {
	validators = append(validators, validator)
}
