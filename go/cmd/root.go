package cmd

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/chaitin/libveinmind/go/plugin"
)

// Root is the search root for each commands.
//
// When executing a plugin command, we will first specify mode
// and its related flags, which creates a search root for
// entities waiting to come, then specify the identifiers for
// entities (e.g. api.Image).
type Root interface {
	Mode() string
	Options() plugin.ExecOption
}

// UniqueRoot is a root which can be uniquely identified.
//
// The identifier should be a hashable object and we will use
// sync.Map.Store to store that object as key.
type UniqueRoot interface {
	Root
	ID() interface{}
}

// Partitioner is the function that decompose concrete type
// into root object and arguments.
//
// When the partitioner is for the root type itself, the
// form of "func(Type) Root" is just ok, and no argument will
// be generated this case. When the partitioner is for some
// object beneath the root, then "func(Type) (Root, string)"
// is required. The Type will be used for registering directly.
type Partitioner interface{}

var typeString = reflect.TypeOf("")

var typeRootObject = reflect.TypeOf((*Root)(nil)).Elem()

type partitioner func(reflect.Value) (Root, []string)

func newPartitioner(p Partitioner) (reflect.Type, partitioner) {
	val := reflect.ValueOf(p)
	typ := val.Type()
	if typ.Kind() != reflect.Func {
		panic("partitioner must be a function")
	}
	if typ.NumIn() != 1 {
		panic("partitioner must have one eact input of object")
	}
	if typ.NumOut() < 1 || typ.Out(0) != typeRootObject {
		panic("partitioner must have at least one output of Root")
	}
	if typ.NumOut() >= 2 && typ.Out(1) != typeString {
		panic("partitioner must yield a string as second result")
	}
	if typ.NumOut() >= 3 {
		panic("partitioner has too many result")
	}
	return typ.In(0), func(in reflect.Value) (Root, []string) {
		out := val.Call([]reflect.Value{in})
		root := out[0].Interface().(Root)
		var ids []string
		if len(out) > 1 {
			ids = append(ids, out[1].Interface().(string))
		}
		return root, ids
	}
}

var partitioners sync.Map

func RegisterPartition(f Partitioner) {
	partitioners.Store(newPartitioner(f))
}

type partition struct {
	root Root
	ids  []string
}

// Scan attempt to partition the objects and invoke plugin
// commands by itering the exec iterator.
func Scan(
	ctx context.Context, iter plugin.ExecIterator,
	objs interface{}, opts ...plugin.ExecOption,
) error {
	objVals := reflect.ValueOf(objs)
	length := objVals.Len()
	var result sync.Map
	for i := 0; i < length; i++ {
		objVal := reflect.ValueOf(objVals.Index(i).Interface())
		objTyp := objVal.Type()
		val, ok := partitioners.Load(objTyp)
		if !ok {
			panic(fmt.Sprintf("undefined partition %q", objTyp))
		}
		root, ids := val.(partitioner)(objVal)
		var rootID interface{} = root
		if uniq, ok := root.(UniqueRoot); ok {
			rootID = uniq.ID()
		}
		val, _ = result.LoadOrStore(rootID, &partition{root: root})
		p := val.(*partition)
		p.ids = append(p.ids, ids...)
	}
	var err error
	result.Range(func(_, val interface{}) bool {
		// Reset iterator for next objects
		defer iter.Reset()

		p := val.(*partition)
		if err = plugin.Exec(ctx, iter, p.ids,
			plugin.WithPrependArgs("--mode", p.root.Mode()),
			p.root.Options(),
			plugin.WithExecOptions(opts...)); err != nil {
			return false
		}
		return true
	})
	return err
}

// ScanIDs attempts to load a root object while passing a list
// of IDs that will be acceptable.
func ScanIDs(
	ctx context.Context, iter plugin.ExecIterator,
	obj interface{}, ids []string, opts ...plugin.ExecOption,
) error {
	objVal := reflect.ValueOf(obj)
	objTyp := objVal.Type()
	val, ok := partitioners.Load(objTyp)
	if !ok {
		panic(fmt.Sprintf("undefined partition %q", objTyp))
	}
	root, pids := val.(partitioner)(objVal)
	if len(pids) > 0 {
		panic("invalid root object with ID")
	}
	if len(ids) == 0 {
		return nil
	}
	return plugin.Exec(ctx, iter, ids,
		plugin.WithPrependArgs("--mode", root.Mode()),
		root.Options(), plugin.WithExecOptions(opts...))
}
