package behaviour

import (
	"time"

	api "github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/pkg/binding"
)

// Psutil specifies the behaviour exhibited by those
// objects implementing api.Psutil interface.
type Psutil struct {
	h *binding.Handle
}

type Process struct {
	h binding.Handle
}

func (p *Psutil) Pids() ([]int32, error) {
	return p.h.PsutilPids()
}

func (p *Psutil) PidExists(pid int32) (bool, error) {
	return p.h.PsutilPidExists(pid)
}

func (p *Psutil) NewProcess(pid int32) (api.Process, error) {
	h, err := p.h.PsutilNewProcess(pid)
	if err != nil {
		return nil, err
	}

	return &Process{
		h: h,
	}, nil
}

func NewPsutil(h *binding.Handle) Psutil {
	return Psutil{h: h}
}

func (p *Process) Parent() (api.Process, error) {
	h, err := p.h.ProcessParent()
	if err != nil {
		return nil, err
	}

	return &Process{
		h: h,
	}, nil
}

func (p *Process) Children() ([]api.Process, error) {
	children, err := p.h.ProcessChildren()
	if err != nil {
		return nil, err
	}

	var result []api.Process
	length := children.Length()
	for i := 0; i < length; i++ {
		func() {
			h := children.Index(i)
			result = append(result, &Process{
				h: h,
			})
		}()
	}
	return result, nil
}

func (p *Process) Cmdline() (string, error) {
	return p.h.ProcessCmdline()
}

func (p *Process) Environ() ([]string, error) {
	return p.h.ProcessEnviron()
}

func (p *Process) Cwd() (string, error) {
	return p.h.ProcessCwd()
}

func (p *Process) Exe() (string, error) {
	return p.h.ProcessExe()
}

func (p *Process) Gids() ([]int32, error) {
	return p.h.ProcessGids()
}

func (p *Process) Uids() ([]int32, error) {
	return p.h.ProcessUids()
}

func (p *Process) Pid() (int32, error) {
	return p.h.ProcessPid()
}

func (p *Process) HostPid() (int32, error) {
	return p.h.ProcessHostPid()
}

func (p *Process) Ppid() (int32, error) {
	return p.h.ProcessPpid()
}

func (p *Process) Name() (string, error) {
	return p.h.ProcessName()
}

func (p *Process) Status() (string, error) {
	return p.h.ProcessStatus()
}

func (p *Process) CreateTime() (time.Time, error) {
	timestamp, err := p.h.ProcessCreateTime()
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(0, timestamp*int64(time.Millisecond)), nil
}

func (p *Process) Close() {
	p.h.Free()
}
