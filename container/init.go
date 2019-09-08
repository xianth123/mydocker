package container

import (
	"fmt"
	"os/exec"
	"path/filepath"

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

	setUpMount()
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

// Init 挂载点
func setUpMount()  {
	pwd ,err := os.Getwd()
	if err != nil {
		fmt.Println("Get current location error %v", err)
		return
	}
	fmt.Println("current location is %s", pwd)
	pivotRoot(pwd)

	// mount proc
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}


func pivotRoot(root string) error {
	// 将 root 重新 mount 一次， 然后，将 root 使用 bind 换一个挂载的地点。然后可以解除 root 的挂载
	// func Mount(source string, target string, fstype string, flags uintptr, data string) (err error)
	// #define MS_BIND      4096
	///* 把一个已挂载文件系统移动到另一个挂载点，相当于先执行卸载，然后将文件系统挂载在另外的一个目录下 */
	// #define MS_REC       16384 /* 为目录子树递归的创建绑定挂载 */
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC); err != nil {
		return fmt.Errorf("Mount rootfs to itself error: %v", err)
	}
	// 创建 rootfs/.pivot_root 存储 old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return nil
	}
	// pivot_root 到显得 rootfs ，old_root 是挂载到 rootfs/.pivot_root
	// 挂载点现在可以在 mount 命令中看到
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	// umount rootfs/.pivot_root
	if err := syscall.Umount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("umount pivot_root fail %v", err)
	}
	return os.Remove(pivotDir)


}
