package service

import (
	"context"
	"encoding/json"
	"io"
	"reflect"
	"sync"

	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"
)

type serviceInvoke struct {
	req    serviceRequest
	reply  serviceResponse
	err    error
	doneCh chan struct{}
}

type serviceClient struct {
	ctx      context.Context
	invokeCh chan *serviceInvoke
}

func (s *serviceClient) runReaderThread(
	r io.ReadCloser, readerCh chan<- serviceResponse,
) error {
	defer func() { _ = r.Close() }()
	d := json.NewDecoder(r)
	for {
		var response serviceResponse
		if err := d.Decode(&response); err != nil {
			if err == io.EOF {
				err = nil
			}
			return err
		}
		select {
		case <-s.ctx.Done():
			return nil
		case readerCh <- response:
		}
	}
}

func (s *serviceClient) runMasterThread(
	w io.WriteCloser, readerCh <-chan serviceResponse,
) (rerr error) {
	defer func() { _ = w.Close() }()
	var current uint64
	pending := make(map[uint64]*serviceInvoke)
	defer func() {
		err := rerr
		if err == nil {
			err = xerrors.New("client closed")
		}
		for _, invoke := range pending {
			invoke.err = err
			close(invoke.doneCh)
		}
	}()
	e := json.NewEncoder(w)
	for {
		select {
		case <-s.ctx.Done():
			return nil
		case invoke := <-s.invokeCh:
			if err := func() error {
				allocated := false
				for next := current + 1; next != current; next++ {
					if _, ok := pending[next]; !ok {
						current = next
						allocated = true
						break
					}
				}
				if !allocated {
					invoke.err = xerrors.New("cannot allocate sequence")
					close(invoke.doneCh)
					return nil
				}
				invoke.req.Sequence = current
				pending[current] = invoke
				return e.Encode(invoke.req)
			}(); err != nil {
				return err
			}
		case reply := <-readerCh:
			if invoke, ok := pending[reply.Sequence]; ok {
				delete(pending, reply.Sequence)
				invoke.reply = reply
				if reply.ErrMsg != nil {
					invoke.err = xerrors.New(*reply.ErrMsg)
				}
				close(invoke.doneCh)
			}
		}
	}
}

func (s *serviceClient) call(
	ns, name string, data json.RawMessage,
) (json.RawMessage, error) {
	invoke := &serviceInvoke{
		req: serviceRequest{
			Namespace: ns,
			Type:      serviceTypeCall,
			Name:      name,
			Args:      data,
		},
		doneCh: make(chan struct{}),
	}
	select {
	case <-s.ctx.Done():
		return nil, s.ctx.Err()
	case s.invokeCh <- invoke:
		<-invoke.doneCh
		return invoke.reply.Reply, invoke.err
	}
}

func (s *serviceClient) getManifest(ns string) (json.RawMessage, error) {
	invoke := &serviceInvoke{
		req: serviceRequest{
			Namespace: ns,
			Type:      serviceTypeGetManifest,
		},
		doneCh: make(chan struct{}),
	}
	select {
	case <-s.ctx.Done():
		return nil, s.ctx.Err()
	case s.invokeCh <- invoke:
		<-invoke.doneCh
		return invoke.reply.Reply, invoke.err
	}
}

func (s *serviceClient) listServices(ns string) ([]string, error) {
	invoke := &serviceInvoke{
		req: serviceRequest{
			Namespace: ns,
			Type:      serviceTypeListServices,
		},
		doneCh: make(chan struct{}),
	}
	select {
	case <-s.ctx.Done():
		return nil, s.ctx.Err()
	case s.invokeCh <- invoke:
		<-invoke.doneCh
		return invoke.reply.Services, invoke.err
	}
}

func (s *serviceClient) hasNamespace(ns string) (bool, error) {
	invoke := &serviceInvoke{
		req: serviceRequest{
			Namespace: ns,
			Type:      serviceTypeHasNamespace,
		},
		doneCh: make(chan struct{}),
	}
	select {
	case <-s.ctx.Done():
		return false, s.ctx.Err()
	case s.invokeCh <- invoke:
		<-invoke.doneCh
		return invoke.reply.Ok, invoke.err
	}
}

