package plugin

import (
	"golang.org/x/xerrors"
)

// ExecIterator is an iterator providing items for executing, and
// thus allow user to provide the list of plugin to execute in an
// independent way while visiting.
//
// The Next() method might returns triple nil when it is wrapping
// and filtering another iterator. This is an expected condition
// and the caller should handle it.
type ExecIterator interface {
	HasNext() bool
	Next() (*Plugin, *Command, error)
	Done()
}

type pluginIterator struct {
	plugins []*Plugin
	i, j    int
}

func (p *pluginIterator) HasNext() bool {
	return p.i < len(p.plugins)
}

func (p *pluginIterator) Next() (*Plugin, *Command, error) {
	plug := p.plugins[p.i]
	var cmd *Command
	if p.j < len(plug.Commands) {
		cmd = &plug.Commands[p.j]
	}
	p.j++
	if p.j >= len(plug.Commands) {
		p.j = 0
		p.i++
	}
	if cmd == nil {
		plug = nil
	}
	return plug, cmd, nil
}

func (p *pluginIterator) Done() {
}

type filterIterator struct {
	iter ExecIterator
	typ  string
}

func (f *filterIterator) HasNext() bool {
	return f.iter.HasNext()
}

func (f *filterIterator) Next() (*Plugin, *Command, error) {
	plug, cmd, err := f.iter.Next()
	if err != nil || plug == nil || cmd == nil {
		return plug, cmd, err
	}
	if cmd.Type == f.typ {
		return plug, cmd, nil
	}
	return nil, nil, nil
}

func (f *filterIterator) Done() {
	f.iter.Done()
}

// ExecRange specifies a range of plugins to execute.
//
// The type itself is defined as a general interface, but actually
// accepts only the following types:
//
//   - *Plugin
//   - []*Plugin
//   - ExecIterator
//   - func() ExecIterator
//   - func() (ExecIterator, error)
//   - interface{ All() ExecIterator }
//   - interface{ All() (ExecIterator, error) }
//   - interface{ ExecIterator Typed(string) } (IterateTyped)
//
// All these types can be unified into ExecIterator by calling
// NewExecIterator so the caller don't have to concern about the
// difference of acceptable types.
type ExecRange interface{}

// IterateAll iterates all executable functions in the command.
func IterateAll(rang ExecRange) (ExecIterator, error) {
	if iter, ok := rang.(ExecIterator); ok {
		return iter, nil
	}
	if reg, ok := rang.(interface{ All() ExecIterator }); ok {
		return reg.All(), nil
	}
	if reg, ok := rang.(interface{ All() (ExecIterator, error) }); ok {
		return reg.All()
	}
	switch obj := rang.(type) {
	case *Plugin:
		return &pluginIterator{
			plugins: []*Plugin{obj},
		}, nil
	case []*Plugin:
		return &pluginIterator{
			plugins: obj,
		}, nil
	case func() ExecIterator:
		return obj(), nil
	case func() (ExecIterator, error):
		return obj()
	default:
		return nil, xerrors.Errorf("unknown exec range %T", rang)
	}
}

// IterateTyped iterates for commands with type in the function.
func IterateTyped(rang ExecRange, typ string) (ExecIterator, error) {
	if reg, ok := rang.(interface{ Typed(string) ExecIterator }); ok {
		return reg.Typed(typ), nil
	}
	iter, err := IterateAll(rang)
	if err != nil {
		return nil, err
	}
	return &filterIterator{iter: iter, typ: typ}, nil
}
