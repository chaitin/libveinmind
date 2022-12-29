from . import binding as binding
from . import runtime as runtime
from . import image as image
from . import filesystem as filesystem
import ctypes as C


class Tarball(runtime.Runtime):
    _new = binding.lookup(b"veinmind_TarballNew", b"VEINMIND_1.3")

    def __init__(self, root):
        with binding.new_str(root) as hstr:
            handle = binding.Handle()
            binding.handle_error(Tarball._new(
                handle.ptr(), hstr.val()))
            super(Tarball, self).__init__(handle=handle)

    _open_image_by_id = binding.lookup(
        b"veinmind_RuntimeOpenImageByID", b"VEINMIND_1.0")

    def open_image_by_id(self, image_id):
        with binding.new_str(image_id) as hstr:
            handle = binding.Handle()
            binding.handle_error(Tarball._open_image_by_id(
                handle.ptr(), self.__handle__().val(), hstr.val()))
            return Image(handle)

    _remove_image_by_id = binding.lookup(
        b"veinmind_TarballRemoveImageByID", b"VEINMIND_1.3")

    def remove_image_by_id(self, image_id):
        with binding.new_str(image_id) as hstr:
            return binding.handle_error(Tarball._remove_image_by_id(
                self.__handle__().val(), hstr.val()))

    _load = binding.lookup(
        b"veinmind_TarballLoad", b"VEINMIND_1.3")

    def load(self, path):
        with binding.Handle() as handle:
            with binding.new_str(path) as hstr:
                binding.handle_error(Tarball._load(
                    handle.ptr(), self.__handle__().val(), hstr.val()))
                return handle.str_list()


class Layer(filesystem.FileSystem):
    "Layer refers to a layer in docker image."

    # Initialize the docker layer object.
    def __init__(self, handle):
        super(Layer, self).__init__(handle=handle)

    _id = binding.lookup(b"veinmind_TarballLayerID", b"VEINMIND_1.3")

    def id(self):
        "Retrieve the diff ID of the docker layer."

        handle = binding.Handle()
        binding.handle_error(Layer._id(
            handle.ptr(), self.__handle__().val()))
        with handle as handle:
            return handle.str()


class Image(image.Image):
    _open_layer = binding.lookup(
        b"veinmind_TarballImageOpenLayer", b"VEINMIND_1.3")
    def open_layer(self, i):
        "Open specified layer in the docker image."

        handle = binding.Handle()
        binding.handle_error(Image._open_layer(
            handle.ptr(), self.__handle__().val(), C.c_size_t(i)))
        return Layer(handle)

    _num_layers = binding.lookup(
        b"veinmind_TarballImageNumLayers", b"VEINMIND_1.3")
    def num_layers(self):
        "Return the number of layers in the tarball image."

        result = C.c_size_t()
        binding.handle_error(Image._num_layers(
            C.pointer(result), self.__handle__().val()))
        return result.value

    _opaques = binding.lookup(
        b"veinmind_TarballLayerOpaques", b"VEINMIND_1.5")
	def opaques(self):
		"Retrieve the opaques of the tarball layer."

		handle = binding.Handle()
		binding.handle_error(Layer._opaques(
			handle.ptr(), self.__handle__().val()))
		with handle as handle:
			return handle.str_list()
	
	_whiteouts = binding.lookup(
        b"veinmind_TarballLayerWhiteouts", b"VEINMIND_1.5")
	def whiteouts(self):
		"Retrieve the whiteouts of the tarball layer."

		handle = binding.Handle()
		binding.handle_error(Layer._whiteouts(
			handle.ptr(), self.__handle__().val()))
		with handle as handle:
			return handle.str_list()