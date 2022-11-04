from . import binding as binding
from . import runtime as runtime
from . import image as image


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
            return image.Image(handle)

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
