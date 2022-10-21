package cmd

import (
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
	return plugin.WithExecOptions(
		plugin.WithPrependArgs(
			"--kube-config", r.k.ConfigPath()),
		plugin.WithPrependArgs(
			"--namespace", r.k.CurrentNamespace()),
		plugin.WithPrependArgs(
			"--in-cluster", func() string {
				if r.k.InCluster() {
					return "true"
				} else {
					return "false"
				}
			}()))
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
			kubernetes.WithKubeConfig(path))
		return nil
	}, "kube-config",
		`flag "--kube-config" specified kube config`)
	pflagext.StringVarF(fset, func(namespace string) error {
		kubernetesFlags = append(kubernetesFlags,
			kubernetes.WithNamespace(namespace))
		return nil
	}, "namespace",
		`flag "--namespace" specified namespace`)
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
