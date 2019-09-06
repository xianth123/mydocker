package main

import (
	log "github.com/sirupsen/logrus"
	"mydocker/container"
	"fmt"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name: "run",
	Usage: `Create a container whit namespace and cgroups limit
			mydocker run -it [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name: "ti",
			Usage: "enable tty",
		},
	},

	Action: func(context *cli.Context) error{
		if len(context.Args()) < 1{
			return fmt.Errorf("miss container command")
		}
		var cmdArray []string
		for _, arg := range context.Args(){
			cmdArray = append(cmdArray, arg)
		}
		tty := context.Bool("ti")
		Run(tty, cmdArray)
		return nil
	},
}

var initCommand = cli.Command{
	Name:         "init",
	Usage:        `Init container process`,
	Action: func(ctx *cli.Context) error {
		log.Infof("init container")
		err := container.RunContainerInitProcess()
		return err
	},
}

