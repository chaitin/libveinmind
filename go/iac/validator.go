package iac

import (
	"bufio"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/chaitin/libveinmind/go/pkg/vfs"
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

		file, err := vfs.Open(path)
		if err != nil {
			return false
		}
		defer file.Close()
		data, err := io.ReadAll(file)
		if err != nil {
			return false
		}
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

	// properties 3: a dockerfile must contains FROM cmd
	file, err := vfs.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	// read file first line which not start with '#'/'ARG' and check is 'FROM' or not
	// see https://docs.docker.com/engine/reference/builder/#from
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		if strings.HasPrefix(strings.TrimSpace(string(line)), "#") || strings.HasPrefix(strings.TrimSpace(string(line)), "ARG") {
			continue
		} else if strings.HasPrefix(strings.TrimSpace(string(line)), "FROM") {
			return true
		} else {
			return false
		}
	}

	return false
}

type KubernetesValidator struct{}

func (k *KubernetesValidator) ID() IACType {
	return Kubernetes
}

func (k *KubernetesValidator) Validate(path string, info fs.FileInfo) bool {
	k8s := struct {
		ApiVersion string      `yaml:"apiVersion"`
		Path       string      `yaml:"path"`
		Kind       string      `yaml:"kind"`
		Meta       interface{} `yaml:"metadata"`
		Spec       interface{} `yaml:"spec"`
		Status     interface{} `yaml:"status"`
	}{}

	ext := filepath.Ext(info.Name())
	if ext == ".yaml" || ext == ".yml" {
		file, err := vfs.Open(path)
		if err != nil {
			return false
		}
		defer file.Close()
		data, err := io.ReadAll(file)
		if err != nil {
			return false
		}
		err = yaml.Unmarshal(data, &k8s)
		if err != nil {
			return false
		}
		if k8s.ApiVersion != "" && k8s.Kind != "" {
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
