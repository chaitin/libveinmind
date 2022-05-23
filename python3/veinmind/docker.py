from . import binding as binding
from . import runtime as runtime
from . import filesystem as filesystem
from . import image as image
import ctypes as C

class Docker(runtime.Runtime):
	"Docker refers to a parsed docker application object."

	_with_config_path_fn = binding.lookup(
		b"veinmind_DockerWithConfigPath", b"VEINMIND_1.1")
	def _with_config_path(path):
		if not isinstance(path, str):
			raise TypeError("path must be str")
		with binding.new_str(path) as hstr:
			hopt = binding.Handle()
			binding.handle_error(Docker._with_config_path_fn(
				hopt.ptr(), hstr.val()))
			return hopt

	_with_data_root_dir_fn = binding.lookup(
		b"veinmind_DockerWithDataRootDir", b"VEINMIND_1.1")
	def _with_data_root_dir(path):
		if not isinstance(path, str):
			raise TypeError("path must be str")
		with binding.new_str(path) as hstr:
			hopt = binding.Handle()
			binding.handle_error(Docker._with_data_root_dir_fn(
				hopt.ptr(), hstr.val()))
			return hopt

	_with_unique_desc_fn = binding.lookup(
		b"veinmind_DockerWithUniqueDesc", b"VEINMIND_1.1")
	def _with_unique_desc(desc):
		if not isinstance(desc, str):
			raise TypeError("desc must be str")
		with binding.new_str(desc) as hstr:
			hopt = binding.Handle()
			binding.handle_error(Docker._with_unique_desc_fn(
				hopt.ptr(), hstr.val()))
			return hopt

	# Initialize the docker object, with specified arguments.
	_make_new_option_list = binding.lookup(
		b"veinmind_DockerMakeNewOptionList", b"VEINMIND_1.1")
	_new = binding.lookup(b"veinmind_DockerNew", b"VEINMIND_1.1")
	def __init__(self, **kwargs):
		hopts = binding.Handle()
		binding.handle_error(
			Docker._make_new_option_list(hopts.ptr()))
		with hopts as hopts:
			config_path = kwargs.pop("config_path", None)
			if config_path is not None:
				with Docker._with_config_path(config_path) as hopt:
					hopts.append(hopt)
			data_root_dir = kwargs.pop("data_root_dir", None)
			if data_root_dir is not None:
				with Docker._with_data_root_dir(data_root_dir) as hopt:
					hopts.append(hopt)
			unique_desc = kwargs.pop("unique_desc", None)
			if unique_desc is not None:
				with Docker._with_unique_desc(unique_desc) as hopt:
					hopts.append(hopt)

			handle = binding.Handle()
			binding.handle_error(Docker._new(
				handle.ptr(), hopts.val()))
			super(Docker, self).__init__(handle=handle)

	# Open a image by its ID and return the image object.
	_open_image_by_id = binding.lookup(
		b"veinmind_RuntimeOpenImageByID", b"VEINMIND_1.0")
	def open_image_by_id(self, image_id):
		"Open a docker image by its specified ID."

		with binding.new_str(image_id) as hstr:
			handle = binding.Handle()
			binding.handle_error(Docker._open_image_by_id(
				handle.ptr(), self.__handle__().val(), hstr.val()))
			return Image(handle)

	_unique_desc = binding.lookup(
		b"veinmind_DockerUniqueDesc", b"VEINMIND_1.1")
	def unique_desc(self):
		"Retrieve the correlated unique descriptor."

		handle = binding.Handle()
		binding.handle_error(Docker._unique_desc(
			handle.ptr(), self.__handle__().val()))
		with handle as handle:
			return handle.str()

class Layer(filesystem.FileSystem):
	"Layer refers to a layer in docker image."

	# Initialize the docker layer object.
	def __init__(self, handle):
		super(Layer, self).__init__(handle=handle)

	_id = binding.lookup(b"veinmind_DockerLayerID", b"VEINMIND_1.0")
	def id(self):
		"Retrieve the diff ID of the docker layer."

		handle = binding.Handle()
		binding.handle_error(Layer._id(
			handle.ptr(), self.__handle__().val()))
		with handle as handle:
			return handle.str()

class Image(image.Image):
	"Image refers to a docker specific image."

	# Initialize the docker image object.
	def __init__(self, handle):
		super(Image, self).__init__(handle=handle)

	_num_layers = binding.lookup(
		b"veinmind_DockerImageNumLayers", b"VEINMIND_1.0")
	def num_layers(self):
		"Return the number of layers in the docker image."

		result = C.c_size_t()
		binding.handle_error(Image._num_layers(
			C.pointer(result), self.__handle__().val()))
		return result.value

	_open_layer = binding.lookup(
		b"veinmind_DockerImageOpenLayer", b"VEINMIND_1.0")
	def open_layer(self, i):
		"Open specified layer in the docker image."

		handle = binding.Handle()
		binding.handle_error(Image._open_layer(
			handle.ptr(), self.__handle__().val(), C.c_size_t(i)))
		return Layer(handle)

	_get_layer_diff_id = binding.lookup(
		b"veinmind_DockerImageGetLayerDiffID", b"VEINMIND_1.0")
	def get_layer_diff_id(self, i):
		"Retrieve the diff ID the of docker layer without opening it."

		handle = binding.Handle()
		binding.handle_error(Image._get_layer_diff_id(
			handle.ptr(), self.__handle__().val(), C.c_size_t(i)))
		with handle as handle:
			return handle.str()
