package container

import (
	"fmt"
	"os/exec"

	//	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"syscall"
	"io/ioutil"
//	log "github.com/sirupsen/logrus"
)

// 设置 mount 参数， 进行 mount 文件的挂载。

// 通过 syscall 执行 命令， 传入 参数和环境.

// 先启动进程，然后进行挂载。

func RunContainerInitProcess() error {
	cmdArray := readUserCommand()
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("Run container get user command error, cmdArray is nil")
	}
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	// 寻找命令的路径
	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
	//	log.Errorf("exec loop path error %v", err)
	}
//	log.Info("Find path %s", path)
	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
	//	log.Errorf(err.Error())
	}
	return nil
}

// 从pipe 中读取命令
func readUserCommand() []string  {
	// 保证管道文件是进程的第四个文件描述符
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
	//	log.Errorf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}
