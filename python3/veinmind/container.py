from . import binding
from . import filesystem
from . import process
import io
import json


class Container(filesystem.FileSystem, process.Psutil):
    def __init__(self, handle):
        super(Container, self).__init__(handle=handle)

    # Retrieve the container ID by its open handle.
    _id = binding.lookup(b"veinmind_ContainerID", b"VEINMIND_1.2")
    def id(self):
        "Retrieve runtime specific container ID."

        hstr = binding.Handle()
        binding.assert_no_error(Container._id(
            hstr.ptr(), self.__handle__().val()))
        with hstr as hstr:
            return hstr.str()

    # Retrieve the container Name by its open handle.
    _name = binding.lookup(b"veinmind_ContainerName", b"VEINMIND_1.2")
    def name(self):
        "Retrieve runtime specific container ID."

        hstr = binding.Handle()
        binding.assert_no_error(Container._name(
            hstr.ptr(), self.__handle__().val()))
        with hstr as hstr:
            return hstr.str()

    # Retrieve the image ID associated with the container by its open handle
    _image_id = binding.lookup(b"veinmind_ContainerImageID", b"VEINMIND_1.2")
    def image_id(self):
        "Retrieve runtime specific image ID associated with the container."

        hstr = binding.Handle()
        binding.assert_no_error(Container._image_id(
            hstr.ptr(), self.__handle__().val()
        ))
        with hstr as hstr:
            return hstr.str()

    # Retrieve the parsed json of runtime OCI Spec format.
    _ocispec_marshal_json = binding.lookup(b"veinmind_ContainerOCISpecMarshalJSON", b"VEINMIND_1.2")
    def ocispec(self):
        "Retrieve the runtime OCI Specification information."

        hspec = binding.Handle()
        binding.handle_error(Container._ocispec_marshal_json(
            hspec.ptr(), self.__handle__().val()))
        hstr = None
        with hspec as hspec:
            hstr = hspec.bytes_to_str()
        with hstr as hstr:
            return json.load(io.StringIO(hstr.str()))

    # Retrieve the parsed json of runtime OCI State format.
    _ocistate_marshal_json = binding.lookup(b"veinmind_ContainerOCIStateMarshalJSON", b"VEINMIND_1.2")
    def ocistate(self):
        "Retrieve the runtime OCI Specification information."

        hspec = binding.Handle()
        binding.handle_error(Container._ocistate_marshal_json(
            hspec.ptr(), self.__handle__().val()))
        hstr = None
        with hspec as hspec:
            hstr = hspec.bytes_to_str()
        with hstr as hstr:
            return json.load(io.StringIO(hstr.str()))