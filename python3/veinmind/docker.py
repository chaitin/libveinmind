from . import binding as binding
from . import runtime as runtime
from . import filesystem as filesystem
from . import image as image
import ctypes as C

class Docker(runtime.Runtime):
	"Docker refers to a parsed docker application object."

	# Initialize the docker object, assuming it is defaultly
	# installed in the path '/var/lib/docker'.
	_new = binding.lookup(b"veinmind_DockerNew", b"VEINMIND_1.0")
	def __init__(self):
		handle = binding.Handle()
		binding.handle_error(Docker._new(handle.ptr()))
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
