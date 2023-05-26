// Package specflags provides the flag for specifying plugin
// specific flags, so that extra arguments might be passed
// to the matched plugin.
//
// Currently, we require the caller to provide a list of
// strings that is either received through a string array
// flag of cobra, or some kind of configuration files.
//
// We matches the format as below for each visited plugins:
//   <pluginName>[:commandPath].<flagName>[=<value>]
// Which will be converted to "--flagName value" or "-f value"
// when we call specified plugins.
//
// This package is also a good example of showing how to
// utilize primitive flags from package plugin to complete
// complex tasks.
package specflags

import (
	"context"
	"path"
	"strings"

	"github.com/chaitin/libveinmind/go/plugin"
)

// WithSpecFlags supplies raw plugin specific flags to
// create an plugin.ExecOption that supplies the flags to
// each executed plugins.
func WithSpecFlags(flags []string) plugin.ExecOption {
	m := make(map[string][]string)
	for _, flag := range flags {
		if !strings.Contains(flag, ".") {
			// There must be at least one dot according
			// to our specified example.
			continue
		}
		var value []string
		if idx := strings.Index(flag, "="); idx >= 0 {
			value = append(value, flag[idx+1:])
			flag = flag[:idx]
		}
		current := strings.Index(flag, ":") + 1
		for {
			idx := strings.Index(flag[current+1:], ".")
			if idx < 0 {
				break
			}
			current = current + 1 + idx
			pattern := flag[:current]
			key := flag[current+1:]
			if len(key) == 1 {
				key = "-" + key
			} else {
				key = "--" + key
			}
			args := m[pattern]
			args = append(args, key)
			args = append(args, value...)
			m[pattern] = args
		}
	}
	return plugin.WithExecInterceptor(func(
		ctx context.Context, plug *plugin.Plugin, cmd *plugin.Command,
		next func(context.Context, ...plugin.ExecOption) error,
	) error {
		var allArgs []string
		if args, ok := m[plug.Name]; ok {
			allArgs = append(allArgs, args...)
		}
		if args, ok := m[plug.Name+":"+path.Join(cmd.Path...)]; ok {
			allArgs = append(allArgs, args...)
		}
		return next(ctx, plugin.WithPrependArgs(allArgs...))
	})
}
