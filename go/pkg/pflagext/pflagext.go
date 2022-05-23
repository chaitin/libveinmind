// Package pflagext is the extension of pflag which is part
// of the cobra framework.
package pflagext

import (
	"strconv"

	"github.com/spf13/pflag"
)

type stringFuncValue func(string) error

func (f stringFuncValue) Set(sval string) error {
	return f(sval)
}

func (stringFuncValue) String() string {
	return ""
}

func (stringFuncValue) Type() string {
	return "string"
}

func StringVarFP(
	fset *pflag.FlagSet, f func(string) error,
	name, shorthand, usage string,
) {
	fset.VarP(stringFuncValue(f), name, shorthand, usage)
}

func StringVarF(
	fset *pflag.FlagSet, f func(string) error,
	name, usage string,
) {
	StringVarFP(fset, f, name, "", usage)
}

type intFuncValue func(int) error

func (f intFuncValue) Set(sval string) error {
	val, err := strconv.Atoi(sval)
	if err != nil {
		return err
	}
	return f(val)
}

func (intFuncValue) String() string {
	return ""
}

func (intFuncValue) Type() string {
	return "int"
}

func IntVarFP(
	fset *pflag.FlagSet, f func(int) error,
	name, shorthand string, usage string,
) {
	fset.VarP(intFuncValue(f), name, shorthand, usage)
}

func IntVarF(
	fset *pflag.FlagSet, f func(int) error,
	name string, usage string,
) {
	IntVarFP(fset, f, name, "", usage)
}
