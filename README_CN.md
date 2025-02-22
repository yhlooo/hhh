**[简体中文](README_CN.md)** | [English](README.md)

---

# hhh - Human readable HTML Helper

`hhh` 是一个 CLI 工具，用于通过 HTTP 、读取本地文件等方式获取 HTML ，降噪并提取 HTML 页面中的文章内容，转换为人类可读的 Markdown 等格式文本。

**🚧 该项目尚未完成 🚧**

## 安装

### Docker

直接使用镜像 [`ghcr.io/yhlooo/hhh`](https://github.com/yhlooo/hhh/pkgs/container/hhh) docker run 即可：

```bash
docker run -it --rm ghcr.io/yhlooo/hhh:latest --help
```

### 通过二进制安装

通过 [Releases](https://github.com/yhlooo/hhh/releases) 页面下载可执行二进制，解压并将其中 `hhh` 文件放置到任意 `$PATH` 目录下。

### 从源码编译

要求 Go 1.23 ，执行以下命令下载源码并构建：

```bash
go install github.com/yhlooo/hhh/cmd/hhh@latest
```

构建的二进制默认将在 `${GOPATH}/bin` 目录下，需要确保该目录包含在 `$PATH` 中。

## 使用

```bash
hhh https://github.com/yhlooo/hhh
```
