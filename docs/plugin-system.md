# 插件系统

## 为什么我们需要一套插件系统

当我们开发了一个容器安全工具，我们自然而然地会想在各种场景下被复用该工具。譬如，当我们写好了一个镜像后门扫描工具，会想在以下场景使用它：

1. 本地扫描：当我们直接执行这个工具时，它能够对本地镜像进行扫描，并且输出镜像中的后门的检测结果；
2. 操作授权：我们可以编写一个 [Docker 授权插件](https://docs.docker.com/engine/extend/plugins_authorization/)，在拉取镜像和启动容器时，自动执行该工具进行镜像后门检测，并阻断风险操作；
3. 远程扫描：我们希望使用这个工具对远程仓库中的镜像进行扫描，发现存在后门的镜像，以为后续的修复、溯源和封禁操作提供支持。

除此以外，我们还能想到很多其他的使用场景，譬如集成到 CI/CD 流程中进行扫描等。但我们不妨先以上述范围作为典型展开讨论。

一般情况下，为了方便开发和调试，我们编写的安全工具都会先支持本地扫描的功能。事实上我们大多数时候都“超量完成”了这个目标：我们会花很多时间在装饰这个工具的输出功能上，可能这个工具一半的代码都花在了命令行配色、输出格式化和各式的报告生成等输出功能的优化上。

现在让我们开始思考如何支持操作授权。显然可以在镜像后门扫描工具的基础上添加一个子命令，并把它封装为一个 Docker 授权插件，可是这要处理很多繁复细节。因此我们会希望存在一个开源的 Docker 授权插件项目，它以子进程的方式调用我们的工具，并且依据工具的检测结果决定是否要阻断用户的操作。如果这样的 Docker 授权插件存在（或者若不存在我们可以自己做一个），那必然会减少很多的我们工作量。在我们开始寻找或编写这样的 Docker 授权插件之前，不妨先思考一个问题：它应该如何获取我们的检测结果呢？

先考虑不对工具本身进行任何修改的方案。能不能直接读取工具的输出和报告，利用正则表达式等匹配出检测结果呢？除却使用正则表达式进行模式匹配本身可能会产生的错误匹配，可以相信的一点是，如果放任不同的安全工具作者对输出结果进行自由发挥，一千个工具可能就有一千种输出模式。因此要想基于正则匹配生成一个通用方案，只能让工具作者提供匹配其工具输出的正则表达式。而相信大部分人都会觉得这个方案处处散发着坏味道，且宁可对工具进行一些修改也不愿意接受这个方案。

如果允许对工具本身进行一些修改，又会有哪些可行的方案呢？既然我们依靠了子进程的方式执行工具，一个最简单的做法是判断命令的退出状态码。譬如状态码为 0 时认为当前检查没有发现威胁，否则认为发现了威胁并进行阻断。这确实是一个可行的而且容易接受的方案，但仍然存在一些问题。

首先状态码常用于表达程序是否成功执行并退出，那么当一个工具执行失败并返回非零状态码 [^1] 时，是否蕴含了应该阻断用户操作的语义呢？有的用户会认为所有安全工具都应该正常执行完毕且报告无威胁后才应该放行用户操作，这样才能确保不会放过任何一个安全问题；有的用户则会认为仍在试用阶段的安全工具，若不能正常工作则应该忽略其结果，只需要确保处于生产阶段的安全工具正常工作即可。因此一刀切地处理非零状态码，并不一定能满足所有用户的使用需求。

其次在发生阻断时，很多时候都希望把完整的阻断理由（譬如在哪个路径下发现了哪种后门，哪个软件资产有哪个重大漏洞）呈现给用户，但是子进程的状态码是不包含这些信息的。而如果工具没有其他渠道报告具体的检测结果的话，只能去尝试解析工具的输出或报告，并重蹈我们先前认为充满坏味道的正则匹配老路。

因此若想以一种合理且可维护的方法提供检测结果给前文所述的 Docker 授权插件，工具应该输出一个计算机可读、格式预定义的检测结果，以供 Docker 授权插件进行解析和处理。

让我们再考虑如何支持远程扫描。我们会希望编写的工具能提供一个子命令，接收远程仓库 URL 及其认证信息，而该子命令只需要把相应的镜像 tar 下载下来，执行原来的扫描代码即可。

在只有单个安全工具需要执行时，这个思路很自然，但是很多时候我们会不止关心一种威胁，因此会同时执行多个检测工具。在这种情况下，若各个工具之间都各自独立下载远程镜像，除了会因为重复下载镜像 tar 文件（注意到镜像中一个 Layer 的大小经常在几百 MB 量级），占用大量网络带宽并且扫描效率低下外，部分公有仓库如 Docker Hub 还有[下载频率限制](https://docs.docker.com/docker-hub/download-rate-limit/)，对于 Docker Hub 的非付费用户而言，其能扫描的镜像个数随着其使用的工具个数呈反比例关系迅速下降，用户体验极差 [^2]。

一个简单可行的解决方法是，执行每一个工具之前，先把待扫描镜像下载到本地（不管是直接下载 tar 还是使用已经存在的容器运行的 Pull 指令），然后执行各工具进行扫描，待扫描完成后再移除镜像释放资源。值得注意的是，这样就把针对远程仓库的扫描转化为本地扫描了，而进行扫描时只需要为各工具指定下载好的镜像，而无需让工具感知当前是否在扫描远程仓库。这也就意味着工具即使只支持本地扫描，也可以在远程扫描的场景下复用。同时为了方便使用，我们往往会编写一个远程扫描入口程序，接收扫描指令，并依此完成镜像下载、工具调用和镜像卸载的工作。

除了操作授权和远程扫描外，我们还能针对很多其他的使用场景展开实现细节上的讨论，并且总能发现在新的场景下，工具总需要作出或多或少的调整才能完美满足需求。事实上，使用场景的种类之多和各场景之间的差异之大，已然让在工具或者 SDK 中通过堆砌代码来应对成为不可能完成的任务。因此，我们会通过抽象和适配等手段，想方设法将新的场景处理为已知的场景（如远程扫描中将扫描远程仓库处理为扫描本地仓库），这样才能达成在不同的使用场景下复用已有工具。

至此，一个插件系统的想法呼之欲出：如前文所述的 Docker 授权插件、远程扫描的入口程序，我们将这样的程序称为宿主程序（Host Program），它们负责处理各使用场景下的具体细节，将其转化为容器、镜像等具体实体的扫描问题；与之相对的，如前文所述的镜像后门扫描工具、镜像漏洞扫描工具，我们将这样的程序称为插件（Plugin），它们则能对容器、镜像等具体实体进行扫描，并发现其中存在的安全问题。

针对某一具体的使用场景，其相关细节的处理是相对固定的，往往有一个与之相对应的宿主程序；而对于某一具体的容器、镜像等实体，往往需要检查多个方面的安全问题，每个方面的安全问题往往有与之对应的插件可以进行检查。因此宿主程序和插件属于一对多的关系，在实际使用中我们往往先配置好了具体的宿主程序，然后依据我们所关心的安全问题插拔安全插件。

而宿主进程和插件之间并非只有简单的而单向的父子进程的调用关系。在前文讨论如何支持操作授权时，宿主进程就需要收集插件输出的检测结果并处理，生成操作授权响应。在本文未讨论到的场景中，也存在其他采集检测结果并进行定制化处理的需求，如生成检测报告文件，生成 Syslog 并转发到 SIEM 等。同理，针对插件的日志输出，会希望支持进行采集和过滤等操作，以支持不同场景下的日志持久化、工具调试和界面展示等需求。为了处理这些插件产生的、需要可定制化处理的数据，会需要打通从插件到宿主进程的通路，并通过宿主进程向插件提供服务（Service）的方式，在插件中通过调用服务产生数据，然后在宿主进程接收这些数据并进行具体处理 [^3]。

在问脉 SDK 中，我们除了针对容器安全相关实体设计了相应的 API 外，还设计并实现了一套插件系统。基于插件系统编写的安全工具只需一次编写，并在本地扫描的场景中验证，便可集成到不同使用场景下的宿主程序中得到复用；同样地，而对于新的使用场景，只需基于插件系统编写新的宿主程序，便可复用现有的已经编写好的插件。

## 从零开始的插件系统

### 你好，插件

我们使用 Python 语言来编写我们的第一个插件，在此前请先[安装问脉 SDK 软件包](../README.md#快速开始)，并通过 `sudo pip3 install veinmind` [^4] 安装问脉 SDK 的 Python Binding。

创建文件 `hello-plugin`，向文件内写入以下内容，并通过 `chmod a+x ./hello-plugin` 赋予该文件可执行权限：

```python
#!/usr/bin/env python3
from veinmind import *
from os.path import join
from stat import *

command.set_manifest(name="hello-plugin", version="1.0.0")

@command.image_command()
def scan(image):
    """Find executables inside images"""
    reporefs = image.reporefs()
    name = reporefs[0] if len(reporefs) > 0 else image.id()
    log.info('image %s scan start', name)
    for root, _, filenames in image.walk('/'):
        for filename in filenames:
            filepath = join(root, filename)
            mode = image.lstat(filepath).st_mode
            if S_ISREG(mode) and (S_IMODE(mode) & 0o111) != 0:
                log.info('image %s has executable: %s', name, filepath)
    log.info('image %s scan done', name)

if __name__ == '__main__':
    command.main()
```

这可能看上去比一般的 Hello World 程序要复杂一些，不过我相信，一个包含容器操作的样例比起简单地打印一串 Hello World 文本更适合作为一个容器安全 SDK 的 Hello World。

上述样例除了一些结果输出的美化处理外，最核心的地方是一个对镜像内文件进行处理的循环。事实上，这个循环只不过是下述代码的“容器化版本”：

```python
#!/usr/bin/env python3
import os
from os.path import join
from stat import *

for root, _, filenames in os.walk('/'):
    for filename in filenames:
        filepath = join(root, filename)
        mode = os.lstat(filepath).st_mode
        if S_ISREG(mode) and (S_IMODE(mode) & 0o111) != 0:
            print('found executable %s' % (filepath))
```

如果读者比较熟悉 Python 的话，不难发现这段代码的功能是从宿主机根目录开始遍历，并打印其中的可执行文件。而 `hello-plugin` 中的代码不过是把 `os` 替换为了当前正在处理的镜像 `image`，因此不难猜出这段代码的功能是从镜像 `image` 的根目录开始遍历，并打印其中的可执行文件。

通过执行 `sudo ./hello-plugin scan` 指令即可印证我们的猜测：执行该指令时，`hello-plugin` 对本机上所有镜像进行了扫描，并打印了每个镜像中的可执行文件。我们也可以在其后添加一个或多个镜像名称或镜像 ID 等参数，如 `sudo ./hello-plugin scan nginx 1a2d3c`，来限定所需扫描的镜像范围。

我们详细考察一下 `hello-plugin` 的代码细节，注意到通过 `from veinmind import *` 我们导入了问脉 SDK 所提供的 `command` 和 `log` 模块，而 `command` 模块提供了 `command.image_command` 和 `command.main` 等函数，为我们处理了 `hello-plugin` 指令执行的诸多细节，简化了我们的开发。

不过相信看到这里，有很多读者会有疑问：迄今为止所展示的样例中，`hello-plugin` 行为怎么看都像是一个可以独立执行的工具，难以看出为何会被称为插件，以及其与本文所谓插件系统之间的联系。

事实上，在插件系统中，插件和宿主程序是成对的概念，脱离宿主程序是无法清晰地解释插件的概念的。为了简单地解释这个问题，我们不妨马上进入第一个宿主程序的编写。

### 你好，宿主程序

我们不妨针对前文所述的远程扫描的场景编写一个简单的宿主程序。依据前文的思路，这样的宿主程序将接收一个远程仓库列表作为参数，调用宿主机上的容器运行时下载相应的远程仓库，并调用插件对仓库进行扫描。

我们选择 Containerd 作为容器运行时，利用它提供的命名空间功能，创建命名空间并在其中进行操作，以避免我们影响宿主机上其他容器的运行。一般情况下，Docker 都会以 Containerd 作为所依赖的下层容器运行时，这就意味着在安装了 Docker 的机器上我们不必额外安装 Containerd。

我们使用 Go 语言来编写宿主程序，以便使用 Containerd 提供的 API。创建文件 `hello-host.go`，并写入以下内容：

```go
package main

import (
	"context"
	"os"
	"path"

	"github.com/containerd/containerd"
	"github.com/distribution/distribution/reference"
	"github.com/spf13/cobra"

	"github.com/chaitin/libveinmind/go/cmd"
	veinmindContainerd "github.com/chaitin/libveinmind/go/containerd"
	"github.com/chaitin/libveinmind/go/plugin"
	"github.com/chaitin/libveinmind/go/plugin/log"
	"github.com/chaitin/libveinmind/go/plugin/service"
)

func scanImages(ctx context.Context, reporefs []string) error {
	client, err := containerd.New(
		"/run/containerd/containerd.sock",
		containerd.WithDefaultNamespace("hello-host"))
	if err != nil {
		return err
	}
	defer client.Close()
	var imageIDs []string
	for _, reporef := range reporefs {
		if named, err := reference.ParseDockerRef(reporef); err == nil {
			reporef = named.String()
		}
		log.Infof("pulling image %q", reporef)
		image, err := client.Pull(ctx, reporef, containerd.WithPullUnpack)
		if err != nil {
			log.Errorf("cannot pull image %q: %v", reporef, err)
			continue
		}
		log.Infof("pulled image %q", reporef)
		imageID := "hello-host/" + string(image.Target().Digest)
		imageIDs = append(imageIDs, imageID)
	}

	plugins, err := plugin.DiscoverPlugins(ctx, ".")
	if err != nil {
		return err
	}
	c, err := veinmindContainerd.New()
	if err != nil {
		return err
	}
	defer c.Close()
	return cmd.ScanImageIDs(ctx, plugins, c, imageIDs,
		plugin.WithExecInterceptor(func(
			ctx context.Context, plug *plugin.Plugin, c *plugin.Command,
			next func(context.Context, ...plugin.ExecOption) error,
		) error {
			reg := service.NewRegistry()
			reg.AddServices(log.WithFields(log.Fields{
				"plugin":  plug.Name,
				"command": path.Join(c.Path...),
			}))
			return next(ctx, reg.Bind())
		}),
	)
}

var rootCmd = &cobra.Command{
	Use:   "hello-host",
	Short: "Our first host program to scan images from repositories",
}

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "scan",
		Short: "scan images from repositories",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return scanImages(c.Context(), args)
		},
	})
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
```

执行指令 `go mod init hello-host && go build -o hello-host ./hello-host.go` 进行编译 [^5]，得到我们的第一个宿主程序 `hello-host`。

将 `hello-host` 和 `hello-plugin` 置于同一目录，然后执行 `sudo ./hello-host scan` 加镜像名，譬如 `sudo ./hello-host scan nginx redis golang:1.16`，即可调用 `hello-plugin` 插件对远程仓库进行扫描。

直接调用 Containerd 的 API 清空命名空间的代码比较复杂，因此为了我们理解和说明的方便，没有在宿主程序中包含扫描完毕后清空命名空间的代码。读者可以直接执行以下代码来清空 `hello-host` 创建的命名空间：

```bash
sudo ctr -n hello-host images ls -q | xargs sudo ctr -n hello-host images rm
sudo ctr -n hello-host content ls -q | xargs sudo ctr -n hello-host content rm
sudo ctr namespace rm hello-host
```

尽管 `hello-host` 的代码比起 `hello-plugin` 的代码来说长多了，但是我们可以很容易发现其核心逻辑位于 `scanImages` 函数中。`scanImages` 函数包含两部分，前一部分依据参数列表中指定的镜像名称，使用 Containerd 的 API 拉取镜像，并记录镜像 ID 到列表 `imageIDs` 中；后一部分则发现当前工作目录中存在的插件，并调用问脉 SDK 的 `ScanImageIDs` 函数驱动这些插件对 `imageIDs` 中指定的镜像进行扫描。让我们把注意力集中到与插件相关的后一部分。

首先要解决的一个问题是，尽管我们发现了 `hello-plugin` 插件并成功执行了其中的 `scan` 函数，但是我们尚不清楚什么样的可执行文件中，什么样的函数会被成功识别并执行。譬如将 `hello-host` 与 `hello-plugin` 置于当前工作目录下时，`hello-host` 中的 `scanImages` 函数就没有被执行。而如果尝试将 `hello-plugin` 中的 `scan` 函数重命名为 `scan_images` 函数，再以同样的指令执行 `hello-host`，重命名后的 `scan_images` 函数依然会被执行。这说明我们执行插件中的函数并不是依靠简单的名字匹配，而是有更高阶的规则去发现待执行的候选函数。

直接执行 `hello-plugin`，我们可以发现其下有两个子命令：`info` 和 `scan`（重命名为 `scan_images` 函数后为 `scan-images` 子命令，以下略）。其中 `info` 并没有被我们手动定义过，而执行 `info` 子命令可以看到以下输出：

```bash
$ sudo ./hello-plugin info
{"manifestVersion": 1, "name": "hello-plugin", "version": "1.0.0", "author": "", "description": "", "tags": [], "commands": [{"type": "image", "data": {}, "path": ["scan"]}]}
```

在这里我们可以明显地看到先前通过 `@command.image_command` 对函数进行标注的成果：插件通过问脉 SDK 对函数进行标注，除了简化代码以外，更重要的是在代码层面确定了每个函数的功能语义。如在 `hello-plugin` 中定义的 `scan` 函数被 `@command.image_command` 标注确定为是针对镜像的扫描函数，那么在 `hello-host` 执行 `cmd.ScanImageIDs` 函数扫描指定镜像时理应调用该函数，这也是我们最终观察到的运行结果。

宿主程序 `hello-host` 与插件 `hello-plugin` 分别使用 Go 和 Python 语言编写，具有不同的语言运行时。而利用子命令的方式可以将语言的差异封装起来，对外仅在进程级别暴露相互调用的接口，对内通过问脉 SDK 的代码约束了各子命令的参数传递与解析方式，最终不同语言编写的宿主程序和插件得以协同运行。

至此，宿主程序发现并调用插件中的函数的方式就一目了然了：

1. 插件编写不同类型的扫描函数并对它们进行分别标注，而 `info` 子命令利用标注信息生成元数据；
2. 宿主程序对每一个放置在插件目录中的可执行文件，调用 `info` 子命令获取插件的元数据并解析，生成插件列表；
3. 宿主程序识别并生成相应的扫描对象，并依据插件元数据调用插件的对应函数，实现可扩展的扫描。

然后另一个问题是，在宿主程序 `hello-host` 调用 `cmd.ScanImageIDs` 时，指定了一个显眼到不可忽略的 `plugin.WithExecInterceptor` 参数，其作用是什么，是否每次调用都需要？

在问脉 SDK 中，`plugin.WithExecInterceptor` 参数允许我们通过类似职责链的方式为每次插件的执行添加一些额外行为，如初始化某一资源，配置运行环境等，其风格和原理类似于 gRPC 的 [`WithUnaryInterceptor`](https://pkg.go.dev/google.golang.org/grpc#WithUnaryInterceptor) 函数。

在 `hello-host` 的例子中，我们在 `plugin.WithExecInterceptor` 使用 `service.NewRegistry` 初始化了一个 Registry，并且往其中添加了一个通过 `log.WithFields` 初始化的服务。而在运行时我们也发现，`hello-plugin` 通过 `log.info` 记录的日志被转发到了 `hello-host` 中，并且在往每条日志添加 `log.Fields` 中的指定的字段后，作为 `hello-host` 本身的日志打印出来。

这便是问脉 SDK 插件系统的的另一重要机制：服务。在问脉 SDK 中，宿主程序只需要将能访问的服务注册到 Registry 中进行索引，即可向插件提供服务，问脉 SDK 为宿主程序和插件处理了与服务相关的诸多细节，如进程间通信、路由、服务调度等；同时，作为常用服务之一，问脉 SDK 也提供了日志服务，以便宿主接收各插件等运行日志，并且在宿主程序中进行过滤、归档和轮转等处理，它也为其他类型服务的编写提供了参考。

[^1]: 事实上很多语言的标准库在遇到运行时异常的情况下，会直接以非零状态码终止当前进程，并且何时终止进程和以何种状态码终止进程，大多数情况下都不受用户控制。
[^2]: 基于这个原因，问脉 SDK 计划上也不会提供远程仓库扫描的功能，用户应该使用问脉开源工具集中的 `veinmind-runner` 指令驱动其他工具 / 插件进行远程仓库扫描。
[^3]: 显然地，通过提供服务，我们除了可以接收和处理数据外还能进行很多别的操作，包括在插件中受控地获取和修改某项系统配置等，从语义上也应该接受这样的用法，但为了说明方便对此不作过多赘述。
[^4]: 问脉 SDK 工作时需要有访问 Docker 根目录等关键目录和文件的权限，为了简化说明，我们直接以 root 用户执行插件程序，在执行前需要为 root 用户安装 Python Binding。
[^5]: 自 containerd v1.6.0 开始使用了 Go 1.17 的 API，若本地安装的 Go 版本过低可能会导致编译错误，此时将本地的 `go.mod` 中 containerd 的版本修改为 v1.5.0 方可编译通过。
