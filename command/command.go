package command

import (
	"Socker/container"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var RunCommand = cli.Command{
	Name:	"run",
	Usage:	`create a new container with namespace and cgroups limit: socker run -ti [command]`,
	Flags:	[]cli.Flag{
		&cli.BoolFlag{
			Name:  "ti",	//open stdin/stdout tunnel
			Usage: "enable tty",
		},
	},

	//get command behind -ti if there is
	//call Run function to build a container
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		cmd := context.Args().Get(0)
		tty := context.Bool("ti")
		Run(tty, cmd)
		return nil
	},

}

var InitCommand = cli.Command{
	Name:	"init",
	Usage:	``,

	Action: func(context *cli.Context) error{
		log.Infof("init come on")
		cmd := context.Args().Get(0)
		err := container.InitProcess(cmd ,context.Args())
		return err
	},
}
