// Package cmd defines the concrete protocol between host and
// plugins based on libVeinMind plugin system.
//
// Each plugins by themselves are standalone programs that can
// be run individually, but they must conform to some protocol
// between them and their host program if they want to produce
// meaningful output that can be interpreted and handled.
//
// The plugins' command line layout should form a directory
// hierarchy with with an info command on the root. The info
// command reflects callable subcommands in the plugin,
// providing information and constrains of each subcommands.
package cmd

import (
	"github.com/spf13/cobra"
)

// Command is libVeinMind plugin command that is built upon
// the widely used cobra framework.
type Command = cobra.Command
