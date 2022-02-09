from . import binding
from . import filesystem
import io
import json

class Image(filesystem.FileSystem):
	def __init__(self, handle):
		super(Image, self).__init__(handle=handle)

	# Retrieve the image ID by its open handle.
	_id = binding.lookup(b"veinmind_ImageID", b"VEINMIND_1.0")
	def id(self):
		"Retrieve runtime specific image ID."

		hstr = binding.Handle()
		binding.assert_no_error(Image._id(
			hstr.ptr(), self.__handle__().val()))
		with hstr as hstr:
			return hstr.str()

	# Retrieves the correlated repos of the image.
	_repos = binding.lookup(b"veinmind_ImageRepos", b"VEINMIND_1.0")
	def repos(self):
		"Retrieve all correlated repos."

		hrepos = binding.Handle()
		binding.handle_error(Image._repos(
			hrepos.ptr(), self.__handle__().val()))
		with hrepos as hrepos:
			return hrepos.str_list()

	# Retrieves the correlated references of the image.
	_repo_refs = binding.lookup(b"veinmind_ImageRepoRefs", b"VEINMIND_1.0")
	def reporefs(self):
		"Retrieve all correlated repo references."

		hrefs = binding.Handle()
		binding.handle_error(Image._repo_refs(
			hrefs.ptr(), self.__handle__().val()))
		with hrefs as hrefs:
			return hrefs.str_list()

	# Retrieve the parsed json of image OCI Spec format.
	_ocispec_v1_marshal_json = binding.lookup(
		b"veinmind_ImageOCISpecV1MarshalJSON", b"VEINMIND_1.0")
	def ocispec_v1(self):
		"Retrieve the image OCI Specification information."

		hspec = binding.Handle()
		binding.handle_error(Image._ocispec_v1_marshal_json(
			hspec.ptr(), self.__handle__().val()))
		hstr = None
		with hspec as hspec:
			hstr = hspec.bytes_to_str()
		with hstr as hstr:
			return json.load(io.StringIO(hstr.str()))
