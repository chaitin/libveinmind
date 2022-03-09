import ctypes as C
import errno
import os

# Load dlvsym for accessing specific version of symbol from
# libveinmind.so, emulating the compiled one's behaviour.
_dlso = C.CDLL("libdl.so")
_dlvsym = _dlso["dlvsym"]
_dlvsym.restype = C.c_void_p

# Load the libveinmind SDK library for futher function lookup.
_libveinmind = C.CDLL("libveinmind.so")

def lookup(name, version, restype=C.c_size_t):
	"Lookup a versioned function from libveinmind SDK library."

	sym = _dlvsym(C.c_void_p(_libveinmind._handle), name, version)
	if sym == 0:
		raise RuntimeError(
			"unsatisfied link {name}@{version}".format(
				name=name, version=version))
	result = C.CFUNCTYPE(restype)(sym)
	result.restype = restype
	return result

# Pre-defined error enumerations for throwing out.
_error_msg = [
	None,
	"out of memory",
	"null pointer receiver",
	"invalid ID specified",
	"invalid OP specified",
	"panic in call",
	"race detected",
]

# BindingRuntimeError is the runtime error specified by binding.
class BindingRuntimeError(RuntimeError):
	"RuntimeError dedicated to wrong API binding usage."

	def __init__(self, code):
		message = "unknown error"
		if code > 0 and code < len(_error_msg):
			message = _error_msg[code]
		RuntimeError.__init__(self, message)

# assert_no_error checks the result from the functions that
# is not expected to generate exception.
def assert_no_error(result):
	"Check and raise error from exception-less context."

	if result == 0:
		return
	raise BindingRuntimeError(result)

# Handle object which can be wrapped to ensure it is freed
# automatically once it is not required.
class Handle:
	"Handle to an allocated resource in the API binding."

	def __init__(self, init=0):
		self._value = C.c_size_t(init) # veinmind_id_t

	_free = lookup(b"veinmind_Free", b"VEINMIND_1.0", None)
	def free(self):
		"Recycle the resource corresponding to the handle."

		if Handle._free is None:
			return
		Handle._free(self._value)
		self._value.value = 0

	def __del__(self):
		self.free()

	def __enter__(self):
		return self

	def __exit__(self, exc_type, exc_value, traceback):
		self.free()

	def ptr(self):
		"Return a modifiable pointer to the handle."

		return C.pointer(self._value)

	def val(self):
		"Return the underlying value of the handle."

		return C.c_size_t(self._value.value)

	# Length function applicable to array like objects.
	_len = lookup(b"veinmind_Length", b"VEINMIND_1.0")
	def len(self):
		"When object is array like, return its length."

		length = C.c_size_t()
		assert_no_error(Handle._len(C.pointer(length), self.val()))
		return length.value

	def __len__(self):
		return self.len()

	# Indexing inside an array like object.
	_index = lookup(b"veinmind_Index", b"VEINMIND_1.0")
	def index(self, i):
		"When object is array like, return item at index."

		if not isinstance(i, int):
			raise TypeError("index must be int")
		result = Handle()
		assert_no_error(Handle._index(
			result.ptr(), self.val(), C.c_size_t(i)))
		return result

	# Raw pointer applicable to bytes object.
	_raw_pointer = lookup(b"veinmind_RawPointer", b"VEINMIND_1.0")
	def rawptr(self):
		"When object is bytes, return its address pointer."

		result = C.c_void_p()
		assert_no_error(Handle._raw_pointer(C.pointer(result), self.val()))
		return result

	# Actually retrieve a bytes object, might be either binary or text.
	def bytes(self):
		"When object is bytes, return its binary buffer."

		return bytes(C.string_at(self.rawptr(), self.len()))

	# Conversion between string and bytes.
	_bytes_to_str = lookup(b"veinmind_BytesToString", b"VEINMIND_1.0")
	def bytes_to_str(self):
		"When object is bytes, convert it to a string object."

		result = Handle()
		assert_no_error(Handle._bytes_to_str(result.ptr(), self.val()))
		return result

	_str_to_bytes = lookup(b"veinmind_StringToBytes", b"VEINMIND_1.0")
	def str_to_bytes(self):
		"When object is string, convert it to a bytes object."

		result = Handle()
		assert_no_error(Handle._str_to_bytes(result.ptr(), self.val()))
		return result

	# Retrieve the content of a string.
	def str(self):
		"When object is string, convert it to a Python string."

		with self.str_to_bytes() as hbytes:
			return str(hbytes.bytes(), "utf-8")

	# Zip the string array into python list.
	def str_list(self):
		"When object is string array, convert it to a Python string list."

		result = list()
		for i in range(len(self)):
			with self.index(i) as item:
				result.append(item.str())
		return result

_bytes = lookup(b"veinmind_Bytes", b"VEINMIND_1.0")
def new_bytes(value):
	"Create a new bytes buffer to be used inside the binding."

	if not isinstance(value, bytes):
		raise TypeError("value must be bytes")
	bytesbuf = C.c_buffer(value)
	result = Handle()
	assert_no_error(_bytes(result.ptr(), bytesbuf, len(value)))
	return result

