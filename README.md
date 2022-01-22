# libVeinMind: 问脉容器感知与安全 SDK

<p>
  <img src="https://img.shields.io/github/release/chaitin/libveinmind.svg" />
  <img src="https://img.shields.io/github/release-date/chaitin/libveinmind.svg?color=blue&label=update" />
</p>

> 容器安全见筋脉，望闻问切治病害。

问脉 (TM) SDK 提供了容器安全领域所关心的容器、镜像和运行时等对象的信息获取、镜像内容器访问等相关操作。并进行合理的抽象以提高容器安全工具对 Docker、Containerd 和 Kubernetes 等不同容器产品的兼容性，简化容器安全工具的开发和维护。

问脉 SDK 是问脉 (TM) 容器安全开源工具箱编译和运行所需的依赖，您也可以基于问脉 SDK 开发符合自己需求的容器安全工具。

问脉 SDK 中的接口部分在本仓库中进行了开源，实现部分采取免费闭源的方式提供给用户。

## 快速开始

要使用问脉 SDK 需要先安装对应平台下的 SDK 软件包。以 Ubuntu 和 Debian 为例，SDK 的安装和配置步骤如下:

1. 添加 libVeinMind 的 APT 包源并更新 APT 包列表和索引，通过执行以下指令：

    ```bash
    echo 'deb [trusted=yes] https://download.veinmind.tech/libveinmind/ ./' | sudo tee /etc/apt/sources.list.d/libveinmind.list
    sudo apt-get update
    ```

2. 执行指令 `sudo apt-get install libveinmind-dev` 以开始安装问脉。

3. 安装过程中，需阅读并同意问脉 SDK 的用户协议，方可完成 SDK 的安装与配置。

## 联系我们

1. 您可以通过 GitHub Issue 直接进行 Bug 反馈和功能建议。

2. 您扫描下方二维码可以通过添加问脉小助手，以加入问脉用户讨论群进行详细讨论：

![](docs/veinmind-group-qrcode.jpg)
