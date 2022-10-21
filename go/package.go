// Package api defines the API outline for working with
// different container runtimes.
//
// The user should first establish some kind of connection
// as client with their desired container runtime. The
// client configuraton could be either by specifying it
// manually, or by recognizing them on the host first.
//
// After establishing a client connection with container
// runtime, the user could invoke the client API to enumerate
// containers and images by their IDs, and open one of these
// entities furtherly.
package api

import (
	"context"

	imageV1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/opencontainers/runtime-spec/specs-go"
)

// Image is the open image object from a runtime.
type Image interface {
	FileSystem

	Close() error
	ID() string

	Repos() ([]string, error)
	RepoRefs() ([]string, error)

	OCISpecV1() (*imageV1.Image, error)
}

type Container interface {
	FileSystem
	Psutil

	Close() error
	ID() string
	Name() string
	ImageID() string

	OCISpec() (*specs.Spec, error)
	OCIState() (*specs.State, error)
}

// Runtime is the connection established with a specific
// container runtime, depending on the implementation and
// container runtime internal nature.
type Runtime interface {
	Close() error

	// ListImageIDs attempt to enumerate the images by their
	// IDs managed by the container runtime, which could be
	// used to open the image.
	ListImageIDs() ([]string, error)

	// FindImageIDs attempt to match image ID by specifying
	// their human readable identifiers. It must follow the
	// following rules.
	//
	// 1. When pattern is image ID recognizable by this
	//    container runtime, it will be searched first.
	// 2. When pattern is pure hexadecimal, the digest value
	//    portion will be matched.
	// 3. When pattern is a single identifier, all images
	//    with the specified identifier will be matched.
	// 4. When pattern is a repository path, all images with
	//    the specified repository but different versions
	//    will be matched.
	// 5. When pattern is a named tagged or canonical
	//    reference, the whole portion will be matched.
	FindImageIDs(pattern string) ([]string, error)

	// OpenImageByID attempt to open a image by its ID.
	OpenImageByID(id string) (Image, error)

	// ListContainerIDs attempt to open a container by its ID.
	ListContainerIDs() ([]string, error)

	// FindContainerIDs attempt to match container ID by specifying
	// their human readable identifiers. It must follow the
	// following rules.
	FindContainerIDs(pattern string) ([]string, error)

	// OpenContainerByID attempt to open a container by its ID.
	OpenContainerByID(id string) (Container, error)
}

type ClusterResource interface {
	Close() error

	// Kind return resource kind
	Kind() string

	// Get attempt to get resource raw bytes from cluster
	Get(ctx context.Context, name string) ([]byte, error)

	// List attempts to list resources from cluster
	// return resource name for Get method
	List(ctx context.Context) ([]string, error)

	// Create attempts to create resource in cluster
	Create(ctx context.Context, resource []byte) error

	// Update attempts to update resource in cluster
	Update(ctx context.Context, resource []byte) error
}

// Cluster is the connection established with a specific
// specs cluster
type Cluster interface {
	Close() error

	// ConfigPath return config path of cluster
	ConfigPath() string

	// ListNamespaces attempt to list all namespaces in cluster
	ListNamespaces() ([]string, error)

	// CurrentNamespace return current namespace of cluster
	CurrentNamespace() string

	// InCluster return kubernetes client whether in cluster
	InCluster() bool

	// Namespace attempt to switch namespace
	Namespace(namespace string) Cluster

	// Resource attempt to open ClusterResource
	// accord schema.GroupVersionResource
	Resource(kind string) (ClusterResource, error)
}
