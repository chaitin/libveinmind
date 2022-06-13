package plugin

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"
)

// DiscoverHandler is the handler for plugin errors in discover.
type DiscoverHandler func(*Plugin, error) error

func ignoreDiscoverHandler(_ *Plugin, _ error) error {
	return nil
}

// DefaultDiscoverHandler just simply ignore all errors.
var DefaultDiscoverHandler = DiscoverHandler(ignoreDiscoverHandler)

// discoverOption is the internal data to modify the way
// of discovering plugins.
type discoverOption struct {
	parallelism int
	exec        Executor
	pattern     string
	errHandler  DiscoverHandler
}

// DiscoverOption specifies how to find and validate plugins.
type DiscoverOption func(*discoverOption)

// WithExecutor specifies that the plugin will be executed
// with the specified option.
func WithExecutor(exec Executor) DiscoverOption {
	return func(p *discoverOption) {
		p.exec = exec
	}
}

// WithGlob specifies the pattern of file to find.
//
// While using, the path for matching will be first trimmed
// relative to the traversal root.
//
// This option is ignored when used with DiscoverPlugin, since
// it has already specified a file to validate.
func WithGlob(pattern string) DiscoverOption {
	return func(p *discoverOption) {
		p.pattern = pattern
	}
}

// WithDiscoverHandler specifies the handler of error.
//
// When unspecified, those plugin with error would be ignored,
// and only successfully discovered ones will be returned.
//
// This option is ignored when used with NewPlugin, since
// it should return the error when plugin object is not returned.
func WithDiscoverHandler(f DiscoverHandler) DiscoverOption {
	return func(p *discoverOption) {
		p.errHandler = f
	}
}

// WithDiscoverParallelism specifies how many info commands will
// be executed in parallel.
//
// Leaving parallelism unspecified or setting it to 0 will cause
// up to runtime.GOMAXPROCS(0) plugins to be executed in parallel.
// Setting it to 1 will disable parallel execution.
//
// This option is ignored when calling NewPlugin directly.
func WithDiscoverParallelism(n int) DiscoverOption {
	return func(p *discoverOption) {
		p.parallelism = n
	}
}

// newDiscoverOption creates the discover option object.
func newDiscoverOption(opts ...DiscoverOption) *discoverOption {
	result := &discoverOption{
		exec:       DefaultExecutor,
		errHandler: DefaultDiscoverHandler,
	}
	for _, opt := range opts {
		opt(result)
	}
	return result
}

// fillPlugin attempt to fill the information of plugin.
func (opt *discoverOption) fillPlugin(plug *Plugin) {
	plug.executor = opt.exec
}

// discover is the internal method of plugin to perform discover.
func (plug *Plugin) discover(ctx context.Context) error {
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}
	defer func() {
		_ = r.Close()
		_ = w.Close()
	}()
	grp, errCtx := errgroup.WithContext(ctx)
	grp.Go(func() error {
		defer func() { _ = w.Close() }()
		return plug.exec(errCtx, []string{"info"}, &os.ProcAttr{
			Files: []*os.File{
				nil, w,
			},
		})
	})
	grp.Go(func() error {
		defer func() { _ = r.Close() }()
		decoder := json.NewDecoder(r)
		return decoder.Decode(&plug.Manifest)
	})
	if err := grp.Wait(); err != nil {
		return err
	}
	if plug.Manifest.ManifestVersion != CurrentManifestVersion {
		return xerrors.New("incompatible plugin")
	}
	return nil
}

// NewPlugin attempt to create and verify a plugin.
func NewPlugin(
	ctx context.Context, path string, opts ...DiscoverOption,
) (*Plugin, error) {
	plug := &Plugin{path: path}
	option := newDiscoverOption(opts...)
	option.fillPlugin(plug)
	if err := plug.discover(ctx); err != nil {
		return nil, err
	}
	return plug, nil
}

// DiscoverPlugins discover plugins by a traversal under
// specified directory recursively.
func DiscoverPlugins(
	ctx context.Context, root string, opts ...DiscoverOption,
) ([]*Plugin, error) {
	option := newDiscoverOption(opts...)
	n := option.parallelism
	if n <= 0 {
		n = runtime.GOMAXPROCS(0)
	}
	if n <= 0 {
		n = 1
	}
	results := make([][]*Plugin, n)
	discoverCh := make(chan *Plugin)
	grp, errCtx := errgroup.WithContext(ctx)
	for i := 0; i < n; i++ {
		j := i
		grp.Go(func() error {
			for {
				select {
				case <-errCtx.Done():
					return nil
				case plug, ok := <-discoverCh:
					if !ok {
						return nil
					}
					if err := plug.discover(errCtx); err != nil {
						err = option.errHandler(plug, err)
						if err != nil {
							return err
						}
					} else {
						results[j] = append(results[j], plug)
					}
				}
			}
		})
	}
	grp.Go(func() error {
		defer close(discoverCh)
		return filepath.Walk(root, func(
			path string, info os.FileInfo, err error,
		) error {
			// Ignore all errors while walking the file path.
			if err != nil {
				return nil
			}
			if option.pattern != "" {
				rel, err := filepath.Rel(root, path)
				if err != nil {
					return err
				}
				ok, err := filepath.Match(
					option.pattern, filepath.Join("/", rel))
				if err != nil {
					return err
				}
				if !ok {
					return nil
				}
			}
			mode := info.Mode()
			if !mode.IsRegular() || ((mode.Perm() & 0550) != 0550) {
				return nil
			}
			plug := &Plugin{path: path}
			option.fillPlugin(plug)
			select {
			case <-errCtx.Done():
				return nil
			case discoverCh <- plug:
			}
			return nil
		})
	})
	if err := grp.Wait(); err != nil {
		return nil, err
	}
	var result []*Plugin
	for i := 0; i < n; i++ {
		result = append(result, results[i]...)
	}
	return result, nil
}
