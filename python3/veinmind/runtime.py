from . import binding as binding
from . import container

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

	# Retrieve all container IDs of the container runtime.
	_list_container_ids = binding.lookup(b"veinmind_RuntimeListContainerIDs", b"VEINMIND_1.2")
	def list_container_ids(self):
		"List all container IDs managed by current runtime."

		handle = binding.Handle()
		binding.handle_error(Runtime._list_container_ids(
			handle.ptr(), self.__handle__().val()
		))
		with handle as handle:
			return handle.str_list()

	# Find container IDs by implementation specific characteristics.
	_find_container_ids = binding.lookup(
		b"veinmind_RuntimeFindContainerIDs", b"VEINMIND_1.2")
	def find_container_ids(self, pattern):
		"Attempt to match and find container IDs by pattern."

		with binding.new_str(pattern) as hstr:
			handle = binding.Handle()
			binding.handle_error(Runtime._find_container_ids(
				handle.ptr(), self.__handle__().val(), hstr.val()))
			with handle as handle:
				return handle.str_list()

	# Open a container by its ID and return the container object.
	_open_container_by_id = binding.lookup(
		b"veinmind_RuntimeOpenContainerByID", b"VEINMIND_1.2")
	def open_container_by_id(self, container_id):
		"Open a container by its specified ID."

		with binding.new_str(container_id) as hstr:
			handle = binding.Handle()
			binding.handle_error(Runtime._open_container_by_id(
				handle.ptr(), self.__handle__().val(), hstr.val()))
			return container.Container(handle)