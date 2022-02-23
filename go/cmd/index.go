package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/xerrors"

	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/libveinmind/go/plugin/service"
)

// Index for mapping user defined commands that is compatible
// with libVeinMind plugin system into command information.
type Index struct {
	info map[*cobra.Command]plugin.Command
}

var defaultIndex = Index{
	info: map[*cobra.Command]plugin.Command{},
}

func (idx *Index) traverseInfo(
	visited map[*cobra.Command]struct{},
	path []string, c *cobra.Command,
) []plugin.Command {
	if _, ok := visited[c]; ok {
		return nil
	}
	visited[c] = struct{}{}
	if children := c.Commands(); children != nil {
		var result []plugin.Command
		var next []string
		next = append(next, path...)
		next = append(next, "")
		for _, child := range children {
			next[len(next)-1] = child.Name()
			result = append(result, idx.traverseInfo(
				visited, next, child)...)
		}
		return result
	} else if info, ok := idx.info[c]; ok {
		current := info
		current.Path = append(current.Path, path...)
		return []plugin.Command{current}
	}
	return nil
}

// NewInfoCommand creates an info command node.
func (idx *Index) NewInfoCommand(m plugin.Manifest) *Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Describe libVeinMind plugin command entrypoints",
		RunE: func(c *cobra.Command, _ []string) error {
			if !c.HasParent() {
				return xerrors.New("missing parent command")
			}
			result := m
			result.ManifestVersion = plugin.CurrentManifestVersion
			result.Commands = idx.traverseInfo(
				make(map[*cobra.Command]struct{}), nil, c.Parent())
			data, err := json.Marshal(result)
			if err != nil {
				return err
			}
			_, err = os.Stdout.Write(data)
			return err
		},
	}
}

// NewInfoCommand issues defaultIndex.NewInfoCommand.
func NewInfoCommand(m plugin.Manifest) *Command {
	return defaultIndex.NewInfoCommand(m)
}

// PluginHandler is the handler for valid plugin command.
type PluginHandler func(*Command, []string) error

// MapPluginCommand creates a plugin command node.
//
// The function sets up and handle initialization of some
// plugin components so that the plugin command can work
// as is expected.
//
// The provided command must have its Run and RunE function
// vacant when caling the function, otherwise the function
// raises a panic.
//
// The obj specified will be marshaled into json and supplied
// to plugin.Command field. If it cannot be marshaled this
// function will also raises a panic.
//
// The result is just the same as input command, but we will
// still return that for the package document to aggregate.
func (idx *Index) MapPluginCommand(
	c *Command, typ string, obj interface{}, f PluginHandler,
) *Command {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(fmt.Sprintf(
			"cannot marshal command data: %v", err))
	}
	if c.Run != nil || c.RunE != nil {
		panic("command.Run and command.RunE must be vacant")
	}
	service.AddHostFlags(c.PersistentFlags())
	c.RunE = func(c *Command, args []string) error {
		if service.Hosted() {
			err := service.InitServiceClient(c.Context())
			if err != nil {
				return err
			}
		}
		defer log.Destroy()
		return f(c, args)
	}
	idx.info[c] = plugin.Command{
		Type: typ,
		Data: b,
	}
	return c
}

// MapPluginCommand issues defaultIndex.MapPluginCommand.
func MapPluginCommand(
	c *Command, typ string, obj interface{}, f PluginHandler,
) *Command {
	return defaultIndex.MapPluginCommand(c, typ, obj, f)
}
