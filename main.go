package main

import (
	"Socker/command"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

func main(){
	//use opensource project 'cli' to define app and parse flags
	app := cli.NewApp()
	app.Name = "socker"

	app.Commands = []*cli.Command{
		&command.RunCommand,
		&command.InitCommand,
		&command.CommitCommand,
		&command.ListCommand,
		&command.LogCommand,
		&command.ExecCommand,
		&command.StopCommand,
		&command.NetworkCommand,
	}

	//init logrus
	app.Before = func(context *cli.Context) error {
		//set log as JSON formatter
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		return nil
	}
	//output args you had just type-in
	sockerCommand := os.Args
	log.Printf("args: %s", sockerCommand)


	if err := app.Run(os.Args); err != nil{
		log.Fatal(err)
	}
}