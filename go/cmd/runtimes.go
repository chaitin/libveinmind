package cmd

import (
	"context"

	"github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/plugin"
)

// RuntimeHandler is the handler for specified runtimes.
type RuntimeHandler func(*Command, api.Runtime) error

// ScanRuntimes scans runtime provided by runtime list.
func ScanRuntimes(
	ctx context.Context, rang plugin.ExecRange,
	runtimes []api.Runtime, opts ...plugin.ExecOption,
) error {
	iter, err := plugin.IterateTyped(rang, "runtime")
	if err != nil {
		return err
	}
	return Scan(ctx, iter, runtimes, opts...)
}

// ScanRuntime scans a runtime provided.
func ScanRuntime(
	ctx context.Context, rang plugin.ExecRange,
	runtime api.Runtime, opts ...plugin.ExecOption,
) error {
	return ScanRuntimes(ctx, rang, []api.Runtime{runtime}, opts...)
}

// MapRuntimeCommand attempts to create a runtime command.
//
// The command will attempt to initialize the runtime object
// from specified mode with flags, and then invoke the function
// specified by caller.
func (idx *Index) MapRuntimeCommand(
	c *Command, f RuntimeHandler,
) *Command {
	return idx.MapModeCommand(c, "runtime", struct{}{}, func(
		c *Command, _ []string, root interface{},
	) error {
		r, ok := root.(api.Runtime)
		if !ok {
			return IncompatibleMode()
		}
		return f(c, r)
	})
}

// MapRuntimeCommand issues defaultIndex.MapRuntimeCommand.
func MapRuntimeCommand(
	c *Command, f RuntimeHandler,
) *Command {
	return defaultIndex.MapRuntimeCommand(c, f)
}
