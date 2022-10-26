package iac

type IAC struct {
	Path string
	Type IACType
}

type IACType string

func (t IACType) String() string {
	return string(t)
}
