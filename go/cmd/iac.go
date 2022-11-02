package cmd

import (
	"context"
	"os"
	"sync"

	"github.com/chaitin/libveinmind/go/iac"
	"github.com/chaitin/libveinmind/go/plugin"
)

// ScanIACs scans iac provided by iac list.
func ScanIACs(
	ctx context.Context, rang plugin.ExecRange,
	IACs []iac.IAC, opts ...plugin.ExecOption,
) error {
	iter, err := plugin.IterateTyped(rang, "iac")
	if err != nil {
		return err
	}

	uniq := sync.Map{}
	for _, i := range IACs {
		actual, loaded := uniq.LoadOrStore(i.Type, []string{i.Path})
		if loaded {
			if paths, ok := actual.([]string); ok {
				paths = append(paths, i.Path)
				uniq.Store(i.Type, paths)
			}
		}
	}

	uniq.Range(func(key, value interface{}) bool {
		t := key.(iac.IACType)
		if paths, ok := value.([]string); ok {
			if err = plugin.Exec(ctx, iter, paths,
				plugin.WithPrependArgs("--iac-type", t.String()),
				plugin.WithExecOptions(opts...)); err != nil {
				return false
			}
		}
		return true
	})

	return nil
}

// ScanIAC scans iac provided by iac.
func ScanIAC(
	ctx context.Context, rang plugin.ExecRange,
	IAC iac.IAC, opts ...plugin.ExecOption,
) error {
	return ScanIACs(ctx, rang, []iac.IAC{IAC}, opts...)
}

// IACHandler is the handler for specified iac.
type IACHandler func(*Command, iac.IAC) error

// MapIACCommand attempts to create a iac command.
func (idx *Index) MapIACCommand(
	c *Command, f IACHandler,
) *Command {
	idx.MapPluginCommand(c, "iac", struct{}{}, func(command *Command, args []string) error {
		t, err := command.Flags().GetString("iac-type")
		if err != nil {
			return err
		}

		var iacs []iac.IAC
		for _, path := range args {
			fi, err := os.Stat(path)
			if err != nil {
				continue
			}

			if fi.IsDir() {
				discovered, err := iac.DiscoverIACs(path)
				if err != nil {
					continue
				}

				iacs = append(iacs, discovered...)
			} else {
				if iac.IsIACType(t) {
					iacs = append(iacs, iac.IAC{
						Path: path,
						Type: iac.IACType(t),
					})
				} else {
					discovered, err := iac.DiscoverType(path)
					if err != nil {
						continue
					}

					if discovered == iac.Unknown {
						continue
					}

					iacs = append(iacs, iac.IAC{
						Path: path,
						Type: discovered,
					})
				}
			}
		}

		for _, i := range iacs {
			_ = f(c, i)
		}

		return nil
	})
	c.Flags().String("iac-type", "dockerfile", "dedicate iac type for iac files")
	return c
}

func MapIACCommand(c *Command, f IACHandler) *Command {
	return defaultIndex.MapIACCommand(c, f)
}
