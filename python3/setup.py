import os
import os.path
import re
import subprocess
import setuptools

# Generate version based on git when it is managed there.
version = "unknown"

# If the file is managed by git, always fetch the version from
# executing "git describe".
if os.path.exists("../.git"):
	# Retrieve and split version output from "git describe".
	output = subprocess.getoutput("git describe --tags")
	m = re.match(r"^v(\d+[.]\d+[.]\d+)(-(\d+)-(\w+))?$", output)
	if m is None:
		raise ValueError("invalid version: {!r}".format(output))
	tag, extra, rev, commit = m.groups()

	# Generate version conforming to PEP-440 to play with pip.
	version = tag
	if extra is not None:
		version = "{tag}.dev{rev}+{commit}".format(
			tag=tag, rev=rev, commit=commit)

# If it is currently inside some python virtualenv, attempt to
# bail out by inspecting the veinmind.egg-info/PKG-INFO file.
elif os.path.exists("veinmind.egg-info/PKG-INFO"):
	with open("veinmind.egg-info/PKG-INFO") as pkginfo:
		for line in pkginfo.readlines():
			m = re.match(r"^Version: (.*)$", line)
			if m is not None:
				version = m.groups()[0]

# Invoke the setuptools for setting up module.
with open("README.md", "r") as readme:
	long_description = readme.read()
setuptools.setup(
	name="veinmind", version=version,
	url="https://github.com/chaitin/libveinmind",
	author="Haoran Luo",
	author_email="haoran.luo@chaitin.com",
	description="libVeinMind API python binding",
	long_description=long_description,
	long_description_content_type="text/markdown",
	project_urls={
		"Bug Tracker": "https://github.com/chaitin/libveinmind/issues",
	},
	classifiers=[
		"Programming Language :: Python :: 3",
		"License :: Free To Use But Restricted",
		"Operating System :: POSIX :: Linux",
	],
	package_dir={"": "."},
	packages=setuptools.find_packages(),
	install_requires=[
		"click==7.1.2",
	],

	# PEP-344: raise from: >=2.5
	# PEP-380: yield from: >=3.3
	python_requires=">=3.3")
