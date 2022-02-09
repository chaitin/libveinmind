from . import binding as binding
from . import runtime as runtime
from . import image as image

class Containerd(runtime.Runtime):
	"Containerd refers to a parsed containerd application object."

	# Initialize the docker object, assuming it is defaultly
	# installed in the path '/var/lib/containerd'.
	_new = binding.lookup(
		b"veinmind_ContainerdNew", b"VEINMIND_1.0")
	def __init__(self):
		handle = binding.Handle()
		binding.handle_error(Containerd._new(handle.ptr()))
		super(Containerd, self).__init__(handle=handle)

	# Open a image by its ID and return the image object.
	_open_image_by_id = binding.lookup(
		b"veinmind_RuntimeOpenImageByID", b"VEINMIND_1.0")
	def open_image_by_id(self, image_id):
		"Open a containerd image by its specified ID."

		with binding.new_str(image_id) as hstr:
			handle = binding.Handle()
			binding.handle_error(Containerd._open_image_by_id(
				handle.ptr(), self.__handle__().val(), hstr.val()))
			return Image(handle)

class Image(image.Image):
	"Image refers to a contaierd specific image."

	# Initialize the docker image object.
	def __init__(self, handle):
		super(Image, self).__init__(handle=handle)
