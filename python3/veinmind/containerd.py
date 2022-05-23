from . import binding as binding
from . import runtime as runtime
from . import image as image

class Containerd(runtime.Runtime):
	"Containerd refers to a parsed containerd application object."

	_with_config_path_fn = binding.lookup(
		b"veinmind_ContainerdWithConfigPath", b"VEINMIND_1.1")
	def _with_config_path(path):
		if not isinstance(path, str):
			raise TypeError("path must be str")
		with binding.new_str(path) as hstr:
			hopt = binding.Handle()
			binding.handle_error(Containerd._with_config_path_fn(
				hopt.ptr(), hstr.val()))
			return hopt

	_with_root_dir_fn = binding.lookup(
		b"veinmind_ContainerdWithRootDir", b"VEINMIND_1.1")
	def _with_root_dir(path):
		if not isinstance(path, str):
			raise TypeError("path must be str")
		with binding.new_str(path) as hstr:
			hopt = binding.Handle()
			binding.handle_error(Containerd._with_root_dir_fn(
				hopt.ptr(), hstr.val()))
			return hopt

	_with_unique_desc_fn = binding.lookup(
		b"veinmind_ContainerdWithUniqueDesc", b"VEINMIND_1.1")
	def _with_unique_desc(desc):
		if not isinstance(desc, str):
			raise TypeError("desc must be str")
		with binding.new_str(desc) as hstr:
			hopt = binding.Handle()
			binding.handle_error(Containerd._with_unique_desc_fn(
				hopt.ptr(), hstr.val()))
			return hopt

	# Initialize the docker object, with specified arguments.
	_make_new_option_list = binding.lookup(
		b"veinmind_ContainerdMakeNewOptionList", b"VEINMIND_1.1")
	_new = binding.lookup(
		b"veinmind_ContainerdNew", b"VEINMIND_1.1")
	def __init__(self, **kwargs):
		hopts = binding.Handle()
		binding.handle_error(
			Containerd._make_new_option_list(hopts.ptr()))
		with hopts as hopts:
			config_path = kwargs.pop("config_path", None)
			if config_path is not None:
				with Containerd._with_config_path(config_path) as hopt:
					hopts.append(hopt)
			root_dir = kwargs.pop("root_dir", None)
			if root_dir is not None:
				with Containerd._with_root_dir(root_dir) as hopt:
					hopts.append(hopt)
			unique_desc = kwargs.pop("unique_desc", None)
			if unique_desc is not None:
				with Containerd._with_unique_desc(unique_desc) as hopt:
					hopts.append(hopt)

			handle = binding.Handle()
			binding.handle_error(Containerd._new(
				handle.ptr(), hopts.val()))
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
