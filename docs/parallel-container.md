# 平行容器

平行容器是指将云原生安全（如容器安全、镜像安全等）应用在其目标云原生平台进行容器化部署的方案。将安全应用进行容器化部署，可以充分利用云原生平台提供的集群化部署、容器编排、高可用、自动恢复等特性，提升安全产品的鲁棒性，降低用户的使用难度和运维成本。

问脉 SDK 在设计时充分考虑到了平行容器的场景，并且提供了能覆盖绝大多数使用场景的基础镜像。在绝大多数情况下，问脉 SDK 及其基础镜像对宿主机和平行容器的运行环境差异进行了统一化处理，因此问脉 SDK 的用户编写应用时无需感知当前是否处于平行容器内。

与之相对，大部分云原生平台运行容器时都需要对容器的挂载卷、命名空间和 Capabilities 等参数进行详细配置，并通过 docker 命令行、`docker-compose.yml` 文件和 Kubernetes Pod 的 yaml 文件等告知云原生平台如何创建和配置容器。而基于问脉 SDK 的基础镜像创建平行容器时，若未对平行容器的运行环境进行正确配置，则会导致平行容器内的问脉 SDK 无法正常工作，从而影响上层应用的运行。

以 python 为例，基于问脉 SDK 的基础镜像创建平行容器并使用问脉 SDK 的示例如下：

```bash
$ docker run --rm -it --mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave' veinmind/python3:1.0.2-stretch python3
Python 3.5.10 (default, Sep 10 2020, 18:47:38) 
[GCC 6.3.0 20170516] on linux
Type "help", "copyright", "credits" or "license" for more information.
>>> from veinmind import *
>>> d = docker.Docker()
>>> d.find_image_ids('python3')
['sha256:ff2047e9fe4b823a15d98d1d0622a235caad42dc0384518e351b9ba920d3fd39', 'sha256:2048b775ce1f9b53d91d637bcb0bfe34d8e463cb07de9bc113c4604035e3a2fd']
```

倘若把其中的 `--mount` 参数去掉，则在平行容器内运行的问脉 SDK 将因为配置错误而无法工作（输出有删减）：

```bash
$ docker run --rm -it veinmind/python3:1.0.2-stretch python3
Python 3.5.10 (default, Sep 10 2020, 18:47:38) 
[GCC 6.3.0 20170516] on linux
Type "help", "copyright", "credits" or "license" for more information.
>>> from veinmind import *
>>> d = docker.Docker()
Traceback (most recent call last):
...
  File "/usr/local/lib/python3.5/site-packages/veinmind/binding.py", line 223, in _handle_syscall_error
    raise OSError(errvalue, os.strerror(errvalue))
FileNotFoundError: [Errno 2] No such file or directory: '/host/var/lib/docker'
...
```

诚然要提供一个“一把梭”命令/配置来最简化平行容器的创建是很容易的，只需要尽可能消除容器内外的区别并给予足够高的访问权限即可。譬如在创建 docker 容器时通过 `--privileged` 给予全部 Capabilities，将平行容器的除了 mount 以外的命名空间全部指定为 `host`，再加上本文描述的和问脉 SDK 运行环境相关的一些配置，即可简单地创建这样的容器。然而这样的容器被创建出来，将在用户的宿主机上展开一个巨大的暴露面，入侵者只需要成功利用 Linux 内核、云原生平台或容器内资产中的漏洞或配置错误，即可以获得当前宿主机上的 root 权限并“入驻”当前机器，这显然是我们安全从业人员需要尽力避免的。

因此比起直接给一个“一把梭”指令创建平行容器，我们更希望使用问脉 SDK 的应用在部署应用给终端用户时，其平行容器的创建配置是经过审慎编写的，在满足其实际运行需求的同时，遵循最小权限原则。本文将对现阶段问脉 SDK 和基础镜像中与平行容器有关的配置、含义及其安全影响进行尽可能详细的阐述，以供应用的编写者在充分理解其含义和后果的前提下作出准确的判断。

## 容器配置

### 宿主机根文件系统

```
--mount 'type=bind,source=/,target=/host,readonly,bind-propagation=rslave'
```

我们应该非常熟悉通过 Dockerfile 等构建方式构建镜像，然后通过构建好的镜像创建容器的使用方式：Dockerfile 构建时会添加应用运行所需的依赖和配置等文件，它们逐步组成了支持应用运行的根文件系统。根文件系统包含在镜像中进行持久化和分发，并且在容器运行时加载以运行容器。而由于应用运行所需的依赖与配置均包含在了镜像中，因此目标宿主机只需要安装了相应的容器运行时 [^1] 即可运行应用，无需再在宿主机上安装对应的依赖或进行配置。

在基于问脉 SDK 构建安全应用的场景，应用的创建者可以选用合适的基础镜像并安装其应用所需的依赖，以为其应用提供运行环境，这与常规的实践是一致的。而问脉 SDK 工作时，需要访问宿主机上相应容器运行时和云原生平台的本地文件，读取相关数据和配置，以构建容器运行时、镜像、容器等 API 对象供 SDK 用户访问和操作。因此在平行容器内，除了访问容器根文件系统外，还需要能访问宿主机根文件系统。访问容器根文件系统是容器技术常规实践的要求，而访问宿主机根文件系统是问脉 SDK 正常运行的需求。

在 Linux 下，根文件系统是以 `/` 为根节点组织的单一树状结构，因此容器的根文件系统和宿主机的根文件系统不可能同时存在于 `/` 下。但是可以利用 Linux 的挂载点机制，创建 `/host` 目录作为挂载点并将宿主机的文件系统挂载于其下，即可实现同时访问宿主机和容器的根文件系统。

以 Docker 为例，依据其[用户文档](https://docs.docker.com/storage/bind-mounts/)所述，在创建容器时指定 `--mount` 参数创建 bind 挂载点，将宿主机上的 `/` 路径挂载到容器内的 `/host`，即通过 `--mount 'type=bind,source=/,target=/host'` 可完成宿主机文件系统的到容器内的映射。

宿主机的根文件系统上本来就存在许多挂载点，若只挂载 `/` 处的挂载点所包含的文件系统，可能会导致需要访问的关键文件或目录未被映射到容器内。因此挂载时需要指定 `bind-propagation` 选项为 `rslave`，使 `/` 目录树下的其他挂载点也被递归地映射到容器内。

而 `readonly` 选项则限制了应用通过 `/host` 对宿主根文件系统的修改操作。

## 容器镜像

| 容器名称 | 基础镜像 | 描述 |
|----------|----------|------|
| [veinmind/base:\*-stretch](https://hub.docker.com/repository/docker/veinmind/base) | [buildpack-deps:stretch-scm](https://hub.docker.com/_/buildpack-deps) | 安装了问脉 SDK 的基础镜像 |
| [veinmind/python3:\*-stretch](https://hub.docker.com/repository/docker/veinmind/python3) | [python:3.5-stretch](https://hub.docker.com/_/buildpack-deps) | 安装了问脉 SDK 和 Python 3.5 基础镜像 |
| [veinmind/go1.16:\*-stretch](https://hub.docker.com/repository/docker/veinmind/python3) | [golang:1.16-stretch](https://hub.docker.com/_/buildpack-deps) | 安装了问脉 SDK 和 Go 1.16 基础镜像 |

[^1]: 此处为了易于描述和理解，隐去了 CPU 架构匹配、内核版本不能过低等细节。