func newServiceClient(ctx context.Context) (*serviceClient, error) {
	r, w, err := openHostFiles()
	if err != nil {
		return nil, err
	}
	grp, errCtx := errgroup.WithContext(ctx)
	client := &serviceClient{
		ctx:      errCtx,
		invokeCh: make(chan *serviceInvoke),
	}
	readerCh := make(chan serviceResponse)
	grp.Go(func() error {
		return client.runReaderThread(r, readerCh)
	})
	grp.Go(func() error {
		return client.runMasterThread(w, readerCh)
	})
	return client, nil
}

var (
	getOnce   sync.Once
	clientObj *serviceClient
	clientErr error
)

func getServiceClient(ctx context.Context) (*serviceClient, error) {
	getOnce.Do(func() {
		clientObj, clientErr = newServiceClient(ctx)
	})
	return clientObj, clientErr
}

// InitServiceClient attempts to initialize the service
// client when it is executed from a plugin.
//
// Though it is dispensible, it is the only chance for
// setting up context for service client, so it should be
// executed before any service has been called.
func InitServiceClient(ctx context.Context) error {
	_, err := getServiceClient(ctx)
	return err
}

// GetService will attempt to retrieve a service from host.
//
// The service argument must be a pointer to function, and a
// common usage for it might be shown as below:
//
//     var target func(x, y int) (z int, err error)
//     plugin.GetService("my-package", "add", &target)
//     target(1, 2) // returns (3, nil) on success
//
// This function panics when the provided service argument
// is not acceptable. While other errors should be returned
// as the error in provided function.
func GetService(namespace, name string, service Service) {
	val := reflect.ValueOf(service)
	panicPointerToFunc := func() {
		panic("service must be pointer to function")
	}
	if val.Type().Kind() != reflect.Ptr {
		panicPointerToFunc()
	}
	typ := val.Type().Elem()
	if typ.Kind() != reflect.Func {
		panicPointerToFunc()
	}
	numOut := typ.NumOut()
	if numOut < 0 || typ.Out(numOut-1) != typeError {
		panic("service must have error as last argument")
	}
	f := func(args []reflect.Value) []reflect.Value {
		result := make([]reflect.Value, numOut)
		for i := 0; i < numOut-1; i++ {
			result[i] = reflect.New(typ.Out(i)).Elem()
		}
		result[numOut-1] = reflect.Zero(typeError)

		if err := func() error {
			client, err := getServiceClient(context.Background())
			if err != nil {
				return err
			}
			input := make([]interface{}, typ.NumIn())
			for i := 0; i < len(input); i++ {
				input[i] = args[i].Interface()
			}
			inputData, err := json.Marshal(&input)
			if err != nil {
				return err
			}
			outputData, err := client.call(
				namespace, name, inputData)
			if err != nil {
				return err
			}
			output := make([]interface{}, numOut-1)
			for i := 0; i < len(output); i++ {
				output[i] = result[i].Addr().Interface()
			}
			return json.Unmarshal(outputData, &output)
		}(); err != nil {
			result[numOut-1] = reflect.ValueOf(err)
		}
		return result
	}
	val.Elem().Set(reflect.MakeFunc(typ, f))
}

// GetManifest attempt to retrieve manifest for namespace.
func GetManifest(namespace string, v interface{}) error {
	client, err := getServiceClient(context.Background())
	if err != nil {
		return err
	}
	data, err := client.getManifest(namespace)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// ListServices attempt to list services in a namespace.
func ListServices(namespace string) ([]string, error) {
	client, err := getServiceClient(context.Background())
	if err != nil {
		return nil, err
	}
	return client.listServices(namespace)
}

// HasNamespace attempt to verify whether namespace exists.
func HasNamespace(namespace string) (bool, error) {
	client, err := getServiceClient(context.Background())
	if err != nil {
		return false, err
	}
	return client.hasNamespace(namespace)
}
