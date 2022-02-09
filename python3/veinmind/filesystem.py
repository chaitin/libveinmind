from . import binding as binding
from . import structseq as structseq
import ctypes as C
import io
import platform
import stat
import struct
import os.path as filepath

class FileInfo(binding.Object):
	def __init__(self, handle):
		super(FileInfo, self).__init__(handle)

	_name = binding.lookup(
		b"veinmind_FileInfoName", b"VEINMIND_1.0")
	def name(self):
		"Retrieve the name of the file info object."

		handle = binding.Handle()
		binding.handle_error(FileInfo._name(
			handle.ptr(), self.__handle__().val()))
		with handle as handle:
			return handle.str()

	_size = binding.lookup(
		b"veinmind_FileInfoSize", b"VEINMIND_1.0")
	def size(self):
		"Retrieve the size of the file info object."

		result = C.c_size_t()
		binding.handle_error(FileInfo._size(
			C.pointer(result), self.__handle__().val()))
		return result.value

	_mode = binding.lookup(
		b"veinmind_FileInfoMode", b"VEINMIND_1.0")
	def mode(self):
		"Retrieve the file mode of the file info object."

		result = C.c_uint32()
		binding.handle_error(FileInfo._mode(
			C.pointer(result), self.__handle__().val()))
		return result.value

	_mtime = binding.lookup(
		b"veinmind_FileInfoModTime", b"VEINMIND_1.0")
	def mtime(self):
		"Retrieve the modified time of the file info object."

		result = C.c_int64()
		binding.handle_error(FileInfo._mtime(
			C.pointer(result), self.__handle__().val()))
		return result.value

	# Retrieve the raw block of platform specific data.
	_sys = binding.lookup(
		b"veinmind_FileInfoSys", b"VEINMIND_1.0")
	def sys(self):
		"Retrieve the raw sys data of the file info object."

		hbytes = binding.Handle()
		binding.handle_error(FileInfo._sys(
			hbytes.ptr(), self.__handle__().val()))
		with hbytes as hbytes:
			return hbytes.bytes()

	def _statof_unknown(self):
		raise NotImplementedError()

	_fields_posix = None
	def _fields_statof_posix():
		import posix
		if FileInfo._fields_posix is None:
			FileInfo._fields_posix = structseq.describe(
				posix.stat_result, [
				"st_dev", "st_ino", "st_nlink", "st_mode",
				"st_uid", "st_gid", "st_rdev", "st_size",
				"st_blksize", "st_blocks",
				"st_atime", "st_atime_ns",
				"st_mtime", "st_mtime_ns",
				"st_ctime", "st_ctime_ns"])
		return FileInfo._fields_posix

	def _statof_linux_amd64(self):
		b = self.sys()
		pattern = "@QLLIIIQnnnqqqqqqqqq"
		assert struct.calcsize(pattern) == len(b)
		t = struct.unpack(pattern, b)
		import posix
		fields = FileInfo._fields_statof_posix()
		result = [0]*posix.stat_result.n_fields
		result[fields["st_dev"]] = t[0]
		result[fields["st_ino"]] = t[1]
		result[fields["st_nlink"]] = t[2]
		result[fields["st_mode"]] = t[3]
		result[fields["st_uid"]] = t[4]
		result[fields["st_gid"]] = t[5]
		if "st_rdev" in fields:
			result[fields["st_rdev"]] = t[6]
		result[fields["st_size"]] = t[7]
		if "st_blksize" in fields:
			result[fields["st_blksize"]] = t[8]
		if "st_blocks" in fields:
			result[fields["st_blocks"]] = t[9]
		result[fields["st_atime"]-3] = t[10]
		result[fields["st_atime"]] = t[10]+t[11]*1e-9
		result[fields["st_atime_ns"]] = t[10]*int(1e9)+t[11]
		result[fields["st_mtime"]-3] = t[12]
		result[fields["st_mtime"]] = t[12]+t[13]*1e-9
		result[fields["st_mtime_ns"]] = t[12]*int(1e9)+t[13]
		result[fields["st_ctime"]-3] = t[14]
		result[fields["st_ctime"]] = t[14]+t[15]*1e-9
		result[fields["st_ctime_ns"]] = t[14]*int(1e9)+t[15]
		return posix.stat_result(result)

	def _statof_linux_arm64(self):
		b = self.sys()
		pattern = "@QLIIIIQQniinqqqqqqii"
		assert struct.calcsize(pattern) == len(b)
		t = struct.unpack(pattern, b)
		import posix
		fields = FileInfo._fields_statof_posix()
		result = [0]*posix.stat_result.n_fields
		result[fields["st_dev"]] = t[0]
		result[fields["st_ino"]] = t[1]
		result[fields["st_mode"]] = t[2]
		result[fields["st_nlink"]] = t[3]
		result[fields["st_uid"]] = t[4]
		result[fields["st_gid"]] = t[5]
		if "st_rdev" in fields:
			result[fields["st_rdev"]] = t[6]
		result[fields["st_size"]] = t[8]
		if "st_blksize" in fields:
			result[fields["st_blksize"]] = t[9]
		if "st_blocks" in fields:
			result[fields["st_blocks"]] = t[11]
		result[fields["st_atime"]-3] = t[12]
		result[fields["st_atime"]] = t[12]+t[13]*1e-9
		result[fields["st_atime_ns"]] = t[12]*int(1e9)+t[13]
		result[fields["st_mtime"]-3] = t[14]
		result[fields["st_mtime"]] = t[14]+t[15]*1e-9
		result[fields["st_mtime_ns"]] = t[14]*int(1e9)+t[15]
		result[fields["st_ctime"]-3] = t[16]
		result[fields["st_ctime"]] = t[16]+t[17]*1e-9
		result[fields["st_ctime_ns"]] = t[16]*int(1e9)+t[17]
		return posix.stat_result(result)

	def _statof_linux_i386(self):
		b = self.sys()
		pattern = "@QHIIIIIQHQnQiiiiiiQ"
		assert struct.calcsize(pattern) == len(b)
		t = struct.unpack(pattern, b)
		import posix
		fields = FileInfo._fields_statof_posix()
		result = [0]*posix.stat_result.n_fields
		result[fields["st_dev"]] = t[0]
		result[fields["st_mode"]] = t[3]
		result[fields["st_nlink"]] = t[4]
		result[fields["st_uid"]] = t[5]
		result[fields["st_gid"]] = t[6]
		if "st_rdev" in fields:
			result[fields["st_rdev"]] = t[7]
		result[fields["st_size"]] = t[9]
		if "st_blksize" in fields:
			result[fields["st_blksize"]] = t[10]
		if "st_blocks" in fields:
			result[fields["st_blocks"]] = t[11]
		result[fields["st_atime"]-3] = t[12]
		result[fields["st_atime"]] = t[12]+t[13]*1e-9
		result[fields["st_atime_ns"]] = t[12]*int(1e9)+t[13]
		result[fields["st_mtime"]-3] = t[14]
		result[fields["st_mtime"]] = t[14]+t[15]*1e-9
		result[fields["st_mtime_ns"]] = t[14]*int(1e9)+t[15]
		result[fields["st_ctime"]-3] = t[16]
		result[fields["st_ctime"]] = t[16]+t[17]*1e-9
		result[fields["st_ctime_ns"]] = t[16]*int(1e9)+t[17]
		result[fields["st_ino"]] = t[18]
		return posix.stat_result(result)

	# Platform dependent conversion of file info into
	# os.stat_result object in python.
	_statof_platforms = {
		"Linux|x86_64":  _statof_linux_amd64,
		"Linux|amd64":   _statof_linux_amd64,
		"Linux|aarch64": _statof_linux_arm64,
		"Linux|i386":    _statof_linux_i386,
		"Linux|i686":    _statof_linux_i386,
	}
	_statof_platform_name = platform.system()+'|'+platform.machine()
	_statof_platform = _statof_unknown
	if _statof_platform_name in _statof_platforms:
		_statof_platform = _statof_platforms[_statof_platform_name]
	def stat(self):
		"Retrieve the os.stat_result of the file info object."

		return FileInfo._statof_platform(self)

