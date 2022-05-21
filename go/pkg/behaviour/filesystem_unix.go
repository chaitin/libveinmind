// +build !windows
//go:build !windows
// +build !windows

package behaviour

import "syscall"

type fileStat = syscall.Stat_t
