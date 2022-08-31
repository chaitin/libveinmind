import ctypes as C
from . import binding


class Process(binding.Object):
    def __init__(self, handle):
        super(Process, self).__init__(handle=handle)

    _parent = binding.lookup(b"veinmind_ProcessParent", b"VEINMIND_1.2")
    def parent(self):
        "Retrieves parent process from current process"

        h_process = binding.Handle()
        binding.handle_error(Process._parent(
            h_process.ptr(), self.__handle__().val()
        ))
        return Process(h_process)

    _children = binding.lookup(b"veinmind_ProcessChildren", b"VEINMIND_1.2")
    def children(self):
        "Retrieves children process from current process"

        h_procsses = binding.Handle()
        binding.handle_error(Process._children(
            h_procsses.ptr(), self.__handle__().val()
        ))

        result = list()
        for i in range(h_procsses.len()):
            result.append(Process(h_procsses.index(i)))
        return result

    _cmdline = binding.lookup(b"veinmind_ProcessCmdline", b"VEINMIND_1.2")
    def cmdline(self):
        "Retrieves cmdline from process"

        h_str = binding.Handle()
        binding.assert_no_error(Process._cmdline(h_str.ptr(), self.__handle__().val()))
        return h_str.str()

    _environ = binding.lookup(b"veinmind_ProcessEnviron", b"VEINMIND_1.2")
    def environ(self):
        "Retrieves environ from process"

        h_str_list = binding.Handle()
        binding.assert_no_error(Process._environ(h_str_list.ptr(), self.__handle__().val()))
        return h_str_list.str_list()

    _cwd = binding.lookup(b"veinmind_ProcessCwd", b"VEINMIND_1.2")
    def cwd(self):
        "Retrieves cwd from process"

        h_str = binding.Handle()
        binding.assert_no_error(Process._cwd(h_str.ptr(), self.__handle__().val()))
        return h_str.str()

    _exe = binding.lookup(b"veinmind_ProcessExe", b"VEINMIND_1.2")
    def exe(self):
        "Retrieves exe from process"

        h_str = binding.Handle()
        binding.assert_no_error(Process._exe(h_str.ptr(), self.__handle__().val()))
        return h_str.str()

    _gids = binding.lookup(b"veinmind_ProcessGids", b"VEINMIND_1.2")
    def gids(self):
        "Retrieves gids from process"

        h_int32_array = binding.Handle()
        binding.assert_no_error(Process._gids(h_int32_array.ptr(), self.__handle__().val()))
        return h_int32_array.int32_list()

    _uids = binding.lookup(b"veinmind_ProcessUids", b"VEINMIND_1.2")
    def uids(self):
        "Retrieves uids from process"

        h_int32_array = binding.Handle()
        binding.assert_no_error(Process._uids(h_int32_array.ptr(), self.__handle__().val()))
        return h_int32_array.int32_list()

    _pid = binding.lookup(b"veinmind_ProcessPid", b"VEINMIND_1.2")
    def pid(self):
        "Retrieves pid from process"

        h_res = C.c_int32()
        binding.assert_no_error(Process._pid(C.pointer(h_res), self.__handle__().val()))
        return h_res.value

    _host_pid = binding.lookup(b"veinmind_ProcessHostPid", b"VEINMIND_1.2")
    def host_pid(self):
        "Retrieves host_pid from process"

        h_res = C.c_int32()
        binding.assert_no_error(Process._host_pid(C.pointer(h_res), self.__handle__().val()))
        return h_res.value

    _ppid = binding.lookup(b"veinmind_ProcessPpid", b"VEINMIND_1.2")
    def ppid(self):
        "Retrieves ppid from process"

        h_res = C.c_int32()
        binding.handle_error(Process._ppid(C.pointer(h_res), self.__handle__().val()))
        return h_res.value

    _name = binding.lookup(b"veinmind_ProcessName", b"VEINMIND_1.2")
    def name(self):
        "Retrieves name from process"

        h_str = binding.Handle()
        binding.assert_no_error(Process._name(h_str.ptr(), self.__handle__().val()))
        return h_str.str()

    _status = binding.lookup(b"veinmind_ProcessStatus", b"VEINMIND_1.2")
    def status(self):
        "Retrieves status from process"

        h_str = binding.Handle()
        binding.assert_no_error(Process._status(h_str.ptr(), self.__handle__().val()))
        return h_str.str()

    _create_time = binding.lookup(b"veinmind_ProcessCreateTime", b"VEINMIND_1.2")
    def create_time(self):
        "Retrieves create_time from process"

        h_res = C.c_int64()
        binding.assert_no_error(Process._create_time(C.pointer(h_res), self.__handle__().val()))
        return h_res.value

    def close(self):
        binding.assert_no_error(self.__handle__().free())

class Psutil(binding.Object):
    def __init__(self, handle):
        super(Psutil, self).__init__(handle=handle)

    _pids = binding.lookup(b"veinmind_PsutilPids", b"VEINMIND_1.2")
    def pids(self):
        "List pid in container"

        h_int32_array = binding.Handle()
        binding.handle_error(Psutil._pids(
            h_int32_array.ptr(), self.__handle__().val()
        ))
        return h_int32_array.int32_list()

    _pid_exists = binding.lookup(b"veinmind_PsutilPidExists", b"VEINMIND_1.2")
    def pid_exist(self, pid):
        "Determine if PID exists"

        h_res = C.c_int32()
        binding.assert_no_error(Psutil._pid_exists(
            C.pointer(h_res), self.__handle__().val(), C.c_int32(pid)
        ))
        if h_res.value == 0:
            return False
        else:
            return True

    _process = binding.lookup(b"veinmind_PsutilNewProcess", b"VEINMIND_1.2")
    def Process(self, pid):
        "Create process from psutil"

        h_process = binding.Handle()
        binding.handle_error(Psutil._process(
            h_process.ptr(), self.__handle__().val(), C.c_int32(pid)
        ))
        return Process(h_process)
