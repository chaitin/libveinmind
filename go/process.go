package api

type Process interface {
	Close()
	Children() ([]Process, error)
	Cmdline() (string, error)
	Cwd() (string, error)
	Environ() ([]string, error)
	Exe() (string, error)
	Gids() ([]int32, error)
	Parent() (Process, error)
	Ppid() (int32, error)
	Pid() (int32, error)
	Uids() ([]int32, error)
	Name() (string, error)
}

type Psutil interface {
	Pids() ([]int32, error)
	PidExists(pid int32) (bool, error)
	NewProcess(pid int32) (Process, error)
}
