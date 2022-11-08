// Package plugin defines the plugin system built in with the
// libveinmind SDK, allowing easy integration and composition
// of hosts and plugins.
//
// This package provides abstraction for callers of plugins,
// including discovery of plugins, definition of plugin object
// model and interface to issue plugin commands.
package plugin

import (
	"context"
	"encoding/json"
	"os"

	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"
)

// CurrentManifestVersion is the current version of manifest.
//
// The version is maintained by the SDK itself and increments
// on the schema of plugin.Manifest changes. The version will
// increment when the newer scheme has newer indispensible
// items or deleted fields, making it incompatible with older
// ones. And hosts with different manifest version will not
// be able to work together.
const CurrentManifestVersion = 1

// Manifest is the information describing the plugin itself.
type Manifest struct {
	Name        string   `json:"name,omitempty"`
	Version     string   `json:"version,omitempty"`
	Author      string   `json:"author,omitempty"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`

	// Auto generated fields that user written values will
	// be overwritten when return.
	ManifestVersion int       `json:"manifestVersion"`
	Commands        []Command `json:"commands,omitempty"`
}

// Command describes a callable subcommand with the path to
// call and its use circumstance.
type Command struct {
	Path []string        `json:"path"`
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`

	// TODO: also interpret the options and arguments.
}

// Executor will attempt to execute the plugin with
// given configuration.
//
// Sometimes it takes extra steps to execute the plugin,
// like switching to specific namespace and so on, and
// these works will be done with executor.
type Executor func(
	ctx context.Context, plugin *Plugin,
	path string, argv []string, attr *os.ProcAttr,
) error

func ExecuteStartProcessWithContext(
	ctx context.Context, _ *Plugin,
	path string, argv []string, attr *os.ProcAttr) error {
	proc, err := os.StartProcess(path, argv, attr)
	if err != nil {
		return err
	}

	wait := make(chan error, 1)
	errG, _ := errgroup.WithContext(ctx)
	errG.Go(func() error {
		state, err := proc.Wait()
		if err != nil {
			return err
		}
		if !state.Success() {
			return xerrors.New(state.String())
		}

		return nil
	})
	go func() {
		wait <- errG.Wait()
	}()

	select {
	case <-ctx.Done():
		err = proc.Kill()
		if err != nil {
			return err
		}
	case err = <-wait:
		return err
	}
	return nil
}

func executeStartProcess(
	ctx context.Context, _ *Plugin,
	path string, argv []string, attr *os.ProcAttr,
) error {
	proc, err := os.StartProcess(path, argv, attr)
	if err != nil {
		return err
	}
	state, err := proc.Wait()
	if err != nil {
		return err
	}
	if !state.Success() {
		return xerrors.New(state.String())
	}
	return nil
}

var DefaultExecutor = Executor(executeStartProcess)

// Plugin is the parsed and verified plugin, ready for
// issuing commands to it.
//
// Upon discovery, the process will attempt to call the
// "info" subcommand to verify whether the specified
// executable file is callable, and produces compatible
// manifest. How to use the plugin depends on the users.
type Plugin struct {
	Manifest

	path     string
	executor Executor
}

// exec the plugin with arguments directly.
func (plugin *Plugin) exec(
	ctx context.Context, args []string, attr *os.ProcAttr,
) error {
	var argv []string
	argv = append(argv, plugin.path)
	argv = append(argv, args...)
	return plugin.executor(ctx, plugin, plugin.path, argv, attr)
}
