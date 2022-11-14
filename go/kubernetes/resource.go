package kubernetes

import (
	"context"
	"encoding/json"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
)

type Resource struct {
	kind string

	resourceInterface dynamic.ResourceInterface
}

func (r Resource) List(ctx context.Context) ([]string, error) {
	resList, err := r.resourceInterface.List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	names := make([]string, 0)

	for _, item := range resList.Items {
		names = append(names, item.GetName())
	}

	return names, nil
}

func (r Resource) Get(ctx context.Context, name string) ([]byte, error) {
	res, err := r.resourceInterface.Get(ctx, name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return resBytes, nil
}

func (r Resource) Create(ctx context.Context, resource []byte) error {
	object := make(map[string]interface{})
	err := json.Unmarshal(resource, &object)
	if err != nil {
		return err
	}

	unstructuredObject := &unstructured.Unstructured{Object: object}

	_, err = r.resourceInterface.Create(ctx, unstructuredObject, v1.CreateOptions{})
	return err
}

func (r Resource) Update(ctx context.Context, resource []byte) error {
	object := make(map[string]interface{})
	err := json.Unmarshal(resource, &object)
	if err != nil {
		return err
	}

	unstructuredObject := &unstructured.Unstructured{Object: object}

	_, err = r.resourceInterface.Update(ctx, unstructuredObject, v1.UpdateOptions{})
	return err
}

func (r Resource) Watch(ctx context.Context) (watch.Interface, error) {
	return r.resourceInterface.Watch(ctx, v1.ListOptions{})
}

func (r Resource) Kind() string {
	return r.kind
}

func (r Resource) Close() error {
	return nil
}
