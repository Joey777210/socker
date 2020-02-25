package command

import (
	"Socker/cgroup"
	"Socker/container"
	"Socker/overlay2"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

//called by runCommand
func Run(tty bool, command []string, resourceConfig cgroup.ResourceConfig){
	//gets the command
	parent, writePipe:= container.NewParentProcess(tty)
	//parent command starts to manipulate
	if err := parent.Start(); err != nil{
		log.Error(err)
	}

	//use socker-cgroup as cgroup name
	//create cgroup manager and set res
	cgroupManager := cgroup.NewCgroupManager("socker-cgroup")
	//defer cgroupManager.Destroy()
	//set res
	cgroupManager.Set(&resourceConfig)
	//set container into cgroup
	cgroupManager.Apply(parent.Process.Pid)
	//init container
	sendInitCommand(command, writePipe)
	log.Print("socker: exit container")
	parent.Wait()
	overlay2.DeleteWorkSpace("/root", "/root/mergeDir")
	os.Exit(-1)
}

func sendInitCommand(command []string, writePipe *os.File) {
	cmdStr := strings.Join(command, " ")
	log.Infof("command init write is %s", cmdStr)
	writePipe.WriteString(cmdStr)
	writePipe.Close()
}
