import click.core as core
import click.decorators as decorators

# manifest version that current library supports.
_current_manifest_version = 1

class Manifest(object):
	"""
	Manifest is the plugin manifest whose content is essential
	for its hosting command.

	You might use the decorator "command.manifest" to annotate
	either a command, group and so on to specify the manifest.
	And it is usually annotated within "__veinmind_manifest__" attr.
	"""
	def __init__(self, **kwargs):
		self.name = kwargs.pop("name", "")
		self.version = kwargs.pop("version", "")
		self.author = kwargs.pop("author", "")
		self.description = kwargs.pop("description", "")
		self.tags = kwargs.pop("tags", list())

	def serialize(self, commands=None):
		"""
		Create the dict conforming to the format specified
		in github.com/chaitin/libveimind/go/plugin.Manifest,
		so that it can be serialized and transferred back.
		"""
		return {
			"manifestVersion": _current_manifest_version,
			"name":            self.name,
			"version":         self.version,
			"author":          self.author,
			"description":     self.description,
			"tags":            self.tags,
			"commands":        commands,
		}

def manifest(**kwargs):
	"""
	Manifest annotates the current function or command object
	with a manifest, so that it can be used as the manifest to
	present back to the host once matched.
	"""
	def g(f):
		f.__veinmind_manifest__ = Manifest(**kwargs)
		return f
	return g

# default manifest is set through set_manifest function, which
# creates a global manifest that will override the search of info
# command (if it does not have an internal manifest).
_default_manifest = None

def set_manifest(**kwargs):
	"""Specify the manifest of the python program as plugin."""
	global _default_manifest
	_default_manifest = Manifest(**kwargs)

def _generate_default_manifest():
	"""Attempt to generate the default manifest."""

	# We'll gladly return the user-specified manifest.
	if _default_manifest is not None:
		return _default_manifest

	# We need to generate the default manifest ourselves, just
	# like what we have done in the golang counterpart.
	#
	# XXX: merely using the __main__ will not yield enough
	# information for us, so we returns the path of the main
	# package for identifying them. The best recommended way is
	# that user specifies themselves.
	name = None
	import __main__
	import os.path
	if hasattr(__main__, '__file__'):
		name = os.path.dirname(os.path.realpath(__main__.__file__))
	return Manifest(name=name,
		description="a plugin powered by libVeinMind proudly")

def _search_manifest(cmd, ctx=None):
	"""Search for manifest in the tree of command."""
	if hasattr(cmd, "__veinmind_manifest__"):
		return cmd.__veinmind_manifest__
	if isinstance(cmd, core.Group):
		for name in cmd.list_commands(ctx):
			manifest = _search_manifest(
				cmd.get_command(ctx, name), ctx=ctx)
			if manifest is not None:
				return manifest
	return None

def _aggregate_commands(cmd, *path, ctx=None):
	"""Search for commands and aggregate them into list."""
	result = list()
	if isinstance(cmd, core.Group):
		for name in cmd.list_commands(ctx):
			result.extend(_aggregate_commands(
				cmd.get_command(ctx, name), *path, name, ctx=ctx))
	elif hasattr(cmd, "__veinmind_command__"):
		if not isinstance(cmd.__veinmind_command__, dict):
			raise ValueError("__veinmind_command__ must be a dict")
		info = cmd.__veinmind_command__.copy()
		info["path"] = [*path]
		result.append(info)
	return result

class InfoCommand(core.Command):
	"""
	InfoCommand will attempt to lookup the command in its current
	group, and searches whether the attr of "__veinmind_command__"
	presents. The attribute will be used as the content to present
	in the info command.
	"""

	def __init__(self, manifest=None, **kwargs):
		"""Info command will either use the specified manifest, or
		just find it while visiting all groups and commands."""
		kwargs.setdefault("help",
			"Describe libVeinMind plugin command entrypoints")
		super(InfoCommand, self).__init__("info", **kwargs)
		self.manifest = manifest

	def invoke(self, ctx):
		"""Generate manifest for other libveinmind commands."""
		if ctx.parent is None:
			ctx.fail("InfoCommand must be attached to certain command")
		cmd = ctx.parent.command

		# Aggregate and create the manifest from parent command.
		#
		# XXX: unlike the golang counterpart, the manifest can be
		# attached to any children commands (that can be found by
		# _search_manifest). So we must attempt to find it before
		# we are attempting to use the default manifest.
		#
		# TODO: maybe we should unify the behaviour of manifest
		# fetching in both golang and python3.
		manifest = None
		if self.manifest is not None:
			manifest = self.manifest
		else:
			manifest = _search_manifest(cmd)
		if manifest is None:
			manifest = _generate_default_manifest()

		# Aggregate and collect the commands for manifest.
		commands = _aggregate_commands(cmd)

		# Write the manifest result back to the caller.
		import json
		print(json.dumps(manifest.serialize(commands=commands)))

