import io
import os
import urllib.parse as urlparse
import json

# service_openers is the registry for the service openers
# that can be used
_service_openers = dict()

def service_opener(name=None):
	"""
	Define a new service opener with its specified name.

	Upon the protocol match, the service opener will receive
	a urlparse.ParseResult, whose path will be escaped upon
	invoking, plus an indicator for whether it is open for
	read or write. The service opener should return a file
	pointer like object that support reading or writing.
	"""
	def g(f):
		key = name if name is not None else f.__name__
		_service_openers[key] = f
		return f
	return g

def _open_service(path, rw):
	url = urlparse.urlparse(path)
	escaped_path = urlparse.unquote(url.path)
	if url.scheme not in _service_openers:
		raise KeyError("unsupported scheme {scheme}"
			.format(scheme=url.scheme))
	return _service_openers[url.scheme](urlparse.ParseResult(
			url.scheme, url.netloc, escaped_path, url.params,
			url.query, url.fragment), rw)

@service_opener("file")
def _open_file(url, rw):
	fd = os.open(url.path, os.O_RDWR if rw else os.O_RDONLY)
	return os.fdopen(fd, 'wb' if rw else 'rb', 0)

# _is_hosted indicator stores whether current application
# has been put into hosted mode.
_is_hosted = False

def is_hosted():
	return _is_hosted

def assert_hosted():
	assert _is_hosted, "unhosted condition not handled"

# _infile and _outfile are the files open for communicating
# to the plugin host.
_infile = None
_outfile = None

def init_service_client(*args):
	"""
	Initialize the client for accessing service.

	When there's no argument passed in, it is equivalent
	to running in standalone mode, with its output truncated
	to empty file.

	When either one or two arguments are specified, it will
	attempt to open the stream for reading and writing.
	"""
	if len(args) == 0:
		return
	elif len(args) > 2:
		raise ValueError("too many arguments for opening file")
	else:
		global _infile, _outfile, _is_hosted
		if len(args) > 1:
			_infile = _open_service(args[0], False)
			_outfile = _open_service(args[1], True)
		else:
			_infile = _open_service(args[0], True)
			_outfile = _infile
		_is_hosted = True

# _initer is the input iterator that attempts to decode json
# data from the input stream. It returns whenever a newer
# json has been decoded.
_initer = None

def _input_iter():
	assert_hosted()
	dec = json.JSONDecoder()
	remaining = b''
	cursor = 0
	depth = 0
	inside_quote = False

	while True:
		while cursor >= len(remaining):
			current = _infile.read(io.DEFAULT_BUFFER_SIZE)
			if current == None:
				raise EOFError()
			remaining = remaining+current
		assert cursor < len(remaining)

		char = remaining[cursor:cursor+1]
		if inside_quote:
			if char == b'"':
				inside_quote = False
			elif char == b'\\':
				# Skip the next character of backslash.
				cursor = cursor+1
		elif depth <= 0:
			if char == b'{':
				depth = depth + 1
			elif cursor == 0 and char.isspace():
				# Trim leading spaces at depth 0.
				cursor = cursor - 1
				remaining = remaining[1:]
		else:
			if char == b'"':
				inside_quote = True
			elif char == b'{':
				depth = depth + 1
			elif char == b'}':
				depth = depth - 1
		cursor = cursor+1

		# XXX: the root object transferred will be
		# always the JSON object, so we need not to
		# consider the depth of '[' and ']' pairs.
		if depth <= 0 and char == b'}':
			result, _ = dec.raw_decode(
				str(remaining[:cursor], 'utf-8'))
			remaining = remaining[cursor:]
			cursor = 0
			yield result


# Request attempts to perform a request to host.
_id = 0
def _request(payload):
	assert_hosted()
	global _id
	_id = _id + 1
	sequence = _id
	payload["sequence"] = sequence
	_outfile.write(bytes(json.dumps(payload), 'utf-8'))
	global _initer
	if _initer is None:
		_initer = _input_iter()
	response = next(_initer)
	if response["sequence"] != sequence:
		# TODO: we only supports single threaded
		# request of services, but we might consider
		# supporting blocking services as there might
		# be some in the future.
		raise RuntimeError("simultaneous multiple requests")
	if "error" in response:
		err = response["error"]
		if err is not None:
			raise RuntimeError(err)
	return response

# User specific functions to play with services.
def has_namespace(namespace):
	"""Tell whether a namespace exists under a namespace."""
	response = _request({
		"type": "hasNamespace",
		"namespace": namespace,
	})
	return response.get("ok", False)

def get_manifest(namespace):
	"""Retrieve the manifest of a namespace."""
	response = _request({
		"type": "getManifest",
		"namespace": namespace,
	})
	return response.get("reply", None)

def list_services(namespace):
	"""Retrieve the list of services under a namespace."""
	response = _request({
		"type": "listServices",
		"namespace": namespace,
	})
	return response.get("services", list())

def rawcall(namespace, name, *args):
	"""Raw invoke the service specified by namespace and name."""
	response = _request({
		"type": "call",
		"namespace": namespace,
		"name": name,
		"args": [*args],
	})
	return tuple(response.get("reply", list()))

def service(namespace, name=None):
	"""
	Decorator for wrapping service with its fallback behaviour.

	Either the service of the host will be called, with arguments
	passed to it be returned, or the underlying function will be
	called when it is not hosted (or service not found).
	"""
	def g(f):
		_name = name
		if _name is None:
			_name = f.__name__
		def h(*args):
			if not hasattr(f, "__veinmind_serviced__"):
				if not is_hosted():
					f.__veinmind_serviced__ = False
				elif not has_namespace(namespace):
					f.__veinmind_serviced__ = False
				else:
					services = list_services(namespace)
					f.__veinmind_serviced__ = _name in services

			if f.__veinmind_serviced__:
				return rawcall(namespace, name, *args)
			else:
				return f(*args)
		return h
	return g
