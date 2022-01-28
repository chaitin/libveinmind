from . import binding as binding

class Runtime(binding.Object):
	def __init__(self, handle):
		super(Runtime, self).__init__(handle=handle)

	# Retrieve all image IDs of the docker object.
	_list_image_ids = binding.lookup(
		b"veinmind_RuntimeListImageIDs", b"VEINMIND_1.0")
	def list_image_ids(self):
		"List all image IDs managed by current runtime."

		handle = binding.Handle()
		binding.handle_error(Runtime._list_image_ids(
			handle.ptr(), self.__handle__().val()))
		with handle as handle:
			return handle.str_list()

	# Find image IDs by implementation specific characteristics.
	_find_image_ids = binding.lookup(
		b"veinmind_RuntimeFindImageIDs", b"VEINMIND_1.0")
	def find_image_ids(self, pattern):
		"Attempt to match and find image IDs by pattern."

		with binding.new_str(pattern) as hstr:
			handle = binding.Handle()
			binding.handle_error(Runtime._find_image_ids(
				handle.ptr(), self.__handle__().val(), hstr.val()))
			with handle as handle:
				return handle.str_list()
