package container

import (
//	"github.com/sirupsen/logrus"
	"os"
	"syscall"
)

// 设置 mount 参数， 进行 mount 文件的挂载。

// 通过 syscall 执行 命令， 传入 参数和环境.

// 先启动进程，然后进行挂载。
func RunContainerInitProcess(command string, args []string) error {
//	logrus.Infof("command %s", command)

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
//		logrus.Errorf(err.Error())
	}
	return nil
}
