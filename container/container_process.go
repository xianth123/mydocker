package container

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
)

// 返回一个包装了 namespace 的命令描述，和往这个命令写入参数的 write pipe 文件。
func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Errorf("New pipe error %v", err)
		return nil, nil
	}
	cmd := exec.Command("/proc/self/exe", "init") //返回一个命令
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.ExtraFiles = []*os.File{readPipe}  // 子进程继承读 pipe 文件。
	                                      // 标准输入，标准输出，错误不能读取，所以，readPipe 的文件描述符是3
	return cmd, writePipe				// 返回一个 写 pipe 文件进行命令的输入
}

func NewPipe() (*os.File, *os.File, error)  {
	read, write, err := os.Pipe()
	if err != nil{
		return nil, nil, err
	}
	return read, write, err
}
