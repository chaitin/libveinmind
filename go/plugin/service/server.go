package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"

	"github.com/chaitin/libveinmind/go/plugin"
)

// serviceFunc is the wrapped namespace function.
type serviceFunc func(json.RawMessage) (json.RawMessage, error)

func newServiceFunc(service Service) serviceFunc {
	val := reflect.ValueOf(service)
	typ := val.Type()
	if typ.Kind() != reflect.Func {
		panic("service must be function")
	}
	numOut := typ.NumOut()
	lastError := false
	if numOut > 0 && typ.Out(numOut-1) == typeError {
		lastError = true
	}
	return func(input json.RawMessage) (_ json.RawMessage, rerr error) {
		defer func() {
			if err := recover(); err != nil {
				rerr = xerrors.Errorf(
					"panic in service: %s", err)
			}
		}()
		numIn := typ.NumIn()
		callArgs := make([]reflect.Value, numIn)
		jsonArgs := make([]interface{}, numIn)
		for i := 0; i < numIn; i++ {
			callArgs[i] = reflect.New(typ.In(i)).Elem()
			jsonArgs[i] = callArgs[i].Addr().Interface()
		}
		if err := json.Unmarshal(input, &jsonArgs); err != nil {
			return nil, err
		}
		callReply := val.Call(callArgs)
		jsonReply := make([]interface{}, numOut)
		var callErr error
		if lastError {
			if errReply := callReply[numOut-1]; !errReply.IsNil() {
				callErr = errReply.Interface().(error)
			}
			jsonReply = jsonReply[:numOut-1]
		}
		for i := 0; i < len(jsonReply); i++ {
			jsonReply[i] = callReply[i].Interface()
		}
		output, err := json.Marshal(&jsonReply)
		if err != nil {
			return nil, err
		}
		return output, callErr
	}
}

type namespace struct {
	manifest json.RawMessage
	services map[string]serviceFunc
}

// Registry records the marshaled services that will be
// provided to the plugins.
type Registry struct {
	parent     *Registry
	namespaces map[string]*namespace
}

func (r *Registry) find(ns string) *namespace {
	if r == nil {
		return nil
	}
	if obj, ok := r.namespaces[ns]; ok {
		return obj
	}
	return r.parent.find(ns)
}

// Define a new namespace in the current registry.
func (r *Registry) Define(ns string, manifest interface{}) {
	if _, ok := r.namespaces[ns]; ok {
		panic(fmt.Sprintf("conflict namespace %q", ns))
	}
	data, err := json.Marshal(manifest)
	if err != nil {
		panic(fmt.Sprintf("marshal manifest %q: %v", ns, err))
	}
	r.namespaces[ns] = &namespace{
		manifest: json.RawMessage(data),
		services: make(map[string]serviceFunc),
	}
}

// Inherit attempt to inherit and add their own services.
func (r *Registry) Inherit() *Registry {
	return &Registry{
		parent:     r,
		namespaces: make(map[string]*namespace),
	}
}

// NewRegistry creates a root registry object.
func NewRegistry() *Registry {
	return &Registry{
		namespaces: make(map[string]*namespace),
	}
}

// Services inverse the control of registering operations
// and simplifies the interface provided to user.
type Services interface {
	Add(*Registry)
}

// AddService attempt to add the service to namespace.
//
// The namespace must be defined in current registry, which
// means the caller cannot modifies namespaces in the parent
// layer of current registry.
func (r *Registry) AddService(ns, name string, svc Service) {
	n, ok := r.namespaces[ns]
	if !ok {
		panic(fmt.Sprintf("undefined namespace %q", ns))
	}
	if _, ok := n.services[name]; ok {
		panic(fmt.Sprintf("conflict service %q", name))
	}
	n.services[name] = newServiceFunc(svc)
}

// AddServices attempt to add the services to the registry.
func (r *Registry) AddServices(svcs Services) {
	svcs.Add(r)
}

type serviceServer struct {
	ctx   context.Context
	group *errgroup.Group
}

