package main

import (
	log "github.com/sirupsen/logrus"
	"mydocker/cgroups/subsystems"
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
		cli.StringFlag{
			Name: "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name: "cpushare",
			Usage: "cpushare limit",
		},
		cli.StringFlag{
			Name: "cpuset",
			Usage: "cpuset limit",
		},
		cli.StringFlag{
			Name: "v",
			Usage: "volume",
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
		resConf := &subsystems.ResourceCofing{
			MemoryLimit: context.String("m"),
			CpuSet: context.String("cpuset"),
			CpuShare: context.String("cpuset"),
		}
		log.Infof(" 666666666666 %v", resConf)
		volume := context.String("v")
		Run(tty, cmdArray, resConf, volume)
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

