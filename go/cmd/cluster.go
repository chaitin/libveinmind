package cmd

import (
	"context"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/plugin"
)

// ScanClusters scans cluster list.
func ScanClusters(
	ctx context.Context, rang plugin.ExecRange,
	cluster []api.Cluster, opts ...plugin.ExecOption,
) error {
	iter, err := plugin.IterateTyped(rang, "cluster")
	if err != nil {
		return err
	}
	return Scan(ctx, iter, cluster,
		plugin.WithExecOptions(opts...))
}

// ClusterHandler is the handler for specified cluster.
type ClusterHandler func(*Command, api.Cluster) error

// MapClusterCommand issues defaultIndex.MapClusterCommand.
func MapClusterCommand(
	c *Command, f ClusterHandler,
) *Command {
	return defaultIndex.MapClusterCommand(c, f)
}

// MapClusterCommand attempts to create a cluster command.
//
// The command will attempt to initialize the cluster object
// from specified mode with flags, and then invoke the function
// specified by caller.
func (idx *Index) MapClusterCommand(
	c *Command, f ClusterHandler,
) *Command {
	return idx.MapModeCommand(c, "cluster", struct{}{}, func(
		c *Command, _ []string, root interface{},
	) error {
		r, ok := root.(api.Cluster)
		if !ok {
			return IncompatibleMode()
		}
		return f(c, r)
	})
}