class FileInfoList(binding.Object):
	def __init__(self, handle):
		super(FileInfoList, self).__init__(handle=handle)

	def __len__(self):
		return self.__handle__().len()

	def __getitem__(self, index):
		return FileInfo(self.__handle__().index(index))

	def __iter__(self):
		for i in range(len(self)):
			yield self[i]

class File(binding.Object, io.RawIOBase):
	def __init__(self, handle):
		super(File, self).__init__(handle=handle)

	_close = binding.lookup(
		b"veinmind_Close", b"VEINMIND_1.0")
	def close(self):
		binding.handle_error(File._close(
			self.__handle__().val()))

	_read = binding.lookup(
		b"veinmind_Read", b"VEINMIND_1.0")
	def readinto(self, b):
		try:
			with binding.new_bytes(bytes(b)) as hbytes:
				nread = C.c_size_t()
				binding.handle_error(File._read(C.pointer(nread),
					self.__handle__().val(), hbytes.val()))
				b[:nread.value] = hbytes.bytes()[:nread.value]
				return nread.value
		except EOFError:
			return 0

	def readable(self):
		return True

	_write = binding.lookup(
		b"veinmind_Write", b"VEINMIND_1.0")
	def write(self, b):
		try:
			with binding.new_bytes(bytes(b)) as hbytes:
				nwrite = C.c_size_t()
				binding.handle_error(File._write(C.pointer(nwrite),
					self.__handle__().val(), hbytes.val()))
				return nwrite.value
		except EOFError:
			return 0

	def writable(self):
		return True

