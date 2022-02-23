package plugin

import (
	"context"
	"os"
	"runtime"

	"golang.org/x/sync/errgroup"
)

// ExecHandler is the handler for plugin errors in exec.
type ExecHandler func(*Plugin, *Command, error) error

func ignoreExecHandler(_ *Plugin, _ *Command, _ error) error {
	return nil
}

// DefaultExecHandler just simply ignore all errors.
var DefaultExecHandler = ExecHandler(ignoreExecHandler)

type execOption struct {
	parallelism  int
	errHandler   ExecHandler
	args         []string
	generators   []ExecGenerator
	interceptors []ExecInterceptor
}

// clone creates a copy of current options so that we can reuse
// the options by each contexts easily.
func (e *execOption) clone() *execOption {
	result := &execOption{
		parallelism: e.parallelism,
		errHandler:  e.errHandler,
	}
	result.args = append(result.args, e.args...)
	result.generators = append(result.generators, e.generators...)
	result.interceptors = append(result.interceptors, e.interceptors...)
	return result
}

// ExecOption specifies works to perform when calling a
// plugin command.
type ExecOption func(*execOption)

// WithExecOptions creates a composite of exec option.
//
// This is useful when external packages want to specify their
// own exec options and represent them as if only one option
// dedicated to their command is specified.
func WithExecOptions(opts ...ExecOption) ExecOption {
	return func(p *execOption) {
		for _, opt := range opts {
			opt(p)
		}
	}
}

// WithExecParallelism specifies how many commands will be
// executed in parallel.
//
// Leaving paralellism unspecified or setting it to 0 will cause
// up to runtime.GOMAXPROCS(0) plugins to be executed in parallel.
// Setting it to 1 will disable parallel execution.
//
// This option is ignored when calling plugin.Exec directly.
func WithExecParallelism(n int) ExecOption {
	return func(p *execOption) {
		p.parallelism = n
	}
}

// WithExecHandler specifies the handler of error.
//
// When unspecified, those plugin with error would be ignored,
// and other plugins will be execute on.
//
// Unlike WithDiscoverHandler, the option is also used even if
// Plugin.Exec is called directly.
func WithExecHandler(f ExecHandler) ExecOption {
	return func(p *execOption) {
		p.errHandler = f
	}
}

// WithPrependArgs specifies a portion of arguments to prepend
// to the arguments passed to Plugin.Exec function.
//
// This allows external packages to specifies their own
// arguments and flags in determined context.
func WithPrependArgs(args ...string) ExecOption {
	return func(p *execOption) {
		p.args = append(p.args, args...)
	}
}

// ExecGenerator is the function that generates exec options
// for each plugin and command that is visited.
type ExecGenerator func(plug *Plugin, c *Command) []ExecOption

// WithExecGenerator specifies a generator for each command.
func WithExecGenerator(generator ExecGenerator) ExecOption {
	return func(p *execOption) {
		p.generators = append(p.generators, generator)
	}
}

// ExecInterceptor is the function that intercepts how to
// execute the plugin.
//
// This function creates a responsibility chain, and it should
// invoke the next function to move on. The arguments generated
// can also be passed to the next function.
type ExecInterceptor func(
	ctx context.Context, plug *Plugin, c *Command,
	next func(context.Context, ...ExecOption) error,
) error

// WithExecInterceptor specifies a interceptor to execute.
func WithExecInterceptor(interceptor ExecInterceptor) ExecOption {
	return func(p *execOption) {
		p.interceptors = append(p.interceptors, interceptor)
	}
}

// newExecOption creates the exec option object.
func newExecOption(opts ...ExecOption) *execOption {
	result := &execOption{
		errHandler: DefaultExecHandler,
	}
	for _, opt := range opts {
		opt(result)
	}
	return result
}

func (p *execOption) exec(
	ctx context.Context, args []string, plug *Plugin, cmd *Command,
) error {
	for len(p.generators) > 0 {
		generator := p.generators[0]
		p.generators = p.generators[1:]
		for _, f := range generator(plug, cmd) {
			f(p)
		}
	}
	if len(p.interceptors) > 0 {
		interceptor := p.interceptors[0]
		p.interceptors = p.interceptors[1:]
		return interceptor(ctx, plug, cmd,
			func(ctx context.Context, opts ...ExecOption) error {
				for _, f := range opts {
					f(p)
				}
				return p.exec(ctx, args, plug, cmd)
			})
	}
	var execArgs []string
	execArgs = append(execArgs, cmd.Path...)
	execArgs = append(execArgs, p.args...)
	execArgs = append(execArgs, args...)
	return plug.exec(ctx, execArgs, &os.ProcAttr{
		Files: []*os.File{nil, nil, nil},
	})
}

type execItem struct {
	plug *Plugin
	cmd  *Command
}

// Exec the plugins specified by iterator with options.
func Exec(
	ctx context.Context, iter ExecIterator,
	args []string, opts ...ExecOption,
) error {
	option := newExecOption(opts...)
	n := option.parallelism
	if n <= 0 {
		n = runtime.GOMAXPROCS(0)
	}
	if n <= 0 {
		n = 1
	}
	grp, errCtx := errgroup.WithContext(ctx)
	execCh := make(chan execItem)
	for i := 0; i < n; i++ {
		grp.Go(func() error {
			for {
				select {
				case <-errCtx.Done():
					return nil
				case item, ok := <-execCh:
					if !ok {
						return nil
					}
					if err := option.clone().exec(ctx, args,
						item.plug, item.cmd); err != nil {
						err = option.errHandler(
							item.plug, item.cmd, err)
						if err != nil {
							return err
						}
					}
				}
			}
		})
	}
	grp.Go(func() error {
		defer iter.Done()
		defer close(execCh)
		for iter.HasNext() {
			plug, cmd, err := iter.Next()
			if err != nil {
				return err
			}
			if plug == nil || cmd == nil {
				continue
			}
			item := execItem{
				plug: plug,
				cmd:  cmd,
			}
			select {
			case <-errCtx.Done():
				return nil
			case execCh <- item:
			}
		}
		return nil
	})
	return grp.Wait()
}
