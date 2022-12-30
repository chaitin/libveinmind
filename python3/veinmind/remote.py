from . import binding as binding
from . import runtime as runtime
from . import image as image
from . import filesystem as filesystem
import ctypes as C


class Remote(runtime.Runtime):
    _new = binding.lookup(b"veinmind_RemoteNew", b"VEINMIND_1.4")

    def __init__(self, root):
        with binding.new_str(root) as hstr:
            handle = binding.Handle()
            binding.handle_error(Remote._new(
                handle.ptr(), hstr.val()))
            super(Remote, self).__init__(handle=handle)

    _open_image_by_id = binding.lookup(
        b"veinmind_RuntimeOpenImageByID", b"VEINMIND_1.0")

    def open_image_by_id(self, image_id):
        with binding.new_str(image_id) as hstr:
            handle = binding.Handle()
            binding.handle_error(Remote._open_image_by_id(
                handle.ptr(), self.__handle__().val(), hstr.val()))
            return Image(handle)

    _load = binding.lookup(
        b"veinmind_RemoteLoad", b"VEINMIND_1.4")

    def load(self, image_ref,username,password):
        with binding.Handle() as handle:
            with binding.new_str(image_ref) as hstr:
                with binding.new_str(username) as ustr:
                    with binding.new_str(password) as pstr:
                        binding.handle_error(Remote._load(
                            handle.ptr(), self.__handle__().val(), hstr.val(),ustr.val(),pstr.val()))
                        return handle.str_list()


class Layer(filesystem.FileSystem):
    "Layer refers to a layer in docker image."

    # Initialize the docker layer object.
    def __init__(self, handle):
        super(Layer, self).__init__(handle=handle)

    _id = binding.lookup(b"veinmind_RemoteLayerID", b"VEINMIND_1.4")
    def id(self):
        "Retrieve the diff ID of the remote layer."

        handle = binding.Handle()
        binding.handle_error(Layer._id(
            handle.ptr(), self.__handle__().val()))
        with handle as handle:
            return handle.str()

    _opaques = binding.lookup(b"veinmind_RemoteLayerOpaques", b"VEINMIND_1.5")
    def opaques(self):
        "Retrieve the opaques of the remote layer."

        handle = binding.Handle()
        binding.handle_error(Layer._opaques(
            handle.ptr(), self.__handle__().val()))
        with handle as handle:
            return handle.str_list()
	
    _whiteouts = binding.lookup(b"veinmind_RemoteLayerWhiteouts", b"VEINMIND_1.5")
    def whiteouts(self):
        "Retrieve the whiteouts of the remote layer."

        handle = binding.Handle()
        binding.handle_error(Layer._whiteouts(
            handle.ptr(), self.__handle__().val()))
        with handle as handle:
            return handle.str_list()


class Image(image.Image):
    _open_layer = binding.lookup(
        b"veinmind_RemoteImageOpenLayer", b"VEINMIND_1.4")
    def open_layer(self, i):
        "Open specified layer in the docker image."

        handle = binding.Handle()
        binding.handle_error(Image._open_layer(
            handle.ptr(), self.__handle__().val(), C.c_size_t(i)))
        return Layer(handle)

    _num_layers = binding.lookup(
        b"veinmind_RemoteImageNumLayers", b"VEINMIND_1.4")
    def num_layers(self):
        "Return the number of layers in the tarball image."

        result = C.c_size_t()
        binding.handle_error(Image._num_layers(
            C.pointer(result), self.__handle__().val()))
        return result.value
    
    _get_layer_diff_id = binding.lookup(
        b"veinmind_RemoteImageGetLayerDiffID", b"VEINMIND_1.5")
    def get_layer_diff_id(self, i):
        "Retrieve the diff ID the of remote layer without opening it."
        
        handle = binding.Handle()
        binding.handle_error(Image._get_layer_diff_id(
            handle.ptr(), self.__handle__().val(), C.c_size_t(i)))
        with handle as handle:
            return handle.str()