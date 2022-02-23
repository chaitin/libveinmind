使用方式
=============

示例1: 获取当前机器 ``docker`` 镜像 ID 列表
----------------------------------------------

调用 ``docker.Docker`` 类默认的构造函数创建一个客户端，该客户端用于与容器运行时进行交互。

通过客户端调用 ``list_image_ids`` 函数，获得当前机器 ``docker`` 镜像 ID 列表。

.. code:: python

    >>> from veinmind import docker
    >>> client = docker.Docker()
    >>> client.list_image_ids()
    ['sha256:b8604a3fe8543c9e6afc29550de05b36cd162a97aa9b2833864ea8a5be11f3e2', 'sha256:9fb2db53533ec20ca178e7e64ce04fba2a9b0280dd0c65cf5908e667acf77f61']

示例2: 获取名称为 ``ubuntu`` 的镜像中的镜像基本信息
--------------------------------------------------------------------

客户端通过 ``open_image_by_id`` ，可以根据镜像 ID 获得对应的 ``Image`` 实例。

``Image`` 类提供了一系列获取基本信息的函数，如 ``reporefs`` 则是获取镜像对应的 ``reference`` 。

.. code:: python

    from veinmind import docker

    client = docker.Docker()
    ids = client.find_image_ids("ubuntu")
    for id in ids:
        image = client.open_image_by_id(id)
        print("image id: " + image.id())
        for ref in image.reporefs():
            print("image ref: " + ref)
        for repo in image.repos():
            print("image repo: " + repo)
        print("image ocispec: " + str(image.ocispec_v1()))


示例3: 读取名称为 ``ubuntu`` 的镜像中的 ``/etc/passwd`` 文件
------------------------------------------------------------------

``Image`` 类继承了 ``FileSystem`` 相关方法，可以进行文件操作。

通过调用 ``image.open("/etc/passwd")`` 可以打开 ``/etc/passwd`` 文件，使用方式和 Python 中的 ``open`` 函数大体一致。

.. code:: python

    from veinmind import docker

    client = docker.Docker()
    ids = client.find_image_ids("ubuntu")
    for id in ids:
        image = client.open_image_by_id(id)
        try:
            with image.open("/etc/passwd") as f:
                print(f.read())
        except FileNotFoundError as e:
            print("/etc/passwd is not found in " + image.id())

示例4: 遍历 ``redis`` 的镜像中每一层中的所有文件
------------------------------------------------------------

通过 ``Image`` 类，可以获取镜像的 ``Layer`` 实例，通过调用 ``open_layer`` 并传入对应下标即可。

``Layer`` 类继承了 ``FileSystem`` 相关方法，可以进行文件操作，其中 ``layer.walk`` 和 Python 中 ``os.walk`` 函数使用方式相同。

.. code:: python

    from veinmind import docker
    import os

    client = docker.Docker()
    ids = client.find_image_ids("redis")
    for id in ids:
        image = client.open_image_by_id(id)
        for layer_index in range(image.num_layers()):
            layer = image.open_layer(layer_index)
            for root, dirs, files in layer.walk("/"):
                for file in files:
                    filepath = os.path.join(root, file)
                    print(filepath)