def new_str(value):
	"Create a new string to be used inside the binding."

	if not isinstance(value, str):
		raise TypeError("value must be str")
	with new_bytes(bytes(value, "utf-8")) as hbytes:
		return hbytes.bytes_to_str()

# Lookup and fetch the max runtime error value.
_runtime_error_max_value = lookup(
	b"veinmind_RuntimeErrorMaxValue", b"VEINMIND_1.0")()

# Arbitrary error handler as a fallback option.
_err_string = lookup(b"veinmind_ErrorString", b"VEINMIND_1.0")
def _handle_unknown_error(herr):
	handle = Handle()
	assert_no_error(_err_string(handle.ptr(), herr.val()))
	with handle as hstr:
		raise Exception(hstr.str())

# Wrapped or chained error which should be handled separatedly.
_err_message = lookup(b"veinmind_ErrorMessage", b"VEINMIND_1.0")
_err_unwrap = lookup(b"veinmind_ErrorUnwrap", b"VEINMIND_1.0")
def _handle_wrap_error(herr):
	hmsg = Handle()
	assert_no_error(_err_message(hmsg.ptr(), herr.val()))
	msg = None
	with hmsg as hmsg:
		msg = hmsg.str()
	hnext = Handle()
	assert_no_error(_err_unwrap(hnext.ptr(), herr.val()))
	if hnext.val().value == 0:
		raise Exception(msg)
	else:
		with hnext as hnext:
			try:
				_handle_error_internal(hnext)
			except Exception as exc:
				raise Exception(msg) from exc

# Raw syscall error from a POSIX error number.
_syscall_err_get_errno = lookup(
	b"veinmind_SyscallErrorGetErrno", b"VEINMIND_1.0")
def _handle_syscall_error(herr):
	errval = C.c_size_t()
	_syscall_err_get_errno(C.pointer(errval), herr.val())
	errvalue = errval.value
	raise OSError(errvalue, os.strerror(errvalue))

# OS path error handler which attach op and path information.
# op code will be omitted for such error, adapted for python.
_os_path_err_get_op = lookup(
	b"veinmind_OSPathErrorGetOp", b"VEINMIND_1.0")
_os_path_err_get_path = lookup(
	b"veinmind_OSPathErrorGetPath", b"VEINMIND_1.0")
def _handle_os_path_error(herr):
	hnext = Handle()
	assert_no_error(_err_unwrap(hnext.ptr(), herr.val()))
	if hnext.val().value == 0:
		_handle_unknown_error(herr)
	else:
		with hnext as hnext:
			filename = None
			hpath = Handle()
			assert_no_error(_os_path_err_get_path(
				hpath.ptr(), herr.val()))
			with hpath as hpath:
				filename = hpath.str()
			try:
				_handle_error_internal(hnext)
			except OSError as osexc:
				osexc.filename = filename
				raise osexc
			except Exception as exc:
				hop = Handle()
				assert_no_error(_os_path_err_get_op(
					hop.ptr(), herr.val()))
				op = hop.str()
				raise Exception('cannot {op} {path!r}'.format(
					op = op, path = filename))

# Error type judging functions for casting them into correct types.
_is_wrap_err = lookup(b"veinmind_IsWrapError", b"VEINMIND_1.0")
_is_syscall_err = lookup(b"veinmind_IsSyscallError", b"VEINMIND_1.0")
_is_os_path_err = lookup(b"veinmind_IsOSPathError", b"VEINMIND_1.0")
_is_eof_err = lookup(b"veinmind_IsEOFError", b"VEINMIND_1.0")

# _handle_error_internal handles and rethrows the error recursively,
# while the error resource (herr) is managed by its callers.
def _handle_error_internal(herr):
	if _is_wrap_err(herr.val()) != 0:
		_handle_wrap_error(herr)
	elif _is_eof_err(herr.val()) != 0:
		raise EOFError()
	elif _is_syscall_err(herr.val()) != 0:
		_handle_syscall_error(herr)
	elif _is_os_path_err(herr.val()) != 0:
		_handle_os_path_error(herr)
	else:
		_handle_unknown_error(herr)

def handle_error(result):
	"Handles the error result of returned functions."

	if result == 0:
		return
	elif result < _runtime_error_max_value:
		raise BindingRuntimeError(result)
	else:
		with Handle(result) as herr:
			_handle_error_internal(herr)

# Object imbues the sematic that a handle is held by this super
# class and its subclasses can add their object specific behavoiurs.
class Object:
	"Object which holds and manages the handle resource."

	def __init__(self, handle):
		if not isinstance(handle, Handle):
			raise TypeError('handle must be binding.Handle')
		self.handle = handle

	def __handle__(self):
		return self.handle

	def __enter__(self):
		return self

	def __exit__(self, exc_type, exc_value, traceback):
		self.handle.free()