class PluginCommand(core.Command):
	"""
	PluginCommand is a service aware command which will initialize
	the service submodule before entering any user code. It is
	designed to be equivalent to cmd.PluginCommand in golang API.
	"""

	def __init__(self, name=None, **kwargs):
		super(PluginCommand, self).__init__(name, **kwargs)
		self.params.append(core.Option(["--host"],
			type=str, multiple=True,
			help="the URL of host communication file"))
		f = self.callback
		if f is None:
			return
		if hasattr(f, "__veinmind_manifest__"):
			self.__veinmind_manifest__ = f.__veinmind_manifest__
		if hasattr(f, "__veinmind_command__"):
			self.__veinmind_command__ = f.__veinmind_command__

	def invoke(self, ctx):
		from . import service
		host = ctx.params.pop("host", [])
		service.init_service_client(*host)
		return self.invoke_plugin(ctx)

	def invoke_plugin(self, ctx):
		return ctx.invoke(self.callback, root, **ctx.params)

def _plugin_command(name=None, **kwargs):
	return decorators.command(name=name, cls=PluginCommand, **kwargs)

# Mode registers which will be used in the mode command.
_mode_registry = dict()

def mode(name=None, **kwargs):
	"""
	Define a mode specifed by name and its mode handler.

	Please notice that multiple root objects might be initialized
	in a single mode, so a mode command must initialize and manage
	the lifecycle of root objects, while calling the callback with
	respect to each root object.
	"""
	def g(f):
		mode = decorators.command(name=name, **kwargs)(f)
		_mode_registry[mode.name] = mode
	return g

@mode(name="docker")
def _docker_mode(callback, **kwargs):
	from . import docker
	with docker.Docker() as d:
		callback(d)

@mode(name="containerd")
def _containerd_mode(callback, **kwargs):
	from . import containerd
	with containerd.Containerd() as c:
		callback(c)

class ModeCommand(PluginCommand):
	"""ModeCommand is the command that will accept in a mode flag
	along side with mode parameters. It will initialize the mode
	object and specify it as the first parameter for invoke."""

	def __init__(self, **kwargs):
		super(ModeCommand, self).__init__(**kwargs)
		self.params.append(core.Option(["--mode"], default="",
			help="select mode to retrieve root object"))
		for key, mode in _mode_registry.items():
			self.params.append(core.Option(["--"+key],
				is_flag=True, default=False,
				help="specify {mode} as the mode in use"
					.format(mode=key)))
			self.params.extend(mode.params)

	def invoke_plugin(self, ctx):
		mode = ctx.params["mode"]
		if mode == "":
			for key in _mode_registry.keys():
				if ctx.params[key]:
					mode = key
					break
		# TODO: support auto-recognizer mode for scanning.
		if mode == "":
			mode = "docker"

		if mode not in _mode_registry:
			ctx.fail("unknown mode {mode}".format(mode=mode))
		def g(root):
			# Remove mode specific options from the ctx.params.
			ctx.params.pop("mode", None)
			for key, mode in _mode_registry.items():
				ctx.params.pop(key, None)
				for param in mode.params:
					ctx.params.pop(param.name, None)
			self.invoke_mode(ctx, root)

		return ctx.invoke(_mode_registry[mode], g, **ctx.params)

	def invoke_mode(self, ctx, root):
		ctx.invoke(self.callback, root, **ctx.params)

def _mode_command(name=None, **kwargs):
	return decorators.command(name=name, cls=ModeCommand, **kwargs)

class RuntimeCommand(ModeCommand):
	"""RuntimeCommand is the command that accepts in an argument
	of the runtime object."""

	def __init__(self, **kwargs):
		super(RuntimeCommand, self).__init__(**kwargs)
		self.__veinmind_command__ = {
			"type": "runtime",
			"data": {},
		}

	def invoke_mode(self, ctx, root):
		from . import runtime
		if not isinstance(root, runtime.Runtime):
			ctx.fail("incompatible mode")
		ctx.invoke(self.callback, root, **ctx.params)

def _runtime_command(name=None, **kwargs):
	return decorators.command(name=name, cls=RuntimeCommand, **kwargs)

