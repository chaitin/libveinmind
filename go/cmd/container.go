package cmd

import (
	"context"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/plugin"
)

// ScanAllContainers scans container provided by runtime list.
func ScanAllContainers(
	ctx context.Context, rang plugin.ExecRange,
	runtime []api.Runtime, opts ...plugin.ExecOption,
) error {
	iter, err := plugin.IterateTyped(rang, "container")
	if err != nil {
		return err
	}
	return Scan(ctx, iter, runtime, opts...)
}

// ScanContainers scans container provided by container list.
func ScanContainers(
	ctx context.Context, rang plugin.ExecRange,
	containers []api.Container, opts ...plugin.ExecOption,
) error {
	iter, err := plugin.IterateTyped(rang, "container")
	if err != nil {
		return err
	}
	return Scan(ctx, iter, containers,
		plugin.WithPrependArgs("--id"),
		plugin.WithExecOptions(opts...))
}

// ScanContainer scan an container provided.
func ScanContainer(
	ctx context.Context, rang plugin.ExecRange,
	container api.Container, opts ...plugin.ExecOption,
) error {
	return ScanContainers(ctx, rang, []api.Container{container}, opts...)
}

// ScanContainerIDs with a runtime and a list of IDs provided.
func ScanContainerIDs(
	ctx context.Context, rang plugin.ExecRange,
	runtime api.Runtime, ids []string, opts ...plugin.ExecOption,
) error {
	iter, err := plugin.IterateTyped(rang, "container")
	if err != nil {
		return err
	}
	return ScanIDs(ctx, iter, runtime, ids,
		plugin.WithPrependArgs("--id"),
		plugin.WithExecOptions(opts...))
}

// containerExactIDs specifies whether the argument list specifies
// ID instead of searchable names.
var containerExactIDs bool

// ContainerIDsHandler is the handler for current list of containers.
type ContainerIDsHandler func(*Command, api.Runtime, []string) error

// MapContainerIDsCommand attempts to create an container IDs command.
//
// The command will attempt to initialize the runtime object
// from specified mode with flags, scan and match containers in
// the runtime, and collect those qualified container IDs.
func (idx *Index) MapContainerIDsCommand(
	c *Command, f ContainerIDsHandler,
) *Command {
	c = idx.MapModeCommand(c, "container", struct{}{}, func(
		c *Command, args []string, root interface{},
	) error {
		r, ok := root.(api.Runtime)
		if !ok {
			return IncompatibleMode()
		}
		var containerIDs []string
		if len(args) == 0 {
			ids, err := r.ListContainerIDs()
			if err != nil {
				return err
			}
			containerIDs = append(containerIDs, ids...)
		} else if containerExactIDs {
			containerIDs = append(containerIDs, args...)
		} else {
			for _, arg := range args {
				ids, err := r.FindContainerIDs(arg)
				if err != nil {
					return err
				}
				containerIDs = append(containerIDs, ids...)
			}
		}
		return f(c, r, containerIDs)
	})
	flags := c.PersistentFlags()
	flags.BoolVar(&containerExactIDs, "id", false,
		"whether fully qualified ID is specified")
	return c
}

// ContainerHandler is the handler for specified containers.
type ContainerHandler func(*Command, api.Container) error

// MapContainerCommand attempts to create a container command.
//
// The command will attempt to initialize the runtime object
// from specified mode with flags, scan and match containers in
// the runtime, and open matched containers, one at once.
func (idx *Index) MapContainerCommand(
	c *Command, f ContainerHandler,
) *Command {
	return idx.MapContainerIDsCommand(c, func(
		c *Command, r api.Runtime, containerIDs []string,
	) error {
		for _, containerID := range containerIDs {
			if err := func() error {
				container, err := r.OpenContainerByID(containerID)
				if err != nil {
					return err
				}
				defer func() { _ = container.Close() }()
				return f(c, container)
			}(); err != nil {
				return err
			}
		}
		return nil
	})
}

// AddContainerCommand invokes MapContainerCommand with no return.
func (idx *Index) AddContainerCommand(
	c *Command, f ContainerHandler,
) {
	_ = idx.MapContainerCommand(c, f)
}

// AddContainerIDsCommand invokes MapContainerCommand with no return.
func (idx *Index) AddContainerIDsCommand(
	c *Command, f ContainerIDsHandler,
) {
	_ = idx.MapContainerIDsCommand(c, f)
}

// MapContainerIDsCommand issues defaultIndex.MapContainerIDsCommand.
func MapContainerIDsCommand(
	c *Command, f ContainerIDsHandler,
) *Command {
	return defaultIndex.MapContainerIDsCommand(c, f)
}

// AddContainerCommand issues defaultIndex.AddContainerIDsCommand.
func AddContainerIDsCommand(
	c *Command, f ContainerIDsHandler,
) {
	defaultIndex.AddContainerIDsCommand(c, f)
}

// MapContainerCommand issues defaultIndex.MapContainerCommand.
func MapContainerCommand(
	c *Command, f ContainerHandler,
) *Command {
	return defaultIndex.MapContainerCommand(c, f)
}

// AddContainerCommand issues defaultIndex.AddContainerCommand.
func AddContainerCommand(
	c *Command, f ContainerHandler,
) {
	defaultIndex.AddContainerCommand(c, f)
}
