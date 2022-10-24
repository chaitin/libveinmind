package iac

import (
	"io/fs"
	"strings"
	"sync"
)

var validators = sync.Map{}

type Validator interface {
	ID() IACType
	Validate(path string, info fs.FileInfo) bool
}

func RegisterValidator(validator Validator) {
	validators.Store(validator.ID(), validator)
}

type DockerfileValidator struct{}

func (v *DockerfileValidator) ID() IACType {
	return Dockerfile
}

func (v *DockerfileValidator) Validate(path string, info fs.FileInfo) bool {
	if strings.ToLower(info.Name()) == "dockerfile" {
		return true
	} else {
		return false
	}
}

type DockerComposeValidator struct{}

func (v *DockerComposeValidator) ID() IACType {
	return DockerCompose
}

func (v *DockerComposeValidator) Validate(path string, info fs.FileInfo) bool {
	if strings.ToLower(info.Name()) == "docker-compose.yml" {
		return true
	} else {
		return false
	}
}

func init() {
	RegisterValidator(&DockerfileValidator{})
	RegisterValidator(&DockerComposeValidator{})
}
