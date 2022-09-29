package kubernetes

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"

	api "github.com/chaitin/libveinmind/go"
)

type Kubernetes struct {
	// kubernetes cluster namespace
	namespace string

	// kube kubeConfig path for cluster
	kubeConfig string

	// dynamicClient reference dynamic.Interface
	// return data use map[string]interface{} format
	dynamicClient dynamic.Interface

	// restMapper reference mete.RESTMapper
	// used to fetch schema.GroupVersionResource from kind
	restMapper meta.RESTMapper
}

type NewOption func(kubernetes *Kubernetes) error

func WithNamespace(namespace string) NewOption {
	return func(kubernetes *Kubernetes) error {
		kubernetes.namespace = namespace
		return nil
	}
}

func WithKubeConfig(path string) NewOption {
	return func(kubernetes *Kubernetes) error {
		kubernetes.kubeConfig = path
		return nil
	}
}

func New(options ...NewOption) (*Kubernetes, error) {
	k := new(Kubernetes)

	for _, opt := range options {
		err := opt(k)
		if err != nil {
			continue
		}
	}

	if k.kubeConfig == "" {
		if os.Getenv("KUBECONFIG") == "" {
			return nil, errors.New("kubernetes: can't find kube config path")
		} else {
			k.kubeConfig = os.Getenv("KUBECONFIG")
		}
	}

	if k.namespace == "" {
		k.namespace = "default"
	}

	// init dynamic client config
	config := genericclioptions.NewConfigFlags(true)
	*config.KubeConfig = k.kubeConfig
	configLoader := config.ToRawKubeConfigLoader()
	restConfig, err := configLoader.ClientConfig()
	if err != nil {
		return nil, errors.Wrap(err, "kubernetes: can't get rest config")
	}

	// init dynamic client
	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "kubernetes: can't init dynamic client")
	}
	k.dynamicClient = dynamicClient

	// init rest mapper
	mapper, err := config.ToRESTMapper()
	if err != nil {
		return nil, errors.Wrap(err, "kubernetes: can't init rest mapper")
	}
	k.restMapper = mapper

	return k, nil
}

func (k *Kubernetes) ListNamespaces() ([]string, error) {
	namespaceResource, err := k.Resource(Namespaces.String())
	if err != nil {
		return nil, err
	}

	return namespaceResource.List(context.Background())
}

func (k *Kubernetes) CurrentNamespace() string {
	return k.namespace
}

func (k *Kubernetes) ConfigPath() string {
	return k.kubeConfig
}

func (k *Kubernetes) Namespace(namespace string) api.Cluster {
	k.namespace = namespace
	return k
}

func (k *Kubernetes) Resource(kind string) (api.ClusterResource, error) {
	gvr, err := k.restMapper.ResourceFor(schema.GroupVersionResource{Resource: kind})
	if err != nil {
		return nil, err
	}

	// cluster resource can't use namespace (namespaced is false)
	if IsClusterKind(kind) {
		return Resource{kind, k.dynamicClient.Resource(gvr)}, nil
	} else if IsNamespaceKind(kind) {
		return Resource{kind, k.dynamicClient.Resource(gvr).Namespace(k.namespace)}, nil
	} else {
		return nil, errors.New("kubernetes: not support resource kind for cluster")
	}
}

func (k *Kubernetes) Close() error {
	return nil
}
