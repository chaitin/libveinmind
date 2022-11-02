package iac

type IAC struct {
	Path string
	Type IACType
}

type IACType string

const (
	Unknown       IACType = "unknown"
	Dockerfile    IACType = "dockerfile"
	DockerCompose IACType = "docker-compose"
	Kubernetes    IACType = "kubernetes"
)

func (t IACType) String() string {
	return string(t)
}

func IsIACType(t string) bool {
	switch t {
	case Dockerfile.String(), DockerCompose.String(), Kubernetes.String():
		return true
	default:
		return false
	}
}
