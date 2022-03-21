import enum
import time
import sys
from . import service

_namespace = "github.com/chaitin/libveinmind/logging"

# Normalize timezone and format into RFC3339 format.
_timezone = time.strftime('%z')
assert len(_timezone) == 5
if _timezone == "+0000" or _timezone == "-0000":
	_timezone = "Z"
else:
	_timezone = _timezone[0:3] + ':' + _timezone[3:5]
_format = "%Y-%m-%dT%H:%M:%S" + _timezone

# Log level enumerations for comparation and filtering.
@enum.unique
class Level(enum.Enum):
	PANIC = 0
	FATAL = 1
	ERROR = 2
	WARN  = 3
	INFO  = 4
	DEBUG = 5
	TRACE = 6

@service.service(_namespace, "getConfig")
def _get_config():
	return ({
		"level": Level.INFO.value,
	},)

@service.service(_namespace, "log")
def _log(logs):
	pass

_log_level = None
def log(level, msg, *args, **kwargs):
	global _log_level
	if _log_level is None:
		config = _get_config()[0]
		_log_level = config.get("level", 0)
	if level.value > _log_level:
		return
	message = msg % args
	fields = kwargs.copy()
	if len(fields) == 0:
		fields = None
	if service.is_hosted():
		_log([{
			"time": time.strftime(_format),
			"level": level.value,
			"fields": fields,
			"msg": message,
		}])
	else:
		payload = "{asctime} - {level} - {message}".format(
			asctime=time.asctime(), level=level.name,
			message=message)
		if fields is not None:
			payload = payload + "\t" + str(fields)
		print(payload, file=sys.stderr)

	# Handle extra behaviour for fatal and panic level.
	if level == Level.PANIC:
		raise Exception(message)
	elif level == Level.FATAL:
		raise SystemExit(1)
		

def trace(msg, *args, **kwargs):
	log(Level.TRACE, msg, *args, **kwargs)

def debug(msg, *args, **kwargs):
	log(Level.DEBUG, msg, *args, **kwargs)

def info(msg, *args, **kwargs):
	log(Level.INFO, msg, *args, **kwargs)

def warn(msg, *args, **kwargs):
	log(Level.WARN, msg, *args, **kwargs)

def warning(msg, *args, **kwargs):
	warn(msg, *args, **kwargs)

def error(msg, *args, **kwargs):
	log(Level.ERROR, msg, *args, **kwargs)

def fatal(msg, *args, **kwargs):
	log(Level.FATAL, msg, *args, **kwargs)

def panic(msg, *args, **kwargs):
	log(Level.PANIC, msg, *args, **kwargs)

def critical(msg, *args, **kwargs):
	panic(msg, *args, **kwargs)

class Entry:
	def __init__(self, **kwargs):
		self.fields = kwargs.copy()

	def _combine_fields(self, **kwargs):
		passing = kwargs.copy()
		for key, value in self.fields.items():
			passing.setdefault(key, value)
		return passing

	def log(self, level, msg, *args, **kwargs):
		log(level, msg, *args, **(self._combine_fields(**kwargs)))

	def trace(self, msg, *args, **kwargs):
		self.log(Level.TRACE, msg, *args, **kwargs)

	def debug(self, msg, *args, **kwargs):
		self.log(Level.DEBUG, msg, *args, **kwargs)

	def info(self, msg, *args, **kwargs):
		self.log(Level.INFO, msg, *args, **kwargs)

	def warn(self, msg, *args, **kwargs):
		self.log(Level.WARN, msg, *args, **kwargs)

	def warning(self, msg, *args, **kwargs):
		self.warn(msg, *args, **kwargs)

	def error(self, msg, *args, **kwargs):
		self.log(Level.ERROR, msg, *args, **kwargs)

	def fatal(self, msg, *args, **kwargs):
		self.log(Level.FATAL, msg, *args, **kwargs)

	def panic(self, msg, *args, **kwargs):
		self.log(Level.PANIC, msg, *args, **kwargs)

	def critical(self, msg, *args, **kwargs):
		self.panic(msg, *args, **kwargs)
