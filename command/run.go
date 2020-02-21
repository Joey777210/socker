package command

import (
	"Socker/container"
	log "github.com/sirupsen/logrus"
	"os"
)

//called by runCommand
func Run(tty bool, command string){
	//gets the command
	parent := container.NewParentProcess(tty, command)
	//parent command starts to manipulate
	if err := parent.Start(); err != nil{
		log.Error(err)
	}
	log.Print("socker: exit container")
	parent.Wait()
	os.Exit(-1)
}