func (s *serviceServer) runMasterThread(
	r io.ReadCloser, registry *Registry,
	writerCh chan<- serviceResponse,
) error {
	defer func() { _ = r.Close() }()
	d := json.NewDecoder(r)
	for {
		var request serviceRequest
		if err := d.Decode(&request); err != nil {
			if err == io.EOF {
				err = nil
			}
			return err
		}

		// Attempt to serve the client request right now,
		// we may create new thread here when a call has
		// been issued, and the result will be sent back
		// on executor thread.
		response, err := func() (*serviceResponse, error) {
			n := registry.find(request.Namespace)
			if request.Type == serviceTypeHasNamespace {
				return &serviceResponse{
					Ok: n != nil,
				}, nil
			}
			if n == nil {
				return nil, xerrors.Errorf(
					"undefined namespace %q", request.Namespace)
			}
			switch request.Type {
			case serviceTypeCall:
				f, ok := n.services[request.Name]
				if !ok {
					return nil, xerrors.Errorf(
						"undefined service %q", request.Name)
				}
				s.group.Go(func() error {
					return s.runExecutorThread(f, request, writerCh)
				})
				return nil, nil
			case serviceTypeGetManifest:
				return &serviceResponse{
					Reply: n.manifest,
				}, nil
			case serviceTypeListServices:
				var services []string
				for name := range n.services {
					services = append(services, name)
				}
				return &serviceResponse{
					Services: services,
				}, nil
			default:
				return nil, xerrors.Errorf(
					"invalid request type %q", request.Type)
			}
		}()
		if err != nil {
			if response == nil {
				response = &serviceResponse{}
			}
			response.ErrMsg = new(string)
			*response.ErrMsg = err.Error()
		}
		if response != nil {
			response.Sequence = request.Sequence
			select {
			case <-s.ctx.Done():
				return nil
			case writerCh <- *response:
			}
		}
	}
}

func (s *serviceServer) runWriterThread(
	w io.WriteCloser, writerCh <-chan serviceResponse,
) error {
	defer func() { _ = w.Close() }()
	e := json.NewEncoder(w)
	for {
		select {
		case <-s.ctx.Done():
			return nil
		case response := <-writerCh:
			if err := e.Encode(response); err != nil {
				if err == io.EOF {
					err = nil
				}
				return err
			}
		}
	}
}

func (s *serviceServer) runExecutorThread(
	f serviceFunc, request serviceRequest,
	writerCh chan<- serviceResponse,
) error {
	var response serviceResponse
	result, err := f(request.Args)
	response.Sequence = request.Sequence
	response.Reply = result
	if err != nil {
		response.ErrMsg = new(string)
		*response.ErrMsg = err.Error()
	}
	select {
	case <-s.ctx.Done():
		return nil
	case writerCh <- response:
	}
	return nil
}

func (r *Registry) startServiceServer(
	ctx context.Context, group *errgroup.Group,
	reader io.ReadCloser, writer io.WriteCloser,
) {
	server := serviceServer{
		ctx:   ctx,
		group: group,
	}
	writerCh := make(chan serviceResponse)
	group.Go(func() error {
		return server.runWriterThread(writer, writerCh)
	})
	group.Go(func() error {
		return server.runMasterThread(reader, r, writerCh)
	})
}

type bindOption struct {
	bind BindFunc
}

// BindOption is the option that could be use when bind.
type BindOption func(*bindOption)

// BindFunc binds running service to specified file.
//
// When bind function is called, the server corresponding to
// specified registry is already running. And bind function
// should create communication pipes, redirect input and output
// of the pipes to the stream provided by the server, and
// finally invoke the next function to go on.
type BindFunc func(
	ctx context.Context, plug *plugin.Plugin, cmd *plugin.Command,
	reader io.ReadCloser, writer io.WriteCloser,
	next func(context.Context, ...plugin.ExecOption) error,
) error

func WithBindFunc(f BindFunc) BindOption {
	if f == nil {
		panic("invalid nil argument")
	}
	return func(option *bindOption) {
		option.bind = f
	}
}

// Bind the registry into a running service.
func (r *Registry) Bind(opts ...BindOption) plugin.ExecOption {
	option := newDefaultBindOption()
	for _, f := range opts {
		f(option)
	}
	return plugin.WithExecInterceptor(func(
		ctx context.Context, plug *plugin.Plugin, cmd *plugin.Command,
		next func(context.Context, ...plugin.ExecOption) error,
	) (rerr error) {
		cancelCtx, cancel := context.WithCancel(ctx)
		group, groupCtx := errgroup.WithContext(cancelCtx)
		defer func() {
			if err := group.Wait(); err != nil {
				rerr = err
			}
		}()

		// Content from the service server will be written to
		// outputWriter, and can be fetched from outputReader.
		// Vice versa for inputReader and inputWriter.
		inputReader, inputWriter := io.Pipe()
		outputReader, outputWriter := io.Pipe()
		defer func() {
			_ = inputReader.Close()
			_ = inputWriter.Close()
			_ = outputReader.Close()
			_ = outputWriter.Close()
		}()

		// Start the service server and delegate to bind function.
		defer cancel()
		r.startServiceServer(groupCtx, group, inputReader, outputWriter)
		return option.bind(groupCtx, plug, cmd, outputReader, inputWriter, next)
	})
}
