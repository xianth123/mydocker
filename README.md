# MyDocker

用 Go 从零实现的简化版容器运行时，用于学习 Linux 容器背后的核心机制。

项目通过 Linux namespace 隔离进程环境，使用 cgroups v1 限制资源，通过 AUFS 构造可写文件系统，并在容器进程中执行 `pivot_root` 和 `/proc` 挂载。它更接近一个容器原理实验，而不是 Docker 的生产级替代品。

> [!WARNING]
> 本项目基于较早期的 Linux 容器技术栈，目标环境为 Go 1.12、cgroups v1 和 AUFS。代码仅供学习与实验，请勿用于生产环境，也不要在包含重要数据的主机上直接运行。

## 已实现功能

- UTS、PID、Mount、Network 和 IPC namespace 隔离
- 父进程通过管道向容器 init 进程传递命令
- BusyBox rootfs 与 AUFS 读写层
- `pivot_root` 切换容器根文件系统
- 挂载独立的 `/proc`
- cgroups v1 内存、CPU 和 CPUSet 子系统
- 宿主机目录挂载到容器
- 容器退出后的挂载点和可写层清理

## 工作流程

```text
mydocker run
    ├── 创建 namespace 隔离的子进程
    ├── 准备 BusyBox + AUFS 工作目录
    ├── 创建并应用 cgroups 资源限制
    ├── 通过管道发送用户命令
    └── mydocker init
          ├── pivot_root
          ├── 挂载 /proc
          └── syscall.Exec 执行容器命令
```

## 环境要求

- Linux（依赖 Linux 专有的 namespace、mount 和 cgroup API）
- root 权限
- Go 1.12 左右的历史环境
- cgroups v1
- 支持 AUFS 的内核，并已安装 AUFS 挂载工具
- 一个解压到 `/root/busybox` 的 BusyBox rootfs

现代发行版通常默认使用 cgroups v2 和 OverlayFS，可能无法直接运行本项目。

## 构建

源码中的包路径采用早期 GOPATH 风格。建议在兼容的 Go 1.12 环境中将仓库放到 `$GOPATH/src/mydocker`：

```bash
mkdir -p "$GOPATH/src"
git clone https://github.com/xianth123/mydocker.git "$GOPATH/src/mydocker"
cd "$GOPATH/src/mydocker"

GO111MODULE=off go get github.com/sirupsen/logrus
GO111MODULE=off go get github.com/urfave/cli
GO111MODULE=off go build -o mydocker .
```

## 准备 rootfs

项目运行时会直接使用 `/root/busybox` 作为只读基础层。请先准备与宿主机架构匹配的 BusyBox rootfs：

```bash
sudo mkdir -p /root/busybox
sudo tar -xf busybox.tar -C /root/busybox
```

确认 rootfs 中至少存在准备执行的命令，例如：

```bash
sudo test -x /root/busybox/bin/sh
```

## 使用方法

查看命令帮助：

```bash
./mydocker --help
./mydocker run --help
```

启动交互式容器：

```bash
sudo ./mydocker run --ti /bin/sh
```

运行单条命令：

```bash
sudo ./mydocker run /bin/echo hello-from-mydocker
```

限制内存和 CPU 核心：

```bash
sudo ./mydocker run --ti --m 100m --cpuset 0 /bin/sh
```

挂载宿主机目录：

```bash
sudo mkdir -p /tmp/mydocker-data
sudo ./mydocker run --ti --v /tmp/mydocker-data:/data /bin/sh
```

## 命令参数

| 参数 | 说明 | 示例 |
| --- | --- | --- |
| `--ti` | 将当前终端连接到容器进程 | `--ti` |
| `--m` | cgroups v1 内存限制 | `--m 100m` |
| `--cpuset` | 指定容器可使用的 CPU 核心 | `--cpuset 0` |
| `--cpushare` | CPU shares 参数；当前实现仍需进一步修正和验证 | `--cpushare 512` |
| `--v` | 挂载目录，格式为 `宿主机路径:容器路径` | `--v /tmp/data:/data` |

## 项目结构

```text
.
├── main.go                         # CLI 入口
├── main_command.go                 # run 和 init 子命令
├── run.go                          # 容器启动、cgroup 应用和清理流程
├── container/
│   ├── container_process.go        # namespace、AUFS、volume 和工作目录
│   └── init.go                     # pivot_root、/proc 挂载和命令执行
├── cgroups/
│   ├── cgroup_manager.go           # cgroup 生命周期管理
│   └── subsystems/                 # memory、cpu、cpuset 子系统
└── src/ch2/                        # namespace、mount 和 cgroup 独立实验
```

## 已知限制

- 仅支持 Linux，并且必须以 root 权限运行。
- 文件系统路径写死为 `/root/busybox`、`/root/writeLayer` 和 `/root/mnt`。
- 使用 AUFS，未适配现代发行版常见的 OverlayFS。
- cgroup 实现面向 v1，未适配 cgroups v2。
- 没有镜像管理、网络配置、容器持久化、日志驱动或安全策略。
- 命令通过空格拆分，不支持完整的 shell 引号与转义语义。
- `--cpushare` 的参数传递和 CPU 子系统写入逻辑仍需修正。
- 错误处理和资源清理仍有实验性质，不保证异常退出后的主机状态。

## 学习建议

建议按照下面的顺序阅读源码：

1. `main.go`：了解 CLI 入口。
2. `main_command.go` 和 `run.go`：跟踪 `run` 到 `init` 的调用过程。
3. `container/container_process.go`：理解 namespace 和 AUFS 工作目录。
4. `container/init.go`：理解 `pivot_root`、`/proc` 和 `syscall.Exec`。
5. `cgroups/`：理解资源限制如何应用到容器进程。
6. `src/ch2/`：分别运行 namespace 与 cgroup 小实验。

## 后续改进方向

- 修复 Go Modules 包路径并升级依赖
- 支持 cgroups v2
- 使用 OverlayFS 替代 AUFS
- 移除固定目录，改为可配置的运行时目录
- 补充单元测试和集成测试
- 完善错误传播和异常清理
- 增加容器状态与生命周期管理

## 致谢

该项目用于学习 Go 与 Linux 容器运行时原理。欢迎通过 Issue 或 Pull Request 交流改进建议。
