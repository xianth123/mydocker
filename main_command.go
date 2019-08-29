package main

import (
	log "github.com/Sirupsen/logrus"
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
		cmd := context.Args().Get(0)
		tty := context.Bool("ti")
		Run(tty, cmd)
		return nil
	},
}

var initCommand = cli.Command{
	Name:         "init",
	Usage:        `Init container process`,
	Action: func(ctx *cli.Context) error {
		log.Infof("init container")
		cmd := ctx.Args().Get(0)
		log.Infof("command %s", cmd)
		err := container.RunContainerInitProcess(cmd, nil)
		return err
	},
}