class ImageCommand(ModeCommand):
	"""ImageCommand is the command that accepts in either an
	argument of the image object, or arguments of a image object
	plus the list of image IDs."""

	def __init__(self, pass_image_id=False, **kwargs):
		super(ImageCommand, self).__init__(**kwargs)
		self.params.append(core.Option(["--id"],
			default=False, is_flag=True,
			help="whether fully qualified ID is specified"))
		self.pass_image_id = pass_image_id
		self.allow_extra_args=True
		self.no_args_is_help=False
		self.__veinmind_command__ = {
			"type": "image",
			"data": {},
		}

	def invoke_mode(self, ctx, root):
		from . import runtime
		if not isinstance(root, runtime.Runtime):
			ctx.fail("incompatible mode")

		# Retrieve the list of image IDs for invocation.
		image_ids = list()
		fully_qualified_id = ctx.params.pop("id", False)
		if len(ctx.args) == 0:
			image_ids.extend(root.list_image_ids())
		if fully_qualified_id:
			image_ids.extend(ctx.args)
		else:
			for arg in ctx.args:
				image_ids.extend(root.find_image_ids(arg))

		# Invoke by passing image IDs or objects.
		if self.pass_image_id:
			ctx.invoke(self.callback, root, image_ids, **ctx.params)
		else:
			for image_id in image_ids:
				with root.open_image_by_id(image_id) as image:
					ctx.invoke(self.callback, image, **ctx.params)

def _image_command(name=None, **kwargs):
	return decorators.command(name=name, cls=ImageCommand, **kwargs)

def _image_id_command(name=None, **kwargs):
	kwargs["pass_image_id"] = True
	return _image_command(name=name, **kwargs)

class Group(core.Group):
	"""Command group augmented with our commands."""

	def add_info_command(self, manifest=None, **kwargs):
		"""Create an info command under current group."""
		self.add_command(InfoCommand(manifest=manifest, **kwargs))

	def plugin_command(self, *args, **kwargs):
		"""Decorator creating a plugin command under current group."""
		def g(f):
			cmd = _plugin_command(*args, **kwargs)(f)
			self.add_command(cmd)
			return cmd
		return g

	def mode_command(self, *args, **kwargs):
		"""Decorator creating a mode command under current group."""
		def g(f):
			cmd = _mode_command(*args, **kwargs)(f)
			self.add_command(cmd)
			return cmd
		return g

	def runtime_command(self, *args, **kwargs):
		"""Decorator creating a runtime command under current group."""
		def g(f):
			cmd = _runtime_command(*args, **kwargs)(f)
			self.add_command(cmd)
			return cmd
		return g

	def image_command(self, *args, **kwargs):
		"""Decorator creating an image command under current group."""
		def g(f):
			cmd = _image_command(*args, **kwargs)(f)
			self.add_command(cmd)
			return cmd
		return g

	def image_id_command(self, *args, **kwargs):
		"""Decorator creating an image ID command under current group."""
		def g(f):
			cmd = _image_id_command(*args, **kwargs)(f)
			self.add_command(cmd)
			return cmd
		return g

def group(name=None, attach_main=True, **kwargs):
	"""Create a command group with parameter."""
	kwargs.setdefault("cls", Group)
	f = main.group if attach_main else decorators.group
	return f(name=name, **kwargs)

# While decorating functions with groupless function like
# "command.runtime" and "command.image", it is also registered to
# the group "main" here.
@group(attach_main=False)
def main(*args, **kwargs):
	pass
main.add_info_command()

def command(name=None, attach_main=True, **kwargs):
	"""Create a callable command with parameters."""
	f = main.command if attach_main else decorators.command
	return f(name=name, **kwargs)

def option(*args, **kwargs):
	"""Invokes click.decorators.option."""
	return decorators.option(*args, **kwargs)

def plugin_command(name=None, attach_main=True, **kwargs):
	"""Create a plugin command with parameters."""
	f = main.plugin_command if attach_main else _plugin_command
	return f(name=name, **kwargs)

def mode_command(name=None, attach_main=True, **kwargs):
	"""Create a mode command with parameters."""
	f = main.mode_command if attach_main else _mode_command
	return f(name=name, **kwargs)

def runtime_command(name=None, attach_main=True, **kwargs):
	"""Create a runtime command with parameters."""
	f = main.runtime_command if attach_main else _runtime_command
	return f(name=name, **kwargs)

def image_command(name=None, attach_main=True, **kwargs):
	"""Create an image command with parameters."""
	f = main.image_command if attach_main else _image_command
	return f(name=name, **kwargs)

def image_id_command(name=None, attach_main=True, **kwargs):
	"""Create an image ID command with parameters."""
	f = main.image_id_command if attach_main else _image_id_command
	return f(name=name, **kwargs)
