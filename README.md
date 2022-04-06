# libVeinMind: 问脉容器感知与安全 SDK

<p>
  <img src="https://img.shields.io/github/release/chaitin/libveinmind.svg" />
  <img src="https://img.shields.io/github/release-date/chaitin/libveinmind.svg?color=blue&label=update" />
</p>

> 容器安全见筋脉，望闻问切治病害。

问脉 (TM) SDK 提供了容器安全领域所关心的容器、镜像和运行时等对象的信息获取、镜像内容器访问等相关操作。并进行合理的抽象以提高容器安全工具对 Docker、Containerd 和 Kubernetes 等不同容器产品的兼容性，简化容器安全工具的开发和维护。

问脉 SDK 是[问脉 (TM) 容器安全开源工具箱](https://github.com/chaitin/veinmind-tools)编译和运行所需的依赖，您也可以基于问脉 SDK 开发符合自己需求的容器安全工具。

问脉 SDK 中的接口部分在本仓库中进行了开源，实现部分采取免费闭源的方式提供给用户。

## 快速开始

一般情况下，应用应以平行容器的方式发布和部署，用户无需额外安装依赖。如何以平行容器方式构建与发布应用，详见[平行容器](docs/parallel-container.md)说明文档。

若应用只支持本地运行，或需要搭建本地环境进行开发，则需要先安装对应平台下的问脉 SDK 软件包。

软件安装包元信息中包含问脉 SDK 的相关许可协议，在开发和使用时请遵守许可协议。当您下载并安装 SDK 软件包后即视为您已同意问脉 SDK 使用协议。

在 Ubuntu 和 Debian 平台下，添加问脉 SDK 的 APT 仓库即可安装所需软件包：

```bash
echo 'deb [trusted=yes] https://download.veinmind.tech/libveinmind/apt/ ./' | sudo tee /etc/apt/sources.list.d/libveinmind.list
sudo apt-get update
sudo apt-get install libveinmind-dev
```

在 RedHat 和 CentOS 平台下，添加问脉 SDK 的 yum 仓库即可安装所需软件包：

```bash
sudo cat > /etc/yum.repos.d/libveinmind.repo << ==EOF==
[libveinmind]
name=libVeinMind SDK yum repository
baseurl=https://download.veinmind.tech/libveinmind/yum/
enabled=1
gpgcheck=0
==EOF==
sudo yum install libveinmind-devel
```

## 开发指南

- API 文档及使用样例（[Golang](https://pkg.go.dev/github.com/chaitin/libveinmind)、[Python3](docs/python-usage.rst)）
- [插件系统](docs/plugin-system.md)（如何开发可复用的容器安全工具）
- [平行容器](docs/parallel-container.md)（如何容器化部署容器安全工具）

## 联系我们

1. 您可以通过 GitHub Issue 直接进行 Bug 反馈和功能建议。

2. 您扫描下方二维码可以通过添加问脉小助手，以加入问脉用户讨论群进行详细讨论：

![](docs/veinmind-group-qrcode.jpg)
