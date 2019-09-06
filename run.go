package main

import (
	"mydocker/container"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

// 是否改变输出
// 命令

// init 做一些初始化工作
// 给执行的命令 包装 namespace
func Run(tty bool, comArray []string)  {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil{
		log.Error(err)
	}
	sendInitCommand(comArray, writePipe)
	parent.Wait()
	os.Exit(-1)
}

func sendInitCommand(comArray []string, writePipe *os.File)  {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}

