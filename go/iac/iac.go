package iac

type IAC struct {
	Path string
	Type IACType
}

type IACType string

const (
	Dockerfile    IACType = "dockerfile"
	DockerCompose IACType = "docker-compose"
	Kubernetes    IACType = "kubernetes"
)

func (t IACType) String() string {
	return string(t)
}
