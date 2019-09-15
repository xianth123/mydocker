package container

import (
	"fmt"
	"strings"

	//	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
)

// 返回一个包装了 namespace 的命令描述，和往这个命令写入参数的 write pipe 文件。
func NewParentProcess(tty bool, volume string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
//		log.Errorf("New pipe error %v", err)
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
	cmd.Dir = "/root/busybox"             // 设置 cmd 的dir
	mntURL := "/root/mnt/"
	rootURL := "/root/"
	NewWorkSpace(rootURL, mntURL, volume)
	cmd.Dir = mntURL
	return cmd, writePipe				// 返回一个 写 pipe 文件进行命令的输入
}

func NewPipe() (*os.File, *os.File, error)  {
	read, write, err := os.Pipe()
	if err != nil{
		return nil, nil, err
	}
	return read, write, err
}

// 创建 AUFS 联合文件系统
func NewWorkSpace(rootURL string, mntURL string, volume string)  {
//	CreateReadOnlyLayer(rootURL)
	CreateWriteLayer(rootURL)
	CreateMountPoint(rootURL, mntURL)
	if volume != ""{
		volumeURLs := volumeUrlExtract(volume)
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			// 根据输入的 volume 宿主机url docker url 分别进行挂载
			MountVolume(rootURL, mntURL, volumeURLs)
			println("%q", volumeURLs)
		}else {
			println("volume error")
		}
	}
}

func volumeUrlExtract(volume string) ([]string){
	var volumeURLs [] string
	volumeURLs = strings.Split(volume, ":")
	return volumeURLs

}

func MountVolume(rootURL string, mntURL string, volumeURLs []string) {
	// 创建宿主机文件目录
	parentUrl := volumeURLs[0]
	exist, _ := PathExists(parentUrl)
	if exist == false{
		if err := os.Mkdir(parentUrl, 0777); err != nil {
			println(" mkdir parent dir error, %s    %v", parentUrl, err)
		}
	}
	containerUrl := volumeURLs[1]
	containerVolumeURL := mntURL + containerUrl
	if err := os.Mkdir(containerVolumeURL, 0777); err != nil {
		println(" mkdir container dir error, %s    %v", parentUrl, err)
	}
	dirs := "dirs=" + parentUrl
	println("mount dir %s", dirs)
	println("mount cmd mount -t aufs -o %s none %s", dirs, containerVolumeURL)
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerVolumeURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		println("mount volume error , %v", err)
	}
}

// 获取 busybox.tar 的路径，将其解压到，解压后的 busybox 就是只读层
func CreateReadOnlyLayer(rootURL string)  {
	busyboxURL := rootURL + "busybox/"
	busyboxTarURL := rootURL + "busybox.tar"
	exist, err := PathExists(busyboxTarURL)
	if err != nil {
		fmt.Println("Fail to judeg whether dir %s exitst, %v", busyboxTarURL, err)
	}
	if exist == false {
		if err := os.Mkdir(busyboxURL, 0777); err != nil {
			fmt.Println("mkdir %s error, %v", busyboxURL, err)
		}
		if _, err := exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxURL).CombinedOutput(); err != nil {
			fmt.Println("untar %s error, %v", busyboxURL, err)
		}
	}

}

// 创建一个文件夹作为只写层，将这个文件夹
func CreateWriteLayer(rootURL string) {
	writeURL := rootURL + "writeLayer/"
	if err := os.Mkdir(writeURL, 0777); err != nil {
		fmt.Println("create write layer error, %v", err)

	}
}

//创建一个挂载点
func CreateMountPoint(rootURL string, mntURL string){
	if err := os.Mkdir(mntURL, 0777); err != nil {
		fmt.Println("mkdir  error, %v", err)
	}
	dirs := "dirs=" + rootURL + "writeLayer:" + rootURL + "busybox"
	println("mount dir %s", dirs)
	println("mount cmd mount -t aufs -o %s none %s", dirs, mntURL)
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("cmd run error, %v", err)
	}
}

func DeleteWordSpace(rootURL string, mntURL string, volume string) {
	if (volume != ""){
		volumeUrls := volumeUrlExtract(volume)
		length := len(volumeUrls)
		if length == 2 && volumeUrls[0] != "" && volumeUrls[1] != "" {
			DeleteMountPointWithVolume(rootURL, mntURL, volumeUrls)
		} else {
			DeleteMountPoint(rootURL, mntURL)
		}
	} else {
		DeleteMountPoint(rootURL, mntURL)
	}
	DeleteWriteLayer(rootURL)
}

func DeleteMountPointWithVolume(rootURL string, mntURL string, volumeUrls []string)  {
	containerUrl := mntURL + volumeUrls[1]
	cmd := exec.Command("umount", containerUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		println("umont volume failed %v", err)
	}
	DeleteMountPoint(rootURL, mntURL)
}

func DeleteMountPoint(rootURL string, mntURL string)  {
	cmd := exec.Command("umount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("cmd run error, %v", err)
	}
	if err := os.RemoveAll(mntURL); err != nil {
		fmt.Println("remove error, %v", err)
	}
}

func DeleteWriteLayer(rootURL string) {
	writeURL := rootURL + "writeLayer/"
	if err := os.RemoveAll(writeURL); err != nil {

	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err

}
