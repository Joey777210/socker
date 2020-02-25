package command

import (
	"Socker/cgroup"
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
		&cli.StringFlag{
			Name:        "m",
			Usage:       "limit memory usage",
		},
		&cli.StringFlag{
			Name:        "cpushare",
			Usage:       "limit cpushare usage",
		},
		&cli.StringFlag{
			Name:        "cpuset",
			Usage:       "limit cpuset usage",
		},
	},

	//get command behind -ti if there is
	//call Run function to build a container
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}
		tty := context.Bool("ti")
		resourceConfig := cgroup.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuShare:    context.String("cpuset"),
			CpuSet:      context.String("cpushare"),
		}

		Run(tty, cmdArray, resourceConfig)
		return nil


	},

}

var InitCommand = cli.Command{
	Name:	"init",
	Usage:	``,

	Action: func(context *cli.Context) error{
		log.Infof("init come on")
		err := container.InitProcess()
		return err


	},
}
