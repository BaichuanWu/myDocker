package main

import (
	"fmt"

	"github.com/BaichuanWu/myDocker/container"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name:  "run",
	Usage: `myDocker run -it [cmd]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
	},
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("empty command")
		}
		var cmds []string
		for _, c := range ctx.Args() {
			cmds = append(cmds, c)
		}
		tty := ctx.Bool("it")
		Run(tty, cmds)
		return nil
	},
}

var initCommand = cli.Command{
	Name: "init",

	Action: func(ctx *cli.Context) error {
		log.Infof("Init come on 999999")
		err := container.RunContainerInitProcess()
		return err
	},
}
