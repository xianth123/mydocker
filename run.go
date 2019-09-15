package main

import (
	"mydocker/cgroups"
	"mydocker/cgroups/subsystems"
	"mydocker/container"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

// 是否改变输出
// 命令

// init 做一些初始化工作
// 给执行的命令 包装 namespace
func Run(tty bool, comArray []string, res *subsystems.ResourceCofing, volume string)  {
	parent, writePipe := container.NewParentProcess(tty, volume)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil{
		log.Error(err)
	}
	//
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(comArray, writePipe)
	parent.Wait()
	mntURL := "/root/mnt"
	rootURL := "/root"
	container.DeleteWordSpace(rootURL, mntURL, volume)
	os.Exit(-1)
}

func sendInitCommand(comArray []string, writePipe *os.File)  {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}

