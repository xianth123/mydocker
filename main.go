package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

const usage = `This is a simple docker container runtime implementation
				enjoy it, good luck!`

// 设置一个 app 应用，初始化两条命令


func main()  {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = usage
	app.Commands = []cli.Command{
		initCommand,
		runCommand,
	}
	app.Before = func(context *cli.Context) error {
		log.SetFormatter(&log.JSONFormatter{}) // 设置日志为 json 类型
		log.SetOutput(os.Stdout)	// 设置输出为os 输出
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
