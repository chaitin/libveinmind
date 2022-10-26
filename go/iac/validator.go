package iac

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"gopkg.in/yaml.v3"
)

var validators = sync.Map{}

type Validator interface {
	ID() IACType
	Validate(path string, info fs.FileInfo) bool
}

func RegisterValidator(validator Validator) {
	validators.Store(validator.ID(), validator)
}

type DockerComposeValidator struct {
}

func (d *DockerComposeValidator) ID() IACType {
	return DockerCompose
}

func (d *DockerComposeValidator) Validate(path string, info fs.FileInfo) bool {
	compose := struct {
		Version  string                 `yaml:"version"`
		Services map[string]interface{} `yaml:"services"`
		Volumes  map[string]interface{} `yaml:"volumes,omitempty"`
		Networks map[string]interface{} `yaml:"networks,omitempty"`
		Secrets  map[string]interface{} `yaml:"secrets,omitempty"`
	}{}

	ext := filepath.Ext(info.Name())
	if ext == ".yaml" || ext == ".yml" {
		// properties 1: file name is docker-compose.yml
		if strings.ToLower(info.Name()) == "docker-compose.yml" || strings.ToLower(info.Name()) == "docker-compose.yaml" {
			return true
		}

		file, err := os.Open(path)
		if err != nil {
			return false
		}
		defer file.Close()
		data, err := io.ReadAll(file)
		err = yaml.Unmarshal(data, &compose)
		if err != nil {
			return false
		}
		if compose.Version != "" && compose.Services != nil {
			return true
		}
	}
	return false
}

type DockerfileValidator struct{}

func (v *DockerfileValidator) ID() IACType {
	return Dockerfile
}

func (v *DockerfileValidator) Validate(path string, info fs.FileInfo) bool {
	// properties 1: file name is dockerfile/containerfile
	if strings.ToLower(info.Name()) == "dockerfile" || strings.ToLower(info.Name()) == "containerfile" {
		return true
	}
	// properties 2: file ext is dockerfile/containerfile
	ext := filepath.Ext(info.Name())
	if strings.ToLower(ext) == "dockerfile" || strings.ToLower(ext) == "containerfile" {
		return true
	}

	// properties 3: try parse file
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	dockerfile, err := parser.Parse(file)
	if err != nil {
		return false
	}

	for _, child := range dockerfile.AST.Children {
		_, err := instructions.ParseInstruction(child)
		if err != nil {
			return false
		}
	}

	return true
}

type KubernetesValidator struct{}

func (k *KubernetesValidator) ID() IACType {
	return Kubernetes
}

func (k *KubernetesValidator) Validate(path string, info fs.FileInfo) bool {
	compose := struct {
		ApiVersion string      `yaml:"apiVersion"`
		Path       string      `yaml:"path"`
		Kind       string      `yaml:"kind"`
		Meta       interface{} `yaml:"metadata"`
		Spec       interface{} `yaml:"spec"`
		Status     interface{} `yaml:"status"`
	}{}

	ext := filepath.Ext(info.Name())
	if ext == ".yaml" || ext == ".yml" {
		file, err := os.Open(path)
		if err != nil {
			return false
		}
		defer file.Close()
		data, err := io.ReadAll(file)
		err = yaml.Unmarshal(data, &compose)
		if err != nil {
			return false
		}
		if compose.ApiVersion != "" && compose.Kind != "" {
			return true
		}
	}
	return false
}

func init() {
	// Register validator
	RegisterValidator(&DockerfileValidator{})
	RegisterValidator(&DockerComposeValidator{})
	RegisterValidator(&KubernetesValidator{})
}
