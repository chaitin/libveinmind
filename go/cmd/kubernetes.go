package cmd

import (
	"encoding/base64"
	"strings"

	"github.com/spf13/pflag"

	"github.com/chaitin/libveinmind/go/kubernetes"
	"github.com/chaitin/libveinmind/go/pkg/pflagext"
	"github.com/chaitin/libveinmind/go/plugin"
)

type kubernetesRoot struct {
	k *kubernetes.Kubernetes
}

func (r kubernetesRoot) ID() interface{} {
	return r.k
}

func (r kubernetesRoot) Mode() string {
	return "kubernetes"
}

func (r kubernetesRoot) Options() plugin.ExecOption {
	var opts []plugin.ExecOption

	if r.k.ConfigPath() != "" {
		opts = append(opts, plugin.WithPrependArgs(
			"--kube-config-path", r.k.ConfigPath()))
	}

	if r.k.ConfigBytes() != nil {
		opts = append(opts, plugin.WithPrependArgs("--kube-config-bytes",
			base64.StdEncoding.EncodeToString(r.k.ConfigBytes())))
	}

	opts = append(opts, plugin.WithPrependArgs(
		"--in-cluster", func() string {
			if r.k.InCluster() {
				return "true"
			} else {
				return "false"
			}
		}()))
	return plugin.WithExecOptions(
		opts...,
	)
}

var kubernetesFlags []kubernetes.NewOption

type kubernetesMode struct {
}

func (kubernetesMode) Name() string {
	return "kubernetes"
}

func (kubernetesMode) AddFlags(fset *pflag.FlagSet) {
	pflagext.StringVarF(fset, func(path string) error {
		kubernetesFlags = append(kubernetesFlags,
			kubernetes.WithKubeConfigPath(path))
		return nil
	}, "kube-config-path",
		`flag "--kube-config-path" specified kube config`)
	pflagext.StringVarF(fset, func(config string) error {
		b, err := base64.StdEncoding.DecodeString(config)
		if err != nil {
			return err
		}
		kubernetesFlags = append(kubernetesFlags,
			kubernetes.WithKubeConfigBytes(b))
		return nil
	}, "kube-config-bytes",
		`flag "--kube-config-bytes" specified kube config bytes`)
	pflagext.StringVarF(fset, func(inCluster string) error {
		if strings.ToLower(inCluster) == "true" {
			kubernetesFlags = append(kubernetesFlags,
				kubernetes.WithInCluster())
		}
		return nil
	}, "in-cluster",
		`flag "--in-cluster" specified in-cluster`)
}

func (kubernetesMode) Invoke(c *Command, args []string, m ModeHandler) error {
	k, err := kubernetes.New(kubernetesFlags...)
	if err != nil {
		return err
	}
	defer func() { _ = k.Close() }()
	return m(c, args, k)
}

func init() {
	RegisterPartition(func(k *kubernetes.Kubernetes) Root {
		return kubernetesRoot{k: k}
	})
	RegisterMode(&kubernetesMode{})
}
