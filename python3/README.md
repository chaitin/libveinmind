# 问脉 Python Binding

问脉 Python Binding 定义了问脉 SDK 在 Python 语言中接口，使得在 Python 语言中使用问脉 SDK 成为可能。

您需要先安装问脉 SDK 以开始使用 Python Binding，安装方式可参考[本仓库的 README 文档](https://github.com/chaitin/libveinmind)。

问脉 Python Binding 已经上传到 pypi.org，因此可以直接执行以下指令拉取最新版本的 Python 包：

```bash
pip install veinmind
```

或者您也可以通过手动构建的方式，在本目录执行以下指令手动构建并安装问脉 Python Binding（需安装对应的 Python 构建工具）：

```bash
python -m build
pip install dist/veinmind-*.whl
```

安装完成后即可在 Python 代码中通过 `import veinmind` 的方式导入并开始使用问脉 SDK。