class FileSystem(binding.Object):
	def __init__(self, handle):
		super(FileSystem, self).__init__(handle=handle)

	_open = binding.lookup(
		b"veinmind_Open", b"VEINMIND_1.0")
	def openraw(self, path):
		"Open the raw file at specified path for operating."

		with binding.new_str(path) as hstr:
			handle = binding.Handle()
			binding.handle_error(FileSystem._open(
				handle.ptr(), self.__handle__().val(), hstr.val()))
			return File(handle)

	_open_mode = {
		'r':  lambda f: io.TextIOWrapper(io.BufferedReader(f)),
		'rt': lambda f: io.TextIOWrapper(io.BufferedReader(f)),
		'rb': lambda f: io.BufferedReader(f),
	}
	def open(self, path, mode='r'):
		"Open the file at specified path with specified mode."

		# XXX: we support only read only mode by now, so if the
		# caller attempt to open with mode like '[wx][abt+U]',
		# it will be receiving invalid mode now.
		if mode not in FileSystem._open_mode:
			raise ValueError(
				'invalid mode: {mode!r}'.format(mode=mode))
		return FileSystem._open_mode[mode](self.openraw(path))

	_readlink = binding.lookup(
		b"veinmind_Readlink", b"VEINMIND_1.0")
	def readlink(self, path):
		"Read the target of specified link file at path."

		with binding.new_str(path) as hstr:
			handle = binding.Handle()
			binding.handle_error(FileSystem._readlink(
				handle.ptr(), self.__handle__().val(), hstr.val()))
			with handle as handle:
				return handle.str()

	_evalsymlink = binding.lookup(
		b"veinmind_EvalSymlink", b"VEINMIND_1.0")
	def evalsymlink(self, path):
		"Eval the final target of specified link file at path."

		with binding.new_str(path) as hstr:
			handle = binding.Handle()
			binding.handle_error(FileSystem._evalsymlink(
				handle.ptr(), self.__handle__().val(), hstr.val()))
			with handle as handle:
				return handle.str()

	_stat = binding.lookup(
		b"veinmind_Stat", b"VEINMIND_1.0")
	def stat_fileinfo(self, path):
		"Retreive the file stat represented by FileInfo at path."

		with binding.new_str(path) as hstr:
			hstat = binding.Handle()
			binding.handle_error(FileSystem._stat(
				hstat.ptr(), self.__handle__().val(), hstr.val()))
			return FileInfo(hstat)

	def stat(self, path):
		"Retrieve the file stat represented by os.stat_result at path."

		with self.stat_fileinfo(path) as stat:
			return stat.stat()

	_lstat = binding.lookup(
		b"veinmind_Lstat", b"VEINMIND_1.0")
	def lstat_fileinfo(self, path):
		"Retrieve the file lstat represented by FileInfo at path."

		with binding.new_str(path) as hstr:
			hstat = binding.Handle()
			binding.handle_error(FileSystem._lstat(
				hstat.ptr(), self.__handle__().val(), hstr.val()))
			return FileInfo(hstat)

	def lstat(self, path):
		"Retrieve the file lstat represented by os.stat_result at path."

		with self.lstat_fileinfo(path) as stat:
			return stat.stat()

	_readdir = binding.lookup(
		b"veinmind_Readdir", b"VEINMIND_1.0")
	def listdir_fileinfo(self, path):
		"Read directory content by FileInfo list at path."

		with binding.new_str(path) as hstr:
			handle = binding.Handle()
			binding.handle_error(FileSystem._readdir(
				handle.ptr(), self.__handle__().val(), hstr.val()))
			return FileInfoList(handle)

	# Anonymous function for receiving info and returning
	# its file name portion.
	def _listdir_fileinfo_name(info):
		with info as info:
			return info.name()

	def listdir(self, path='/'):
		"Read directory content by filename list at path."

		# If the specified path is bytes, convert the argument
		# and result between bytes and str to conform to the
		# interface convention of os.listdir.
		if isinstance(path, bytes):
			p = str(path, "utf-8")
			result = self.listdir(p)
			return list(map(lambda x: bytes(x, "utf-8"), result))

		with self.listdir_fileinfo(path) as result:
			return list(map(FileSystem._listdir_fileinfo_name, result))

	def walk(self, top, onerror=None):
		"Walk the specified directory hierarchy."

		# Fetch the current directory entries and report first.
		dirnames = list()
		filenames = list()
		try:
			with self.listdir_fileinfo(top) as dirents:
				for dirent in dirents:
					with dirent as dirent:
						name = dirent.name()
						if stat.S_ISDIR(dirent.stat().st_mode):
							dirnames.append(name)
						else:
							filenames.append(name)
		except OSError as e:
			if onerror is not None:
				onerror(e)
		yield top, dirnames, filenames

		# Recursively walk into the subdirectories retrieved.
		for dirname in dirnames:
			yield from self.walk(filepath.join(top, dirname),
				onerror=onerror)
