package command

import (
	"Socker/cgroup"
	//"Socker/cgroup"
	"Socker/container"
	"Socker/network"
	"Socker/overlay2"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

//attention! this is v1 version of cli
var RunCommand = cli.Command{
	Name:  "run",
	Usage: `create a new container with namespace and cgroups limit: socker run -ti [command]`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "ti", //open stdin/stdout tunnel
			Usage: "enable tty",
		},
		&cli.StringFlag{
			Name:  "m",
			Usage: "limit memory usage",
		},
		&cli.StringFlag{
			Name:  "cpushare",
			Usage: "limit cpushare usage",
		},
		&cli.StringFlag{
			Name:  "cpuset",
			Usage: "limit cpuset usage",
		},
		&cli.BoolFlag{
			Name:  "d",
			Usage: "detach container",
		},
		&cli.BoolFlag{
			Name:  "mqtt",
			Usage: "open mqtt sub and pubsss",
		},
		&cli.StringFlag{
			Name:  "name",
			Usage: "container name",
		},
		&cli.StringFlag{
			Name:  "net",
			Usage: "container network",
		},
		&cli.StringSliceFlag{
			Name: "p",
			Usage: "port mapping",
		},
		&cli.StringSliceFlag{
			Name: "e",
			Usage: "set environment",
		},
	},
	Action: func(context *cli.Context) error {
		if context.Args().Len() < 1 {
			return fmt.Errorf("Missing container command")
		}
		var cmdArray []string
		for i := 0; i < context.Args().Len(); i++ {
			arg := context.Args().Get(i)
			cmdArray = append(cmdArray, arg)
		}
		imageName := cmdArray[0]
		cmdArray = cmdArray[1:]

		createTty := context.Bool("ti")
		detach := context.Bool("d")
		mqtt := context.Bool("mqtt")
		envSlice := context.StringSlice("e")

		if createTty && detach {
			return fmt.Errorf("ti and d paramter can not both provided")
		}

		resConf := &cgroup.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuSet:      context.String("cpuset"),
			CpuShare:    context.String("cpushare"),
		}
		log.Infof("createTty %v", createTty)
		containerName := context.String("name")
		network := context.String("net")

		portmapping := context.StringSlice("p")

		Run(createTty, cmdArray, resConf, containerName, network, portmapping, mqtt, imageName, envSlice)
		return nil
	},
}

var InitCommand = cli.Command{
	Name:  "init",
	Usage: `Init comtainer`,

	Action: func(context *cli.Context) error {
		log.Infof("init come on")
		err := container.InitProcess()
		return err
	},
}

var CommitCommand = cli.Command{
	Name:  "commit",
	Usage: `commit a container into image`,

	Action: func(context *cli.Context) error {
		if context.Args().Len() < 1 {
			return fmt.Errorf("Missing image name when commit")
		}
		imageName := context.Args().Get(0)
		overlay2.CommitContainer(imageName)
		return nil
	},
}

var ListCommand = cli.Command{
	Name:  "ps",
	Usage: `list all the containers`,
	Action: func(context *cli.Context) error {
		container.ListContainers()
		return nil
	},
}

var LogCommand = cli.Command{
	Name:  "logs",
	Usage: `print logs of a container`,
	Action: func(context *cli.Context) error {
		if context.Args().Len() < 1 {
			return fmt.Errorf("Please input your container name")
		}
		containerName := context.Args().Get(0)
		container.LogContainer(containerName)
		return nil
	},
}

var ExecCommand = cli.Command{
	Name:  "exec",
	Usage: `exec a command into container`,
	Action: func(context *cli.Context) error {
		//the second call
		if os.Getenv(container.ENV_EXEC_PID) != "" {
			log.Infof("%d", os.Getgid())
			log.Infof("pid callback pid %d", os.Getgid())
			return nil
		}

		//	./socker exec containerName command
		if context.Args().Len() < 2 {
			return fmt.Errorf("Missing container name or command")
		}
		containerName := context.Args().Get(0)
		var commandArray []string
		for _, arg := range context.Args().Tail() {
			commandArray = append(commandArray, arg)
		}
		container.ExecContainer(containerName, commandArray)
		return nil
	},
}

var StopCommand = cli.Command{
	Name:  "stop",
	Usage: `stop a container`,
	Action: func(context *cli.Context) error {
		if context.Args().Len() < 1 {
			return fmt.Errorf("Missing container name")
		}
		containerName := context.Args().Get(0)
		container.StopContainer(containerName)
		return nil
	},
}

var RemoveCommand = cli.Command{
	Name:  "rm",
	Usage: `remove unused containers`,
	Action: func(context *cli.Context) error {
		if context.Args().Len() < 1 {
			return fmt.Errorf("Missing container name")
		}
		containerName := context.Args().Get(0)
		container.RemoveContainer(containerName)
		return nil
	},
}

var NetworkCommand = cli.Command{
	Name:  "network",
	Usage: `set network for a container`,
	Subcommands: []*cli.Command{
		{
			Name:  "create",
			Usage: "create a container network",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "driver",
					Usage: "make a network driver",
				},
				&cli.StringFlag{
					Name:  "subnet",
					Usage: "set subnet IP and mask e.g. 192.168.0.1/24",
				},
			},

			Action: func(context *cli.Context) error {
				if context.Args().Len() < 1 {
					return fmt.Errorf("Missing network command")
				}
				driverName := context.String("driver")
				subnet := context.String("subnet")
				networkName := context.Args().Get(0)
				err := network.Init()
				if err != nil {
					log.Errorf("init network %s error %v", networkName, err)
				}
				err = network.CreateNetwork(driverName, subnet, networkName)
				if err != nil {
					log.Errorf("create network %s error %v", networkName, err)
				}
				return nil
			},
		},		{
			Name: "list",
			Usage: "list container network",
			Action:func(context *cli.Context) error {
				network.Init()
				network.ListNetwork()
				return nil
			},
		},
		{
			Name: "remove",
			Usage: "remove container network",
			Action:func(context *cli.Context) error {
				if context.Args().Len() < 1 {
					return fmt.Errorf("Missing network name")
				}
				network.Init()
				err := network.DeleteNetwork(context.Args().Get(0))
				if err != nil {
					return fmt.Errorf("remove network error: %+v", err)
				}
				return nil
			},
		},
	},
}

var ImageCommand = cli.Command{
	Name:  "image",
	Usage: `show imgages and delete image`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "ls",
			Usage: "list all images",
		},
		&cli.StringFlag{
			Name:  "rm",
			Usage: "remove a image by name",
		},
	},
	Action: func(context *cli.Context) error{
		ls := context.Bool("ls")
		if ls {
			err := container.ImageLs()
			return err
		}
		imageName := context.String("rm")
		err := container.ImageRemove(imageName)
		if err != nil {
			return err
		}
		return nil
	},
}
