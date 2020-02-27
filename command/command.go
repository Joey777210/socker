package command

import (
	"Socker/cgroup"
	"Socker/container"
	"Socker/overlay2"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

//attention! this is v1 version of cli
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
		&cli.BoolFlag{
			Name:        "d",
			Usage:       "detach container",
		},
		&cli.StringFlag{
			Name:        "name",
			Usage:       "container name",
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
		detach := context.Bool("d")

		//tty and detach can not both exist
		if tty && detach {
			return fmt.Errorf("ti and d paramter can not both provideed")
		}

		resourceConfig := cgroup.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuShare:    context.String("cpuset"),
			CpuSet:      context.String("cpushare"),
		}

		log.Infof("createTty %v", tty)
		//get container name and pass on
		containerName := context.String("name")
		Run(tty, cmdArray, resourceConfig, containerName)
		return nil
	},

}

var InitCommand = cli.Command{
	Name:	"init",
	Usage:	`Init comtainer`,

	Action: func(context *cli.Context) error{
		log.Infof("init come on")
		err := container.InitProcess()
		return err


	},
}

var CommitCommand = cli.Command{
	Name:	"commit",
	Usage:	`commit a container into image`,

	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing image name when commit")
		}
			imageName := context.Args().Get(0)
			overlay2.CommitContainer(imageName)
			return nil
	},
}

var ListCommand = cli.Command{
	Name:	"ps",
	Usage:	`list all the containers`,
	Action:	func(context *cli.Context) error {
		container.ListContainers()
		return nil
	},
}

var LogCommand = cli.Command{
	Name:	"logs",
	Usage:	`print logs of a container`,
	Action:	func(context *cli.Context) error {
		if len (context.Args()) < 1{
			return fmt.Errorf("Please input your container name")
		}
		containerName := context.Args().Get(0)
		container.LogContainer(containerName)
		return nil
	},
}