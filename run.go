package main

import (
	"log"
	"mydocker/container"
	"os"
)

// 是否改变输出
// 命令

// init 做一些初始化工作
// 给执行的命令 包装 namespace
func Run(tty bool, command string)  {
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil{
		log.Fatal(err)
	}
	parent.Wait()
	os.Exit(-1)

}
