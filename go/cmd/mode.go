package cmd

import (
	"fmt"
	"sync"

	"github.com/spf13/pflag"
	"golang.org/x/xerrors"
)

// modes is the registry of all modes.
var modes sync.Map

// ModeHandler is the handler for the selected mode.
type ModeHandler func(*Command, []string, interface{}) error

// Mode is a delegated object that initializes root object(s),
// manage their lifecycle(s) and pass it to subcommand.
//
// Mode are arranged on the same level, being selected by
// "--mode" flag and detail specified through flags. But
// different modes might possesses completely different root
// object provision traits.
type Mode interface {
	Name() string
	AddFlags(*pflag.FlagSet)
	Invoke(*Command, []string, ModeHandler) error
}

// RegisterMode to be used in subcommands.
func RegisterMode(mode Mode) {
	modes.Store(mode.Name(), mode)
}

// modeName is the flag that is used for selecting mode.
var modeName string

// IncompatibleMode is the helper for reporting the user
// specified mode is not compatible with the command.
func IncompatibleMode() error {
	return xerrors.Errorf("incompatible mode %q", modeName)
}

// modeFlag is the helper flag allowing user to specify the
// name of each mode as flag directly.
type modeFlag string

func (m modeFlag) String() string {
	return ""
}

func (m modeFlag) Set(_ string) error {
	modeName = string(m)
	return nil
}

func (m modeFlag) Type() string {
	return ""
}

// MapModeCommand attempts to create a mode command.
//
// The command will lookup the specified mode, initialize
// root objects, then invoke the function provided by caller.
func (idx *Index) MapModeCommand(
	c *Command, typ string, obj interface{}, f ModeHandler,
) *Command {
	c = idx.MapPluginCommand(c, typ, obj, func(
		c *Command, args []string,
	) error {
		mode, ok := modes.Load(modeName)
		if !ok {
			return xerrors.Errorf("unknown mode %q", modeName)
		}
		return mode.(Mode).Invoke(c, args, f)
	})
	flags := c.PersistentFlags()
	// TODO: create a mode "cognitive" that automatically
	// discovers and recognizes containers.
	flags.StringVarP(&modeName, "mode", "m", "docker",
		"select mode to retrieve root object")
	modes.Range(func(key, value interface{}) bool {
		name := key.(string)
		mode := value.(Mode)
		mode.AddFlags(flags)
		flag := flags.VarPF(modeFlag(name), name, "",
			fmt.Sprintf("specify %q as the mode in use", name))
		flag.NoOptDefVal = "true"
		return true
	})
	return c
}

// AddModeCommand invokes MapModeCommand with no return.
func (idx *Index) AddModeCommand(
	c *Command, typ string, obj interface{}, f ModeHandler,
) {
	_ = idx.MapModeCommand(c, typ, obj, f)
}

// MapModeCommand issues defaultIndex.MapModeCommand.
func MapModeCommand(
	c *Command, typ string, obj interface{}, f ModeHandler,
) *Command {
	return defaultIndex.MapModeCommand(c, typ, obj, f)
}

// MapModeCommand issues defaultIndex.AddModeCommand.
func AddModeCommand(
	c *Command, typ string, obj interface{}, f ModeHandler,
) {
	defaultIndex.AddModeCommand(c, typ, obj, f)
}
