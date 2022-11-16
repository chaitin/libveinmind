// Package kubernetes is the API implementation on kubernetes.
package kubernetes

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"

	api "github.com/chaitin/libveinmind/go"
)

type Kubernetes struct {
	// kubernetes cluster namespace
	namespace string

	// kubeConfigPath path for cluster
	kubeConfigPath string

	// kubeConfigBytes for cluster
	kubeConfigBytes []byte

	// dynamicClient reference dynamic.Interface
	// return data use map[string]interface{} format
	dynamicClient dynamic.Interface

	// restMapper reference mete.RESTMapper
	// used to fetch schema.GroupVersionResource from kind
	restMapper meta.RESTMapper

	// inCluster dedicate whether kubernetes client in cluster
	inCluster bool
}

type NewOption func(kubernetes *Kubernetes) error

func WithNamespace(namespace string) NewOption {
	return func(kubernetes *Kubernetes) error {
		kubernetes.namespace = namespace
		return nil
	}
}

func WithKubeConfigPath(path string) NewOption {
	return func(kubernetes *Kubernetes) error {
		kubernetes.kubeConfigPath = path
		return nil
	}
}

func WithKubeConfigBytes(config []byte) NewOption {
	return func(kubernetes *Kubernetes) error {
		kubernetes.kubeConfigBytes = config
		return nil
	}
}

func WithInCluster() NewOption {
	return func(kubernetes *Kubernetes) error {
		kubernetes.inCluster = true
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

	var (
		restConfig *rest.Config
		err        error
	)

	// init namespace
	if k.namespace == "" {
		k.namespace = "default"
	}

	// init rest config
	if k.inCluster {
		restConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		if k.kubeConfigPath == "" {
			if os.Getenv("KUBECONFIG") != "" {
				k.kubeConfigPath = os.Getenv("KUBECONFIG")
			}
		}

		if k.kubeConfigBytes != nil && len(k.kubeConfigBytes) > 0 {
			restConfig, err = clientcmd.RESTConfigFromKubeConfig(k.kubeConfigBytes)
			if err != nil {
				return nil, errors.Wrap(err, "kubernetes: can't get rest config")
			}
		} else if k.kubeConfigPath != "" {
			config := genericclioptions.NewConfigFlags(true)
			*config.KubeConfig = k.kubeConfigPath
			configLoader := config.ToRawKubeConfigLoader()
			restConfig, err = configLoader.ClientConfig()
			if err != nil {
				return nil, errors.Wrap(err, "kubernetes: can't get rest config")
			}
		} else {
			return nil, errors.New("kubernetes: can'f find kube config path or bytes")
		}
	}

	// init rest mapper
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	grs, err := restmapper.GetAPIGroupResources(clientset.Discovery())
	if err != nil {
		return nil, err
	}

	k.restMapper = restmapper.NewDiscoveryRESTMapper(grs)

	// init dynamic client
	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "kubernetes: can't init dynamic client")
	}
	k.dynamicClient = dynamicClient

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
	return k.kubeConfigPath
}

func (k *Kubernetes) ConfigBytes() []byte {
	return k.kubeConfigBytes
}

func (k *Kubernetes) InCluster() bool {
	return k.inCluster
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
